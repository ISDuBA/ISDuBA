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
	"github.com/jackc/pgx/v5/pgconn"
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
	var (
		ok   bool
		name string
		url  string
		id   int64
	)
	if name, ok = parse(ctx, notEmpty, ctx.PostForm("name")); !ok {
		return
	}
	if url, ok = parse(ctx, endsWith("/aggregator.json"), ctx.PostForm("url")); !ok {
		return
	}
	const sql = `INSERT INTO aggregators (name, url) VALUES ($1, $2) RETURNING id`
	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			return conn.QueryRow(rctx, sql, name, url).Scan(&id)
		}, 0,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("not a unique value: %v", err.Error()),
			})
		} else {
			slog.Error("inserting aggregator failed", "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"id": id})
}

func (c *Controller) deleteAggregator(ctx *gin.Context) {
	id, ok := parse(ctx, toInt64, ctx.Param("id"))
	if !ok {
		return
	}
	const sql = `DELETE FROM aggregators WHERE id = $1`
	var deleted bool
	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			tag, err := conn.Exec(rctx, sql, id)
			deleted = tag.RowsAffected() > 1
			return err
		}, 0,
	); err != nil {
		slog.Error("delete aggregator failed", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if deleted {
		ctx.JSON(http.StatusOK, gin.H{"message": "deleted"})
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	}
}
