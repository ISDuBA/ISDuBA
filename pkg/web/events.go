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

// overviewEvents is an endpoint that returns a list of events.
//
//	@Summary		Returns a list of events.
//	@Description	Returns all events that match the specified query.
//	@Param			query	query	string	false	"Event query"
//	@Produce		json
//	@Success		200	{object}	web.overviewEvents.events
//	@Failure		400	{object}	models.Error
//	@Failure		401
//	@Failure		404	{object}	models.Error
//	@Failure		500	{object}	models.Error
//	@Router			/events [get]
func (c *Controller) overviewEvents(ctx *gin.Context) {
	parser := query.Parser{
		Mode:            query.EventMode,
		MinSearchLength: MinSearchLength,
		Me:              ctx.GetString("uid"),
	}

	// The query to filter the documents.
	expr, ok := parse(ctx, parser.Parse, ctx.DefaultQuery("query", "true"))
	if !ok {
		return
	}

	// Filter the allowed
	expr = c.andTLPExpr(ctx, expr)

	builder := query.SQLBuilder{Mode: query.EventMode}
	builder.CreateWhere(expr)

	fields := strings.Fields(
		ctx.DefaultQuery("columns", "event event_state time actor comments_id message id"))

	if err := builder.CheckProjections(fields); err != nil {
		models.SendError(ctx, http.StatusBadRequest, err)
		return
	}

	orderFields := strings.Fields(ctx.DefaultQuery("orders", "-time"))
	order, err := builder.CreateOrder(orderFields)
	if err != nil {
		models.SendError(ctx, http.StatusBadRequest, err)
		return
	}

	var (
		calcCount     bool
		count         int64
		limit, offset int64 = -1, -1
	)

	calcCount = ctx.Query("count") != ""

	if lim := ctx.Query("limit"); lim != "" {
		if limit, ok = parse(ctx, toInt64, lim); !ok {
			return
		}
	}

	if ofs := ctx.Query("offset"); ofs != "" {
		if offset, ok = parse(ctx, toInt64, ofs); !ok {
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
				slog.Debug("events", "SQL", query.InterpolateSQLqnd(sql, builder.Replacements))
			}
			rows, err := conn.Query(rctx, sql, builder.Replacements...)
			if err != nil {
				return fmt.Errorf("cannot fetch results: %w", err)
			}
			defer rows.Close()
			filtered := builder.RemoveIgnoredFields(fields)
			if results, err = scanRows(rows, filtered); err != nil {
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

	type events struct {
		Events []map[string]any `json:"events"`
		Count  int64            `json:"count,omitempty"`
	}
	h := events{}
	if calcCount {
		h.Count = count
	}
	if len(results) > 0 {
		h.Events = results
	}
	ctx.JSON(http.StatusOK, h)
}

// viewEvents is an endpoint that returns the events of the specified advisory.
//
//	@Summary		Returns all events.
//	@Description	Returns all events from the specified advisory.
//	@Param			publisher	path	string	true	"Publisher"
//	@Param			trackingid	path	string	true	"Tracking ID"
//	@Produce		json
//	@Success		200	{array}		web.viewEvents.event
//	@Failure		400	{object}	models.Error
//	@Failure		401
//	@Failure		404	{object}	models.Error
//	@Failure		500	{object}	models.Error
//	@Router			/events/{publisher}/{trackingid} [get]
func (c *Controller) viewEvents(ctx *gin.Context) {
	var key models.AdvisoryKey
	if err := ctx.ShouldBindUri(&key); err != nil {
		models.SendError(ctx, http.StatusBadRequest, err)
		return
	}

	if key.Publisher == "" || key.TrackingID == "" {
		models.SendErrorMessage(ctx, http.StatusBadRequest, "missing publisher or tracking_id")
		return
	}

	expr := c.andTLPExpr(ctx,
		query.FieldEqString("tracking_id", key.TrackingID).And(
			query.FieldEqString("publisher", key.Publisher)))

	builder := query.SQLBuilder{}
	builder.CreateWhere(expr)

	type event struct {
		Event      models.Event    `json:"event_type"`
		State      models.Workflow `json:"state"`
		Time       time.Time       `json:"time"`
		Actor      *string         `json:"actor,omitempty"`
		DocumentID int64           `json:"document_id"`
		CommentID  *int64          `json:"comment_id,omitempty"`
		PrevSSVC   *string         `json:"prev_ssvc,omitempty"`
	}

	var events []event
	var exists bool

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			existsSQL := `SELECT EXISTS(` +
				`SELECT FROM documents JOIN advisories ON documents.advisories_id = advisories.id ` +
				`WHERE ` + builder.WhereClause + `)`
			if err := conn.QueryRow(
				rctx, existsSQL, builder.Replacements...).Scan(&exists); err != nil {
				return err
			}
			if !exists {
				return nil
			}
			fetchSQL := `SELECT event, documents_id, time, actor, state, comments_id, prev_ssvc FROM events_log ` +
				`WHERE documents_id in (` +
				`SELECT documents.id ` +
				`FROM documents JOIN advisories ON documents.advisories_id = advisories.id ` +
				`WHERE ` + builder.WhereClause + `) ORDER BY time DESC`
			rows, _ := conn.Query(rctx, fetchSQL, builder.Replacements...)
			var err error
			events, err = pgx.CollectRows(
				rows,
				func(row pgx.CollectableRow) (event, error) {
					var ev event
					var act sql.NullString
					var prevSSVC sql.NullString
					err := row.Scan(&ev.Event, &ev.DocumentID, &ev.Time, &ev.Actor, &ev.State, &ev.CommentID, &prevSSVC)
					ev.Time = ev.Time.UTC()
					if act.Valid {
						ev.Actor = &act.String
					}
					if prevSSVC.Valid {
						ev.PrevSSVC = &prevSSVC.String
					}
					return ev, err
				})
			return err
		}, 0,
	); err != nil {
		slog.Error("database error", "err", err)
		models.SendError(ctx, http.StatusInternalServerError, err)
		return
	}

	if !exists {
		models.SendErrorMessage(ctx, http.StatusNotFound, "advisory not found")
		return
	}

	ctx.JSON(http.StatusOK, events)
}
