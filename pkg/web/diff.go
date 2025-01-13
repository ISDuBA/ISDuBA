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
	"github.com/ISDuBA/ISDuBA/pkg/models"
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

// viewDiff is an endpoint that returns diff between two documents.
//
//	@Summary		Returns a diff.
//	@Description	Returns a diff between two documents.
//	@Param			document1	path	string	true	"Document 1 ID"
//	@Param			document2	path	string	true	"Document 2 ID"
//	@Produce		json
//	@Success		200	{object}	any
//	@Failure		400	{object}	models.Error
//	@Failure		404	{object}	models.Error
//	@Failure		500	{object}	models.Error
//	@Router			/diff/{document1}/{document2} [get]
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
			models.SendError(ctx, http.StatusBadRequest, err)
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
		tlps := c.tlps(ctx)
		if len(tlps) == 0 {
			models.SendErrorMessage(ctx, http.StatusNotFound, "document not found")
			return
		}
		tlpExpr := tlps.AsExpr()
		if err := c.db.Run(
			ctx.Request.Context(),
			func(rctx context.Context, conn *pgxpool.Conn) error {
				for _, f := range fromDB {
					expr := query.FieldEqInt("id", f.id).And(tlpExpr)
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
				models.SendErrorMessage(ctx, http.StatusNotFound, "document not found")
			} else {
				slog.Error("database error", "err", err)
				models.SendError(ctx, http.StatusInternalServerError, err)
			}
			return
		}
	}

	// Do we need to load documents from the temp store?
	for _, f := range fromTempStore {
		user := ctx.GetString("uid")
		r, entry, err := c.ts.Fetch(user, f.id)
		switch {
		case errors.Is(err, tempstore.ErrFileNotFound):
			models.SendError(ctx, http.StatusNotFound, err)
			return
		case err != nil:
			slog.Error("temp store fetch error", "err", err)
			models.SendError(ctx, http.StatusInternalServerError, err)
			return
		}
		data := make([]byte, int(entry.Length))
		if _, err := io.ReadFull(r, data); err != nil {
			slog.Error("temp store read error", "err", err)
			models.SendError(ctx, http.StatusInternalServerError, err)
			return
		}
		*f.doc = data
	}

	// Create the patch.
	patch, err := jsonpatch.CreatePatch(doc[0], doc[1])
	if err != nil {
		slog.Error("creating patch failed", "err", err)
		models.SendError(ctx, http.StatusInternalServerError, err)
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
				models.SendError(ctx, http.StatusInternalServerError, err)
				return
			}
			x, ok := locate(d1, p.Path)
			if !ok {
				models.SendErrorMessage(ctx, http.StatusNotFound, fmt.Sprintf("path %q not found", p.Path))
			} else {
				ctx.JSON(http.StatusOK, x)
			}
			return
		}
		models.SendErrorMessage(ctx, http.StatusNotFound, "path/op not found")
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
		models.SendError(ctx, http.StatusInternalServerError, err)
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
