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
	"database/sql"
	"errors"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ISDuBA/ISDuBA/pkg/database/query"

	"github.com/ISDuBA/ISDuBA/pkg/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gin-gonic/gin"
)

// addJob creates a new job configuration.
func (c *Controller) addJob(ctx *gin.Context) {
	jobConfig := models.JobConfig{
		Worker: 1,
	}

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
	var err error
	if jobConfig.Insecure, err = strconv.ParseBool(ctx.PostForm("insecure")); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "please specify 'insecure' boolean parameter",
		})
		return
	}

	if jobConfig.IgnoreSignatureCheck, err = strconv.ParseBool(ctx.PostForm("ignore_signature_check")); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "please specify 'ignore_signature_check' boolean parameter",
		})
		return
	}

	for _, domain := range jobConfig.Domains {
		if domain == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "don't specify empty domains",
			})
			return
		}
	}

	if clientKey := ctx.PostForm("client_key"); clientKey != "" {
		jobConfig.ClientKey = &clientKey
	}

	if clientPassphrase := ctx.PostForm("client_passphrase"); clientPassphrase != "" {
		jobConfig.ClientPassphrase = &clientPassphrase
	}

	if ignorePattern := ctx.PostForm("ignore_pattern"); ignorePattern != "" {
		jobConfig.IgnorePattern = &ignorePattern
	}

	header := ctx.PostFormArray("header")
	if header == nil {
		header = []string{}
	}

	var file *multipart.FileHeader
	if file, err = ctx.FormFile("client_cert"); err == nil {
		open, err := file.Open()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer open.Close()

		certFile := bytes.NewBuffer(nil)
		if _, err := io.Copy(certFile, open); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		jobConfig.ClientCerts = certFile.Bytes()
	}

	if startRange, ok := ctx.GetPostForm("start_range"); ok {
		startRangeTime, err := time.Parse("2006-01-02", startRange)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			jobConfig.StartRange = &startRangeTime
		}
	}

	if endRange, ok := ctx.GetPostForm("start_range"); ok {
		endRangeTime, err := time.Parse("2006-01-02", endRange)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			jobConfig.EndRange = &endRangeTime
		}
	}

	if rate, ok := ctx.GetPostForm("rate"); ok {
		rateFloat, err := strconv.ParseFloat(rate, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			jobConfig.Rate = &rateFloat
		}
	}

	if worker, ok := ctx.GetPostForm("worker"); ok {
		workerInt, err := strconv.ParseInt(worker, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			jobConfig.Worker = int(workerInt)
		}
	}

	const insertSQL = `INSERT INTO jobs (` +
		`name,` +
		`insecure,` +
		`ignore_signature_check,` +
		`client_key,` +
		`client_passphrase,` +
		`rate,` +
		`worker,` +
		`start_range,` +
		`end_range,` +
		`ignore_pattern,` +
		`domains,` +
		`http_headers,` +
		`client_certs` +
		`) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)` +
		`RETURNING id`

	var jobID int64

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			return conn.QueryRow(rctx, insertSQL,
				jobConfig.Name,
				jobConfig.Insecure,
				jobConfig.IgnoreSignatureCheck,
				jobConfig.ClientKey,
				jobConfig.ClientPassphrase,
				jobConfig.Rate,
				jobConfig.Worker,
				jobConfig.StartRange,
				jobConfig.EndRange,
				jobConfig.IgnorePattern,
				jobConfig.Domains,
				header,
				jobConfig.ClientCerts,
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

// updateJob updates a job configuration.
func (c *Controller) updateJob(ctx *gin.Context) {
	jobConfig := models.JobConfig{}

	var jobIDs string
	if jobIDs = ctx.PostForm("job_id"); jobIDs == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "missing 'job_id'",
		})
		return
	}

	var err error
	jobConfig.ID, err = strconv.ParseInt(jobIDs, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
	if jobConfig.Insecure, err = strconv.ParseBool(ctx.PostForm("insecure")); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "please specify 'insecure' boolean parameter",
		})
		return
	}

	if jobConfig.IgnoreSignatureCheck, err = strconv.ParseBool(ctx.PostForm("ignore_signature_check")); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "please specify 'ignore_signature_check' boolean parameter",
		})
		return
	}

	for _, domain := range jobConfig.Domains {
		if domain == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "don't specify empty domains",
			})
			return
		}
	}

	if clientKey := ctx.PostForm("client_key"); clientKey != "" {
		jobConfig.ClientKey = &clientKey
	}

	if clientPassphrase := ctx.PostForm("client_passphrase"); clientPassphrase != "" {
		jobConfig.ClientPassphrase = &clientPassphrase
	}

	if ignorePattern := ctx.PostForm("ignore_pattern"); ignorePattern != "" {
		jobConfig.IgnorePattern = &ignorePattern
	}

	if rate, ok := ctx.GetPostForm("rate"); ok {
		rateFloat, err := strconv.ParseFloat(rate, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			jobConfig.Rate = &rateFloat
		}
	}

	if worker, ok := ctx.GetPostForm("worker"); ok {
		workerInt, err := strconv.ParseInt(worker, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			jobConfig.Worker = int(workerInt)
		}
	}

	header := ctx.PostFormArray("header")
	if header == nil {
		header = []string{}
	}

	var file *multipart.FileHeader
	if file, err = ctx.FormFile("client_cert"); err == nil {
		open, err := file.Open()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer open.Close()

		certFile := bytes.NewBuffer(nil)
		if _, err := io.Copy(certFile, open); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		jobConfig.ClientCerts = certFile.Bytes()
	}

	if startRange, ok := ctx.GetPostForm("start_range"); ok {
		startRangeTime, err := time.Parse("2006-01-02", startRange)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			jobConfig.StartRange = &startRangeTime
		}
	}

	if endRange, ok := ctx.GetPostForm("start_range"); ok {
		endRangeTime, err := time.Parse("2006-01-02", endRange)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			jobConfig.EndRange = &endRangeTime
		}
	}

	expr := query.FieldEqInt("id", jobConfig.ID)
	builder := query.SQLBuilder{}
	builder.CreateWhere(expr)

	updateSQL := `UPDATE jobs SET ` +
		`name = $1,` +
		`insecure = $2,` +
		`ignore_signature_check = $3,` +
		`rate = $4,` +
		`worker = $5,` +
		`start_range = $6,` +
		`end_range = $7,` +
		`ignore_pattern = $8,` +
		`domains = $9, ` +
		`http_headers = $10,` +
		`client_certs = $11 ` +
		`WHERE ` +
		builder.WhereClause +
		`RETURNING id`

	var jobID int64

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			return conn.QueryRow(rctx, updateSQL,
				jobConfig.Name,
				jobConfig.Insecure,
				jobConfig.IgnoreSignatureCheck,
				jobConfig.Rate,
				jobConfig.Worker,
				jobConfig.StartRange,
				jobConfig.EndRange,
				jobConfig.IgnorePattern,
				jobConfig.Domains,
				header,
				jobConfig.ClientCerts,
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
	ctx.JSON(http.StatusOK, gin.H{
		"id": jobID,
	})
}

