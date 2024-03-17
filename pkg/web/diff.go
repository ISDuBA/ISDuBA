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

	"github.com/ISDuBA/ISDuBA/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"gomodules.xyz/jsonpatch/v2"
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
			ctx.JSON(http.StatusInternalServerError, gin.H{"error:": err.Error()})
			return
		}
		expr1 = expr1.And(tlpExpr)
		expr2 = expr2.And(tlpExpr)
	}
	var (
		where1, replacements1, _ = expr1.Where()
		where2, replacements2, _ = expr2.Where()
		doc1, doc2               []byte
		rctx                     = ctx.Request.Context()
	)
	if err := c.db.Run(rctx, func(conn *pgxpool.Conn) error {
		const fetchSQL = `SELECT original FROM documents WHERE `
		fetch1SQL := fetchSQL + where1
		fetch2SQL := fetchSQL + where2
		if err := conn.QueryRow(rctx, fetch1SQL, replacements1...).Scan(&doc1); err != nil {
			return err
		}
		return conn.QueryRow(rctx, fetch2SQL, replacements2...).Scan(&doc2)
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error:": "document not found"})
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
	}
	ctx.JSON(http.StatusOK, patch)
}
