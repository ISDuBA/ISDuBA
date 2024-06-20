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
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sergi/go-diff/diffmatchpatch"
	"gomodules.xyz/jsonpatch/v2"

	"github.com/ISDuBA/ISDuBA/pkg/database"
)

func (c *Controller) viewDiff(ctx *gin.Context) {
	docID1, err := strconv.ParseInt(ctx.Param("document1"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	docID2, err := strconv.ParseInt(ctx.Param("document2"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expr1 := database.FieldEqInt("id", docID1)
	expr2 := database.FieldEqInt("id", docID2)

	// Filter the allowed
	if tlps := c.tlps(ctx); len(tlps) > 0 {
		conditions := tlps.AsConditions()
		parser := database.Parser{}
		tlpExpr, err := parser.Parse(conditions)
		if err != nil {
			slog.Warn("TLP filter failed", "err", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		expr1 = expr1.And(tlpExpr)
		expr2 = expr2.And(tlpExpr)
	}
	var (
		where1, replacements1, _ = expr1.Where(false)
		where2, replacements2, _ = expr2.Where(false)
		doc1, doc2               []byte
	)
	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			const fetchSQL = `SELECT original FROM documents WHERE `
			fetch1SQL := fetchSQL + where1
			fetch2SQL := fetchSQL + where2
			if err := conn.QueryRow(rctx, fetch1SQL, replacements1...).Scan(&doc1); err != nil {
				return err
			}
			return conn.QueryRow(rctx, fetch2SQL, replacements2...).Scan(&doc2)
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

	// Create the patch.
	patch, err := jsonpatch.CreatePatch(doc1, doc2)
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
			if err := json.Unmarshal(doc1, &d1); err != nil {
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
	if err := json.Unmarshal(doc1, &d1); err != nil {
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
