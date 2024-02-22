// This file is Free Software under the MIT License
// without warranty, see README.md and LICENSES/MIT.txt for details.
//
// SPDX-License-Identifier: MIT
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package web

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/ISDuBA/ISDuBA/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (c *Controller) createComment(ctx *gin.Context) {

	docIDs := ctx.Param("document")
	docID, err := strconv.ParseInt(docIDs, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	query := fmt.Sprintf("$id %d int =", docID)
	expr := database.MustParse(query)

	// Filter the allowed
	if tlps := c.tlps(ctx); len(tlps) > 0 {
		conditions := tlps.AsConditions()
		tlpExpr, err := database.Parse(conditions)
		if err != nil {
			slog.Warn("TLP filter failed", "err", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error:": err})
			return
		}
		expr = expr.And(tlpExpr)
	}

	var (
		where, replacements, _ = expr.Where()
		exists                 bool
		commentator            = ctx.GetString("uid")
		message, _             = ctx.GetPostForm("message")
		now                    = time.Now().UTC()
		commentID              int64
		rctx                   = ctx.Request.Context()
	)

	if err := c.db.Run(rctx, func(conn *pgxpool.Conn) error {
		tx, err := conn.BeginTx(rctx, pgx.TxOptions{})
		if err != nil {
			return err
		}
		defer tx.Rollback(rctx)

		existsSQL := `SELECT exists(SELECT FROM documents ` + where + `)`
		if err := tx.QueryRow(rctx, existsSQL, replacements...).Scan(&exists); err != nil {
			return err
		}
		if !exists {
			return nil
		}

		const insertSQL = `INSERT INTO comments ` +
			`(documents_id, time, commentator, message) ` +
			`VALUES ($1, $2, $3, $4) ` +
			`RETURNING id`

		if err := tx.QueryRow(rctx, insertSQL, docID, now, commentator, message).Scan(&commentID); err != nil {
			return err
		}

		const eventSQL = `INSERT INTO events_log ` +
			`(event, state, time, actor, documents_id) ` +
			`VALUES('add_comment', (SELECT state FROM documents WHERE id = $3), $1, $2, $3)`

		var actor sql.NullString
		if !c.cfg.General.AnonymousEventLogging {
			actor.String = commentator
		}
		if _, err := tx.Exec(rctx, eventSQL, now, actor, message); err != nil {
			return err
		}

		return tx.Commit(rctx)
	}); err != nil {
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"id":          commentID,
		"time":        now,
		"commentator": commentator,
	})
}

func (c *Controller) updateComment(ctx *gin.Context) {
	// TODO: Implement me!
	_ = ctx
}

func (c *Controller) viewComments(ctx *gin.Context) {

	idS := ctx.Param("document")
	id, err := strconv.ParseInt(idS, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	query := fmt.Sprintf("$id %d int =", id)
	expr := database.MustParse(query)

	// Filter the allowed
	if tlps := c.tlps(ctx); len(tlps) > 0 {
		conditions := tlps.AsConditions()
		tlpExpr, err := database.Parse(conditions)
		if err != nil {
			slog.Warn("TLP filter failed", "err", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error:": err})
			return
		}
		expr = expr.And(tlpExpr)
	}

	where, replacements, _ := expr.Where()

	type comment struct {
		ID          int64     `json:"id"`
		Time        time.Time `json:"time"`
		Commentator string    `json:"commentator"`
		Message     string    `json:"message"`
	}

	var comments []comment
	var exists bool

	rctx := ctx.Request.Context()
	if err := c.db.Run(rctx, func(conn *pgxpool.Conn) error {
		existsSQL := `SELECT exists(SELECT FROM documents ` + where + `)`
		if err := conn.QueryRow(rctx, existsSQL, replacements...).Scan(&exists); err != nil {
			return err
		}
		if !exists {
			return nil
		}
		const fetchSQL = `SELECT id, time, commentator, message FROM comments ` +
			`WHERE documents_id = $1 ORDER BY time DESC`
		rows, _ := conn.Query(rctx, fetchSQL, id)
		var err error
		comments, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (comment, error) {
			var com comment
			err := row.Scan(&com.ID, &com.Time, &com.Commentator, &com.Message)
			com.Time = com.Time.UTC()
			return com, err
		})
		return err
	}); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}

	ctx.JSON(http.StatusOK, comments)
}
