// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package web

import (
	"bytes"
	"context"
	"errors"
	"github.com/ISDuBA/ISDuBA/pkg/models"
	"github.com/jackc/pgx/v5/pgxpool"
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
)

func downloadHandler(c *Controller, ctx *gin.Context) func(d downloader.DownloadedDocument) error {
	return func(d downloader.DownloadedDocument) error {

		if d.ValStatus != downloader.ValidValidationStatus {
			slog.Info("Got invalid document")
			return nil
		}

		r := bytes.NewReader(d.Data.Bytes())
		actor := models.Importer

		var id int64
		if err := c.db.Run(ctx.Request.Context(), func(ctx context.Context, conn *pgxpool.Conn) error {
			var err error
			id, err = models.ImportDocument(ctx, conn, r, &actor, nil, false)
			return err
		}, 0); err != nil {
			if errors.Is(err, models.ErrAlreadyInDatabase) {
				slog.Warn("advisory already in database", "file", d.InitialReleaseDate)
				err = nil
			}
			return err
		}
		slog.Info("inserted", "id", id)

		slog.Info("Imported advisory", "release-date", d.InitialReleaseDate)
		return nil
	}
}

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

	cfg := &downloader.Config{
		Worker:                 defaultWorker,
		RemoteValidatorPresets: []string{defaultPreset},
		ForwardQueue:           defaultForwardQueue,
		FailedForwardHandler:   failedForward,
		DownloadHandler:        downloadHandler(c, ctx),
		ValidationMode:         defaultValidationMode,
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
