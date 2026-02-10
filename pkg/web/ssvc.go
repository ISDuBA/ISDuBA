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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ISDuBA/ISDuBA/pkg/models"
)

// changeSSVC is an endpoint that changes the SSVC of the specified document.
//
//	@Summary		Changes the SSVC.
//	@Description	This updates the SSVC of the specified document.
//	@Param			document	path	int		true	"Document ID"
//	@Param			vector		query	string	true	"SSVC vector"
//	@Produce		json
//	@Success		200	{object}	models.Success
//	@Failure		400	{object}	models.Error
//	@Failure		401
//	@Failure		403	{object}	models.Error
//	@Failure		404	{object}	models.Error
//	@Failure		500	{object}	models.Error
//	@Router			/sources/{document} [delete]
func (c *Controller) changeSSVC(ctx *gin.Context) {
	documentID, ok := parse(ctx, toInt64, ctx.Param("document"))
	if !ok {
		return
	}

	vector := ctx.DefaultQuery("vector", "")
	if err := models.ValidateSSVCv2Vector(vector); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	const (
		//		findSSVC = `SELECT docs.ssvc, ads.tracking_id, ads.publisher, docs.tlp, ads.state::text ` +
		//			`FROM documents docs JOIN advisories ads ` +
		//			`ON docs.advisories_id = ads.id ` +
		//			`WHERE docs.id = $1`
		// First part taken from above
		findSSVC = `SELECT sh.ssvc, ads.tracking_id, ads.publisher, docs.tlp, ads.state::text ` +
			`FROM documents docs JOIN advisories ads ` +
			`ON docs.advisories_id = ads.id ` +
			// LEFT JOIN so we just get an empty ssvc if there is none in the history
			// LATERAL so immediately the latest one is taken
			// Find ssvc from ssvc_history instead
			`LEFT JOIN LATERAL ` +
			`(SELECT ssvc FROM ssvc_history ` +
			// find document
			// use latest with change number as last resort tiebreaker (unique)
			`WHERE documents_id = docs.id ORDER BY changedate DESC, change_number DESC LIMIT 1) ` +
			`sh ON true WHERE docs.id = $1;`
		switchToAssessing = `UPDATE advisories SET state = 'assessing' ` +
			`WHERE (tracking_id, publisher) = ($1, $2)`
		insertLog = `INSERT INTO events_log (event, state, actor, documents_id) ` +
			`VALUES ($1::events, $2::workflow, $3, $4)`
		updateSSVC = `INSERT INTO ssvc_history (actor, documents_id, ssvc, change_number) VALUES ` +
			`($1::varchar, $2::integer, $3, ( ` +
			`SELECT COALESCE(MAX(change_number), 0::bigint) + 1 FROM ssvc_history WHERE documents_id = $2 ` +
			`))`
	)

	var forbidden, unchanged, bad bool

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			tx, err := conn.BeginTx(rctx, pgx.TxOptions{})
			if err != nil {
				return err
			}
			defer tx.Rollback(rctx)

			var (
				ssvc       sql.NullString
				trackingID string
				publisher  string
				tlp        string
				state      string
			)
			if err := tx.QueryRow(rctx, findSSVC, documentID).Scan(
				&ssvc, &trackingID, &publisher, &tlp, &state,
			); err != nil {
				return err
			}

			// check if we are allowed to do
			if tlps := c.tlps(ctx); len(tlps) > 0 && !tlps.Allowed(publisher, models.TLP(tlp)) {
				forbidden = true
				return nil
			}

			// check if it's a real change
			if ssvc.Valid && ssvc.String == vector {
				unchanged = true
				return nil
			}

			// check if we are in a state that allows changing
			if st := models.Workflow(state); st != models.ReadWorkflow && st != models.AssessingWorkflow {
				bad = true
				return nil
			}

			actor := c.currentUser(ctx)
			logEvent := func(event models.Event, state models.Workflow) error {
				_, err := tx.Exec(rctx, insertLog, string(event), string(state), actor, documentID)
				return err
			}

			// If we are in the 'read' state switch to 'assessing'.
			if st := models.Workflow(state); st == models.ReadWorkflow {
				// Check if the transition is allowed to user.
				roles := st.TransitionsRoles(models.AssessingWorkflow)
				if len(roles) == 0 || !c.hasAnyRole(ctx, roles...) {
					forbidden = true
					return nil
				}
				// Do the actual state change
				if _, err := tx.Exec(rctx, switchToAssessing, trackingID, publisher); err != nil {
					return err
				}
				// Log the state change.
				if err := logEvent(models.StateChangeEvent, models.AssessingWorkflow); err != nil {
					return err
				}
			}

			// Now do the actual SSVC update.
			if _, err := tx.Exec(rctx, updateSSVC, actor, documentID, vector); err != nil {
				return err
			}

			return tx.Commit(rctx)
		}, 0,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "advisory not found"})
		} else {
			slog.Error("database error", "err", err)
			models.SendError(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	switch {
	case forbidden:
		models.SendErrorMessage(ctx, http.StatusForbidden, "access denied")
	case unchanged:
		models.SendSuccess(ctx, http.StatusOK, "unchanged")
	case bad:
		models.SendErrorMessage(ctx, http.StatusBadRequest, "unsuited state")
	default:
		models.SendSuccess(ctx, http.StatusOK, "changed")
	}
}

// ViewSSVC is an endpoint that returns the SSVC of the specified document.
//
//	@Summary		Returns the SSVC.
// ToDo: Finish me
func (c *Controller) viewSSVC(ctx *gin.Context) {
	documentID, ok := parse(ctx, toInt64, ctx.Param("document"))
	if !ok {
		return
	}

	const findSSVC = `SELECT sh.ssvc, ads.publisher, docs.tlp ` +
		`FROM documents docs JOIN advisories ads ` +
		`ON docs.advisories_id = ads.id ` +
		`LEFT JOIN LATERAL ` +
		`(SELECT ssvc FROM ssvc_history ` +
		`WHERE documents_id = docs.id ORDER BY changedate DESC, change_number DESC LIMIT 1) ` +
		`sh ON true WHERE docs.id = $1;`

	var (
		forbidden bool
		ssvc      string
	)

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {

			var (
				ssvcdb    sql.NullString
				publisher string
				tlp       string
			)
			if err := conn.QueryRow(rctx, findSSVC, documentID).Scan(
				&ssvcdb, &publisher, &tlp,
			); err != nil {
				return err
			}

			// check if we are allowed to do
			if tlps := c.tlps(ctx); len(tlps) > 0 && !tlps.Allowed(publisher, models.TLP(tlp)) {
				forbidden = true
				return nil
			}

			ssvc = ssvcdb.String
			return nil
		}, 0,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "advisory not found"})
		} else {
			slog.Error("database error", "err", err)
			models.SendError(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	switch {
	case forbidden:
		models.SendErrorMessage(ctx, http.StatusForbidden, "access denied")
	default:
		ctx.JSON(http.StatusOK, gin.H{"ssvc": ssvc})
	}
}

// SSVCHistoryEntry represents a singular ssvc change event
type SSVCHistoryEntry struct {
	SSVC         *string   `json:"ssvc"`
	ChangeDate   time.Time `json:"changedate"`
	ChangeNumber int64     `json:"change_number"`
	Actor        *string   `json:"actor"`
}

// viewSSVCHistory is an endpoint that returns the SSVC History of the specified document.
//
//	@Summary		View the SSVC History.
//	@Description	This returns the SSVC History of the specified document.
// ToDo: Finish me
func (c *Controller) viewSSVCHistory(ctx *gin.Context) {
	documentID, ok := parse(ctx, toInt64, ctx.Param("document"))
	if !ok {
		return
	}
	// fetch access data
	const findPublisherTLP = `SELECT ads.publisher, docs.tlp ` +
		`FROM documents docs  JOIN advisories ads ` +
		`ON docs.advisories_id = ads.id ` +
		`WHERE docs.id = $1 `
	// fetch entire history if exists
	const findSSVCHistory = `SELECT ssvc, changedate, change_number, actor ` +
		`FROM ssvc_history ` +
		`WHERE documents_id = $1 ` +
		`ORDER BY changedate DESC, change_number DESC;`

	var (
		forbidden   bool
		ssvcHistory = []SSVCHistoryEntry{}
	)
	tlps := c.tlps(ctx)
	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {

			var publisher, tlp string
			err := conn.QueryRow(rctx, findPublisherTLP, documentID).Scan(&publisher, &tlp)
			if err != nil {
				return err
			}
			if len(tlps) > 0 && !tlps.Allowed(publisher, models.TLP(tlp)) {
				forbidden = true
				return nil
			}

			rows, err := conn.Query(rctx, findSSVCHistory, documentID)
			if err != nil {
				return err
			}

			defer rows.Close()
			for rows.Next() {
				var (
					ssvc       sql.NullString
					changedate time.Time
					changeNum  int64
					actor      sql.NullString
				)
				if err := rows.Scan(
					&ssvc, &changedate, &changeNum, &actor,
				); err != nil {
					return err
				}

				entry := SSVCHistoryEntry{
					ChangeDate:   changedate,
					ChangeNumber: changeNum,
				}
				// If no ssvc was set
				if ssvc.Valid {
					val := ssvc.String
					entry.SSVC = &val
				}
				// If there's no actor. ToDo: Evaluate if that can happen
				if actor.Valid {
					val := actor.String
					entry.Actor = &val
				}

				ssvcHistory = append(ssvcHistory, entry)
			}
			return rows.Err()
		}, 0,
	); err != nil {
		slog.Error("database error", "err", err)
		models.SendError(ctx, http.StatusInternalServerError, err)
		return
	}
	switch {
	case forbidden:
		// Maybe log attempt?: slog.Warn("unauthorized access attempt", "user", c.currentUser(ctx), "doc_id", documentID)
		models.SendErrorMessage(ctx, http.StatusForbidden, "access denied")
	case len(ssvcHistory) == 0:
		models.SendErrorMessage(ctx, http.StatusNotFound, "No History found")
	default:
		ctx.JSON(http.StatusOK, gin.H{"ssvcHistory": ssvcHistory})
	}
}
