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

	"github.com/ISDuBA/ISDuBA/pkg/database/query"
	"github.com/ISDuBA/ISDuBA/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	aggregatedDocument struct {
		id   int64
		rows [][]any
	}
	aggregatedDocuments []aggregatedDocument
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
		ads      aggregatedDocuments
		sql      = builder.CreateQuery(fields, order, -1, -1)
		rctx     = ctx.Request.Context()
		filtered = builder.RemoveIgnoredFields(fields)
		esc      = needsEscaping(filtered, builder.Aliases)
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
			if ads, err = scanAggregatedDocuments(rows, filtered); err != nil {
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
		fmt.Fprint(w, "{")
		if calcCount {
			fmt.Fprintf(w, `"count":%d`, len(ads))
		}
		// Only produce documents when we have them.
		if len(ads) > 0 {
			if calcCount {
				fmt.Fprint(w, ',')
			}
			fmt.Fprint(w, `{"documents":[`)
			firstDocument := true
			data := make(map[string]any, len(fields))
			enc := json.NewEncoder(w)

			idIdx := slices.Index(fields, "id")

			if err := ads.window(limit, offset, func(ad *aggregatedDocument) error {
				if firstDocument {
					firstDocument = false
				} else {
					fmt.Fprint(w, ',')
				}
				fmt.Fprintf(w, `"document":{"id":%d,"data":["`, ad.id)

				for i, row := range ad.rows {
					if i > 0 {
						fmt.Fprint(w, ',')
					}
					clear(data)
					for j, v := range row {
						if j != idIdx {
							data[fields[j]] = esc.escape(j, v)
						}
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
		}
		_, err := fmt.Fprint(w, '}')
		return err
	}))
}

type jsonStream func(http.ResponseWriter) error

func (js jsonStream) Render(w http.ResponseWriter) error { return js(w) }

func (jsonStream) WriteContentType(w http.ResponseWriter) {
	render.JSON{}.WriteContentType(w)
}

// scanAggregatedDocuments returns a result set into a slice of aggregated documents.
func scanAggregatedDocuments(
	rows pgx.Rows,
	fields []string,
) (aggregatedDocuments, error) {
	idIdx := slices.Index(fields, "id")
	if idIdx == -1 {
		return nil, errors.New("missing id column to aggregate")
	}
	values := make([]any, len(fields))
	ptrs := make([]any, len(fields))
	for i := range ptrs {
		ptrs[i] = &values[i]
	}
	ads := make(aggregatedDocuments, 0, 512)
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
			ads = append(ads, aggregatedDocument{
				id:   id,
				rows: [][]any{results},
			})
			lastID = id
		} else {
			last := &ads[len(ads)-1].rows
			*last = append(*last, results)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scanning failed: %w", err)
	}
	return ads, nil
}

func (ads aggregatedDocuments) window(
	limit, offset int64,
	write func(*aggregatedDocument) error,
) error {
	var start, end int
	if offset < 0 {
		start = 0
	} else {
		start = min(int(offset), len(ads))
	}
	if limit < 0 {
		end = len(ads)
	} else {
		end = min(int(offset+limit), len(ads))
	}
	for i := start; i < end; i++ {
		if err := write(&ads[i]); err != nil {
			return err
		}
	}
	return nil
}
