// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

// batchSize to process only a managable amount of documents each iteration
const batchSize = 100

func UpdateDocumentsIndex(
	ctx context.Context,
	tx pgx.Tx,
) error {
	globalInsertLock.Lock()
	defer globalInsertLock.Unlock()

	const selectBatch = `SELECT id, advisories_id, original ` +
		`FROM documents WHERE id > $1 ORDER BY id LIMIT $2`

	type docRow struct {
		id         int64
		advisoryID int64
		original   []byte
	}

	var last, total int64
	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		rows, err := tx.Query(ctx, selectBatch, last, batchSize)
		if err != nil {
			return fmt.Errorf("selecting documents failed: %w", err)
		}
		defer rows.Close()
		var batch []docRow
		for rows.Next() {
			var row docRow
			if err := rows.Scan(*&row.id, *&row.advisoryID, &row.original); err != nil {
				return fmt.Errorf("scanning document failed: %w", err)
			}
			batch = append(batch, row)
		}
		if err := rows.Err(); err != nil {
			return fmt.Errorf("reading documents failed: %w", err)
		}
		if len(batch) == 0 {
			return nil
		}

		for _, row := range batch {
			if err := updateDocumentIndex(ctx, tx, row.id, row.advisoryID, row.original); err != nil {
				return fmt.Errorf("updating document %d indexes failed: %w", row.id, err)
			}
			last = row.id
		}
		total += int64(len(batch))
		slog.Info("updated documents index", "total", total, "last_id", last)
	}
}

// updateDocumentIndex redoes indexing of one document from its original
// and stores the updated texts in the database
func updateDocumentIndex(ctx context.Context, tx pgx.Tx, id, advisoryID int64, original []byte) error {
	const (
		updateDoc  = `UPDATE documents SET document = $1 WHERE id = $2`
		deleteText = `DELETE FROM documents_texts WHERE documents_id = $1`
	)

	var document any
	if err := json.NewDecoder(bytes.NewReader(original)).Decode(&document); err != nil {
		return fmt.Errorf("decoding original failed: %w", err)
	}

	idxer := indexDocument(document)

	if _, err := tx.Exec(ctx, updateDoc, document, id); err != nil {
		return fmt.Errorf("updating document failed: %w", err)
	}

	if _, err := tx.Exec(ctx, deleteText, id); err != nil {
		return fmt.Errorf("deleting old texts failed: %w", err)
	}

	return storeDocumentTexts(ctx, tx, advisoryID, id, idxer)
}
