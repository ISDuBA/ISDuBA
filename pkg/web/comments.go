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
	"errors"
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

func (c *Controller) createComment(ctx *gin.Context) {
	docIDs := ctx.Param("document")
	docID, err := strconv.ParseInt(docIDs, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expr := database.FieldEqInt("id", docID)

	// Filter the allowed
	if tlps := c.tlps(ctx); len(tlps) > 0 {
		conditions := tlps.AsConditions()
		parser := database.Parser{}
		tlpExpr, err := parser.Parse(conditions)
		if err != nil {
			slog.Warn("TLP filter failed", "err", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		expr = expr.And(tlpExpr)
	}

	var (
		where, replacements, _ = expr.Where()
		exists                 bool
		commentingAllowed      bool
		forbidden              bool
		commentator            = ctx.GetString("uid")
		message, _             = ctx.GetPostForm("message")
		now                    = time.Now().UTC()
		commentID              *int64
		rctx                   = ctx.Request.Context()
	)

	if err := c.db.Run(rctx, func(rctx context.Context, conn *pgxpool.Conn) error {
		tx, err := conn.BeginTx(rctx, pgx.TxOptions{})
		if err != nil {
			return err
		}
		defer tx.Rollback(rctx)

		stateSQL := `SELECT state, docs.tracking_id, docs.publisher ` +
			`FROM documents docs JOIN advisories ads ` +
			`ON (docs.tracking_id, docs.publisher) = (ads.tracking_id, ads.publisher) ` +
			` WHERE ` + where

		var (
			stateS     string
			trackingID string
			publisher  string
		)
		if err := tx.QueryRow(rctx, stateSQL, replacements...).Scan(
			&stateS, &trackingID, &publisher,
		); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil
			}
			return err
		}
		exists = true

		// Check if we are in a state in which commenting is allowed.
		state := models.Workflow(stateS)
		commentingAllowed = state == models.ReadWorkflow ||
			state == models.AssessingWorkflow
		if !commentingAllowed {
			return nil
		}

		var actor sql.NullString
		if !c.cfg.General.AnonymousEventLogging {
			actor.String = commentator
			actor.Valid = true
		}

		logEvent := func(event models.Event, state models.Workflow) error {
			const eventSQL = `INSERT INTO events_log ` +
				`(event, state, time, actor, documents_id, comments_id) ` +
				`VALUES($1::events, $2::workflow, $3, $4, $5, $6)`
			_, err := tx.Exec(
				rctx, eventSQL, string(event), string(state), now, actor, docID, commentID)
			return err
		}

		// Switch to assessing state if we are not in.
		if state == models.ReadWorkflow {
			// Check if the transition is allowed to user.
			roles := models.ReadWorkflow.TransitionsRoles(models.AssessingWorkflow)
			if !c.hasAnyRole(ctx, roles...) {
				forbidden = true
				return nil
			}

			// Switch to assessing state.
			const assessingStateSQL = `UPDATE advisories SET state = 'assessing' ` +
				`WHERE (tracking_id, publisher) = ($1, $2)`
			if _, err := tx.Exec(rctx, assessingStateSQL, trackingID, publisher); err != nil {
				return err
			}

			// Log that we switched state.
			if err := logEvent(models.StateChangeEvent, models.AssessingWorkflow); err != nil {
				return err
			}
		}

		// Now insert the comment itself
		const insertSQL = `INSERT INTO comments ` +
			`(documents_id, time, commentator, message) ` +
			`VALUES ($1, $2, $3, $4) ` +
			`RETURNING id`

		if err := tx.QueryRow(
			rctx, insertSQL,
			docID, now, commentator, message,
		).Scan(&commentID); err != nil {
			return err
		}

		// Log that we created a comment
		if err := logEvent(models.AddCommentEvent, models.AssessingWorkflow); err != nil {
			return err
		}

		return tx.Commit(rctx)
	}, 0); err != nil {
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	switch {
	case !exists:
		ctx.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
	case !commentingAllowed:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid state to comment"})
	case forbidden:
		ctx.JSON(http.StatusForbidden, gin.H{"error": "user not allowed to change state"})
	default:
		ctx.JSON(http.StatusCreated, gin.H{
			"id":          commentID,
			"time":        now,
			"commentator": commentator,
		})
	}
}

func (c *Controller) updateComment(ctx *gin.Context) {
	commentIDs := ctx.Param("id")
	commentID, err := strconv.ParseInt(commentIDs, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var (
		exists      bool
		now         = time.Now().UTC()
		commentator = ctx.GetString("uid")
		message, _  = ctx.GetPostForm("message")
		rctx        = ctx.Request.Context()
	)
	if err := c.db.Run(rctx, func(rctx context.Context, conn *pgxpool.Conn) error {
		tx, err := conn.BeginTx(rctx, pgx.TxOptions{})
		if err != nil {
			return err
		}
		defer tx.Rollback(rctx)

		const updateSQL = `UPDATE comments ` +
			`SET message = $1 ` +
			`WHERE id = $2 AND commentator = $3 ` +
			`RETURNING documents_id`

		var docID int64
		switch err := tx.QueryRow(
			rctx, updateSQL, message, commentID, commentator,
		).Scan(&docID); {
		case errors.Is(err, pgx.ErrNoRows):
			exists = false
			return nil
		case err != nil:
			return err
		}
		exists = true

		const eventSQL = `INSERT INTO events_log ` +
			`(event, state, time, actor, documents_id, comments_id) ` +
			`VALUES('change_comment', ` +
			`(SELECT state FROM advisories ads JOIN documents docs ` +
			`ON (ads.tracking_id, ads.publisher) = (docs.tracking_id, docs.publisher) ` +
			`WHERE docs.id = $3), ` +
			`$1, $2, $3, $4)`

		var actor sql.NullString
		if !c.cfg.General.AnonymousEventLogging {
			actor.String = commentator
			actor.Valid = true
		}
		if _, err := tx.Exec(rctx, eventSQL, now, actor, docID, commentID); err != nil {
			return err
		}

		return tx.Commit(rctx)
	}, 0); err != nil {
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}

func (c *Controller) viewComments(ctx *gin.Context) {
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

	type comment struct {
		DocumentID  int64     `json:"document_id"`
		ID          int64     `json:"id"`
		Time        time.Time `json:"time"`
		Commentator string    `json:"commentator"`
		Message     string    `json:"message"`
	}

	var comments []comment
	var exists bool

	rctx := ctx.Request.Context()
	if err := c.db.Run(rctx, func(rctx context.Context, conn *pgxpool.Conn) error {
		existsSQL := `SELECT exists(SELECT FROM documents WHERE ` + where + `)`
		if err := conn.QueryRow(
			rctx, existsSQL, replacements...).Scan(&exists); err != nil {
			return err
		}
		if !exists {
			return nil
		}
		const fetchSQL = `SELECT id, documents_id, time, commentator, message FROM comments ` +
			`WHERE documents_id = $1 ORDER BY time DESC`
		rows, _ := conn.Query(rctx, fetchSQL, id)
		var err error
		comments, err = pgx.CollectRows(
			rows,
			func(row pgx.CollectableRow) (comment, error) {
				var com comment
				err := row.Scan(&com.ID, &com.DocumentID, &com.Time, &com.Commentator, &com.Message)
				com.Time = com.Time.UTC()
				return com, err
			})
		return err
	}, 0); err != nil {
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}

	ctx.JSON(http.StatusOK, comments)
}
