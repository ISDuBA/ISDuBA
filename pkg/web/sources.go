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
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/ISDuBA/ISDuBA/pkg/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ISDuBA/ISDuBA/pkg/database"
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

	// Allow insecure download.
	if insecure, err := strconv.ParseBool(ctx.PostForm("insecure")); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "please specify 'insecure' boolean parameter",
		})
		return
	} else {
		jobConfig.Insecure = insecure
	}

	for _, domain := range jobConfig.Domains {
		if domain == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "don't specify empty domains",
			})
			return
		}
	}

	const insertSQL = `INSERT INTO jobs (` +
		`name,` +
		`insecure,` +
		`ignore_signature_check,` +
		`rate,` +
		`worker,` +
		`start_range,` +
		`end_range,` +
		`ignore_pattern,` +
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
			const fetchSQL = `SELECT id, name, insecure, ignore_signature_check, rate, worker, start_range, end_range, ignore_pattern, domains FROM jobs`
			rows, _ := conn.Query(rctx, fetchSQL)
			var err error
			jobs, err = pgx.CollectRows(
				rows,
				func(row pgx.CollectableRow) (models.JobConfig, error) {
					var j models.JobConfig
					var ignorePattern sql.NullString
					var startRange sql.NullTime
					var endRange sql.NullTime
					err := row.Scan(&j.ID, &j.Name, &j.Insecure, &j.IgnoreSignatureCheck, &j.Rate, &j.Worker, &startRange, &endRange, &ignorePattern, &j.Domains)
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

// addCron creates a new cron job.
func (c *Controller) addCron(ctx *gin.Context) {
	var cron models.Cron

	// We need the name.
	if cron.Name = ctx.PostForm("name"); cron.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "missing 'name'",
		})
		return
	}

	//
	if cron.CronTiming = ctx.PostForm("cron_timing"); cron.CronTiming == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "missing 'cron_timing'",
		})
		return
	}

	var jobIDs string
	if jobIDs := ctx.PostForm("job_id"); jobIDs == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "missing 'job_id'",
		})
		return
	}

	var err error
	cron.JobId, err = strconv.ParseInt(jobIDs, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	const insertSQL = `INSERT INTO cron (` +
		`name,` +
		`job_id,` +
		`cront_timing` +
		`) VALUES ($1, $2, $3)` +
		`RETURNING id`

	var jobID int64

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			return conn.QueryRow(rctx, insertSQL,
				cron.Name,
				cron.JobId,
				cron.CronTiming,
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

// viewCrons returns all cron jobs.
func (c *Controller) viewCrons(ctx *gin.Context) {
	var crons []models.Cron

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			const fetchSQL = `SELECT id, name, job_id, cron_timing FROM jobs`
			rows, _ := conn.Query(rctx, fetchSQL)
			var err error
			crons, err = pgx.CollectRows(
				rows,
				func(row pgx.CollectableRow) (models.Cron, error) {
					var cron models.Cron
					err := row.Scan(&cron.Id, &cron.Name, &cron.JobId, &cron.CronTiming)
					return cron, err
				})
			return err
		}, 0,
	); err != nil {
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, crons)
}

// viewTasks returns all tasks.
func (c *Controller) viewTasks(ctx *gin.Context) {
	var tasks []models.Task

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			const fetchSQL = `SELECT id, created, job_id, status FROM jobs`
			rows, _ := conn.Query(rctx, fetchSQL)
			var err error
			tasks, err = pgx.CollectRows(
				rows,
				func(row pgx.CollectableRow) (models.Task, error) {
					var task models.Task
					err := row.Scan(&task.Id, &task.Created, &task.JobId, &task.Status)
					return task, err
				})
			return err
		}, 0,
	); err != nil {
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}

// runJob runs a configured job
func (c *Controller) runJob(ctx *gin.Context) {
	jobIDs := ctx.Param("id")
	jobID, err := strconv.ParseInt(jobIDs, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expr := database.FieldEqInt("id", jobID)
	builder := database.SQLBuilder{}
	builder.CreateWhere(expr)

	var jobConf models.JobConfig

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			fetchSQL := `SELECT name, insecure, ignore_signature_check, rate, worker, start_range, end_range, ignore_pattern, domains FROM jobs WHERE ` +
				builder.WhereClause
			rows, _ := conn.Query(rctx, fetchSQL)
			var err error
			jobConf, err = pgx.CollectOneRow(
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

	t := time.Now().UTC()

	c.downloadWorker.Enqueue(worker.DownloadJob{
		Config:         jobConf,
		ForwardQueue:   defaultForwardQueue,
		Presets:        c.cfg.Importer.RemoteValidatorPresets,
		ValidationMode: c.cfg.Importer.ValidationMode,
		Db:             c.db,
		LogFile:        c.cfg.Importer.LogPath + t.Format(time.RFC3339) + ".log",
		LogLevel:       c.cfg.Importer.LogLevel,
	})
	slog.Info("Queued download for domains", "domains", jobConf.Domains)

	ctx.JSON(http.StatusOK, gin.H{"msg": "queued import job"})
}
