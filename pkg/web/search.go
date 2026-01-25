// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"text/template"

	"github.com/ISDuBA/ISDuBA/pkg/database/query"
	"github.com/ISDuBA/ISDuBA/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	aggregatedSection struct {
		id      int64
		results [][]any
	}
	aggregatedResult struct {
		fields   []string
		sections []aggregatedSection
	}
)

func (c *Controller) aggregatedResults(
	ctx *gin.Context,
	calcCount bool,
	limit, offset int64,
	fields []string,
	order string,
	builder *query.SQLBuilder,
) {
	var (
		result   *aggregatedResult
		sql      = builder.CreateQuery(fields, order, -1, -1)
		rctx     = ctx.Request.Context()
		filtered = builder.RemoveIgnoredFields(fields)
		escape   = needsEscaping(filtered, builder.Aliases)
	)

	if slog.Default().Enabled(rctx, slog.LevelDebug) {
		slog.Debug("documents", "SQL", query.InterpolateSQLqnd(sql, builder.Replacements))
	}
	if err := c.db.Run(
		rctx,
		func(rctx context.Context, conn *pgxpool.Conn) error {
			rows, err := conn.Query(rctx, sql, builder.Replacements...)
			if err != nil {
				return fmt.Errorf("cannot fetch results: %w", err)
			}
			defer rows.Close()
			if result, err = scanAggregatedRows(rows, filtered); err != nil {
				return fmt.Errorf("loading data failed: %w", err)
			}
			return nil
		},
		c.cfg.Database.MaxQueryDuration, // In case the user provided a very expensive query.
	); err != nil {
		slog.Error("database error", "err", err)
		models.SendError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Render(http.StatusOK, jsonStream(func(w http.ResponseWriter) error {
		// TODO: Produce output.
		_ = calcCount
		fmt.Fprint(w, "{")
		if calcCount {
			fmt.Fprintf(w, `"count": %d`, len(result.sections))
		}
		// Only produce documents when we have them.
		if len(result.sections) > 0 {
			if calcCount {
				fmt.Fprint(w, ',')
			}
			fmt.Fprint(w, `{"documents":[`)
			firstSection := true
			data := make(map[string]any, len(fields))
			enc := json.NewEncoder(w)

			if err := result.window(limit, offset, func(as *aggregatedSection) error {
				if firstSection {
					firstSection = false
				} else {
					fmt.Fprint(w, ',')
				}
				fmt.Fprintf(w, `"document":{"id:%d,"data":["`, as.id)

				for i, result := range as.results {
					if i > 0 {
						fmt.Fprint(w, ',')
					}
					clear(data)
					for j, v := range result {
						name := fields[j]
						if name == "id" {
							continue
						}
						if escape[j] {
							// XXX: A little bit hacky to support client.
							if s, ok := v.(string); ok {
								v = template.HTMLEscapeString(s)
							}
						}
						data[name] = v
					}
					if err := enc.Encode(data); err != nil {
						return err
					}
				}

				_, err := fmt.Fprint(w, ']')
				return err
			}); err != nil {
				slog.Error("writing window failed", "err", err)
				return err
			}
			if _, err := fmt.Fprint(w, `]}`); err != nil {
				return err
			}
			_ = escape
		}
		_, err := fmt.Fprint(w, '}')
		return err
	}))
}

type jsonStream func(http.ResponseWriter) error

var jsonContentType = []string{"application/json; charset=utf-8"}

func (js jsonStream) Render(w http.ResponseWriter) error { return js(w) }

func (jsonStream) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = jsonContentType
	}
}

// scanAggregatedRows turns a result set into an aggregatedResult.
func scanAggregatedRows(
	rows pgx.Rows,
	fields []string,
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
				results: [][]any{results},
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

func (ar *aggregatedResult) window(
	limit, offset int64,
	write func(*aggregatedSection) error,
) error {
	sections := ar.sections
	var start, end int
	if offset < 0 {
		start = 0
	} else {
		start = min(int(offset), len(sections))
	}
	if limit < 0 {
		end = len(sections)
	} else {
		end = min(int(offset+limit), len(sections))
	}
	for i := start; i < end; i++ {
		if err := write(&sections[i]); err != nil {
			return err
		}
	}
	return nil
}
