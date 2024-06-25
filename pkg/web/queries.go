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
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ISDuBA/ISDuBA/pkg/database"
	"github.com/ISDuBA/ISDuBA/pkg/models"
)

func (c *Controller) createStoredQuery(ctx *gin.Context) {

	query := models.StoredQuery{
		Definer: ctx.GetString("uid"),
	}

	// We need the name.
	if query.Name = ctx.PostForm("name"); query.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "missing 'name'",
		})
		return
	}

	// Advisories flag
	if advisories, ok := ctx.GetPostForm("advisories"); ok {
		var err error
		if query.Advisories, err = strconv.ParseBool(advisories); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "bad 'advisories' value: " + err.Error(),
			})
			return
		}
	}

	// Global flag
	if global, ok := ctx.GetPostForm("global"); ok {
		var err error
		if query.Global, err = strconv.ParseBool(global); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "bad 'global' value: " + err.Error(),
			})
			return
		}
	}
	// Global is only for admins.
	if query.Global && !c.hasAnyRole(ctx, models.Admin) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "global flag can only be used by admins",
		})
		return
	}

	parser := database.Parser{
		Advisory:  query.Advisories,
		Languages: c.cfg.Database.TextSearch,
	}

	// The query to filter the documents.
	query.Query = ctx.DefaultPostForm("query", "true")
	expr, err := parser.Parse(query.Query)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "bad 'query' value: " + err.Error(),
		})
		return
	}
	// In advisory mode we only show the latest.
	if query.Advisories {
		expr = expr.And(database.BoolField("latest"))
	}

	// columns are not optional.
	if query.Columns = strings.Fields(ctx.PostForm("columns")); len(query.Columns) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "missing 'columns' value",
		})
		return
	}

	builder := database.SQLBuilder{Advisory: query.Advisories}
	builder.CreateWhere(expr)
	if err := builder.CheckProjections(query.Columns); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "bad 'columns' value: " + err.Error(),
		})
		return
	}

	// Check if we have orders given.
	if orders, ok := ctx.GetPostForm("order"); ok {
		os := strings.Fields(orders)
		if _, err := builder.CreateOrder(os); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "bad 'order' value: " + err.Error(),
			})
			return
		}
		query.Orders = &os
	}

	// Check if we have a description given.
	if description, ok := ctx.GetPostForm("description"); ok {
		query.Description = &description
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

	var queryID, queryNum int64

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			return conn.QueryRow(rctx, insertSQL,
				query.Advisories,
				query.Definer,
				query.Global,
				query.Name,
				query.Description,
				query.Query,
				query.Columns,
				query.Orders,
			).Scan(&queryID, &queryNum)
		}, 0,
	); err != nil {
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

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			definer := ctx.GetString("uid")
			rows, _ := conn.Query(rctx, selectSQL, definer)
			var err error
			queries, err = pgx.CollectRows(rows,
				func(row pgx.CollectableRow) (*models.StoredQuery, error) {
					var query models.StoredQuery
					if err := row.Scan(
						&query.ID,
						&query.Advisories,
						&query.Definer,
						&query.Global,
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
		}, 0,
	); err != nil {
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
		deleteSQLPrefix  = `DELETE FROM stored_queries WHERE id = $1 AND `
		deleteNoAdminSQL = deleteSQLPrefix + `definer = $2`
		deleteAdminSQL   = deleteSQLPrefix + `(definer = $2 OR global)`
	)

	var tag pgconn.CommandTag

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			// Admins are allowed to delete globals.
			var deleteSQL string
			if c.hasAnyRole(ctx, models.Admin) {
				deleteSQL = deleteAdminSQL
			} else {
				deleteSQL = deleteNoAdminSQL
			}
			definer := ctx.GetString("uid")
			var err error
			tag, err = conn.Exec(rctx, deleteSQL, queryID, definer)
			return err
		}, 0,
	); err != nil {
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

func (c *Controller) fetchStoredQuery(ctx *gin.Context) {

	queryID, err := strconv.ParseInt(ctx.Param("query"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	const selectSQL = `SELECT ` +
		`advisories,` +
		`definer,` +
		`global,` +
		`name,` +
		`description,` +
		`query,` +
		`num,` +
		`columns,` +
		`orders ` +
		`FROM stored_queries WHERE id = $1 AND ` +
		`(global OR definer = $2)`

	query := models.StoredQuery{
		ID: queryID,
	}
	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			definer := ctx.GetString("uid")
			return conn.QueryRow(rctx, selectSQL, queryID, definer).Scan(
				&query.Advisories,
				&query.Definer,
				&query.Global,
				&query.Name,
				&query.Description,
				&query.Query,
				&query.Num,
				&query.Columns,
				&query.Orders,
			)
		}, 0,
	); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		default:
			slog.Error("database error", "err", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, &query)
}

func (c *Controller) updateStoredQuery(ctx *gin.Context) {

	queryID, err := strconv.ParseInt(ctx.Param("query"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	const (
		selectSQLPrefix = `SELECT ` +
			`advisories,` +
			`global,` +
			`name,` +
			`description,` +
			`query,` +
			`num,` +
			`columns,` +
			`orders ` +
			`FROM stored_queries WHERE id = $1 AND `
		selectNoAdminSQL = selectSQLPrefix +
			`definer = $2`
		selectAdminSQL = selectSQLPrefix +
			`(global OR definer = $2)`
	)

	var bad string
	var notFound, unchanged bool

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			tx, err := conn.BeginTx(rctx, pgx.TxOptions{})
			if err != nil {
				return err
			}
			defer tx.Rollback(rctx)

			admin := c.hasAnyRole(ctx, models.Admin)

			var query models.StoredQuery

			var selectSQL string
			if admin {
				selectSQL = selectAdminSQL
			} else {
				selectSQL = selectNoAdminSQL
			}
			definer := ctx.GetString("uid")
			if err := tx.QueryRow(rctx, selectSQL, queryID, definer).Scan(
				&query.Advisories,
				&query.Global,
				&query.Name,
				&query.Description,
				&query.Query,
				&query.Num,
				&query.Columns,
				&query.Orders,
			); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					notFound = true
					return nil
				}
				return err
			}
			var fields []string
			var values []any

			// add tracks the real changes to be stored in the database.
			add := func(changed bool, field string, value any) {
				if changed {
					fields = append(fields, field)
					values = append(values, value)
				}
			}

			// Check advisories
			if advs, ok := ctx.GetPostForm("advisories"); ok {
				advisories, err := strconv.ParseBool(advs)
				if err != nil {
					bad = "bad 'advisories' value: " + err.Error()
					return nil
				}
				add(advisories != query.Advisories, "advisories", advisories)
				// Write back as the advisory mode changes the parser behavior.
				query.Advisories = advisories
			}

			parser := database.Parser{
				Advisory:  query.Advisories,
				Languages: c.cfg.Database.TextSearch,
			}

			// Check query
			var expr *database.Expr
			if qs, ok := ctx.GetPostForm("query"); ok {
				expr, err = parser.Parse(qs)
				add(qs != query.Query, "query", qs)
			} else {
				// We need to re-check if the database value is still valid.
				expr, err = parser.Parse(query.Query)
			}
			if err != nil {
				bad = "bad 'query' value: " + err.Error()
				return nil
			}
			if query.Advisories {
				expr = expr.And(database.BoolField("latest"))
			}

			builder := database.SQLBuilder{Advisory: query.Advisories}
			builder.CreateWhere(expr)

			// Check columns
			if cols, ok := ctx.GetPostForm("columns"); ok {
				columns := strings.Fields(cols)
				if err := builder.CheckProjections(columns); err != nil {
					bad = "bad 'columns' value: " + err.Error()
					return nil
				}
				add(!slices.Equal(columns, query.Columns), "columns", columns)
			}

			// Check global
			if glb, ok := ctx.GetPostForm("global"); ok {
				global, err := strconv.ParseBool(glb)
				if err != nil {
					bad = "bad 'global' value: " + err.Error()
					return nil
				}
				// Only admins are allowed to set global
				if !admin && global {
					bad = "none admins are not allowed to set global"
					return nil
				}
				add(global != query.Global, "global", global)
			}

			// Check num
			if glb, ok := ctx.GetPostForm("num"); ok {
				num, err := strconv.ParseInt(glb, 10, 64)
				if err != nil {
					bad = "bad 'num' value: " + err.Error()
					return nil
				}
				add(num != query.Num, "num", num)
			}

			// Check name
			if nm, ok := ctx.GetPostForm("name"); ok {
				if nm == "" {
					bad = "empty name is not allowed"
					return nil
				}
				add(nm != query.Name, "name", nm)
			}

			// Check description
			if desc, ok := ctx.GetPostForm("description"); ok {
				if desc == "" {
					var s *string
					add(query.Description != nil, "description", s)
				} else {
					add(query.Description == nil || *query.Description != desc, "description", &desc)
				}
			}

			// Check orders
			if os, ok := ctx.GetPostForm("orders"); ok {
				orders := strings.Fields(os)
				if len(orders) == 0 {
					var s *[]string
					add(query.Orders != nil, "orders", s)
				} else {
					if _, err := builder.CreateOrder(orders); err != nil {
						bad = "invalid 'orders' value: " + err.Error()
						return nil
					}
					add(query.Orders == nil || !slices.Equal(*query.Orders, orders), "orders", orders)
				}
			}

			// Only try to update if there are real changes.
			if len(fields) == 0 {
				unchanged = true
				return nil
			}

			values = append(values, queryID)

			var placeholders strings.Builder
			for i := range fields {
				if i > 0 {
					placeholders.WriteByte(',')
				}
				placeholders.WriteByte('$')
				placeholders.WriteString(strconv.Itoa(i + 1))
			}

			// Brackets are only allowed if we have more than one argument.
			var op, cl string
			if len(fields) > 1 {
				op, cl = "(", ")"
			}
			updateSQL := fmt.Sprintf(
				"UPDATE stored_queries SET %[1]s%[2]s%[3]s = %[1]s%[4]s%[3]s WHERE id = $%[5]d",
				op, strings.Join(fields, ","), cl,
				placeholders.String(),
				len(values))

			slog.Debug("update statement", "stmt", updateSQL)

			tag, err := tx.Exec(rctx, updateSQL, values...)
			if err != nil {
				return err
			}
			unchanged = tag.RowsAffected() == 0
			return tx.Commit(rctx)
		}, 0,
	); err != nil {
		// As name and num changes can cause unique constraint violations
		// don't report these not as internal server errors as this expected.
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			bad = "not a unique value: %s" + err.Error()
		} else {
			slog.Error("database error", "err", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	switch {
	case bad != "":
		ctx.JSON(http.StatusBadRequest, gin.H{"error": bad})
	case notFound:
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	case unchanged:
		ctx.JSON(http.StatusOK, gin.H{"message": "unchanged"})
	default:
		ctx.JSON(http.StatusOK, gin.H{"message": "changed"})
	}
}
