// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package web

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ISDuBA/ISDuBA/pkg/database/query"
	"github.com/ISDuBA/ISDuBA/pkg/models"
)

// deleteDocument is an end point for deleting a document.
func (c *Controller) deleteDocument(ctx *gin.Context) {
	// Get an ID from context
	idS := ctx.Param("id")
	docID, err := strconv.ParseInt(idS, 10, 64)
	// Error handling for id acquisition
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	// FieldEqInt is a shortcut mainly for building expressions
	// accessing an integer column like 'id's.
	// Expr encapsulates a parsed expression to be converted to an SQL WHERE clause.
	expr := query.FieldEqInt("id", docID)

	// Filter the allowed
	if tlps := c.tlps(ctx); len(tlps) > 0 {
		tlpExpr := tlps.AsExpr()
		expr = expr.And(tlpExpr)
	}

	builder := query.SQLBuilder{}
	builder.CreateWhere(expr)

	deleted := false

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			tx, err := conn.BeginTx(rctx, pgx.TxOptions{})
			if err != nil {
				return err
			}
			defer tx.Rollback(rctx)

			const deletePrefix = `DELETE FROM documents WHERE `
			deleteSQL := deletePrefix + builder.WhereClause
			slog.Debug("delete document", "SQL",
				qndSQLReplace(deleteSQL, builder.Replacements))

			tags, err := tx.Exec(rctx, deleteSQL, builder.Replacements...)
			if err != nil {
				return fmt.Errorf("delete failed: %w", err)
			}

			if deleted = tags.RowsAffected() > 0; deleted {
				actor := c.currentUser(ctx)
				const eventSQL = `INSERT INTO events_log ` +
					`(event, actor) ` +
					`VALUES('delete_document'::events, $1)`
				if _, err := tx.Exec(rctx, eventSQL, actor); err != nil {
					return fmt.Errorf("event logging failed: %w", err)
				}
			}

			return tx.Commit(rctx)
		}, 0,
	); err != nil {
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if deleted {
		ctx.JSON(http.StatusOK, gin.H{"message": "document deleted"})
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
	}
}

// importDocument is an end point to import a document.
func (c *Controller) importDocument(ctx *gin.Context) {

	var actor *string
	if user := c.currentUser(ctx); user.Valid {
		actor = &user.String
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	f, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	limited := http.MaxBytesReader(
		ctx.Writer, f, int64(c.cfg.General.AdvisoryUploadLimit))
	defer limited.Close()

	var id int64

	if err = c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			id, err = models.ImportDocument(
				rctx, conn, limited, actor, c.tlps(ctx), false)
			return err
		}, 0,
	); err != nil {
		switch {
		case errors.Is(err, models.ErrAlreadyInDatabase):
			ctx.JSON(http.StatusConflict, gin.H{"error": "already in database"})
		case errors.Is(err, models.ErrNotAllowed):
			ctx.JSON(http.StatusForbidden, gin.H{"error": "wrong publisher/tlp"})
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": id})
}

// viewDocument is an end point to export a document.
func (c *Controller) viewDocument(ctx *gin.Context) {
	idS := ctx.Param("id")
	id, err := strconv.ParseInt(idS, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	expr := query.FieldEqInt("id", id)

	// Filter the allowed
	if tlps := c.tlps(ctx); len(tlps) > 0 {
		tlpExpr := tlps.AsExpr()
		expr = expr.And(tlpExpr)
	}

	fields := []string{"original"}
	builder := query.SQLBuilder{}
	builder.CreateWhere(expr)
	sql := builder.CreateQuery(fields, "", -1, -1)

	var original []byte

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			return conn.QueryRow(rctx, sql, builder.Replacements...).Scan(&original)
		}, 0,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}
		return
	}

	extraHeaders := map[string]string{
		"Content-Disposition": `attachment; filename="document.json"`,
	}

	ctx.DataFromReader(
		http.StatusOK, int64(len(original)),
		"application/json",
		bytes.NewReader(original),
		extraHeaders)
}

