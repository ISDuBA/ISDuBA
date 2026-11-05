// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

package models

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"

	"github.com/jackc/pgx/v5"
)

var globalInsertLock sync.Mutex

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

func keepByKeys(keys []string) replacer {
	return func(ks []string, v string) (any, bool) {
		if len(ks) == 0 {
			return v, false
		}
		_, found := slices.BinarySearch(keys, ks[len(ks)-1])
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

// transformJSON replaces defined things in-place in the document.
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
)

// indexDocument transforms a given document via replacer functions
// pre replacers are for validation and keeping certain strings (see import.go)
// returns indexer holding all texts replaced by numbers
func indexDocument(document any, pre ...replacer) *indexer[string] {
	idxer := newIndexer[string]()

	reps := append(pre,
		keepAndIndex(idxer.index, "document", "publisher", "name"),
		keepAndIndex(idxer.index, "document", "title"),
		keepAndIndexSuffix(idxer.index, "vulnerabilities", "cve"),
		keepByKeys(excludeKeys),
		replaceByIndex(idxer.index),
	)

	transformJSON(document, chainReplacers(reps...))
	return idxer
}

// storeDocumentTexts resolves idxer texts to unique_texts ids
// and links them to the document via documents_texts.
func storeDocumentTexts(ctx context.Context, tx pgx.Tx, advisoryID, id int64, idxer *indexer[string]) error {
	const (
		queryText     = `SELECT id FROM unique_texts WHERE txt = $1`
		insertText    = `INSERT INTO unique_texts (txt) VALUES ($1) RETURNING id`
		insertDocText = `INSERT INTO documents_texts (documents_id, num, txt_id) VALUES ($1, $2, $3)`
		loadTexts     = `SELECT u.id, txt FROM documents d JOIN documents_texts t ` +
			`ON d.id = t.documents_id JOIN unique_texts u ` +
			`ON t.txt_id = u.id ` +
			`WHERE d.advisories_id = $1`
	)

	txtIDs := make([]int64, len(idxer.elements))
	for i := range txtIDs {
		txtIDs[i] = -1
	}

	// If we already have a document with the given publisher/tracking_id pair
	// it is very likely that they share a lot of the same strings.
	if err := func() error {
		rows, err := tx.Query(ctx, loadTexts, advisoryID)
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
		return fmt.Errorf("loading old texts failed: %w", err)
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
		return fmt.Errorf("finding txt failed: %w", err)
	}

	// We need to insert some
	if insertTextBatch.Len() > 0 {
		if err := tx.SendBatch(ctx, insertTextBatch).Close(); err != nil {
			return fmt.Errorf("inserting txt failed: %w", err)
		}
	}

	batch := &pgx.Batch{}
	for i, txtID := range txtIDs {
		batch.Queue(insertDocText, id, i, txtID)
	}
	if err := tx.SendBatch(ctx, batch).Close(); err != nil {
		return fmt.Errorf("inserting txt failed: %w", err)
	}

	return nil
}
