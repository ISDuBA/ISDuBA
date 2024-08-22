// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package web

import (
	"encoding/pem"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ISDuBA/ISDuBA/pkg/config"
	"github.com/ISDuBA/ISDuBA/pkg/sources"
	"github.com/gin-gonic/gin"
)

type sourceAge struct {
	time.Duration
}

// UnmarshalParam implements [binding.BindUnmarshaler].
func (sa *sourceAge) UnmarshalParam(param string) error {
	duration, err := time.ParseDuration(param)
	if err != nil {
		return err
	}
	*sa = sourceAge{duration}
	return nil
}

// MarshalText implements [encoding.TextMarshaler].
func (sa sourceAge) MarshalText() ([]byte, error) {
	s := sa.String()
	return []byte(s), nil
}

type source struct {
	ID                   int64      `json:"id" form:"id"`
	Name                 string     `json:"name" form:"name" binding:"required,min=1"`
	URL                  string     `json:"url" form:"url" binding:"required,min=1"`
	Active               bool       `json:"active" form:"active"`
	Rate                 *float64   `json:"rate,omitempty" form:"rate" binding:"omitnil,gt=0"`
	Slots                *int       `json:"slots,omitempty" form:"slots" binding:"omitnil,gte=1"`
	Headers              []string   `json:"headers,omitempty" form:"headers"`
	StrictMode           *bool      `json:"strict_mode,omitempty" form:"strict_mode"`
	Insecure             *bool      `json:"insecure,omitempty" form:"insecure"`
	SignatureCheck       *bool      `json:"signature_check,omitempty" form:"signature_check"`
	Age                  *sourceAge `json:"age,omitempty" form:"age"`
	IgnorePatterns       []string   `json:"ignore_patterns,omitempty" form:"ignore_patterns"`
	ClientCertPublic     *string    `json:"client_cert_public,omitempty" form:"client_cert_public"`
	ClientCertPrivate    *string    `json:"client_cert_private,omitempty" form:"client_cert_private"`
	ClientCertPassphrase *string    `json:"client_cert_passphrase,omitempty" form:"client_cert_passphrase"`
}

type feed struct {
	ID       int64               `json:"id"`
	Label    string              `json:"label"`
	URL      string              `json:"url"`
	Rolie    bool                `json:"rolie"`
	LogLevel config.FeedLogLevel `json:"log_level"`
}

var stars = "***"

func threeStars(b bool) *string {
	if b {
		return &stars
	}
	return nil
}

func (c *Controller) viewSources(ctx *gin.Context) {
	srcs := []*source{}
	c.sm.AllSources(func(
		id int64,
		name string,
		url string,
		active bool,
		rate *float64,
		slots *int,
		headers []string,
		strictMode *bool,
		insecure *bool,
		signatureCheck *bool,
		age *time.Duration,
		ignorePatterns []*regexp.Regexp,
		hasClientCertPublic bool,
		hasClientCertPrivate bool,
		hasClientCertPassphrase bool,
	) {
		var sa *sourceAge
		if age != nil {
			sa = &sourceAge{*age}
		}
		srcs = append(srcs, &source{
			ID:                   id,
			Name:                 name,
			URL:                  url,
			Active:               active,
			Rate:                 rate,
			Slots:                slots,
			Headers:              headers,
			StrictMode:           strictMode,
			Insecure:             insecure,
			SignatureCheck:       signatureCheck,
			Age:                  sa,
			IgnorePatterns:       sources.AsStrings(ignorePatterns),
			ClientCertPublic:     threeStars(hasClientCertPublic),
			ClientCertPrivate:    threeStars(hasClientCertPrivate),
			ClientCertPassphrase: threeStars(hasClientCertPassphrase),
		})
	})
	ctx.JSON(http.StatusOK, gin.H{"sources": srcs})
}

