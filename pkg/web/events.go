// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package web

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ISDuBA/ISDuBA/pkg/database/query"
	"github.com/ISDuBA/ISDuBA/pkg/models"
)

func (c *Controller) overviewEvents(ctx *gin.Context) {

	parser := query.Parser{
		Mode:            query.EventMode,
		MinSearchLength: MinSearchLength,
		Me:              ctx.GetString("uid"),
	}

	// The query to filter the documents.
	expr, err := parser.Parse(ctx.DefaultQuery("query", "true"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Filter the allowed
	expr = c.andTLPExpr(ctx, expr)

	builder := query.SQLBuilder{Mode: query.EventMode}
	builder.CreateWhere(expr)

	fields := strings.Fields(
		ctx.DefaultQuery("columns", "event event_state time actor comments_id id"))

	if err := builder.CheckProjections(fields); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderFields := strings.Fields(ctx.DefaultQuery("order", "-time"))
	order, err := builder.CreateOrder(orderFields)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var (
		calcCount, ok bool
		count         int64
		limit, offset int64 = -1, -1
	)

	calcCount = ctx.Query("count") != ""

	if lim := ctx.Query("limit"); lim != "" {
		if limit, ok = parseInt(ctx, lim); !ok {
			return
		}
	}

	if ofs := ctx.Query("offset"); ofs != "" {
		if offset, ok = parseInt(ctx, ofs); !ok {
			return
		}
	}

	var results []map[string]any

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			if calcCount {
				if err := conn.QueryRow(
					rctx,
					builder.CreateCountSQL(),
					builder.Replacements...,
				).Scan(&count); err != nil {
					return fmt.Errorf("cannot calculate count %w", err)
				}
			}
			// Skip fields if they are not requested.
			if len(fields) == 0 {
				return nil
			}

			sql := builder.CreateQuery(fields, order, limit, offset)

			if slog.Default().Enabled(rctx, slog.LevelDebug) {
				slog.Debug("events", "SQL", qndSQLReplace(sql, builder.Replacements))
			}
			rows, err := conn.Query(rctx, sql, builder.Replacements...)
			if err != nil {
				return fmt.Errorf("cannot fetch results: %w", err)
			}
			defer rows.Close()
			if results, err = scanRows(rows, fields, builder.Aliases); err != nil {
				return fmt.Errorf("loading data failed: %w", err)
			}
			return nil
		},
		c.cfg.Database.MaxQueryDuration, // In case the user provided a very expensive query.
	); err != nil {
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h := gin.H{}
	if calcCount {
		h["count"] = count
	}
	if len(results) > 0 {
		h["events"] = results
	}
	ctx.JSON(http.StatusOK, h)
}

func (c *Controller) viewEvents(ctx *gin.Context) {
	id, ok := parseInt(ctx, ctx.Param("document"))
	if !ok {
		return
	}

	expr := c.andTLPExpr(ctx, query.FieldEqInt("id", id))

	builder := query.SQLBuilder{}
	builder.CreateWhere(expr)

	type event struct {
		Event      models.Event    `json:"event_type"`
		State      models.Workflow `json:"state"`
		Time       time.Time       `json:"time"`
		Actor      *string         `json:"actor,omitempty"`
		DocumentID int64           `json:"document_id"`
		CommentID  *int64          `json:"comment_id,omitempty"`
	}

	var events []event
	var exists bool

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			existsSQL := `SELECT exists(SELECT FROM documents WHERE ` +
				builder.WhereClause + `)`
			if err := conn.QueryRow(
				rctx, existsSQL, builder.Replacements...).Scan(&exists); err != nil {
				return err
			}
			if !exists {
				return nil
			}
			const fetchSQL = `SELECT event, documents_id, time, actor, state, comments_id FROM events_log ` +
				`WHERE documents_id = $1 ORDER BY time DESC`
			rows, _ := conn.Query(rctx, fetchSQL, id)
			var err error
			events, err = pgx.CollectRows(

				rows,
				func(row pgx.CollectableRow) (event, error) {
					var ev event
					var act sql.NullString
					err := row.Scan(&ev.Event, &ev.DocumentID, &ev.Time, &ev.Actor, &ev.State, &ev.CommentID)
					ev.Time = ev.Time.UTC()
					if act.Valid {
						ev.Actor = &act.String
					}
					return ev, err
				})
			return err
		}, 0,
	); err != nil {
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}

	ctx.JSON(http.StatusOK, events)
}
