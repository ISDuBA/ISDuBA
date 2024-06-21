// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package worker

import (
	"bytes"
	"context"
	"errors"
	"log/slog"

	"github.com/ISDuBA/ISDuBA/pkg/database"
	"github.com/ISDuBA/ISDuBA/pkg/models"

	"github.com/csaf-poc/csaf_distribution/v3/lib/downloader"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DownloadJob struct {
	Domains        []string
	Worker         int
	ForwardQueue   int
	Preset         string
	ValidationMode downloader.ValidationMode
	Db             *database.DB
}

func downloadHandler(db *database.DB, ctx context.Context) func(d downloader.DownloadedDocument) error {
	return func(d downloader.DownloadedDocument) error {
		if d.ValStatus != downloader.ValidValidationStatus {
			slog.Info("Got invalid document")
			return nil
		}

		r := bytes.NewReader(d.Data.Bytes())
		actor := models.Importer

		var id int64
		if err := db.Run(ctx, func(ctx context.Context, conn *pgxpool.Conn) error {
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

func FailedForwardHandler() func(filename, doc, sha256, sha512 string) error {
	return func(filename, doc, sha256, sha512 string) error {
		return nil
	}
}

func DownloadWorker(ctx context.Context, downloadJobChan <-chan DownloadJob) {
	for {
		select {
		case <-ctx.Done():
			return
		case job := <-downloadJobChan:
			func() {
				cfg := &downloader.Config{
					Worker:                 job.Worker,
					RemoteValidatorPresets: []string{job.Preset},
					ForwardQueue:           job.ForwardQueue,
					FailedForwardHandler:   FailedForwardHandler(),
					DownloadHandler:        downloadHandler(job.Db, ctx),
					ValidationMode:         job.ValidationMode,
					Logger:                 slog.Default(),
				}
				d, err := downloader.NewDownloader(cfg)
				if err != nil {
					slog.Warn("Creating new downloader failed", "err", err)
					return
				}
				defer d.Close()

				err = d.Run(ctx, job.Domains)
				if err != nil {
					slog.Warn("Download failed", "err", err)
				}
			}()
		}
	}
}
