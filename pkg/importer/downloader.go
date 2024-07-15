// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package importer

import (
	"bytes"
	"context"
	"errors"
	"github.com/ISDuBA/ISDuBA/pkg/database"
	"github.com/ISDuBA/ISDuBA/pkg/models"
	"github.com/csaf-poc/csaf_distribution/v3/csaf/filter"
	"log/slog"
	"os"

	csafmodel "github.com/csaf-poc/csaf_distribution/v3/csaf/models"
	"github.com/csaf-poc/csaf_distribution/v3/lib/downloader"
	"github.com/jackc/pgx/v5/pgxpool"
)

// downloadWorker downloads and imports csaf documents.
type downloadWorker struct{}

// DownloadJob describes the download configuration.
type DownloadJob struct {
	Config         models.JobConfig
	ForwardQueue   int
	ValidationMode downloader.ValidationMode
	Presets        []string
	Db             *database.DB
	LogFile        string
	LogLevel       slog.Level
}

func newDownloadWorker() *downloadWorker {
	return &downloadWorker{}
}

func downloadHandler(ctx context.Context, db *database.DB) func(d downloader.DownloadedDocument) error {
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

func failedForwardHandler() func(filename, doc, sha256, sha512 string) error {
	return func(filename, doc, sha256, sha512 string) error {
		return nil
	}
}

func (w *downloadWorker) run(ctx context.Context, job DownloadJob) error {
	logFile, err := os.OpenFile(job.LogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		slog.Error("Couldn't open log file", "file", job.LogFile, "err", err)
	}
	defer logFile.Close()

	slogOptions := slog.HandlerOptions{
		Level: job.LogLevel,
	}
	logger := slog.New(slog.NewJSONHandler(logFile, &slogOptions))

	var ignorePatterns filter.PatternMatcher
	if job.Config.IgnorePattern != nil {
		ignorePatterns, err = filter.NewPatternMatcher([]string{*job.Config.IgnorePattern})
		if err != nil {
			return err
		}
	}

	cfg := &downloader.Config{
		Insecure:               job.Config.Insecure,
		IgnoreSignatureCheck:   job.Config.IgnoreSignatureCheck,
		ClientKey:              job.Config.ClientKey,
		ClientPassphrase:       job.Config.ClientPassphrase,
		Rate:                   job.Config.Rate,
		Worker:                 job.Config.Worker,
		IgnorePattern:          ignorePatterns,
		RemoteValidatorPresets: job.Presets,
		ForwardQueue:           job.ForwardQueue,
		FailedForwardHandler:   failedForwardHandler(),
		DownloadHandler:        downloadHandler(ctx, job.Db),
		ValidationMode:         job.ValidationMode,
		Logger:                 logger,
	}

	if job.Config.StartRange != nil && job.Config.EndRange != nil {
		cfg.Range = &csafmodel.TimeRange{*job.Config.StartRange, *job.Config.EndRange}
	}
	d, err := downloader.NewDownloader(cfg)
	if err != nil {
		slog.Warn("Creating new downloader failed", "err", err)
		return err
	}
	defer d.Close()

	err = d.Run(ctx, job.Config.Domains)
	if err != nil {
		slog.Warn("Download failed", "err", err)
		return err
	}
	return nil
}
