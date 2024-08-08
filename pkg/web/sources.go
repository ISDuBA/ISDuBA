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
	"net/url"
	"time"

	"github.com/ISDuBA/ISDuBA/pkg/config"
	"github.com/ISDuBA/ISDuBA/pkg/sources"
	"github.com/gin-gonic/gin"
)

type source struct {
	ID     int64    `json:"id" form:"id"`
	Name   string   `json:"name" form:"name" binding:"required,min=1"`
	URL    string   `json:"url" form:"url" binding:"required,min=1"`
	Active *bool    `json:"active,omitempty" form:"active"`
	Rate   *float64 `json:"rate,omitempty" form:"rate" binding:"omitnil,gt=0"`
	Slots  *int     `json:"slots,omitempty" form:"slots" binding:"omitnil,gte=1"`
}

type feed struct {
	ID       int64               `json:"id"`
	Label    string              `json:"label"`
	URL      string              `json:"url"`
	Rolie    bool                `json:"rolie"`
	LogLevel config.FeedLogLevel `json:"log_level"`
}

func (c *Controller) viewSources(ctx *gin.Context) {
	var srcs []*source
	c.sm.AllSources(func(
		id int64,
		name string,
		url string,
		active bool,
		rate *float64,
		slots *int,
	) {
		srcs = append(srcs, &source{
			ID:     id,
			Name:   name,
			URL:    url,
			Active: &active,
			Rate:   rate,
			Slots:  slots,
		})
	})
	ctx.JSON(http.StatusOK, gin.H{"sources": srcs})
}

func (c *Controller) createSource(ctx *gin.Context) {
	var src source
	if err := ctx.ShouldBind(&src); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if src.Rate != nil &&
		(c.cfg.Sources.MaxRatePerSource != 0 && *src.Rate > c.cfg.Sources.MaxRatePerSource) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "'rate' out of range"})
		return
	}
	if src.Slots != nil && *src.Slots > c.cfg.Sources.MaxSlotsPerSource {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "'slots' out of range"})
		return
	}
	switch id, err := c.sm.AddSource(
		src.Name,
		src.URL,
		src.Active,
		src.Rate,
		src.Slots,
	); {
	case errors.Is(err, sources.ErrInvalidArgument):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case err != nil:
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	default:
		ctx.JSON(http.StatusCreated, gin.H{"id": id})
	}
}

func (c *Controller) deleteSource(ctx *gin.Context) {
	var input struct {
		ID int64 `uri:"id" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	switch err := c.sm.RemoveSource(input.ID); {
	case errors.Is(err, sources.ErrNoSuchEntry):
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	case err != nil:
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	default:
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
	feeds := []*feed{}
	switch err := c.sm.Feeds(input.SourceID, func(
		id int64,
		label string,
		url *url.URL,
		rolie bool,
		lvl config.FeedLogLevel,
	) {
		feeds = append(feeds, &feed{
			ID:       id,
			Label:    label,
			URL:      url.String(),
			Rolie:    rolie,
			LogLevel: lvl,
		})
	}); {
	case errors.Is(err, sources.ErrNoSuchEntry):
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case err != nil:
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	default:
		ctx.JSON(http.StatusOK, gin.H{"feeds": feeds})
	}
}

func (c *Controller) createFeed(ctx *gin.Context) {
	var input struct {
		SourceID int64  `uri:"id"`
		Label    string `form:"label" binding:"required,min=1"`
		URL      string `form:"url" binding:"required,url"`
		LogLevel string `form:"log_level" binding:"oneof=debug info warn error ''"`
	}
	if err := errors.Join(ctx.ShouldBind(&input), ctx.ShouldBindUri(&input)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var logLevel config.FeedLogLevel
	if input.LogLevel == "" {
		logLevel = c.cfg.Sources.FeedLogLevel
	} else {
		logLevel, _ = config.ParseFeedLogLevel(input.LogLevel)
	}
	parsed, _ := url.Parse(input.URL)
	switch feedID, err := c.sm.AddFeed(
		input.SourceID,
		input.Label,
		parsed,
		logLevel,
	); {
	case errors.Is(err, sources.ErrNoSuchEntry):
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, sources.ErrInvalidArgument):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case err != nil:
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	default:
		ctx.JSON(http.StatusCreated, gin.H{"id": feedID})
	}
}

func (c *Controller) viewFeed(ctx *gin.Context) {
	var input struct {
		FeedID int64 `uri:"id"`
	}
	if err := ctx.ShouldBindUri(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	f := feed{ID: input.FeedID}
	switch err := c.sm.Feed(input.FeedID, func(
		label string,
		url *url.URL,
		rolie bool,
		lvl config.FeedLogLevel,
	) {
		f.Label = label
		f.URL = url.String()
		f.Rolie = rolie
		f.LogLevel = lvl
	}); {
	case errors.Is(err, sources.ErrNoSuchEntry):
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case err != nil:
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	default:
		ctx.JSON(http.StatusOK, &f)
	}
}

func (c *Controller) deleteFeed(ctx *gin.Context) {
	var input struct {
		FeedID int64 `uri:"id"`
	}
	if err := ctx.ShouldBindUri(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	switch err := c.sm.RemoveFeed(input.FeedID); {
	case errors.Is(err, sources.ErrNoSuchEntry):
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case err != nil:
		slog.Error("removing feed failed", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	default:
		ctx.JSON(http.StatusOK, gin.H{"message": "deleted"})
		return
	}
}

func (c *Controller) feedLog(ctx *gin.Context) {
	var input struct {
		FeedID int64 `uri:"id"`
	}
	if err := ctx.ShouldBindUri(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	type entry struct {
		Time    time.Time           `json:"time"`
		Level   config.FeedLogLevel `json:"level"`
		Message string              `json:"msg"`
	}
	entries := []entry{}
	if err := c.sm.FeedLog(input.FeedID, func(
		t time.Time,
		lvl config.FeedLogLevel,
		msg string,
	) {
		entries = append(entries, entry{
			Time:    t,
			Level:   lvl,
			Message: msg,
		})
	}); err != nil {
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, entries)
}

// defaultMessage returns the default message.
func (c *Controller) defaultMessage(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": c.cfg.Sources.DefaultMessage})
}

func (c *Controller) pmd(ctx *gin.Context) {
	var input struct {
		URL string `form:"url" binding:"required,min=1"`
	}
	if err := ctx.ShouldBindQuery(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	lpmd := c.sm.PMD(input.URL)
	if !lpmd.Valid() {
		h := gin.H{}
		if n := len(lpmd.Messages); n > 0 {
			msgs := make([]string, 0, n)
			for i := range lpmd.Messages {
				msgs = append(msgs, lpmd.Messages[i].Message)
			}
			h["messages"] = msgs
		}
		ctx.JSON(http.StatusBadGateway, h)
		return
	}
	ctx.JSON(http.StatusOK, lpmd.Document)
}
