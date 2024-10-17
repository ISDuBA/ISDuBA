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
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	selectImportStatsSQL = `SELECT ` +
		`date_bin($1, time, $2) AS bucket,` +
		`count(*) AS count ` +
		`FROM downloads ` +
		`%s ` + // placeholder for deeper joins.
		`WHERE time BETWEEN $2 AND $3 ` +
		`%s ` + // placeholder for more filters.
		`GROUP BY bucket ` +
		`ORDER BY bucket`
	selectCVEStatsSQL = `SELECT ` +
		`date_bin($1, time, $2) AS bucket,` +
		`count(distinct cve_id) AS count ` +
		`FROM downloads JOIN documents ON downloads.documents_id = documents.id ` +
		`JOIN documents_cves ON documents.id = documents_cves.documents_id ` +
		`%s ` + // placeholder for deeper joins.
		`WHERE time BETWEEN $2 AND $3 ` +
		`%s ` + // placeholder for more filters.
		`GROUP BY bucket ` +
		`ORDER BY bucket`
)

func (c *Controller) cveStatsSource(ctx *gin.Context) {
	c.importStatsSourceTmpl(ctx, selectCVEStatsSQL)
}

func (c *Controller) cveStatsFeed(ctx *gin.Context) {
	c.importStatsFeedTmpl(ctx, selectCVEStatsSQL)
}

func (c *Controller) cveStatsAllSources(ctx *gin.Context) {
	c.importStatsAllSourcesTmpl(ctx, selectCVEStatsSQL)
}

func (c *Controller) importStatsSource(ctx *gin.Context) {
	c.importStatsSourceTmpl(ctx, selectImportStatsSQL)
}

func (c *Controller) importStatsAllSources(ctx *gin.Context) {
	c.importStatsAllSourcesTmpl(ctx, selectImportStatsSQL)
}

func (c *Controller) importStatsFeed(ctx *gin.Context) {
	c.importStatsFeedTmpl(ctx, selectImportStatsSQL)
}

func (c *Controller) importStatsSourceTmpl(ctx *gin.Context, sqlTmpl string) {
	sourcesID, ok := parse(ctx, toInt64, ctx.Param("id"))
	if !ok {
		return
	}
	from, to, step, ok := importStatsInterval(ctx)
	if !ok {
		return
	}
	var cond strings.Builder
	cond.WriteString(`AND feeds.sources_id = $4`)
	if !filterImportStats(ctx, &cond) {
		return
	}
	c.serveImportStats(ctx,
		func(rctx context.Context, conn *pgxpool.Conn) (pgx.Rows, error) {
			const joinFeeds = `JOIN feeds ON downloads.feeds_id = feeds.id`
			sql := fmt.Sprintf(sqlTmpl, joinFeeds, cond.String())
			return conn.Query(rctx, sql, step, from, to, sourcesID)
		})
}

func (c *Controller) importStatsAllSourcesTmpl(ctx *gin.Context, sqlTmpl string) {
	from, to, step, ok := importStatsInterval(ctx)
	if !ok {
		return
	}
	var cond strings.Builder
	if !filterImportStats(ctx, &cond) {
		return
	}
	c.serveImportStats(ctx,
		func(rctx context.Context, conn *pgxpool.Conn) (pgx.Rows, error) {
			sql := fmt.Sprintf(sqlTmpl, "", cond.String())
			return conn.Query(rctx, sql, step, from, to)
		})
}

func (c *Controller) importStatsFeedTmpl(ctx *gin.Context, sqlTmpl string) {
	feedID, ok := parse(ctx, toInt64, ctx.Param("id"))
	if !ok {
		return
	}
	from, to, step, ok := importStatsInterval(ctx)
	if !ok {
		return
	}
	var cond strings.Builder
	cond.WriteString(`AND feeds_id = $4`)
	if !filterImportStats(ctx, &cond) {
		return
	}
	c.serveImportStats(ctx,
		func(rctx context.Context, conn *pgxpool.Conn) (pgx.Rows, error) {
			sql := fmt.Sprintf(sqlTmpl, "", cond.String())
			return conn.Query(rctx, sql, step, from, to, feedID)
		})
}

func (c *Controller) serveImportStats(
	ctx *gin.Context,
	query func(context.Context, *pgxpool.Conn) (pgx.Rows, error),
) {
	var list [][]any
	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			rows, _ := query(rctx, conn)
			var err error
			list, err = collectBuckets(rows)
			return err
		}, 0,
	); err != nil {
		slog.Error("Cannot fetch import stats", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, list)
}

func collectBuckets(rows pgx.Rows) ([][]any, error) {
	return pgx.CollectRows(rows,
		func(row pgx.CollectableRow) ([]any, error) {
			var bucket time.Time
			var count int64
			if err := row.Scan(&bucket, &count); err != nil {
				return nil, err
			}
			return []any{bucket.UTC(), count}, nil
		})
}

func importStatsInterval(ctx *gin.Context) (time.Time, time.Time, time.Duration, bool) {
	var (
		ok       bool
		from, to time.Time
		step     time.Duration
		now      = sync.OnceValue(time.Now)
	)
	if value := ctx.Query("from"); value != "" {
		if from, ok = parse(ctx, parseTime, value); !ok {
			return time.Time{}, time.Time{}, 0, false
		}
	} else {
		from = now().Add(-time.Hour * 3 * 24)
	}

	if value := ctx.Query("to"); value != "" {
		if to, ok = parse(ctx, parseTime, value); !ok {
			return time.Time{}, time.Time{}, 0, false
		}
	} else {
		to = now()
	}

	if to.Before(from) {
		to, from = from, to
	}

	if value := ctx.Query("step"); value != "" {
		if step, ok = parse(ctx, time.ParseDuration, value); !ok {
			return time.Time{}, time.Time{}, 0, false
		}
		step = step.Abs()
	} else {
		step = to.Sub(from)
	}
	return from.UTC(), to.UTC(), step, true
}

func filterImportStats(ctx *gin.Context, cond *strings.Builder) bool {
	have := false
	for _, flag := range []string{
		"download_failed",
		"filename_failed",
		"schema_failed",
		"remote_failed",
		"checksum_failed",
		"signature_failed",
		"duplicate_failed",
	} {
		if value := ctx.Query(flag); value != "" {
			v, ok := parse(ctx, strconv.ParseBool, value)
			if !ok {
				return false
			}
			if have {
				cond.WriteString(" OR ")
			} else {
				have = true
				cond.WriteString(" AND (")
			}
			if !v {
				cond.WriteString("NOT ")
			}
			cond.WriteString(flag)
		}
	}
	if have {
		cond.WriteByte(')')
	}
	return true
}