// hasBlock checks if input has a PEM block.
func hasBlock(data []byte) bool {
	block, _ := pem.Decode(data)
	return block != nil
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
	if err := validateHeaders(src.Headers); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ignorePatterns, err := sources.AsRegexps(src.IgnorePatterns)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var clientCertPublic, clientCertPrivate, clientCertPassphrase []byte
	if src.ClientCertPublic != nil {
		clientCertPublic = []byte(*src.ClientCertPublic)
		if !hasBlock(clientCertPublic) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "client_cert_public has no PEM block"})
			return
		}
	}
	if src.ClientCertPrivate != nil {
		clientCertPrivate = []byte(*src.ClientCertPrivate)
		if !hasBlock(clientCertPrivate) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "client_cert_private has no PEM block"})
			return
		}
	}
	if src.ClientCertPassphrase != nil {
		clientCertPassphrase = []byte(*src.ClientCertPassphrase)
	}

	var age *time.Duration
	if src.Age != nil {
		age = &src.Age.Duration
	}

	switch id, err := c.sm.AddSource(
		src.Name,
		src.URL,
		src.Rate,
		src.Slots,
		src.Headers,
		src.StrictMode,
		src.Insecure,
		src.SignatureCheck,
		age,
		ignorePatterns,
		clientCertPublic,
		clientCertPrivate,
		clientCertPassphrase,
	); {
	case err == nil:
		ctx.JSON(http.StatusCreated, gin.H{"id": id})
	case errors.Is(err, sources.InvalidArgumentError("")):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	case err == nil:
		ctx.JSON(http.StatusOK, gin.H{"message": "source deleted"})
	case errors.Is(err, sources.NoSuchEntryError("")):
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func (c *Controller) updateSource(ctx *gin.Context) {
	var input struct {
		SourceID int64 `uri:"id"`
	}
	if err := ctx.ShouldBindUri(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	switch ur, err := c.sm.UpdateSource(input.SourceID, func(su *sources.SourceUpdater) error {
		// name
		if name, ok := ctx.GetPostForm("name"); ok {
			if err := su.UpdateName(name); err != nil {
				return err
			}
		}
		// rate
		if rate, ok := ctx.GetPostForm("rate"); ok {
			var r *float64
			if rate != "" {
				x, err := strconv.ParseFloat(rate, 64)
				if err != nil {
					return sources.InvalidArgumentError(
						fmt.Sprintf("parsing 'rate' failed: %v", err.Error()))
				}
				r = &x
			}
			if err := su.UpdateRate(r); err != nil {
				return err
			}
		}
		// slots
		if slots, ok := ctx.GetPostForm("slots"); ok {
			var sl *int
			if slots != "" {
				x, err := strconv.Atoi(slots)
				if err != nil {
					return sources.InvalidArgumentError(
						fmt.Sprintf("parsing 'slots' failed: %v", err.Error()))
				}
				sl = &x
			}
			if err := su.UpdateSlots(sl); err != nil {
				return err
			}
		}
		// active
		if active, ok := ctx.GetPostForm("active"); ok {
			act, err := strconv.ParseBool(active)
			if err != nil {
				return sources.InvalidArgumentError(
					fmt.Sprintf("parsing 'active' failed: %v", err.Error()))
			}
			if err := su.UpdateActive(act); err != nil {
				return err
			}
		}
		// headers
		if headers, ok := ctx.GetPostFormArray("headers"); ok {
			if err := validateHeaders(headers); err != nil {
				return err
			}
			if err := su.UpdateHeaders(headers); err != nil {
				return err
			}
		}
		// Little helper function for the otional bool fields.
		optBool := func(option string, update func(*bool) error) error {
			value, ok := ctx.GetPostForm(option)
			if !ok {
				return nil
			}
			var b *bool
			if value != "" {
				v, err := strconv.ParseBool(value)
				if err != nil {
					return sources.InvalidArgumentError(
						fmt.Sprintf("parsing %q failed: %v", option, err.Error()))
				}
				b = &v
			}
			return update(b)
		}
		// strictMode
		if err := optBool("strict_mode", su.UpdateStrictMode); err != nil {
			return err
		}
		// insecure
		if err := optBool("insecure", su.UpdateInsecure); err != nil {
			return err
		}
		// signatureCheck
		if err := optBool("signature_check", su.UpdateSignatureCheck); err != nil {
			return err
		}
		// age
		if value, ok := ctx.GetPostForm("age"); ok {
			var age *time.Duration
			if value != "" {
				d, err := time.ParseDuration(value)
				if err != nil {
					return sources.InvalidArgumentError(
						fmt.Sprintf("parsing 'age' failed: %v", err.Error()))
				}
				age = &d
			}
			if err := su.UpdateAge(age); err != nil {
				return err
			}
		}
		// ignorePatterns
		if patterns, ok := ctx.GetPostFormArray("ignore_patterns"); ok {
			regexps, err := sources.AsRegexps(patterns)
			if err != nil {
				return err
			}
			if err := su.UpdateIgnorePatterns(regexps); err != nil {
				return err
			}
		}
		// client certificate update
		optCert := func(option string, update func([]byte) error) error {
			cert, ok := ctx.GetPostForm(option)
			if !ok {
				return nil
			}
			var data []byte
			if cert != "" {
				data = []byte(cert)
				if !hasBlock(data) {
					return sources.InvalidArgumentError(
						fmt.Sprintf("%q has no PEM block", option))
				}
			}
			return update(data)
		}
		if err := optCert("client_cert_public", su.UpdateClientCertPublic); err != nil {
			return err
		}
		if err := optCert("client_cert_private", su.UpdateClientCertPrivate); err != nil {
			return err
		}
		if passphrase, ok := ctx.GetPostForm("client_cert_passphrase"); ok {
			var data []byte
			if passphrase != "" {
				data = []byte(passphrase)
			}
			if err := su.UpdateClientCertPassphrase(data); err != nil {
				return err
			}
		}
		return nil
	}); {
	case err == nil:
		ctx.JSON(http.StatusOK, gin.H{"message": ur.String()})
	case errors.Is(err, sources.NoSuchEntryError("")):
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	case errors.Is(err, sources.InvalidArgumentError("")):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func validateHeaders(headers []string) error {
	for _, header := range headers {
		if k, _, ok := strings.Cut(header, ":"); !ok || strings.TrimSpace(k) == "" {
			return sources.InvalidArgumentError(
				fmt.Sprintf("header %q is invalid", header))
		}
	}
	return nil
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
	case err == nil:
		ctx.JSON(http.StatusOK, gin.H{"feeds": feeds})
	case errors.Is(err, sources.NoSuchEntryError("")):
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	case err == nil:
		ctx.JSON(http.StatusCreated, gin.H{"id": feedID})
	case errors.Is(err, sources.NoSuchEntryError("")):
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, sources.InvalidArgumentError("")):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	case err == nil:
		ctx.JSON(http.StatusOK, &f)
	case errors.Is(err, sources.NoSuchEntryError("")):
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	case err == nil:
		ctx.JSON(http.StatusOK, gin.H{"message": "deleted"})
	case errors.Is(err, sources.NoSuchEntryError("")):
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		slog.Error("removing feed failed", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
