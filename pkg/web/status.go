// This file is Free Software under the MIT License
// without warranty, see README.md and LICENSES/MIT.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package web

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ISDuBA/ISDuBA/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (c *Controller) changeStatus(ctx *gin.Context) {

	var (
		publisher  = ctx.Param("publisher")
		trackingID = ctx.Param("trackingid")
		state      = models.Workflow(ctx.Param("state"))
	)

	if !state.Valid() {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid state %q", state)})
		return
	}

	var forbidden, noTransition bool

	rctx := ctx.Request.Context()
	if err := c.db.Run(rctx, func(conn *pgxpool.Conn) error {
		tx, err := conn.BeginTx(rctx, pgx.TxOptions{})
		if err != nil {
			return err
		}
		defer tx.Rollback(rctx)

		var (
			documentID int64
			currentS   string
			tlp        string
		)

		const findAdvisory = `SELECT id, state::text, tlp ` +
			`FROM advisories ads JOIN documents docs ON ads.documents_id = docs.id ` +
			`WHERE publisher = $1 AND tracking_id = $2`

		if err := tx.QueryRow(rctx, findAdvisory, publisher, trackingID).Scan(
			&documentID, &state, &tlp,
		); err != nil {
			return err
		}

		// Check if we are allowed to access it.
		if tlps := c.tlps(ctx); len(tlps) > 0 && !tlps.Allowed(publisher, models.TLP(tlp)) {
			forbidden = true
			return nil
		}

		// Check if the transition is allowed to user.
		current := models.Workflow(currentS)
		roles := current.TransitionsRoles(state)
		if len(roles) == 0 {
			noTransition = true
			return nil
		}
		if !c.hasAnyRole(ctx, roles...) {
			forbidden = true
			return nil
		}

		// At this point the state change can be done.
		const updateState = `UPDATE advisories SET state = $1::workflow WHERE documents_id = $2`
		if _, err := tx.Exec(rctx, updateState, string(state), documentID); err != nil {
			return err
		}

		// Log the event
		const insertLog = `INSERT INTO events_log (event, state, actor, documents_id) ` +
			`VALUES ('state_change', $1::workflow, $2, $3)`

		var actor *string
		if !c.cfg.General.AnonymousEventLogging {
			uid := ctx.GetString("uid")
			actor = &uid
		}
		if _, err := tx.Exec(rctx, insertLog, string(state), actor, documentID); err != nil {
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

	if forbidden {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}
	if noTransition {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "state transition not possible"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "transition done"})
}
