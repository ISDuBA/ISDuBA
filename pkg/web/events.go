// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package web

import (
	"database/sql"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/ISDuBA/ISDuBA/pkg/database"
	"github.com/ISDuBA/ISDuBA/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (c *Controller) viewEvents(ctx *gin.Context) {
	idS := ctx.Param("document")
	id, err := strconv.ParseInt(idS, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expr := database.FieldEqInt("id", id)

	// Filter the allowed
	if tlps := c.tlps(ctx); len(tlps) > 0 {
		conditions := tlps.AsConditions()
		parser := database.Parser{}
		tlpExpr, err := parser.Parse(conditions)
		if err != nil {
			slog.Warn("TLP filter failed", "err", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		expr = expr.And(tlpExpr)
	}

	where, replacements, _ := expr.Where()

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

	rctx := ctx.Request.Context()
	if err := c.db.Run(rctx, func(conn *pgxpool.Conn) error {
		existsSQL := `SELECT exists(SELECT FROM documents WHERE ` + where + `)`
		if err := conn.QueryRow(
			rctx, existsSQL, replacements...).Scan(&exists); err != nil {
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
	}); err != nil {
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
