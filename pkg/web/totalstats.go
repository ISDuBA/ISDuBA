// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

// Package web implements the endpoints of the web server.
package web

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (c *Controller) statsTotal(ctx *gin.Context) {
	// Counting imports.
	imports, ok := parse(ctx, strconv.ParseBool, ctx.DefaultQuery("imports", "false"))
	if !ok {
		return
	}
	var when time.Time
	if t := ctx.Query("time"); t != "" {
		var ok bool
		if when, ok = parse(ctx, parseTime, t); !ok {
			return
		}
	} else {
		when = time.Now()
	}
	var (
		documentsSQL, advisoriesSQL string
		numDocuments, numAdvisories int64
	)
	if imports {
		documentsSQL = `SELECT count(*) FROM documents ` +
			`JOIN downloads ON documents.id = downloads.documents_id ` +
			`WHERE time <= $1`
		advisoriesSQL = `SELECT count(DISTINCT (publisher, tracking_id)) FROM documents ` +
			`JOIN downloads ON documents.id = downloads.documents_id ` +
			`WHERE time <= $1`
	} else {
		documentsSQL = `SELECT count(*) FROM documents ` +
			`WHERE least(current_release_date, current_timestamp) <= $1`
		advisoriesSQL = `SELECT count(DISTINCT (publisher, tracking_id)) FROM documents ` +
			`WHERE least(current_release_date, current_timestamp) <= $1`
	}
	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			tx, err := conn.BeginTx(rctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})
			if err != nil {
				return err
			}
			defer tx.Rollback(rctx)
			return errors.Join(
				tx.QueryRow(rctx, documentsSQL, when).Scan(&numDocuments),
				tx.QueryRow(rctx, advisoriesSQL, when).Scan(&numAdvisories))
		}, 0,
	); err != nil {
		slog.Error("counting documents/advisories failed", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"documents":  numDocuments,
		"advisories": numAdvisories,
	})
}
