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
	"os"
	"os/signal"
	"strings"

	"github.com/csaf-poc/csaf_distribution/v3/lib/downloader"
	"github.com/gin-gonic/gin"
)

const (
	defaultWorker         = 2
	defaultPreset         = "mandatory"
	defaultForwardQueue   = 5
	defaultValidationMode = downloader.ValidationStrict
	defaultLogFile        = "downloader.log"
	defaultLogLevel       = slog.LevelInfo
)

// download downloads the advisories from the specified source
func (c *Controller) download(ctx *gin.Context) {
	domainsQuery, ok := ctx.GetQuery("domains")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing domains query parameter"})
		return
	}
	domains := strings.Split(domainsQuery, ",")
	slog.Info("Download for domains", "domains", domains)

	failedForward := func(filename, doc, sha256, sha512 string) error {
		return nil
	}

	downloadHandler := func(d downloader.DownloadedDocument) error {
		slog.Info("Got document", "release-date", d.InitialReleaseDate)
		return nil
	}

	cfg := &downloader.Config{
		Worker:                 defaultWorker,
		RemoteValidatorPresets: []string{defaultPreset},
		ForwardQueue:           defaultForwardQueue,
		FailedForwardHandler:   failedForward,
		DownloadHandler:        downloadHandler,
		Logger:                 slog.Default(),
	}

	d, err := downloader.NewDownloader(cfg)
	if err != nil {
		slog.Warn("Creating new downloader failed", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	defer d.Close()

	downloadCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	downloadCtx, stop := signal.NotifyContext(downloadCtx, os.Interrupt)
	defer stop()

	if cfg.ForwardURL != "" {
		f := downloader.NewForwarder(cfg)
		go f.Run()
		defer func() {
			f.Log()
			f.Close()
		}()
		d.Forwarder = f
	}

	err = d.Run(downloadCtx, domains)
	if err != nil {
		slog.Warn("Download failed", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
}
