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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"text/template"

	"github.com/csaf-poc/csaf_distribution/v3/csaf"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ISDuBA/ISDuBA/pkg/database/query"
	"github.com/ISDuBA/ISDuBA/pkg/models"
)

// MinSearchLength enforces a minimal length of search phrases.
const MinSearchLength = 2 // Makes at least "Go" searchable ;-)

// deleteDocument is an end point for deleting a document.
func (c *Controller) deleteDocument(ctx *gin.Context) {
	// Get an ID from context
	docID, ok := parse(ctx, toInt64, ctx.Param("id"))
	if !ok {
		return
	}

	// FieldEqInt is a shortcut mainly for building expressions
	// accessing an integer column like 'id's.
	// Expr encapsulates a parsed expression to be converted to an SQL WHERE clause.
	expr := c.andTLPExpr(ctx, query.FieldEqInt("id", docID))

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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	f, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	limited := http.MaxBytesReader(
		ctx.Writer, f, int64(c.cfg.General.AdvisoryUploadLimit))
	defer limited.Close()

	var buf bytes.Buffer
	tee := io.TeeReader(limited, &buf)

	var document any
	if err := json.NewDecoder(tee).Decode(&document); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "document is not JSON: " + err.Error()})
		return
	}

	msgs, err := csaf.ValidateCSAF(document)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "schema validation failed: " + err.Error()})
		return
	}
	if len(msgs) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "schema validation failed: " + strings.Join(msgs, ", "),
		})
		return
	}

	// Is remote validator configured?
	if c.val != nil {
		rvr, err := c.val.Validate(document)
		if err != nil {
			slog.Error("remote validation failed", "err", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "remote validation failed: " + err.Error(),
			})
			return
		}
		if !rvr.Valid {
			// XXX: Maybe we should tell, what's exactly wrong?
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "remote validation"})
			return
		}
	}

	var id int64
	switch err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			id, err = models.ImportDocumentData(
				rctx, conn, document, buf.Bytes(), actor, c.tlps(ctx), nil, false)
			return err
		}, 0,
	); {
	case err == nil:
		ctx.JSON(http.StatusCreated, gin.H{"id": id})
	case errors.Is(err, models.ErrAlreadyInDatabase):
		ctx.JSON(http.StatusConflict, gin.H{"error": "already in database"})
	case errors.Is(err, models.ErrNotAllowed):
		ctx.JSON(http.StatusForbidden, gin.H{"error": "wrong publisher/tlp"})
	default:
		slog.Error("storing document failed", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// viewDocument is an end point to export a document.
func (c *Controller) viewDocument(ctx *gin.Context) {
	id, ok := parse(ctx, toInt64, ctx.Param("id"))
	if !ok {
		return
	}

	expr := c.andTLPExpr(ctx, query.FieldEqInt("id", id))

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

func (c *Controller) viewForwardTargets(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, c.fm.GetTargets())
}

// forwardDocument is an end point to forward a document.
func (c *Controller) forwardDocument(ctx *gin.Context) {
	id, ok := parse(ctx, toInt64, ctx.Param("id"))
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "no valid document id specified"})
		return
	}

	targetID, ok := parse(ctx, toInt64, ctx.Param("target"))
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "no target specified"})
		return
	}

	expr := c.andTLPExpr(ctx, query.FieldEqInt("id", id))

	fields := []string{"id"}
	builder := query.SQLBuilder{}
	builder.CreateWhere(expr)
	sql := builder.CreateQuery(fields, "", -1, -1)

	var documentID int64

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			return conn.QueryRow(rctx, sql, builder.Replacements...).Scan(&documentID)
		}, 0,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}
		return
	}

	if err := c.fm.ForwardDocument(ctx.Request.Context(), int(targetID), documentID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"id": documentID})
}

// overviewDocuments is an end point to return an overview document.
func (c *Controller) overviewDocuments(ctx *gin.Context) {
	// Use the advisories.
	advisory, ok := parse(ctx, strconv.ParseBool, ctx.DefaultQuery("advisories", "false"))
	if !ok {
		return
	}

	mode := query.DocumentMode
	if advisory {
		mode = query.AdvisoryMode
	}

	parser := query.Parser{
		Mode:            mode,
		MinSearchLength: MinSearchLength,
		Me:              ctx.GetString("uid"),
	}

	// The query to filter the documents.
	expr, ok := parse(ctx, parser.Parse, ctx.DefaultQuery("query", "true"))
	if !ok {
		return
	}

	// Filter the allowed
	expr = c.andTLPExpr(ctx, expr)

	// In advisory mode we only show the latest.
	if advisory {
		expr = expr.And(query.BoolField("latest"))
	}

	builder := query.SQLBuilder{Mode: mode}
	builder.CreateWhere(expr)

	fields := strings.Fields(
		ctx.DefaultQuery("columns", "id title tracking_id version publisher"))

	if err := builder.CheckProjections(fields); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderFields := strings.Fields(
		ctx.DefaultQuery("orders", "publisher tracking_id -current_release_date -rev_history_length"))
	order, err := builder.CreateOrder(orderFields)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var (
		calcCount     bool
		count         int64
		limit, offset int64 = -1, -1
	)

	calcCount = ctx.Query("count") != ""

	if lim := ctx.Query("limit"); lim != "" {
		if limit, ok = parse(ctx, toInt64, lim); !ok {
			return
		}
	}

	if ofs := ctx.Query("offset"); ofs != "" {
		if offset, ok = parse(ctx, toInt64, ofs); !ok {
			return
		}
	}

	var results []map[string]any

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			if calcCount {
				countSQL := builder.CreateCountSQL()
				if slog.Default().Enabled(rctx, slog.LevelDebug) {
					slog.Debug("count", "SQL", qndSQLReplace(countSQL, builder.Replacements))
				}
				if err := conn.QueryRow(
					rctx,
					countSQL,
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

			if slog.Default().Enabled(rctx, slog.LevelDebug) {
				slog.Debug("documents", "SQL", qndSQLReplace(sql, builder.Replacements))
			}
			rows, err := conn.Query(rctx, sql, builder.Replacements...)
			if err != nil {
				return fmt.Errorf("cannot fetch results: %w", err)
			}
			defer rows.Close()
			if results, err = scanRows(rows, fields, builder.Aliases); err != nil {
				return fmt.Errorf("loading data failed: %w", err)
			}
			return nil
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

// scanRows turns a result set into a slice of maps.
func scanRows(
	rows pgx.Rows,
	fields []string,
	aliases map[string]string,
) ([]map[string]any, error) {
	values := make([]any, len(fields))
	ptrs := make([]any, len(fields))
	for i := range ptrs {
		ptrs[i] = &values[i]
	}
	var results []map[string]any
	for rows.Next() {
		if err := rows.Scan(ptrs...); err != nil {
			return nil, fmt.Errorf("scanning row failed: %w", err)
		}
		result := make(map[string]any, len(fields))
		for i, p := range fields {
			v := values[i]
			// XXX: A little bit hacky to support client.
			if _, ok := aliases[p]; ok {
				if s, ok := v.(string); ok {
					v = template.HTMLEscapeString(s)
				}
			}
			result[p] = v
		}
		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scanning failed: %w", err)
	}
	return results, nil
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
