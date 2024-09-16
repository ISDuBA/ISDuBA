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

	"github.com/ISDuBA/ISDuBA/pkg/database/query"
	"github.com/ISDuBA/ISDuBA/pkg/models"
)

func (c *Controller) createStoredQuery(ctx *gin.Context) {
	sq := models.StoredQuery{
		Definer: ctx.GetString("uid"),
	}

	// We need the name.
	if sq.Name = ctx.PostForm("name"); sq.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "missing 'name'",
		})
		return
	}

	// Advisories flag
	if kind, ok := ctx.GetPostForm("kind"); ok {
		if sq.Kind, ok = parse(ctx, parserMode, kind); !ok {
			return
		}
	}

	// Dashboard flag
	if dashboard, ok := ctx.GetPostForm("dashboard"); ok {
		if sq.Dashboard, ok = parse(ctx, strconv.ParseBool, dashboard); !ok {
			return
		}
	}

	// Role
	if role := ctx.PostForm("role"); role != "" {
		wfr, ok := parse(ctx, models.ParseWorkflowRole, role)
		if !ok {
			return
		}
		sq.Role = &wfr
	}

	// Global flag
	if global, ok := ctx.GetPostForm("global"); ok {
		if sq.Global, ok = parse(ctx, strconv.ParseBool, global); !ok {
			return
		}
	}
	// Global is only for admins.
	if sq.Global && !c.hasAnyRole(ctx, models.Admin) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "global flag can only be used by admins",
		})
		return
	}

	parser := query.Parser{Mode: sq.Kind}

	// The query to filter the documents.
	sq.Query = ctx.DefaultPostForm("query", "true")
	expr, ok := parse(ctx, parser.Parse, sq.Query)
	if !ok {
		return
	}
	// In advisory mode we only show the latest.
	if sq.Kind == query.AdvisoryMode {
		expr = expr.And(query.BoolField("latest"))
	}

	// columns are not optional.
	if sq.Columns = strings.Fields(ctx.PostForm("columns")); len(sq.Columns) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "missing 'columns' value",
		})
		return
	}

	builder := query.SQLBuilder{Mode: sq.Kind}
	builder.CreateWhere(expr)
	if err := builder.CheckProjections(sq.Columns); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "bad 'columns' value: " + err.Error(),
		})
		return
	}

	// Check if we have orders given.
	if orders, ok := ctx.GetPostForm("orders"); ok {
		os := strings.Fields(orders)
		if _, err := builder.CreateOrder(os); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "bad 'orders' value: " + err.Error(),
			})
			return
		}
		sq.Orders = &os
	}

	// Check if we have a description given.
	if description, ok := ctx.GetPostForm("description"); ok {
		sq.Description = &description
	}

	const insertSQL = `INSERT INTO stored_queries (` +
		`kind,` +
		`definer,` +
		`global,` +
		`name,` +
		`description,` +
		`query,` +
		`columns,` +
		`orders,` +
		`dashboard,` +
		`role ` +
		`) VALUES ($1::stored_queries_kind, $2, $3, $4, $5, $6, $7, $8, $9, $10)` +
		`RETURNING id, num`

	var queryID, queryNum int64

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			return conn.QueryRow(rctx, insertSQL,
				sq.Kind.String(),
				sq.Definer,
				sq.Global,
				sq.Name,
				sq.Description,
				sq.Query,
				sq.Columns,
				sq.Orders,
				sq.Dashboard,
				sq.Role,
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

func (c *Controller) updateOrder(ctx *gin.Context) {
	type queryOrder struct {
		ID    int64 `json:"id"`
		Order int64 `json:"order"`
	}
	var orders []queryOrder
	if err := ctx.ShouldBindJSON(&orders); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	const (
		prefix   = `UPDATE stored_queries SET num = $2 WHERE id = $1 AND `
		adminSQL = prefix + `(definer = $3 OR global)`
		userSQL  = prefix + `definer = $3`
	)
	var updateSQL string
	if c.hasAnyRole(ctx, models.Admin) {
		updateSQL = adminSQL
	} else {
		updateSQL = userSQL
	}

	batch := &pgx.Batch{}
	definer := ctx.GetString("uid")
	for _, order := range orders {
		batch.Queue(updateSQL, order.ID, order.Order, definer)
	}

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			tx, err := conn.Begin(rctx)
			if err != nil {
				return err
			}
			defer tx.Rollback(rctx)
			if err := tx.SendBatch(rctx, batch).Close(); err != nil {
				return err
			}
			return tx.Commit(rctx)
		}, 0,
	); err != nil {
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "changed"})
}

