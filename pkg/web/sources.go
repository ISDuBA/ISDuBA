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
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (c *Controller) viewSources(ctx *gin.Context) {

	type source struct {
		ID     int64    `json:"id"`
		Name   string   `json:"name"`
		Domain string   `json:"domain,omitempty"`
		PMD    string   `json:"pmd,omitempty"`
		Active bool     `json:"active"`
		Rate   *float64 `json:"rate,omitempty"`
		Slots  *int     `json:"slots,omitempty"`
	}

	var srcs []*source
	const sql = `SELECT id, name, domain, pmd, active, rate, slots FROM sources`

	if err := c.db.Run(ctx.Request.Context(), func(rctx context.Context, con *pgxpool.Conn) error {
		rows, err := con.Query(rctx, sql)
		if err != nil {
			return fmt.Errorf("failed fetching sources: %w", err)
		}
		srcs, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (*source, error) {
			var src source
			return &src, row.Scan(
				&src.ID,
				&src.Name,
				&src.Domain,
				&src.PMD,
				&src.Active,
				&src.Rate,
				&src.Slots,
			)
		})
		return err
	}, 0); err != nil {
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"sources": srcs})
}

func (c *Controller) createSource(ctx *gin.Context) {
	// TODO: Implement me!
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": "'createSource' not implemented, yet.",
	})
}

func (c *Controller) deleteSource(ctx *gin.Context) {
	// TODO: Implement me!
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": "'deleteSource' not implemented, yet.",
	})
}

func (c *Controller) updateSource(ctx *gin.Context) {
	// TODO: Implement me!
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": "'updateSource' not implemented, yet.",
	})
}

func (c *Controller) viewFeeds(ctx *gin.Context) {
	// TODO: Implement me!
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": "'viewFeeds' not implemented, yet.",
	})
}

func (c *Controller) createFeed(ctx *gin.Context) {
	// TODO: Implement me!
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": "'createFeed' not implemented, yet.",
	})
}

func (c *Controller) viewFeed(ctx *gin.Context) {
	// TODO: Implement me!
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": "'viewFeed' not implemented, yet.",
	})
}

func (c *Controller) deleteFeed(ctx *gin.Context) {
	// TODO: Implement me!
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": "'deleteFeed' not implemented, yet.",
	})
}

func (c *Controller) feedLog(ctx *gin.Context) {
	// TODO: Implement me!
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": "'feedLog' not implemented, yet.",
	})
}

// defaultMessage returns the default message.
func (c *Controller) defaultMessage(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": c.cfg.Sources.DefaultMessage})
}
