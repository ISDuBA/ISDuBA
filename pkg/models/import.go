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

	"github.com/csaf-poc/csaf_distribution/v3/csaf"
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

type replacer func([]string, string) (any, bool)

func chainReplacers(replacers ...replacer) replacer {
	return func(keys []string, value string) (any, bool) {
		for _, rep := range replacers {
			if x, ok := rep(keys, value); ok {
				return x, true
			}
		}
		return value, false
	}
}

type indexer[T comparable] struct {
	elements        []T
	indexToElements map[T]int
}

func newIndexer[T comparable]() *indexer[T] {
	return &indexer[T]{
		indexToElements: make(map[T]int),
	}
}

func (i *indexer[T]) index(t T) int {
	if idx, ok := i.indexToElements[t]; ok {
		return idx
	}
	idx := len(i.elements)
	i.elements = append(i.elements, t)
	i.indexToElements[t] = idx
	return idx
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

func keepByKeys(keys []string) replacer {
	return func(ks []string, v string) (any, bool) {
		if len(ks) == 0 {
			return v, false
		}
		_, found := slices.BinarySearch(keys, ks[len(ks)-1])
		return v, found
	}
}

func keepByValues(values []string) replacer {
	return func(_ []string, v string) (any, bool) {
		_, found := slices.BinarySearch(values, v)
		return v, found
	}
}

func replaceByIndex(index func(string) int) replacer {
	return func(_ []string, v string) (any, bool) {
		return index(v), true
	}
}

func keepAndIndex(index func(string) int, path ...string) replacer {
	found := false
	return func(ks []string, v string) (any, bool) {
		if !found && slices.Equal(path, ks) {
			found = true
			_ = index(v)
			return v, true
		}
		return v, false
	}
}

func keepAndIndexSuffix(index func(string) int, path ...string) replacer {
	return func(ks []string, v string) (any, bool) {
		if len(ks) >= len(path) && slices.Equal(path, ks[len(ks)-len(path):]) {
			_ = index(v)
			return v, true
		}
		return v, false
	}
}

func transformJSON(document any, replace replacer) {
	var (
		array  func(arr []any)
		object func(obj map[string]any)
		keys   []string
	)

	array = func(arr []any) {
		for i, v := range arr {
			_ = i
			switch x := v.(type) {
			case string:
				if y, ok := replace(keys, x); ok {
					arr[i] = y
				}
			case []any:
				array(x)
			case map[string]any:
				object(x)

			}
		}
	}

	object = func(obj map[string]any) {
		for k, v := range obj {
			keys = append(keys, k)
			switch x := v.(type) {
			case string:
				if y, ok := replace(keys, x); ok {
					obj[k] = y
				}
			case []any:
				array(x)
			case map[string]any:
				object(x)
			}
			keys = keys[:len(keys)-1]
		}
	}

	switch x := document.(type) {
	case []any:
		array(x)
	case map[string]any:
		object(x)
	}
}

func sorted(s []string) []string {
	slices.Sort(s)
	return s
}

var (
	excludeKeys = sorted([]string{
		"id",
		"category",
		"csaf_version",
		"date",
		"version",
		"label",
		"lang",
		"status",
		"initial_release_date",
		"current_release_date",
		"release_date",
		"discovery_date",
		"vectorString",
	})
	excludeValues = sorted([]string{
		"HIGH",
		"MEDIUM",
		"LOW",
		"LOW_MEDIUM",
		"MEDIUM_HIGH",
		"CHANGED",
		"UNCHANGED",
		"MULTIPLE",
		"SINGLE",
		"NONE",
		"NETWORK",
		"ADJACENT_NETWORK",
		"LOCAL",
		"PHYSICAL",
		"NOT_DEFINED",
		"PARTIAL",
		"COMPLETE",
		"UNPROVEN",
		"PROOF_OF_CONCEPT",
		"FUNCTIONAL",
		"OFFICIAL_FIX",
		"TEMPORARY_FIX",
		"WORKAROUND",
		"UNAVAILABLE",
		"UNCONFIRMED",
		"UNCORROBORATED",
		"CONFIRMED",
		"UNKNOWN",
		"REASONABLE",
		"REQUIRED",
		"CRITICAL",
	})
)

// ChainInTx executes a list of in transaction functions.
func ChainInTx(
	inTxs ...func(context.Context, pgx.Tx, int64) error,
) func(context.Context, pgx.Tx, int64) error {
	return func(ctx context.Context, tx pgx.Tx, docID int64) error {
		for _, inTx := range inTxs {
			if err := inTx(ctx, tx, docID); err != nil {
				return err
			}
		}
		return nil
	}
}

// ImportDocument imports a given advisory into the database.
func ImportDocument(
	ctx context.Context,
	conn *pgxpool.Conn,
	r io.Reader,
	actor *string,
	pstlps PublishersTLPs,
	inTx func(context.Context, pgx.Tx, int64) error,
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
	inTx func(context.Context, pgx.Tx, int64) error,
	dry bool,
) (int64, error) {

	var (
		tlp, tlpOk               = "", false
		publisher, publisherOK   = "", false
		trackingID, trackingIDOK = "", false
	)

	idxer := newIndexer[string]()

	var reps []replacer

	if pstlps != nil {
		reps = append(reps, storer(&tlp, &tlpOk, "document", "distribution", "tlp", "label"))
	}

	transformJSON(document, chainReplacers(
		append(reps,
			storer(&publisher, &publisherOK, "document", "publisher", "name"),
			storer(&trackingID, &trackingIDOK, "document", "tracking", "id"),
			keepAndIndex(idxer.index, "document", "publisher", "name"),
			keepAndIndex(idxer.index, "document", "title"),
			keepAndIndexSuffix(idxer.index, "vulnerabilities", "cve"),
			keepByKeys(excludeKeys),
			keepByValues(excludeValues),
			replaceByIndex(idxer.index),
		)...))

	if !publisherOK {
		return 0, errors.New("missing /document/publisher/name")
	}

	if !trackingIDOK {
		return 0, errors.New("missing /document/tracking/id")
	}

	if pstlps != nil && (!tlpOk || !pstlps.Allowed(publisher, TLP(tlp))) {
		return 0, ErrNotAllowed
	}

	if dry {
		return 0, nil
	}

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	const (
		insertDoc     = `INSERT INTO documents (document, original) VALUES ($1, $2) RETURNING id`
		insertLog     = `INSERT INTO events_log (event, state, actor, documents_id) VALUES ('import_document', 'new', $1, $2)`
		queryText     = `SELECT id FROM unique_texts WHERE txt = $1`
		insertText    = `INSERT INTO unique_texts (txt) VALUES ($1) RETURNING id`
		insertDocText = `INSERT INTO documents_texts (documents_id, num, txt_id) VALUES ($1, $2, $3)`
		loadTexts     = `SELECT u.id, txt FROM documents d JOIN documents_texts t ` +
			`ON d.id = t.documents_id JOIN unique_texts u ` +
			`ON t.txt_id = u.id ` +
			`WHERE d.publisher = $1 AND d.tracking_id = $2`
	)

	var id int64
	if err := tx.QueryRow(
		ctx, insertDoc,
		document, raw,
	).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		// Unique constraint violation
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, ErrAlreadyInDatabase
		}
		return 0, fmt.Errorf("inserting document failed: %w", err)
	}
	if _, err := tx.Exec(ctx, insertLog, actor, id); err != nil {
		return 0, fmt.Errorf("inserting log failed: %w", err)
	}

	txtIDs := make([]int64, len(idxer.elements))
	for i := range txtIDs {
		txtIDs[i] = -1
	}

	// If we already have a document with the given publisher/tracking_id pair
	// it is very likely that they share a lot of the same strings.
	if err := func() error {
		rows, err := tx.Query(ctx, loadTexts, publisher, trackingID)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var textID int64
			var text string
			if err := rows.Scan(&textID, &text); err != nil {
				return err
			}
			if idx, ok := idxer.indexToElements[text]; ok {
				txtIDs[idx] = textID
			}
		}
		return rows.Err()
	}(); err != nil {
		return 0, fmt.Errorf("loading old texts failed: %w", err)
	}

	insertTextBatch := &pgx.Batch{}

	scanText := func(idx int) func(pgx.Row) error {
		return func(row pgx.Row) error {
			if err := row.Scan(&txtIDs[idx]); err != nil {
				if !errors.Is(err, pgx.ErrNoRows) {
					return fmt.Errorf("finding unique text failed: %w", err)
				}
				insertTextBatch.Queue(insertText, idxer.elements[idx]).QueryRow(
					func(row pgx.Row) error { return row.Scan(&txtIDs[idx]) })
			}
			return nil
		}
	}
	textIDsBatch := &pgx.Batch{}
	for i, txt := range idxer.elements {
		if txtIDs[i] == -1 {
			// Only ask for strings we have not found already.
			textIDsBatch.Queue(queryText, txt).QueryRow(scanText(i))
		}
	}

	if err := tx.SendBatch(ctx, textIDsBatch).Close(); err != nil {
		return 0, fmt.Errorf("finding txt failed: %w", err)
	}

	// We need to insert some
	if insertTextBatch.Len() > 0 {
		if err := tx.SendBatch(ctx, insertTextBatch).Close(); err != nil {
			return 0, fmt.Errorf("inserting txt failed: %w", err)
		}
	}

	batch := &pgx.Batch{}
	for i, txtID := range txtIDs {
		batch.Queue(insertDocText, id, i, txtID)
	}
	if err := tx.SendBatch(ctx, batch).Close(); err != nil {
		return 0, fmt.Errorf("inserting txt failed: %w", err)
	}

	if inTx != nil {
		if err := inTx(ctx, tx, id); err != nil {
			return 0, fmt.Errorf("in transaction failed: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commiting transaction failed: %w", err)
	}
	return id, nil
}
