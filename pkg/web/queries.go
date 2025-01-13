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

// createStoredQuery is an endpoint that creates a stored query.
//
//	@Summary		Creates a stored query.
//	@Description	Creates a stored query with the specified configuration.
//	@Param			inputForm	formData	models.StoredQuery	true	"Query configuration"
//	@Accept			multipart/form-data
//	@Produce		json
//	@Success		201	{object}	web.createStoredQuery.createResult
//	@Failure		400	{object}	models.Error
//	@Failure		404	{object}	models.Error
//	@Failure		500	{object}	models.Error
//	@Router			/queries [post]
func (c *Controller) createStoredQuery(ctx *gin.Context) {
	type createResult struct {
		ID  int64 `json:"id"`
		Num int64 `json:"num"`
	}
	sq := models.StoredQuery{
		Definer: ctx.GetString("uid"),
	}

	// We need the name.
	if sq.Name = ctx.PostForm("name"); sq.Name == "" {
		models.SendErrorMessage(ctx, http.StatusBadRequest, "missing 'name'")
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
	if global := ctx.PostForm("global"); global != "" {
		var ok bool
		if sq.Global, ok = parse(ctx, strconv.ParseBool, global); !ok {
			return
		}
	}
	// Global is only for admins.
	if sq.Global && !c.hasAnyRole(ctx, models.Admin) {
		models.SendErrorMessage(ctx, http.StatusForbidden, "global flag can only used by admins")
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
		models.SendErrorMessage(ctx, http.StatusBadRequest, "missing 'columns' value")
		return
	}

	builder := query.SQLBuilder{Mode: sq.Kind}
	builder.CreateWhere(expr)
	if err := builder.CheckProjections(sq.Columns); err != nil {
		models.SendErrorMessage(ctx,
			http.StatusBadRequest,
			"bad 'columns' value: "+err.Error())
		return
	}

	// Check if we have orders given.
	if orders, ok := ctx.GetPostForm("orders"); ok {
		os := strings.Fields(orders)
		if _, err := builder.CreateOrder(os); err != nil {
			models.SendErrorMessage(ctx,
				http.StatusBadRequest,
				"bad 'orders' value"+err.Error())
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
			models.SendErrorMessage(ctx, http.StatusConflict, "already in database")
			return
		}
		slog.Error("database error", "err", err)
		models.SendError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusCreated, createResult{
		ID:  queryID,
		Num: queryNum,
	})
}

// updateOrder is an endpoint that updates the query order.
//
//	@Summary		Updates the query order.
//	@Description	Updates the query order with the specified ordering.
//	@Param			queryOrder	body	web.updateOrder.queryOrder	true	"Query order"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.Success
//	@Failure		500	{object}	models.Error
//	@Router			/queries/orders [post]
func (c *Controller) updateOrder(ctx *gin.Context) {
	type queryOrder struct {
		ID    int64 `json:"id"`
		Order int64 `json:"order"`
	}
	var orders []queryOrder
	if err := ctx.ShouldBindJSON(&orders); err != nil {
		models.SendError(ctx, http.StatusBadRequest, err)
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
		models.SendError(ctx, http.StatusInternalServerError, err)
		return
	}
	models.SendSuccess(ctx, http.StatusOK, "changed")
}

// listStoredQueries is an endpoint that returns all stored queries.
//
//	@Summary		Returns stored queries.
//	@Description	Returns all configured stored queries.
//	@Produce		json
//	@Success		200	{array}		models.StoredQuery
//	@Failure		400	{object}	models.Error
//	@Failure		500	{object}	models.Error
//	@Router			/queries [get]
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
		models.SendError(ctx, http.StatusInternalServerError, err)
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

// listStoredQueries is an endpoint that deletes the specified stored query.
//
//	@Summary		Deletes the stored query.
//	@Description	Deletes the query with the specified ID.
//	@Param			query	path	int	true	"Query ID"
//	@Produce		json
//	@Success		200	{array}		models.Success
//	@Failure		400	{object}	models.Error
//	@Failure		404	{object}	models.Error
//	@Failure		500	{object}	models.Error
//	@Router			/queries [get]
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
		models.SendError(ctx, http.StatusInternalServerError, err)
		return
	}

	if tag.RowsAffected() != 0 {
		models.SendSuccess(ctx, http.StatusOK, "deleted")
	} else {
		models.SendErrorMessage(ctx, http.StatusNotFound, "query not found")
	}
}

// fetchStoredQuery is an endpoint that returns a stored query.
//
//	@Summary		Updates a stored query.
//	@Description	Updates a feed with the specified configuration.
//	@Param			query	path	int	true	"Query ID"
//	@Produce		json
//	@Success		200	{object}	models.StoredQuery
//	@Failure		400	{object}	models.Error
//	@Failure		404	{object}	models.Error
//	@Failure		500	{object}	models.Error
//	@Router			/queries/{id} [get]
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
			models.SendErrorMessage(ctx, http.StatusNotFound, "not found")
		default:
			slog.Error("database error", "err", err)
			models.SendError(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	ctx.JSON(http.StatusOK, &storedQuery)
}

