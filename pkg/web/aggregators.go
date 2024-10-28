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
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (c *Controller) viewAggregators(ctx *gin.Context) {
	type aggregator struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	var list []aggregator
	const sql = `SELECT id, name, url FROM aggregators ORDER by name`
	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			rows, _ := conn.Query(rctx, sql)
			var err error
			list, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (aggregator, error) {
				var a aggregator
				err := row.Scan(&a.ID, &a.Name, &a.URL)
				return a, err
			})
			return err
		}, 0,
	); err != nil {
		slog.Error("fetching aggregators failed", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, list)
}

func (c *Controller) viewAggregator(ctx *gin.Context) {
	// TODO: Implement me!
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Not implemented, yet!"})
}

func (c *Controller) createAggregator(ctx *gin.Context) {
	// TODO: Implement me!
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Not implemented, yet!"})
}

func (c *Controller) deleteAggregator(ctx *gin.Context) {
	// TODO: Implement me!
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Not implemented, yet!"})
}
