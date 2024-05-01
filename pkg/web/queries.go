// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package web

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/ISDuBA/ISDuBA/pkg/database"
	"github.com/ISDuBA/ISDuBA/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (c *Controller) createStoredQuery(ctx *gin.Context) {

	var (
		name        string
		definer     string
		advisories  bool
		global      bool
		query       string
		columns     []string
		orders      *[]string
		description *string
	)

	// We need the name.
	if name = ctx.PostForm("name"); name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "missing 'name'",
		})
		return
	}

	// Advisories flag
	if advisoriesS, ok := ctx.GetPostForm("advisories"); ok {
		var err error
		if advisories, err = strconv.ParseBool(advisoriesS); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "bad 'advisories' value: " + err.Error(),
			})
			return
		}
	}

	// Global flag
	if globalS, ok := ctx.GetPostForm("global"); ok {
		var err error
		if global, err = strconv.ParseBool(globalS); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "bad 'global' value: " + err.Error(),
			})
			return
		}
	}
	// Global is only for admins.
	if global && !c.hasAnyRole(ctx, models.Admin) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "global flag can only be used by admins",
		})
		return
	}

	parser := database.Parser{
		Advisory:  advisories,
		Languages: c.cfg.Database.TextSearch,
	}

	// The query to filter the documents.
	query = ctx.DefaultPostForm("query", "true")
	expr, err := parser.Parse(query)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "bad 'query' value: " + err.Error(),
		})
		return
	}
	// In advisory mode we only show the latest.
	if advisories {
		expr = expr.And(database.BoolField("latest"))
	}

	// columns are not optional.
	if columns = strings.Fields(ctx.PostForm("columns")); len(columns) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "missing 'columns' value",
		})
		return
	}

	_, _, aliases := expr.Where()
	if err := database.CheckProjections(columns, aliases, advisories); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "bad 'columns' value: " + err.Error(),
		})
		return
	}

	// Check if we have orders given.
	if ordersS, ok := ctx.GetPostForm("order"); ok {
		os := strings.Fields(ordersS)
		if _, err := database.CreateOrder(os, aliases, advisories); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "bad 'order' value: " + err.Error(),
			})
			return
		}
		orders = &os
	}

	// Check if we have a description given.
	if descriptionS, ok := ctx.GetPostForm("description"); ok {
		description = &descriptionS
	}

	const insertSQL = `INSERT INTO stored_queries (` +
		`advisories,` +
		`definer,` +
		`global,` +
		`name,` +
		`description,` +
		`query,` +
		`columns,` +
		`orders` +
		`) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)` +
		`RETURNING id, num`

	definer = ctx.GetString("uid")

	var queryID, queryNum int64

	rctx := ctx.Request.Context()
	if err := c.db.Run(rctx, func(conn *pgxpool.Conn) error {
		return conn.QueryRow(rctx, insertSQL,
			advisories, definer, global, name, description,
			query, columns, orders).Scan(&queryID, &queryNum)
	}); err != nil {
		var pgErr *pgconn.PgError
		// Unique constraint violation
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			ctx.JSON(http.StatusConflict, gin.H{"error": "already in database"})
			return
		}
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"id":  queryID,
		"num": queryNum,
	})
}

func (c *Controller) listStoredQueries(ctx *gin.Context) {

	const selectSQL = `SELECT ` +
		`id,` +
		`advisories,` +
		`definer,` +
		`global,` +
		`name,` +
		`description,` +
		`query,` +
		`num,` +
		`columns,` +
		`orders ` +
		`FROM stored_queries WHERE ` +
		`definer = $1 OR global ` +
		`ORDER BY global desc, definer, num`

	var queries []*models.StoredQuery

	rctx := ctx.Request.Context()
	if err := c.db.Run(rctx, func(conn *pgxpool.Conn) error {
		definer := ctx.GetString("uid")
		rows, _ := conn.Query(rctx, selectSQL, definer)
		var err error
		queries, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (*models.StoredQuery, error) {
			var query models.StoredQuery
			if err := row.Scan(
				&query.ID,
				&query.Advisories,
				&query.Definer,
				&query.Gobal,
				&query.Name,
				&query.Description,
				&query.Query,
				&query.Num,
				&query.Columns,
				&query.Orders,
			); err != nil {
				return nil, err
			}
			return &query, nil
		})
		return err
	}); err != nil {
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, queries)
}

func (c *Controller) deleteStoredQuery(ctx *gin.Context) {

	queryID, err := strconv.ParseInt(ctx.Param("query"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	const (
		deleteNoAdminSQL = `DELETE FROM stored_queries WHERE id = $1 AND definer = $2`
		deleteAdminSQL   = `DELETE FROM stored_queries WHERE id = $1 AND global`
	)

	var tag pgconn.CommandTag

	rctx := ctx.Request.Context()
	if err := c.db.Run(rctx, func(conn *pgxpool.Conn) error {
		var err error
		// Admins are allowed to delete globals.
		if c.hasAnyRole(ctx, models.Admin) {
			tag, err = conn.Exec(rctx, deleteAdminSQL, queryID)
		} else {
			definer := ctx.GetString("uid")
			tag, err = conn.Exec(rctx, deleteNoAdminSQL, queryID, definer)
		}
		return err
	}); err != nil {
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if tag.RowsAffected() != 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": "deleted"})
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "query not found"})
	}
}

func (c *Controller) updateStoredQuery(ctx *gin.Context) {
	// TODO: Implement me!
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Not implemented, yet"})
}
