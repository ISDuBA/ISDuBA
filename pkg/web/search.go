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
	"maps"
	"net/http"
	"reflect"
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
		rows []map[string]any
	}
	aggregatedDocuments []aggregatedDocument
)

type encodedTexts[T comparable] struct {
	txts  map[T]string
	fetch func(T) (string, error)
}

func (et *encodedTexts[T]) resolve(t T) (string, error) {
	if m, ok := et.txts[t]; ok {
		return m, nil
	}
	v, err := et.fetch(t)
	if err != nil {
		return "", err
	}
	m, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	if et.txts == nil {
		et.txts = make(map[T]string)
	}
	x := string(m)
	et.txts[t] = x
	return x, nil
}

func (c *Controller) aggregatedResults(
	ctx *gin.Context,
	calcCount bool,
	limit, offset int64,
	builder *query.AdvancedSQLBuilder,
) {
	var (
		ads  aggregatedDocuments
		sql  = builder.CreateQuery(-1, -1)
		rctx = ctx.Request.Context()
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
			if ads, err = scanAggregatedDocuments(rows, builder.Fields()); err != nil {
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

	// We load texts from the database lazily so we render in another connection.
	// XXX: Think about moving it to the DB stuff above.
	if err := c.db.Run(
		rctx,
		func(rctx context.Context, conn *pgxpool.Conn) error {
			// Load a text only if we don't already have it.
			const uniqueSQL = `` +
				`SELECT` +
				` CASE WHEN length(txt)<= 200 THEN txt` +
				` ELSE substring(txt, 0, 197)||'...' END` +
				` FROM unique_texts WHERE id = $1`
			loadedTxts := encodedTexts[int64]{
				fetch: func(id int64) (string, error) {
					var txt string
					if err := conn.QueryRow(rctx, uniqueSQL, id).Scan(&txt); err != nil {
						return "", fmt.Errorf("fetching unique text failed: %w", err)
					}
					return txt, nil
				},
			}
			// JSON encoding keys will always results in the same values
			// so we cache them, too.
			keys := encodedTexts[string]{
				fetch: func(x string) (string, error) { return x, nil },
			}
			// Catch errors occuring while rendering to logged outside.
			var trackedErr error
			ctx.Render(http.StatusOK, jsonStream(trackError(
				&trackedErr,
				func(w http.ResponseWriter) error {
					fmt.Fprint(w, "{")
					if calcCount {
						fmt.Fprintf(w, `"count":%d`, len(ads))
					}
					// Only produce documents when we have them.
					if len(ads) == 0 {
						_, err := fmt.Fprint(w, "}")
						return err
					}
					if calcCount {
						fmt.Fprint(w, ",")
					}
					fmt.Fprint(w, `"documents":[`)
					firstDocument := true
					enc := json.NewEncoder(w)

					if err := ads.window(
						limit, offset,
						func(ad *aggregatedDocument) error {
							if firstDocument {
								firstDocument = false
							} else {
								fmt.Fprint(w, ",")
							}
							fmt.Fprintf(w, `{"id":%d,"data":[`, ad.id)
							for i, row := range ad.rows {
								if i > 0 {
									fmt.Fprint(w, ",")
								}
								firstEntry := true
								fmt.Fprintf(w, "{")
								for k, v := range row {
									if firstEntry {
										firstEntry = false
									} else {
										fmt.Fprintf(w, ",")
									}
									key, err := keys.resolve(k)
									if err != nil {
										return err
									}
									fmt.Fprintf(w, "%s:", key)
									// If we have an alias we need to load the text from the database.
									if builder.HasAlias(k) {
										id, ok := asInt64(v)
										if !ok {
											return fmt.Errorf("alias %q has not an int value", k)
										}
										txt, err := loadedTxts.resolve(id)
										if err != nil {
											return err
										}
										fmt.Fprint(w, txt)
									} else {
										// A none alias field
										if err := enc.Encode(v); err != nil {
											return err
										}
									}
								}
								fmt.Fprintf(w, "}")
							}
							_, err := fmt.Fprint(w, "]}")
							return err
						}); err != nil {
						return fmt.Errorf("writing window failed %w", err)
					}
					if _, err := fmt.Fprint(w, "]"); err != nil {
						return err
					}
					_, err := fmt.Fprint(w, "}")
					return err
				})))
			return trackedErr
		},
		c.cfg.Database.MaxQueryDuration, // In case the user provided a very expensive query.
	); err != nil {
		slog.Error("database error", "err", err)
		// Too late to send an error to the client.
	}
}

type jsonStream func(http.ResponseWriter) error

func (js jsonStream) Render(w http.ResponseWriter) error { return js(w) }

func (jsonStream) WriteContentType(w http.ResponseWriter) {
	render.JSON{}.WriteContentType(w)
}

func trackError(
	err *error,
	fn func(http.ResponseWriter) error,
) func(http.ResponseWriter) error {
	return func(rw http.ResponseWriter) error {
		*err = fn(rw)
		return *err
	}
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
	have := make(map[string]any)
	for rows.Next() {
		if err := rows.Scan(ptrs...); err != nil {
			return nil, fmt.Errorf("scanning row failed: %w", err)
		}
		id, ok := asInt64(values[idIdx])
		if !ok {
			// XXX: Should we panic here!?
			return nil, fmt.Errorf("id column is not an int: %T", values[idIdx])
		}
		if id != lastID { // A new document bundle.
			lastID = id
			clear(have)
			for j, v := range values {
				// Ignore id as it will be already stored in the aggregatedDocument.
				if name := fields[j]; name != "id" {
					have[name] = v
				}
			}
			ads = append(ads, aggregatedDocument{
				id:   id,
				rows: []map[string]any{maps.Clone(have)},
			})
		} else { // We already have documents for this bundle.
			row := make(map[string]any, 1) // The diff should be small.
			for j, v := range values {
				name := fields[j]
				if name == "id" {
					continue
				}
				if x, ok := have[name]; !ok || !reflect.DeepEqual(x, v) {
					row[name] = v
					have[name] = v
				}
			}
			if len(row) > 0 {
				last := &ads[len(ads)-1].rows
				*last = append(*last, row)
			}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scanning failed: %w", err)
	}
	return ads, nil
}

func asInt64(x any) (int64, bool) {
	switch v := x.(type) {
	case int64:
		return v, true
	case int8:
		return int64(v), true
	case int16:
		return int64(v), true
	case int32:
		return int64(v), true
	case uint64:
		return int64(v), true
	case uint8:
		return int64(v), true
	case uint16:
		return int64(v), true
	case uint32:
		return int64(v), true
	case int:
		return int64(v), true
	case uint:
		return int64(v), true
	}
	return 0, false
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
