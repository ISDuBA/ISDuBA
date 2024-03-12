// This file is Free Software under the MIT License
// without warranty, see README.md and LICENSES/MIT.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package web

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/ISDuBA/ISDuBA/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seal-io/meta-api/ssvc/ssvc2"
)

func (c *Controller) changeSSVC(ctx *gin.Context) {

	documentID, err := strconv.ParseInt(ctx.Param("document"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vector := ctx.DefaultQuery("vector", "")
	if _, err := ssvc2.Parse(vector); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	const (
		findSSVC = `SELECT ssvc, docs.tracking_id, docs.publisher, tlp, state::text ` +
			`FROM documents docs JOIN advisories ads ` +
			`ON (docs.tracking_id, docs.publisher) = (ads.trackingid, ads.publisher)` +
			`WHERE id = $1`
		switchToAssessing = `UPDATE advisories SET state = 'assessing' ` +
			`WHERE (tracking_id, publisher) = ($1, $2)`
		insertLog = `INSERT INTO events_log (event, state, actor, documents_id) ` +
			`VALUES ($1::events, $2::workflow, $3, $4)`
		updateSSVC = `UPDATE documents SET ssvc = $1 WHERE id = $2`
	)

	var forbidden, unchanged, bad bool

	rctx := ctx.Request.Context()
	if err := c.db.Run(rctx, func(conn *pgxpool.Conn) error {
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

		// check if its a real change
		if ssvc.Valid && ssvc.String == vector {
			unchanged = true
			return nil
		}

		// check if we are in a state that allows changing
		if st := models.Workflow(state); st != models.ReadWorkflow && st != models.AssessingWorkflow {
			bad = true
			return nil
		}

		var actor *string
		if !c.cfg.General.AnonymousEventLogging {
			uid := ctx.GetString("uid")
			actor = &uid
		}
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
		if _, err := tx.Exec(rctx, updateSSVC, vector, documentID); err != nil {
			return err
		}

		// Log the SSVC change.
		event := models.ChangeSSVCEvent
		if !ssvc.Valid { // Its new.
			event = models.AddSSVCEvent
		}
		if err := logEvent(event, models.AssessingWorkflow); err != nil {
			return err
		}
		return tx.Commit(rctx)
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "advisory not found"})
		} else {
			slog.Error("state change failed", "err", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	switch {
	case forbidden:
		ctx.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
	case unchanged:
		ctx.JSON(http.StatusOK, gin.H{"message": "unchanged"})
	case bad:
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "unsuited state"})
	default:
		ctx.JSON(http.StatusOK, gin.H{"message": "changed"})
	}
}
