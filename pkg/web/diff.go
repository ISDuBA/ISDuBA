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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sergi/go-diff/diffmatchpatch"
	"gomodules.xyz/jsonpatch/v2"

	"github.com/ISDuBA/ISDuBA/pkg/database/query"
	"github.com/ISDuBA/ISDuBA/pkg/tempstore"
)

var tempdocumentIDRe = regexp.MustCompile(`^tempdocument(\d+)$`)

func parseDiffID(s string) (int64, bool, error) {
	if m := tempdocumentIDRe.FindStringSubmatch(s); m != nil {
		id, err := strconv.ParseInt(m[1], 10, 64)
		return id, false, err
	}
	id, err := strconv.ParseInt(s, 10, 64)
	return id, true, err
}

func (c *Controller) viewDiff(ctx *gin.Context) {
	type idDoc struct {
		id  int64
		doc *[]byte
	}
	var (
		doc                   [2][]byte
		fromDB, fromTempStore []idDoc
	)
	for i := range doc {
		id, inDB, err := parseDiffID(ctx.Param("document" + strconv.Itoa(i+1)))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var from *[]idDoc
		if inDB {
			from = &fromDB
		} else {
			from = &fromTempStore
		}
		*from = append(*from, idDoc{id: id, doc: &doc[i]})
	}

	// Do we need to load docs from database?
	if len(fromDB) > 0 {
		if err := c.db.Run(
			ctx.Request.Context(),
			func(rctx context.Context, conn *pgxpool.Conn) error {
				var tlpExpr *query.Expr
				if tlps := c.tlps(ctx); len(tlps) > 0 {
					tlpExpr = tlps.AsExpr()
				}
				for _, f := range fromDB {
					expr := query.FieldEqInt("id", f.id)
					if tlpExpr != nil {
						expr = expr.And(tlpExpr)
					}
					var b query.SQLBuilder
					b.CreateWhere(expr)
					fetchSQL := `SELECT original FROM documents WHERE ` + b.WhereClause
					if err := conn.QueryRow(rctx, fetchSQL, b.Replacements...).Scan(f.doc); err != nil {
						return fmt.Errorf("fetching data from database failed: %w", err)
					}
				}
				return nil
			}, 0,
		); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
			} else {
				slog.Error("database error", "err", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
	}

	// Do we need to load documents from the temp store?
	for _, f := range fromTempStore {
		user := ctx.GetString("uid")
		r, entry, err := c.tmpStore.Fetch(user, f.id)
		switch {
		case errors.Is(err, tempstore.ErrFileNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		case err != nil:
			slog.Error("temp store fetch error", "err", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		data := make([]byte, int(entry.Length))
		if _, err := io.ReadFull(r, data); err != nil {
			slog.Error("temp store read error", "err", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		*f.doc = data
	}

	// Create the patch.
	patch, err := jsonpatch.CreatePatch(doc[0], doc[1])
	if err != nil {
		slog.Error("creating patch failed", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Deliver a specific operation item.
	if op, path := ctx.Query("item_op"), ctx.Query("item_path"); op != "" && path != "" {
		for i := range patch {
			p := &patch[i]
			if p.Operation != op || p.Path != path {
				continue
			}
			var d1 any
			if err := json.Unmarshal(doc[0], &d1); err != nil {
				slog.Error("unmarshaling failed", "err", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			x, ok := locate(d1, p.Path)
			if !ok {
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": fmt.Sprintf("path %q not found", p.Path),
				})
			} else {
				ctx.JSON(http.StatusOK, x)
			}
			return
		}
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "path/op not found",
		})
		return
	}

	// Check if we have to do word diffing. If not deliver the patch directly.
	wordDiffS := ctx.DefaultQuery("word-diff", "false")
	wordDiff, _ := strconv.ParseBool(wordDiffS)
	if !wordDiff {
		ctx.JSON(http.StatusOK, patch)
		return
	}

	// Calculate word diff for "replace" operations.
	var d1 any
	if err := json.Unmarshal(doc[0], &d1); err != nil {
		slog.Error("unmarshaling failed", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dmp := diffmatchpatch.New()

	for i := range patch {
		p := &patch[i]
		if p.Operation != "replace" {
			continue
		}
		r, ok := p.Value.(string)
		if !ok {
			continue
		}
		x, ok := locate(d1, p.Path)
		if !ok {
			continue
		}
		s, ok := x.(string)
		if !ok {
			continue
		}
		diffs := dmp.DiffMain(s, r, false)
		p.Value = doWordDiff(diffs)
	}

	ctx.JSON(http.StatusOK, patch)
}

type wordDiffCommand struct {
	Mode string `json:"m"`
	Text string `json:"t"`
}

func doWordDiff(diffs []diffmatchpatch.Diff) []wordDiffCommand {
	cmds := make([]wordDiffCommand, 0, len(diffs))
	for _, d := range diffs {
		switch d.Type {
		case diffmatchpatch.DiffInsert:
			cmds = append(cmds, wordDiffCommand{"i", d.Text})
		case diffmatchpatch.DiffDelete:
			cmds = append(cmds, wordDiffCommand{"d", d.Text})
		case diffmatchpatch.DiffEqual:
			cmds = append(cmds, wordDiffCommand{"o", d.Text})
		}
	}
	return cmds
}

func locate(doc any, path string) (any, bool) {
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		return nil, false
	}
	if parts[0] == "" {
		parts = parts[1:]
	}
	var recurse func(any, []string) (any, bool)
	recurse = func(current any, rest []string) (any, bool) {
		if len(rest) == 0 {
			return current, true
		}
		switch x := current.(type) {
		case map[string]any:
			v, ok := x[rest[0]]
			if !ok {
				return nil, false
			}
			return recurse(v, rest[1:])
		case []any:
			idx, err := strconv.Atoi(rest[0])
			if err != nil || idx < 0 || idx >= len(x) {
				return nil, false
			}
			return recurse(x[idx], rest[1:])
		default:
			return nil, false
		}
	}
	return recurse(doc, parts)
}
