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
	var (
		from, to time.Time
		step     time.Duration
	)
	if from, to, step, ok = importStatsInterval(ctx, 0); !ok {
		return
	}

	if !from.Equal(to) {
		if step <= 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "step too small"})
			return
		}
		if steps := to.Sub(from) / step; steps >= 65536/2-2 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "too many steps"})
			return
		}
	} else {
		step = time.Hour
	}

	var documentsSQL, advisoriesSQL string

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

	list := [][]any{}

	queries := func(when time.Time) (func(row pgx.Row) error, func(row pgx.Row) error) {
		var numDocs, numAdvs int64
		queryDocuments := func(row pgx.Row) error {
			return row.Scan(&numDocs)
		}
		queryAdvisories := func(row pgx.Row) error {
			if err := row.Scan(&numAdvs); err != nil {
				return err
			}
			list = append(list, []any{when.UTC(), numDocs, numAdvs})
			return nil
		}
		return queryDocuments, queryAdvisories
	}

	batch := &pgx.Batch{}
	for when := from; !when.After(to); when = when.Add(step) {
		docs, advs := queries(when)
		batch.Queue(documentsSQL, when).QueryRow(docs)
		batch.Queue(advisoriesSQL, when).QueryRow(advs)
	}

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			tx, err := conn.Begin(rctx)
			if err != nil {
				return err
			}
			defer tx.Rollback(rctx)
			return tx.SendBatch(rctx, batch).Close()
		}, 0,
	); err != nil {
		slog.Error("counting documents/advisories failed", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, list)
}
