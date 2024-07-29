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

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type source struct {
	ID     int64    `json:"id" form:"id"`
	Name   string   `json:"name" form:"name" binding:"required,min=1"`
	Domain *string  `json:"domain,omitempty" form:"domain" binding:"omitnil,min=1"`
	PMD    *string  `json:"pmd,omitempty" form:"pmd" binding:"omitnil,url"`
	Active *bool    `json:"active" form:"active"`
	Rate   *float64 `json:"rate,omitempty" form:"rate" binding:"gt=0"`
	Slots  *int     `json:"slots,omitempty" form:"slots" binding:"gte=1"`
}

type feed struct {
	ID    int64  `json:"id"`
	Label string `json:"label"`
	URL   string `json:"url"`
	Rolie bool   `json:"rolie"`
}

func (c *Controller) viewSources(ctx *gin.Context) {

	var srcs []*source
	const sql = `SELECT id, name, domain, pmd, active, rate, slots FROM sources`

	if err := c.db.Run(ctx.Request.Context(), func(rctx context.Context, con *pgxpool.Conn) error {
		rows, err := con.Query(rctx, sql)
		if err != nil {
			return fmt.Errorf("failed fetching sources: %w", err)
		}
		srcs, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (*source, error) {
			var src source
			return &src, row.Scan(
				&src.ID,
				&src.Name,
				&src.Domain,
				&src.PMD,
				&src.Active,
				&src.Rate,
				&src.Slots,
			)
		})
		return err
	}, 0); err != nil {
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"sources": srcs})
}

func (c *Controller) createSource(ctx *gin.Context) {

	var src source
	if err := ctx.ShouldBind(&src); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if src.Domain == nil && src.PMD == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing 'domain' or 'pmd'"})
		return
	}

	if src.Domain != nil && src.PMD != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "'domain' and 'pmd' are exclusive"})
		return
	}

	if src.Rate != nil && (c.cfg.Sources.RatePerSlot != 0 && *src.Rate > c.cfg.Sources.RatePerSlot) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "'rate' out of range"})
		return
	}

	if src.Slots != nil && *src.Slots > c.cfg.Sources.SlotsPerSource {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "'slots' out of range"})
		return
	}

	const sql = `INSERT INTO sources (name, domain, pmd, active, rate, slots) ` +
		`VALUES ($1, $2, $3, $4, $5, $6) ` +
		`RETURNING id`

	var id int64

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, con *pgxpool.Conn) error {
			return con.QueryRow(
				rctx, sql,
				src.Name,
				src.Domain,
				src.PMD,
				src.Active != nil && *src.Active,
				src.Rate,
				src.Slots).Scan(&id)
		}, 0,
	); err != nil {
		// As name can cause an unique constraint violation
		// report this as a bad request as this expected.
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "not a unique value: " + pgErr.Message,
			})
		} else {
			slog.Error("database error", "err", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": id})
}

func (c *Controller) deleteSource(ctx *gin.Context) {
	var input struct {
		ID int64 `uri:"id" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.sm.RemoveSource(input.ID); err != nil {
		slog.Error("removing source failed", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	const sql = `DELETE FROM sources WHERE id = $1`

	notFound := false

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, con *pgxpool.Conn) error {
			tags, err := con.Exec(rctx, sql, input.ID)
			if err != nil {
				return fmt.Errorf("removing source failed: %w", err)
			}
			notFound = tags.RowsAffected() == 0
			return nil
		}, 0,
	); err != nil {
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if notFound {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "source deleted"})
	}
}

func (c *Controller) updateSource(ctx *gin.Context) {
	// TODO: Implement me!
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": "'updateSource' not implemented, yet.",
	})
}

func (c *Controller) viewFeeds(ctx *gin.Context) {
	var input struct {
		SourceID int64 `uri:"id"`
	}
	if err := ctx.ShouldBindUri(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	const sql = `SELECT id, label, url, rolie FROM feeds WHERE sources_id = $1`

	var feeds []*feed
	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, con *pgxpool.Conn) error {
			rows, err := con.Query(rctx, sql, input.SourceID)
			if err != nil {
				return fmt.Errorf("fetching feeds failed: %w", err)
			}
			feeds, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (*feed, error) {
				var f feed
				return &f, row.Scan(&f.ID, &f.Label, &f.URL, &f.Rolie)
			})
			return err
		}, 0,
	); err != nil {
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"feeds": feeds})
}

func (c *Controller) createFeed(ctx *gin.Context) {
	var input struct {
		SourceID int64  `uri:"id"`
		Label    string `form:"label" binding:"required,min=1"`
		URL      string `form:"url" binding:"required,url"`
		Rolie    bool   `form:"rolie"`
	}
	if err := errors.Join(ctx.ShouldBind(&input), ctx.ShouldBindUri(&input)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	const (
		sourceSQL = `SELECT EXISTS(SELECT 1 FROM sources WHERE id = $1)`
		insertSQL = `INSERT INTO feeds (label, sources_id, url, rolie) ` +
			`VALUES ($1, $2, $3, $4) ` +
			`RETURNING id`
	)
	var (
		sourceFound bool
		feedID      int64
	)
	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, con *pgxpool.Conn) error {
			tx, err := con.Begin(rctx)
			if err != nil {
				return fmt.Errorf("starting tx failed: %w", err)
			}
			defer tx.Rollback(rctx)
			if err := tx.QueryRow(rctx, sourceSQL,
				input.SourceID,
			).Scan(&sourceFound); err != nil {
				return fmt.Errorf("checking source id failed: %w", err)
			}
			if !sourceFound {
				return nil
			}
			if err := tx.QueryRow(rctx, insertSQL,
				input.Label,
				input.SourceID,
				input.URL,
				input.Rolie,
			).Scan(&feedID); err != nil {
				return fmt.Errorf("inserting feed failed: %w", err)
			}
			return tx.Commit(rctx)
		}, 0,
	); err != nil {
		// As name can cause an unique constraint violation
		// report this as a bad request as this expected.
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "not a unique value: " + pgErr.Message,
			})
		} else {
			slog.Error("database error", "err", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	if !sourceFound {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "source id not found"})
		return
	}
	// Register feed in source manager.
	if err := c.sm.AddFeed(input.SourceID, feedID); err != nil {
		slog.Error("add feed failed", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"id": feedID})
}

func (c *Controller) viewFeed(ctx *gin.Context) {
	var input struct {
		FeedID int64 `uri:"id"`
	}
	if err := ctx.ShouldBindUri(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	const sql = `SELECT label, url, rolie FROM feeds WHERE id = $1`

	feed := feed{ID: input.FeedID}
	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, con *pgxpool.Conn) error {
			return con.QueryRow(rctx, sql, input.FeedID).Scan(
				&feed.Label, &feed.URL, &feed.Rolie)
		}, 0,
	); err != nil {
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, &feed)
}

func (c *Controller) deleteFeed(ctx *gin.Context) {
	// TODO: Implement me!
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": "'deleteFeed' not implemented, yet.",
	})
}

func (c *Controller) feedLog(ctx *gin.Context) {
	// TODO: Implement me!
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": "'feedLog' not implemented, yet.",
	})
}

// defaultMessage returns the default message.
func (c *Controller) defaultMessage(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": c.cfg.Sources.DefaultMessage})
}