// deleteJob deletes a job configuration.
func (c *Controller) deleteJob(ctx *gin.Context) {
	jobIDs := ctx.Param("id")
	jobID, err := strconv.ParseInt(jobIDs, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expr := query.FieldEqInt("id", jobID)
	builder := query.SQLBuilder{}
	builder.CreateWhere(expr)

	deleteSQL := `DELETE FROM jobs WHERE ` +
		builder.WhereClause

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			if _, err := conn.Exec(rctx, deleteSQL); err != nil {
				return err
			}
			return nil
		}, 0,
	); err != nil {
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id": jobID,
	})
}

// viewJobs returns all configured jobs.
func (c *Controller) viewJobs(ctx *gin.Context) {
	var jobs []models.JobConfig

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			const fetchSQL = `SELECT id, name, insecure, ignore_signature_check, rate, worker, start_range, end_range, ignore_pattern, domains, http_headers FROM jobs`
			rows, _ := conn.Query(rctx, fetchSQL)
			var err error
			jobs, err = pgx.CollectRows(
				rows,
				func(row pgx.CollectableRow) (models.JobConfig, error) {
					var j models.JobConfig
					var ignorePattern sql.NullString
					var startRange sql.NullTime
					var endRange sql.NullTime

					var headers []string
					err := row.Scan(&j.ID, &j.Name, &j.Insecure, &j.IgnoreSignatureCheck, &j.Rate, &j.Worker, &startRange, &endRange, &ignorePattern, &j.Domains, &headers)
					if ignorePattern.Valid {
						j.IgnorePattern = &ignorePattern.String
					}
					if startRange.Valid {
						j.StartRange = &startRange.Time
					}
					if endRange.Valid {
						j.EndRange = &endRange.Time
					}

					for _, h := range headers {
						s := strings.Split(h, ":")
						if len(s) < 2 {
							continue
						}
						j.Headers[s[0]] = s[1]
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
	if jobIDs = ctx.PostForm("job_id"); jobIDs == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "missing 'job_id'",
		})
		return
	}

	var err error
	cron.JobID, err = strconv.ParseInt(jobIDs, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cronID, err := c.scheduler.AddCron(ctx, cron)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"cron_id": cronID,
	})
}

