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
	"database/sql"
	"errors"
	"github.com/ISDuBA/ISDuBA/pkg/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/ISDuBA/ISDuBA/pkg/worker"
	"github.com/gin-gonic/gin"
)

const (
	defaultForwardQueue = 0
)

// addJob creates a new job configuration.
func (c *Controller) addJob(ctx *gin.Context) {
	jobConfig := models.JobConfig{}

	// We need the name.
	if jobConfig.Name = ctx.PostForm("name"); jobConfig.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "missing 'name'",
		})
		return
	}

	// Domains to download from.
	if jobConfig.Domains = ctx.PostFormArray("domains"); len(jobConfig.Domains) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "missing 'domains'",
		})
		return
	}

	const insertSQL = `INSERT INTO job_config (` +
		`name,` +
		`insecure` +
		`ignore_signature_check` +
		`rate` +
		`worker` +
		`start_range` +
		`end_range` +
		`ignore_pattern` +
		`domains` +
		`) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)` +
		`RETURNING id`

	var jobID int64

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			return conn.QueryRow(rctx, insertSQL,
				jobConfig.Name,
				jobConfig.Insecure,
				jobConfig.IgnoreSignatureCheck,
				jobConfig.Rate,
				jobConfig.Worker,
				jobConfig.StartRange,
				jobConfig.EndRange,
				jobConfig.IgnorePattern,
				jobConfig.Domains,
			).Scan(&jobID)
		}, 0,
	); err != nil {
		var pgErr *pgconn.PgError
		// Unique constraint violation
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			ctx.JSON(http.StatusConflict, gin.H{"error": "already in database"})
			return
		}
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"id": jobID,
	})
}

// viewJobs returns all configured jobs.
func (c *Controller) viewJobs(ctx *gin.Context) {
	var jobs []models.JobConfig

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			const fetchSQL = `SELECT name, insecure, ignore_signature_check, rate, worker, start_range, end_range, ignore_pattern, domains FROM job_config`
			rows, _ := conn.Query(rctx, fetchSQL)
			var err error
			jobs, err = pgx.CollectRows(
				rows,
				func(row pgx.CollectableRow) (models.JobConfig, error) {
					var j models.JobConfig
					var ignorePattern sql.NullString
					var startRange sql.NullTime
					var endRange sql.NullTime
					err := row.Scan(&j.Name, &j.Insecure, &j.IgnoreSignatureCheck, &j.Rate, &j.Worker, &startRange, &endRange, &ignorePattern, &j.Domains)
					if ignorePattern.Valid {
						j.IgnorePattern = &ignorePattern.String
					}
					if startRange.Valid {
						j.StartRange = &startRange.Time
					}
					if endRange.Valid {
						j.EndRange = &endRange.Time
					}
					return j, err
				})
			return err
		}, 0,
	); err != nil {
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, jobs)
}

// importProvider downloads the advisories from the specified source
func (c *Controller) importProvider(ctx *gin.Context) {
	domainsQuery, ok := ctx.GetQuery("domains")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing domains query parameter"})
		return
	}
	domains := strings.Split(domainsQuery, ",")

	t := time.Now().UTC()

	c.downloadWorker.Enqueue(worker.DownloadJob{
		Domains:        domains,
		Worker:         2,
		ForwardQueue:   defaultForwardQueue,
		Presets:        c.cfg.Importer.RemoteValidatorPresets,
		ValidationMode: c.cfg.Importer.ValidationMode,
		Db:             c.db,
		LogFile:        c.cfg.Importer.LogPath + t.Format(time.RFC3339),
		LogLevel:       c.cfg.Importer.LogLevel,
	})
	slog.Info("Queued download for domains", "domains", domains)

	ctx.JSON(http.StatusOK, gin.H{"msg": "queued import job"})
}
