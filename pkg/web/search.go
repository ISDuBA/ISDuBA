// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

package web

import (
	"errors"
	"fmt"
	"slices"

	"github.com/jackc/pgx/v5"
)

type (
	aggregatedSection struct {
		id      int64
		results []any
	}
	aggregatedResult struct {
		fields   []string
		escape   []bool
		sections []aggregatedSection
	}
)

// scanAggregatedRows turns a result set into an aggregatedResult.
func scanAggregatedRows(
	rows pgx.Rows,
	fields []string,
	escape []bool,
) (*aggregatedResult, error) {
	idIdx := slices.Index(fields, "id")
	if idIdx == -1 {
		return nil, errors.New("missing id column to aggregate")
	}
	values := make([]any, len(fields))
	ptrs := make([]any, len(fields))
	for i := range ptrs {
		ptrs[i] = &values[i]
	}
	ag := aggregatedResult{
		fields: fields,
		escape: escape,
	}
	lastID := int64(-1)
	for rows.Next() {
		if err := rows.Scan(ptrs...); err != nil {
			return nil, fmt.Errorf("scanning row failed: %w", err)
		}
		results := slices.Clone(values)
		id, ok := values[idIdx].(int64)
		if !ok {
			// XXX: Should we panic here!?
			return nil, errors.New("id column is not an int64")
		}
		if id != lastID {
			ag.sections = append(ag.sections, aggregatedSection{
				id:      id,
				results: results,
			})
			lastID = id
		} else {
			last := &ag.sections[len(ag.sections)-1].results
			*last = append(*last, results)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scanning failed: %w", err)
	}
	return &ag, nil
}