func (c *Controller) listStoredQueries(ctx *gin.Context) {
	const selectSQL = `SELECT ` +
		`id,` +
		`kind::text,` +
		`definer,` +
		`global,` +
		`name,` +
		`description,` +
		`query,` +
		`num,` +
		`columns,` +
		`orders,` +
		`dashboard,` +
		`role ` +
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
					var storedQuery models.StoredQuery
					if err := row.Scan(
						&storedQuery.ID,
						&storedQuery.Kind,
						&storedQuery.Definer,
						&storedQuery.Global,
						&storedQuery.Name,
						&storedQuery.Description,
						&storedQuery.Query,
						&storedQuery.Num,
						&storedQuery.Columns,
						&storedQuery.Orders,
						&storedQuery.Dashboard,
						&storedQuery.Role,
					); err != nil {
						return nil, err
					}
					return &storedQuery, nil
				})
			return err
		}, 0,
	); err != nil {
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// remove queries that should only be viewable by other roles
	queries = slices.DeleteFunc(queries, func(query *models.StoredQuery) bool {
		if !query.Global || query.Role == nil {
			return false
		}
		return !c.hasAnyRole(ctx, *query.Role, models.Admin)
	})

	ctx.JSON(http.StatusOK, queries)
}

func (c *Controller) deleteStoredQuery(ctx *gin.Context) {
	queryID, ok := parse(ctx, toInt64, ctx.Param("query"))
	if !ok {
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
	queryID, ok := parse(ctx, toInt64, ctx.Param("query"))
	if !ok {
		return
	}

	const selectSQL = `SELECT ` +
		`kind::text,` +
		`definer,` +
		`global,` +
		`name,` +
		`description,` +
		`query,` +
		`num,` +
		`columns,` +
		`orders,` +
		`dashboard,` +
		`role ` +
		`FROM stored_queries WHERE id = $1 AND ` +
		`(global OR definer = $2)`

	storedQuery := models.StoredQuery{
		ID: queryID,
	}
	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			definer := ctx.GetString("uid")
			return conn.QueryRow(rctx, selectSQL, queryID, definer).Scan(
				&storedQuery.Kind,
				&storedQuery.Definer,
				&storedQuery.Global,
				&storedQuery.Name,
				&storedQuery.Description,
				&storedQuery.Query,
				&storedQuery.Num,
				&storedQuery.Columns,
				&storedQuery.Orders,
				&storedQuery.Dashboard,
				&storedQuery.Role,
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
	ctx.JSON(http.StatusOK, &storedQuery)
}

func (c *Controller) updateStoredQuery(ctx *gin.Context) {
	queryID, ok := parse(ctx, toInt64, ctx.Param("query"))
	if !ok {
		return
	}

	const (
		selectSQLPrefix = `SELECT ` +
			`kind::text,` +
			`global,` +
			`name,` +
			`description,` +
			`query,` +
			`num,` +
			`columns,` +
			`orders,` +
			`dashboard,` +
			`role ` +
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

			var sq models.StoredQuery

			var selectSQL string
			if admin {
				selectSQL = selectAdminSQL
			} else {
				selectSQL = selectNoAdminSQL
			}
			definer := ctx.GetString("uid")
			if err := tx.QueryRow(rctx, selectSQL, queryID, definer).Scan(
				&sq.Kind,
				&sq.Global,
				&sq.Name,
				&sq.Description,
				&sq.Query,
				&sq.Num,
				&sq.Columns,
				&sq.Orders,
				&sq.Dashboard,
				&sq.Role,
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
			if kind, ok := ctx.GetPostForm("kind"); ok {
				var pm query.ParserMode
				if err := pm.UnmarshalText([]byte(kind)); err != nil {
					bad = "bad 'kind' value: " + err.Error()
					return nil
				}
				add(pm != sq.Kind, "kind", kind)
				// Write back as the advisory mode changes the parser behavior.
				sq.Kind = pm
			}

			parser := query.Parser{Mode: sq.Kind}

			// Check query
			var expr *query.Expr
			if qs, ok := ctx.GetPostForm("query"); ok {
				expr, err = parser.Parse(qs)
				add(qs != sq.Query, "query", qs)
			} else {
				// We need to re-check if the database value is still valid.
				expr, err = parser.Parse(sq.Query)
			}
			if err != nil {
				bad = "bad 'query' value: " + err.Error()
				return nil
			}
			if sq.Kind == query.AdvisoryMode {
				expr = expr.And(query.BoolField("latest"))
			}

			builder := query.SQLBuilder{Mode: sq.Kind}
			builder.CreateWhere(expr)

			// Check columns
			if cols, ok := ctx.GetPostForm("columns"); ok {
				columns := strings.Fields(cols)
				if err := builder.CheckProjections(columns); err != nil {
					bad = "bad 'columns' value: " + err.Error()
					return nil
				}
				add(!slices.Equal(columns, sq.Columns), "columns", columns)
			}

			// Check dashboard
			if dash, ok := ctx.GetPostForm("dashboard"); ok {
				dashboard, err := strconv.ParseBool(dash)
				if err != nil {
					bad = "bad 'global' value: " + err.Error()
					return nil
				}
				add(dashboard != sq.Dashboard, "dashboard", dashboard)
			}

			// Check role
			if role, ok := ctx.GetPostForm("role"); ok {
				if role == "" {
					add(sq.Role != nil, "role", nil)
				} else {
					wfr, err := models.ParseWorkflowRole(role)
					if err != nil {
						bad = "bad 'role' value: " + err.Error()
						return nil
					}
					add(sq.Role == nil || wfr != *sq.Role, "role", role)
				}
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
				add(global != sq.Global, "global", global)
			}

			// Check num
			if glb, ok := ctx.GetPostForm("num"); ok {
				num, err := strconv.ParseInt(glb, 10, 64)
				if err != nil {
					bad = "bad 'num' value: " + err.Error()
					return nil
				}
				add(num != sq.Num, "num", num)
			}

			// Check name
			if nm, ok := ctx.GetPostForm("name"); ok {
				if nm == "" {
					bad = "empty name is not allowed"
					return nil
				}
				add(nm != sq.Name, "name", nm)
			}

			// Check description
			if desc, ok := ctx.GetPostForm("description"); ok {
				if desc == "" {
					var s *string
					add(sq.Description != nil, "description", s)
				} else {
					add(sq.Description == nil || *sq.Description != desc, "description", &desc)
				}
			}

			// Check orders
			if os, ok := ctx.GetPostForm("orders"); ok {
				orders := strings.Fields(os)
				if len(orders) == 0 {
					var s *[]string
					add(sq.Orders != nil, "orders", s)
				} else {
					if _, err := builder.CreateOrder(orders); err != nil {
						bad = "invalid 'orders' value: " + err.Error()
						return nil
					}
					add(sq.Orders == nil || !slices.Equal(*sq.Orders, orders), "orders", orders)
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

func (c *Controller) getDefaultQueryExclusion(ctx *gin.Context) {

	// For which user do we want to get the ignored default queries?
	user := c.currentUser(ctx).String

	const selectSQL = `SELECT ` +
		`id ` +
		`FROM default_query_exclusion WHERE "user" = $1`
	var ignored []int
	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			rows, _ := conn.Query(rctx, selectSQL, user)
			var err error
			ignored, err = pgx.CollectRows(
				rows,
				func(row pgx.CollectableRow) (int, error) {
					var id int
					err := row.Scan(&id)
					return id, err
				})
			return err
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
	ctx.JSON(http.StatusOK, &ignored)
}

func (c *Controller) deleteDefaultQueryExclusion(ctx *gin.Context) {
	queryID, ok := parse(ctx, toInt64, ctx.Param("query"))
	if !ok {
		return
	}

	// For which user do we want to delete the ignored default queries?
	user := c.currentUser(ctx).String

	var deleteSQL = `DELETE FROM default_query_exclusion WHERE "user" = $1 AND id = $2`

	var tag pgconn.CommandTag

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			// Admins are allowed to delete globals.
			var err error
			tag, err = conn.Exec(rctx, deleteSQL, user, queryID)
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
		ctx.JSON(http.StatusNotFound, gin.H{"error": "entry not found"})
	}
}

func (c *Controller) insertDefaultQueryExclusion(ctx *gin.Context) {
	queryID, ok := parse(ctx, toInt64, ctx.Param("query"))
	if !ok {
		return
	}
	user := c.currentUser(ctx).String
	var insertSQL = `INSERT INTO default_query_exclusion ("user", id) VALUES ($1, $2) RETURNING "user", id`

	var insertedUser string
	var insertedID int64

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			return conn.QueryRow(rctx, insertSQL,
				user,
				queryID,
			).Scan(&insertedUser, &insertedID)
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
		"user": insertedUser,
		"id":   insertedID,
	})
}
