// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package models

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
	"unicode"

	"github.com/gocsaf/csaf/v3/csaf"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	// ErrAlreadyInDatabase is returned from ImportDocument if the
	// advisory is already in the database.
	ErrAlreadyInDatabase = errors.New("already in database")
	// ErrNotAllowed is returned from ImportDocument if the
	// TLP restrictions are not met.
	ErrNotAllowed = errors.New("not allowed")
)

// badStrings acts as a detector for decoding errors from UTF-8.
// https://pkg.go.dev/encoding/json#Unmarshal states:
// [...] When unmarshaling quoted strings, invalid UTF-8
// or invalid UTF-16 surrogate pairs are not treated as an error.
// Instead, they are replaced by the Unicode replacement character U+FFFD. [...]
func badStrings(bad *[]string) replacer {
	return func(_ []string, v string) (any, bool) {
		if strings.ContainsRune(v, unicode.ReplacementChar) {
			*bad = append(*bad, v)
		}
		return v, false
	}
}

func storer(value *string, found *bool, path ...string) replacer {
	return func(keys []string, v string) (any, bool) {
		if !*found && slices.Equal(path, keys) {
			*found = true
			*value = v
		}
		return v, false
	}
}

// DocumentStoreChainFunc is a function which is called after the document
// is stored in the database. It is executed in the same transaction
// as the storage of the document itself. Main usage is to store additional
// info about the document.
// If the document is already in the database the function is called with
// an id of 0 and set duplicate flag set.
type DocumentStoreChainFunc func(ctx context.Context, tx pgx.Tx, id int64, duplicate bool) error

// ChainInTx executes a list of in transaction functions.
func ChainInTx(inTxs ...DocumentStoreChainFunc) DocumentStoreChainFunc {
	return func(ctx context.Context, tx pgx.Tx, docID int64, duplicate bool) error {
		for _, inTx := range inTxs {
			if err := inTx(ctx, tx, docID, duplicate); err != nil {
				return err
			}
		}
		return nil
	}
}

// StoreFilename returns a function to store the file name alongside the document.
func StoreFilename(filename string) DocumentStoreChainFunc {
	return func(ctx context.Context, tx pgx.Tx, docID int64, duplicate bool) error {
		if duplicate {
			return nil
		}
		const insertSQL = `UPDATE documents ` +
			`SET filename = $1 ` +
			`WHERE id = $2`
		_, err := tx.Exec(ctx, insertSQL, filename, docID)
		return err
	}
}

// ImportDocument imports a given advisory into the database.
func ImportDocument(
	ctx context.Context,
	conn *pgxpool.Conn,
	r io.Reader,
	actor *string,
	pstlps PublishersTLPs,
	inTx DocumentStoreChainFunc,
	dry bool,
) (int64, error) {
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)

	var document any
	if err := json.NewDecoder(tee).Decode(&document); err != nil {
		return 0, err
	}

	msgs, err := csaf.ValidateCSAF(document)
	if err != nil {
		return 0, fmt.Errorf("schema validation failed: %w", err)
	}
	if len(msgs) > 0 {
		return 0, errors.New("schema validation failed: " + strings.Join(msgs, ", "))
	}
	return ImportDocumentData(ctx, conn, document, buf.Bytes(), actor, pstlps, inTx, dry)
}

