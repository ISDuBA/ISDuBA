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
	"fmt"
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
//	@Failure		401	{object}	models.Error
//	@Failure		403	{object}	models.Error
//	@Failure		404	{object}	models.Error
//	@Failure		500	{object}	models.Error
//	@Router			/ssvc/{document} [put]
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
		updateSSVC = `INSERT INTO ssvc_history (actor, documents_id, ssvc) VALUES ` +
			`($1::varchar, $2::integer, $3)`
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
				&ssvc,
				&trackingID,
				&publisher,
				&tlp,
				&state,
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
			ctx.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
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

// viewSSVC is an endpoint that returns the SSVC of the specified document.
//
//	@Summary		Returns the latest SSVC.
//	@Description	Fetches the most recent SSVC for a specific document.
//	@Produce		json
//	@Param			document	path		int	true	"Document ID"
//	@Success		200			{object}	map[string]string
//	@Failure		403			{object}	models.Error
//	@Failure		404			{object}	models.Error
//	@Failure		500			{object}	models.Error
//	@Router			/ssvc/documents/{document} [get]
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
		`sh ON true WHERE docs.id = $1`

	var (
		forbidden bool
		ssvc      models.SSVCResponse
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
				&ssvcdb,
				&publisher,
				&tlp,
			); err != nil {
				return err
			}

			// check if we are allowed to do
			if tlps := c.tlps(ctx); len(tlps) > 0 && !tlps.Allowed(publisher, models.TLP(tlp)) {
				forbidden = true
				return nil
			}

			if ssvcdb.Valid {
				ssvc.SSVC = &ssvcdb.String
			}
			return nil
		}, 0,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
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
		ctx.JSON(http.StatusOK, ssvc)
	}
}

// viewSSVCHistory is an endpoint that returns the SSVC History of the specified advisory.
//
//	@Summary		View the SSVC History.
//	@Description	Returns a list of all SSVC changes for the specified advisory, ordered by changedate.
//	@Produce		json
//	@Success		200			{object}	map[string][]models.SSVCChange
//	@Param			publisher	path	string	true	"Advisory publisher"
//	@Param			trackingid	path	string	true	"Advisory tracking ID"
//	@Failure		403			{object}	models.Error
//	@Failure		404			{object}	models.Error
//	@Failure		500			{object}	models.Error
//	@Router			/ssvc/history/{publisher}/{trackingid} [get]
func (c *Controller) viewSSVCHistory(ctx *gin.Context) {
	publisherNamespace := ctx.Param("publisher")
	trackingID := ctx.Param("trackingid")

	// fetch access data
	const findPublisherTLP = `WITH advisory_docs AS ( ` +
		`SELECT docs.id, ads.publisher, docs.tlp ` +
		`FROM documents docs ` +
		`JOIN advisories ads ON docs.advisories_id = ads.id ` +
		`WHERE ads.publisher = $1 ` +
		`AND ads.tracking_id = $2 ` +
		`ORDER BY docs.id ` +
		`) ` +
		`SELECT publisher, tlp ` +
		`FROM advisory_docs ` +
		`LIMIT 1;`

	// fetch entire history if exists
	const findSSVCHistory = `WITH advisory_docs AS ( ` +
		`SELECT docs.id, docs.version ` +
		`FROM documents docs ` +
		`JOIN advisories ads ON docs.advisories_id = ads.id ` +
		`WHERE ads.publisher = $1 ` +
		`AND ads.tracking_id = $2 ` +
		`) ` +
		`SELECT h.ssvc, h.changedate, h.change_number, h.actor, h.documents_id, ad.version ` +
		`FROM ssvc_history h ` +
		`JOIN advisory_docs ad ON h.documents_id = ad.id ` +
		`ORDER BY h.documents_id ASC, h.changedate DESC, h.change_number DESC;`

	var (
		forbidden   bool
		ssvcHistory = []models.SSVCHistoryEntry{}
	)
	tlps := c.tlps(ctx)
	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {

			var publisher, tlp string
			err := conn.QueryRow(rctx, findPublisherTLP, publisherNamespace, trackingID).Scan(&publisher, &tlp)
			if err != nil {
				return fmt.Errorf("Failed to find authorization information for queried document: %w", err)
			}
			if len(tlps) > 0 && !tlps.Allowed(publisher, models.TLP(tlp)) {
				forbidden = true
				return nil
			}

			rows, err := conn.Query(rctx, findSSVCHistory, publisherNamespace, trackingID)
			if err != nil {
				return fmt.Errorf("scanning for SSVCHistory failed: %w", err)
			}

			defer rows.Close()

			ssvcHistory, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.SSVCHistoryEntry, error) {
				var entry models.SSVCHistoryEntry
				err := row.Scan(
					&entry.SSVC,
					&entry.ChangeDate,
					&entry.ChangeNumber,
					&entry.Actor,
					&entry.DocumentsID,
					&entry.DocumentsVersion,
				)
				return entry, err
			})
			if err != nil {
				return fmt.Errorf("collecting SSVCHistory failed: %w", err)
			}

			return nil
		}, 0,
	); err != nil {
		slog.Error("database error", "err", err)
		models.SendError(ctx, http.StatusInternalServerError, err)
		return
	}
	switch {
	case forbidden:
		models.SendErrorMessage(ctx, http.StatusForbidden, "access denied")
	case len(ssvcHistory) == 0:
		models.SendErrorMessage(ctx, http.StatusNotFound, "no History found")
	default:
		// also return the previous ssvc for easier consumption
		ssvcChanges := buildSSVCChange(ssvcHistory)
		ctx.JSON(http.StatusOK, gin.H{"ssvcChanges": ssvcChanges})
	}
}

// buildSSVCChange turns an array of SSVCHistoryEntry's into an Array of SSVCChange's by looking up the last ssvc if it's not the oldest entry.
func buildSSVCChange(history []models.SSVCHistoryEntry) []models.SSVCChange {
	var changes []models.SSVCChange
	// Iterate over history
	for i, entry := range history {
		// if there exists an older entry for the same document
		// Only works if history is grouped by documents then sorted by time from new to old
		if i+1 < len(history) && history[i+1].DocumentsID == entry.DocumentsID {
			changes = append(changes, models.SSVCChange{
				SSVCHistoryEntry: entry,
				SSVCPrev:         history[i+1].SSVC,
			})
			// Else -> Oldest entry for this document -> Add SSVC event
		} else {
			changes = append(changes, models.SSVCChange{
				SSVCHistoryEntry: entry,
			})
		}
	}
	return changes
}