// updateStoredQuery is an endpoint that updates a stored query.
//
//	@Summary		Updates a stored query.
//	@Description	Updates a stored query with the specified configuration.
//	@Param			id		path		int					true	"Query ID"
//	@Param			query	formData	models.StoredQuery	true	"Query configuration"
//	@Produce		json
//	@Success		200	{object}	models.Success
//	@Failure		400	{object}	models.Error
//	@Failure		404	{object}	models.Error
//	@Failure		500	{object}	models.Error
//	@Router			/queries/{id} [put]
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
			`role, ` +
			`definer ` +
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
				&sq.Definer,
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
			if dash := ctx.PostForm("dashboard"); dash != "" {
				dashboard, err := strconv.ParseBool(dash)
				if err != nil {
					bad = "bad 'dashboard' value: " + err.Error()
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
			if glb := ctx.PostForm("global"); glb != "" {
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
				// Queries by "system-default" are not allowed to be set to non-global
				// See explanation in docs/developer/queries.md
				if !global && sq.Definer == "system-default" {
					bad = "global dashboard queries are not allowed to be set to non-global"
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
			models.SendError(ctx, http.StatusInternalServerError, err)
			return
		}
	}
	switch {
	case bad != "":
		models.SendErrorMessage(ctx, http.StatusBadRequest, bad)
	case notFound:
		models.SendErrorMessage(ctx, http.StatusNotFound, "not found")
	case unchanged:
		models.SendSuccess(ctx, http.StatusOK, "unchanged")
	default:
		models.SendSuccess(ctx, http.StatusOK, "changed")
	}
}

// getDefaultQueryExclusion is an endpoint that returns the exclusion list of all queries.
//
//	@Summary		Returns query exclusions.
//	@Description	Returns exclusions of all queries.
//	@Produce		json
//	@Success		200	{object}	models.Success
//	@Failure		404	{object}	models.Error
//	@Failure		500	{object}	models.Error
//	@Router			/queries/ignore/{query} [get]
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
			models.SendErrorMessage(ctx, http.StatusNotFound, "not found")
		default:
			slog.Error("database error", "err", err)
			models.SendError(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	ctx.JSON(http.StatusOK, &ignored)
}

// deleteDefaultQueryExclusion is an endpoint that deletes the exclusion of the query with specified ID.
//
//	@Summary		Deletes query exclusion.
//	@Description	Deletes the query exclusion with the specified ID.
//	@Param			query	path	int	true	"Query ID"
//	@Produce		json
//	@Success		200	{object}	models.Success
//	@Failure		400	{object}	models.Error
//	@Failure		404	{object}	models.Error
//	@Failure		500	{object}	models.Error
//	@Router			/queries/ignore/{query} [delete]
func (c *Controller) deleteDefaultQueryExclusion(ctx *gin.Context) {
	queryID, ok := parse(ctx, toInt64, ctx.Param("query"))
	if !ok {
		return
	}

	// For which user do we want to delete the ignored default queries?
	user := c.currentUser(ctx).String

	deleteSQL := `DELETE FROM default_query_exclusion WHERE "user" = $1 AND id = $2`

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
		models.SendError(ctx, http.StatusInternalServerError, err)
		return
	}

	if tag.RowsAffected() != 0 {
		models.SendSuccess(ctx, http.StatusOK, "deleted")
	} else {
		models.SendErrorMessage(ctx, http.StatusNotFound, "entry not found")
	}
}

// insertDefaultQueryExclusion is an endpoint that ignores the query with specified ID.
//
//	@Summary		Ignores a query.
//	@Description	Ignores the query with the specified ID.
//	@Param			query	path	int	true	"Query ID"
//	@Produce		json
//	@Success		200	{object}	web.insertDefaultQueryExclusion.createResult
//	@Failure		400	{object}	models.Error
//	@Failure		409	{object}	models.Error
//	@Failure		500	{object}	models.Error
//	@Router			/queries/ignore/{query} [post]
func (c *Controller) insertDefaultQueryExclusion(ctx *gin.Context) {
	type createResult struct {
		User string `json:"user"`
		ID   int64  `json:"id"`
	}
	queryID, ok := parse(ctx, toInt64, ctx.Param("query"))
	if !ok {
		return
	}
	user := c.currentUser(ctx).String
	insertSQL := `INSERT INTO default_query_exclusion ("user", id) VALUES ($1, $2) RETURNING "user", id`

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
			models.SendErrorMessage(ctx, http.StatusConflict, "already in database")
			return
		}
		slog.Error("database error", "err", err)
		models.SendError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusCreated, createResult{
		User: insertedUser,
		ID:   insertedID,
	})
}
