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
		findSSVC = `SELECT docs.ssvc, ads.tracking_id, ads.publisher, docs.tlp, ads.state::text ` +
			`FROM documents docs JOIN advisories ads ` +
			`ON docs.advisories_id = ads.id ` +
			`WHERE docs.id = $1`
		switchToAssessing = `UPDATE advisories SET state = 'assessing' ` +
			`WHERE (tracking_id, publisher) = ($1, $2)`
		insertLog = `INSERT INTO events_log (event, state, actor, documents_id, prev_ssvc) ` +
			`VALUES ($1::events, $2::workflow, $3, $4, $5)`
		updateSSVC = `UPDATE documents SET ssvc = $1 WHERE id = $2`
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
			logEvent := func(event models.Event, state models.Workflow, prevSSVC *string) error {
				_, err := tx.Exec(rctx, insertLog, string(event), string(state), actor, documentID, prevSSVC)
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
				if err := logEvent(models.StateChangeEvent, models.AssessingWorkflow, nil); err != nil {
					return err
				}
			}

			// Now do the actual SSVC update.
			if _, err := tx.Exec(rctx, updateSSVC, vector, documentID); err != nil {
				return err
			}

			// Log the SSVC change.
			event := models.ChangeSSVCEvent
			var prevSSVC *string
			if !ssvc.Valid { // It's new.
				event = models.AddSSVCEvent
			} else {
				prevSSVC = &ssvc.String
			}
			if err := logEvent(event, models.AssessingWorkflow, prevSSVC); err != nil {
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