// ImportDocumentData imports a given advisory into the database.
func ImportDocumentData(
	ctx context.Context,
	conn *pgxpool.Conn,
	document any,
	raw []byte,
	actor *string,
	pstlps PublishersTLPs,
	inTx DocumentStoreChainFunc,
	dry bool,
) (int64, error) {

	var (
		tlp, tlpOk               = "", false
		publisher, publisherOK   = "", false
		trackingID, trackingIDOK = "", false
	)

	var bad []string
	idxer := indexDocument(document,
		badStrings(&bad),
		storer(&tlp, &tlpOk, "document", "distribution", "tlp", "label"),
		storer(&publisher, &publisherOK, "document", "publisher", "name"),
		storer(&trackingID, &trackingIDOK, "document", "tracking", "id"),
	)

	// Check if there where some string decoding errors.
	if len(bad) > 0 {
		return 0, fmt.Errorf("invalid strings found: %+v", bad)
	}

	if !publisherOK {
		return 0, errors.New("missing /document/publisher/name")
	}

	if !trackingIDOK {
		return 0, errors.New("missing /document/tracking/id")
	}

	if !tlpOk {
		return 0, errors.New("missing /document/distribution/tlp/label")
	}

	if pstlps != nil && !pstlps.Allowed(publisher, TLP(tlp)) {
		return 0, ErrNotAllowed
	}

	if dry {
		return 0, nil
	}

	// Allow only one insert at a time.
	// There are transaction serialization issues with the unique texts.
	// TODO: This has to be investigated!
	globalInsertLock.Lock()
	defer globalInsertLock.Unlock()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	const (
		queryAdvisory        = `SELECT id FROM advisories WHERE (tracking_id, publisher) = ($1, $2)`
		insertAdvisory       = `INSERT INTO advisories (tracking_id, publisher) VALUES ($1, $2) RETURNING id`
		savepointDoc         = `SAVEPOINT insert_document`
		rollbackSavepointDoc = `ROLLBACK TO SAVEPOINT insert_document`
		releaseSavepointDoc  = `RELEASE SAVEPOINT insert_document`
		insertDoc            = `INSERT INTO documents (document, original, advisories_id) VALUES ($1, $2, $3) RETURNING id`
		insertLog            = `INSERT INTO events_log (event, state, actor, documents_id) VALUES ('import_document', 'new', $1, $2)`
	)

	// We need an advisory before we insert a document.
	var (
		advisoryID      int64
		missingAdvisory bool
	)
	switch err := tx.QueryRow(
		ctx, queryAdvisory, trackingID, publisher).Scan(&advisoryID); {
	case errors.Is(err, pgx.ErrNoRows):
		missingAdvisory = true
	case err != nil:
		return 0, fmt.Errorf("querying advisory failed: %w", err)
	}
	if missingAdvisory {
		if err := tx.QueryRow(
			ctx, insertAdvisory, trackingID, publisher).Scan(&advisoryID); err != nil {
			return 0, fmt.Errorf("creating advisory (%q/%s) failed: %w",
				publisher, trackingID, err)
		}
	}

	// Using a savepoint only rolls back the transaction partially.
	if _, err := tx.Exec(ctx, savepointDoc); err != nil {
		return 0, err
	}

	var id int64
	if err := tx.QueryRow(
		ctx, insertDoc,
		document, raw,
		advisoryID,
	).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		// Unique constraint violation
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			// Rolling back to savepoint to enable execution of the
			// remaining book keeping within the inTx callback.
			if _, err2 := tx.Exec(ctx, rollbackSavepointDoc); err2 != nil {
				return 0, errors.Join(ErrAlreadyInDatabase, err2)
			}
			if inTx != nil {
				if err2 := inTx(ctx, tx, 0, true); err2 != nil {
					return 0, errors.Join(ErrAlreadyInDatabase, err2)
				}
				if err2 := tx.Commit(ctx); err2 != nil {
					return 0, errors.Join(ErrAlreadyInDatabase, err2)
				}
			}
			return 0, ErrAlreadyInDatabase
		}
		return 0, fmt.Errorf("inserting document failed: %w", err)
	}

	// Keep the document
	if _, err := tx.Exec(ctx, releaseSavepointDoc); err != nil {
		return 0, err
	}

	if _, err := tx.Exec(ctx, insertLog, actor, id); err != nil {
		return 0, fmt.Errorf("inserting log failed: %w", err)
	}

	if err := storeDocumentTexts(ctx, tx, advisoryID, id, idxer); err != nil {
		return 0, err
	}

	if inTx != nil {
		if err := inTx(ctx, tx, id, false); err != nil {
			return 0, fmt.Errorf("in transaction failed: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commiting transaction failed: %w", err)
	}
	return id, nil
}