// deleteCron returns all cron jobs.
func (c *Controller) deleteCron(ctx *gin.Context) {
	cronIDs := ctx.Param("id")
	var cronID int64
	cronID, err := strconv.ParseInt(cronIDs, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, err := c.scheduler.DeleteCron(ctx, cronID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"cron_id": cronID,
	})
}

// viewCrons returns all cron jobs.
func (c *Controller) viewCrons(ctx *gin.Context) {
	var crons []models.Cron

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			const fetchSQL = `SELECT id, name, job_id, cron_timing FROM cron`
			rows, _ := conn.Query(rctx, fetchSQL)
			var err error
			crons, err = pgx.CollectRows(
				rows,
				func(row pgx.CollectableRow) (models.Cron, error) {
					var cron models.Cron
					err := row.Scan(&cron.ID, &cron.Name, &cron.JobID, &cron.CronTiming)
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

// getTaskLog returns the log file of the task.
func (c *Controller) getTaskLog(ctx *gin.Context) {
	taskIDs := ctx.Param("id")
	var taskID int64
	taskID, err := strconv.ParseInt(taskIDs, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var logFileLocation string

	expr := query.FieldEqInt("id", taskID)
	builder := query.SQLBuilder{}
	builder.CreateWhere(expr)

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			fetchSQL := `SELECT log_file FROM tasks WHERE ` +
				builder.WhereClause
			err := conn.QueryRow(rctx, fetchSQL).Scan(&logFileLocation)
			return err
		}, 0,
	); err != nil {
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.File(logFileLocation)
}

// viewTasks returns all tasks.
func (c *Controller) viewTasks(ctx *gin.Context) {
	var tasks []models.Task

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			const fetchSQL = `SELECT id, created, job_id, status FROM tasks`
			rows, _ := conn.Query(rctx, fetchSQL)
			var err error
			tasks, err = pgx.CollectRows(
				rows,
				func(row pgx.CollectableRow) (models.Task, error) {
					var task models.Task
					err := row.Scan(&task.ID, &task.Created, &task.JobID, &task.Status)
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

	task, err := c.scheduler.AddTask(ctx, jobID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"task_id": task})
}

// abortTask runs a configured job
func (c *Controller) abortTask(ctx *gin.Context) {
	jobIDs := ctx.Param("id")
	jobID, err := strconv.ParseInt(jobIDs, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.scheduler.AbortTask(jobID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
