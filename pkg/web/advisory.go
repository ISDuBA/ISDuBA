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
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ISDuBA/ISDuBA/pkg/models"
)

type advisoryStates []models.AdvisoryState

func (c *Controller) changeStatusAll(ctx *gin.Context, inputs advisoryStates) {
	const (
		findAdvisory = `SELECT id, state::text, tlp ` +
			`FROM advisories ads ` +
			`JOIN documents docs ON (ads.tracking_id, ads.publisher) = (docs.tracking_id, docs.Publisher) ` +
			`WHERE docs.publisher = $1 AND docs.tracking_id = $2 ` +
			`and latest`
		updateState = `UPDATE advisories SET state = $1::workflow WHERE (tracking_id, publisher) = ($2, $3)`
		insertLog   = `INSERT INTO events_log (event, state, actor, documents_id) ` +
			`VALUES ('state_change', $1::workflow, $2, $3)`
	)

	actor := c.currentUser(ctx)
	tlps := c.tlps(ctx)

	var forbidden, noTransition, bad bool

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			tx, err := conn.BeginTx(rctx, pgx.TxOptions{})
			if err != nil {
				return err
			}
			defer tx.Rollback(rctx)

			for i := range inputs {
				var (
					input      = &inputs[i]
					documentID int64
					current    string
					tlp        string
				)

				if input.Publisher == "" || input.TrackingID == "" {
					bad = true
					return nil
				}

				slog.Debug("state change",
					"publisher", input.Publisher,
					"tracking_id", input.TrackingID,
					"state", input.State)

				if err := tx.QueryRow(rctx, findAdvisory, input.Publisher, input.TrackingID).Scan(
					&documentID, &current, &tlp,
				); err != nil {
					return err
				}

				// Check if we are allowed to access it.
				if len(tlps) > 0 && !tlps.Allowed(input.Publisher, models.TLP(tlp)) {
					forbidden = true
					return nil
				}

				slog.Debug("current state", "state", current)

				// Check if the transition is allowed to user.
				roles := models.Workflow(current).TransitionsRoles(input.State)
				if len(roles) == 0 {
					noTransition = true
					return nil
				}
				if !c.hasAnyRole(ctx, roles...) {
					forbidden = true
					return nil
				}

				// At this point the state change can be done.
				if _, err := tx.Exec(rctx, updateState,
					string(input.State), input.TrackingID, input.Publisher,
				); err != nil {
					return err
				}

				// Log the event
				if _, err := tx.Exec(rctx, insertLog, string(input.State), actor, documentID); err != nil {
					return err
				}
			}

			return tx.Commit(rctx)
		}, 0,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			models.SendErrorMessage(ctx, http.StatusNotFound, "advisory not found")
		} else {
			slog.Error("state change failed", "err", err)
			models.SendError(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	switch {
	case bad:
		models.SendErrorMessage(ctx, http.StatusBadRequest, "bad input")
	case forbidden:
		models.SendErrorMessage(ctx, http.StatusForbidden, "access denied")
	case noTransition:
		models.SendErrorMessage(ctx, http.StatusBadRequest, "state transition not possible")
	default:
		models.SendSuccess(ctx, http.StatusOK, "transition done")
	}
}

// changeStatus changes the status of an advisory.
//
//	@Summary		Changes the status of an advisory.
//	@Description	Changes the status of the specified advisory, if allowed.
//	@Param			publisher	path	string	true	"Publisher"
//	@Param			trackingid	path	string	true	"Tracking ID"
//	@Param			state	path	string	true	"Advisory status"
//	@Produce		json
//	@Success		200	{object}		models.Success
//	@Failure		400	{object}	models.Error
//	@Failure		401
//	@Failure		403	{object}	models.Error
//	@Failure		404	{object}	models.Error
//	@Failure		500	{object}	models.Error
//	@Router			/status/{publisher}/{trackingid}/{state} [put]
func (c *Controller) changeStatus(ctx *gin.Context) {
	var input models.AdvisoryState
	if err := ctx.ShouldBindUri(&input); err != nil {
		models.SendError(ctx, http.StatusBadRequest, err)
		return
	}
	c.changeStatusAll(ctx, advisoryStates{input})
}

// changeStatusBulk changes the status of multiple advisories.
//
//	@Summary		Bulk changes status.
//	@Description	Changes the status of multiple advisories, if allowed.
//	@Param			input	body	advisoryStates	true	"Advisory states"
//	@Accept		json
//	@Produce		json
//	@Success		200	{object}		models.Success
//	@Failure		400	{object}	models.Error
//	@Failure		401
//	@Failure		403	{object}	models.Error
//	@Failure		404	{object}	models.Error
//	@Failure		500	{object}	models.Error
//	@Router			/status [put]
func (c *Controller) changeStatusBulk(ctx *gin.Context) {
	var inputs advisoryStates
	if err := ctx.ShouldBindJSON(&inputs); err != nil {
		models.SendError(ctx, http.StatusBadRequest, err)
		return
	}
	c.changeStatusAll(ctx, inputs)
}

// deleteAdvisory deletes a given advisory.
//
//	@Summary		Deletes an advisory.
//	@Description	Deletes the specified advisory.
//	@Param			publisher	path	string	true	"Publisher"
//	@Param			trackingid	path	string	true	"Tracking ID"
//	@Produce		json
//	@Success		200	{array}		web.viewEvents.event
//	@Failure		400	{object}	models.Error
//	@Failure		401
//	@Failure		403	{object}	models.Error
//	@Failure		404	{object}	models.Error
//	@Failure		500	{object}	models.Error
//	@Router			/advisory/{publisher}/{trackingid} [delete]
func (c *Controller) deleteAdvisory(ctx *gin.Context) {
	var key models.AdvisoryKey
	if err := ctx.ShouldBindUri(&key); err != nil {
		models.SendError(ctx, http.StatusBadRequest, err)
		return
	}

	if key.Publisher == "" || key.TrackingID == "" {
		models.SendErrorMessage(ctx, http.StatusBadRequest, "missing publisher or tracking_id")
		return
	}

	const (
		tlpSQL = `SELECT tlp ` +
			`FROM advisories ads ` +
			`JOIN documents docs ON (ads.tracking_id, ads.publisher) = (docs.tracking_id, docs.Publisher) ` +
			`WHERE docs.publisher = $1 AND docs.tracking_id = $2 ` +
			`AND latest`
		deleteSQL = `DELETE FROM documents WHERE publisher = $1 AND tracking_id = $2`
	)

	var forbidden, deleted bool

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			tx, err := conn.BeginTx(rctx, pgx.TxOptions{})
			if err != nil {
				return err
			}
			defer tx.Rollback(rctx)

			var tlp string
			if err := tx.QueryRow(rctx, tlpSQL, key.Publisher, key.TrackingID).Scan(&tlp); err != nil {
				return fmt.Errorf("finding latest tlp failed: %w", err)
			}
			if tlps := c.tlps(ctx); !tlps.Allowed(key.Publisher, models.TLP(tlp)) {
				forbidden = true
				return nil
			}

			tags, err := tx.Exec(rctx, deleteSQL, key.Publisher, key.TrackingID)
			if err != nil {
				return fmt.Errorf("deleting advisory documents failed: %w", err)
			}
			// Log if there were documents deleted.
			if deleted = tags.RowsAffected() > 0; deleted {
				actor := c.currentUser(ctx)
				const eventSQL = `INSERT INTO events_log ` +
					`(event, actor) ` +
					`VALUES('delete_document'::events, $1)`
				if _, err := tx.Exec(rctx, eventSQL, actor); err != nil {
					return fmt.Errorf("event logging failed: %w", err)
				}
			}

			return tx.Commit(rctx)
		}, 0,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			models.SendErrorMessage(ctx, http.StatusNotFound, "advisory not found")
		} else {
			slog.Error("deleting advisory failed", "err", err)
			models.SendError(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	switch {
	case forbidden:
		models.SendErrorMessage(ctx, http.StatusForbidden, "not allowed to delete advisory")
	case !deleted:
		models.SendErrorMessage(ctx, http.StatusNotFound, "advisory not found")
	default:
		models.SendSuccess(ctx, http.StatusOK, "advisory deleted")
	}
}
