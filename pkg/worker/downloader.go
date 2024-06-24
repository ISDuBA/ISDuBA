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
	"os"

	"github.com/ISDuBA/ISDuBA/pkg/database"
	"github.com/ISDuBA/ISDuBA/pkg/models"

	"github.com/csaf-poc/csaf_distribution/v3/lib/downloader"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DownloadWorker struct {
	ctx  context.Context
	jobs chan DownloadJob
}

type DownloadJob struct {
	Presets        []string
	ValidationMode downloader.ValidationMode
	Db             *database.DB
	LogFile        string
	LogLevel       slog.Level
	Domains        []string
	Worker         int
	ForwardQueue   int
}

func NewDownloadWorker(ctx context.Context) *DownloadWorker {
	return &DownloadWorker{
		ctx:  ctx,
		jobs: make(chan DownloadJob),
	}
}

func (w *DownloadWorker) Enqueue(job DownloadJob) {
	w.jobs <- job
}

func downloadHandler(db *database.DB, ctx context.Context) func(d downloader.DownloadedDocument) error {
	return func(d downloader.DownloadedDocument) error {
		if d.ValStatus != downloader.ValidValidationStatus {
			slog.Info("Got invalid document")
			return nil
		}

		r := bytes.NewReader(d.Data)
		actor := models.Importer

		var id int64
		if err := db.Run(ctx, func(ctx context.Context, conn *pgxpool.Conn) error {
			var err error
			id, err = models.ImportDocument(ctx, conn, r, &actor, nil, false)
			return err
		}, 0); err != nil {
			if errors.Is(err, models.ErrAlreadyInDatabase) {
				slog.Warn("advisory already in database")
				err = nil
			}
			return err
		}

		slog.Info("Imported advisory", "id", id)
		return nil
	}
}

func FailedForwardHandler() func(filename, doc, sha256, sha512 string) error {
	return func(filename, doc, sha256, sha512 string) error {
		return nil
	}
}

func (w *DownloadWorker) Run() {
	for {
		select {
		case <-w.ctx.Done():
			return
		case job := <-w.jobs:
			func() {

				logFile, err := os.OpenFile(job.LogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
				if err != nil {
					slog.Error("Couldn't open log file", "file", job.LogFile, "err", err)
				}
				defer logFile.Close()

				slogOptions := slog.HandlerOptions{
					Level: job.LogLevel,
				}
				logger := slog.New(slog.NewJSONHandler(logFile, &slogOptions))

				cfg := &downloader.Config{
					Worker:                 job.Worker,
					RemoteValidatorPresets: job.Presets,
					ForwardQueue:           job.ForwardQueue,
					FailedForwardHandler:   FailedForwardHandler(),
					DownloadHandler:        downloadHandler(job.Db, w.ctx),
					ValidationMode:         job.ValidationMode,
					Logger:                 logger,
				}
				d, err := downloader.NewDownloader(cfg)
				if err != nil {
					slog.Warn("Creating new downloader failed", "err", err)
					return
				}
				defer d.Close()

				err = d.Run(w.ctx, job.Domains)
				if err != nil {
					slog.Warn("Download failed", "err", err)
				}
			}()
		}
	}
}

func (w *DownloadWorker) Close() {
	close(w.jobs)
}