// overviewDocuments is an end point to return an overview document.
func (c *Controller) overviewDocuments(ctx *gin.Context) {

	// Use the advisories.
	advisoryS := ctx.DefaultQuery("advisories", "false")
	advisory, err := strconv.ParseBool(advisoryS)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	parser := query.Parser{
		Advisory:  advisory,
		Languages: c.cfg.Database.TextSearch,
		Me:        ctx.GetString("uid"),
	}

	// The query to filter the documents.
	expr, err := parser.Parse(ctx.DefaultQuery("query", "true"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Filter the allowed
	if tlps := c.tlps(ctx); len(tlps) > 0 {
		tlpExpr := tlps.AsExpr()
		expr = expr.And(tlpExpr)
	}

	// In advisory mode we only show the latest.
	if advisory {
		expr = expr.And(query.BoolField("latest"))
	}

	builder := query.SQLBuilder{Advisory: advisory}
	builder.CreateWhere(expr)

	fields := strings.Fields(
		ctx.DefaultQuery("columns", "id title tracking_id version publisher"))

	if err := builder.CheckProjections(fields); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderFields := strings.Fields(
		ctx.DefaultQuery("order", "publisher tracking_id -current_release_date -rev_history_length"))
	order, err := builder.CreateOrder(orderFields)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var (
		calcCount     bool
		count         int64
		limit, offset int64 = -1, -1
		results       []map[string]any
	)

	if count := ctx.Query("count"); count != "" {
		calcCount = true
	}

	if lim := ctx.Query("limit"); lim != "" {
		limit, err = strconv.ParseInt(lim, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if ofs := ctx.Query("offset"); ofs != "" {
		offset, err = strconv.ParseInt(ofs, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			if calcCount {
				if err := conn.QueryRow(
					rctx,
					builder.CreateCountSQL(),
					builder.Replacements...,
				).Scan(&count); err != nil {
					return fmt.Errorf("cannot calculate count %w", err)
				}
			}
			// Skip fields if they are not requested.
			if len(fields) == 0 {
				return nil
			}

			sql := builder.CreateQuery(fields, order, limit, offset)

			values := make([]any, len(fields))
			ptrs := make([]any, len(fields))
			for i := range ptrs {
				ptrs[i] = &values[i]
			}

			if slog.Default().Enabled(rctx, slog.LevelDebug) {
				slog.Debug("documents", "SQL", qndSQLReplace(sql, builder.Replacements))
			}
			rows, err := conn.Query(rctx, sql, builder.Replacements...)
			if err != nil {
				return fmt.Errorf("cannot fetch results: %w", err)
			}
			defer rows.Close()
			for rows.Next() {
				if err := rows.Scan(ptrs...); err != nil {
					return fmt.Errorf("scan failed: %w", err)
				}
				result := make(map[string]any, len(fields))
				for i, p := range fields {
					result[p] = values[i]
				}
				results = append(results, result)
			}
			return rows.Err()
		},
		c.cfg.Database.MaxQueryDuration, // In case the user provided a very expensive query.
	); err != nil {
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h := gin.H{}
	if calcCount {
		h["count"] = count
	}
	if len(results) > 0 {
		h["documents"] = results
	}
	ctx.JSON(http.StatusOK, h)
}

var (
	dirtyReplace     *regexp.Regexp
	dirtyReplaceOnce sync.Once
)

// qndSQLReplace is a quick and dirty hack to re-substitute strings
// into SQL statements. Warning: USE FOR LOGGING ONLY!
// The separation SQL <-> replacements were done beforehand to
// prevent injections!
func qndSQLReplace(sql string, replacements []any) string {
	dirtyReplaceOnce.Do(func() {
		dirtyReplace = regexp.MustCompile(`\$([\d]+)`)
	})
	sql = dirtyReplace.ReplaceAllStringFunc(sql, func(s string) string {
		m := dirtyReplace.FindStringSubmatch(s)
		return `'%[` + m[1] + `]s'`
	})
	return fmt.Sprintf(sql, replacements...)
}
