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
	"github.com/csaf-poc/csaf_distribution/v3/csaf/filter"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/ISDuBA/ISDuBA/pkg/database/query"

	"github.com/ISDuBA/ISDuBA/pkg/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gin-gonic/gin"
)

func parseJobConfig(ctx *gin.Context, requireID bool) (*models.JobConfig, error) {

	jobConfig := models.JobConfig{
		Worker: 1,
	}
	var err error

	if requireID {
		var jobIDs string
		if jobIDs = ctx.PostForm("job_id"); jobIDs == "" {
			return nil, errors.New("job_id is required")
		}

		jobConfig.ID, err = strconv.ParseInt(jobIDs, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	// We need the name.
	if jobConfig.Name = ctx.PostForm("name"); jobConfig.Name == "" {
		return nil, errors.New("missing 'name'")

	}

	// Domains to download from.
	if jobConfig.Domains = ctx.PostFormArray("domains"); len(jobConfig.Domains) == 0 {
		return nil, errors.New("missing 'domains'")
	}

	// Allow insecure download.
	if jobConfig.Insecure, err = strconv.ParseBool(ctx.PostForm("insecure")); err != nil {
		return nil, errors.New("invalid 'insecure'")
	}

	if jobConfig.IgnoreSignatureCheck, err = strconv.ParseBool(ctx.PostForm("ignore_signature_check")); err != nil {
		return nil, errors.New("invalid 'ignore_signature_check'")
	}

	for _, domain := range jobConfig.Domains {
		if domain == "" {
			return nil, errors.New("missing 'domains'")
		}
	}

	if clientKey := ctx.PostForm("client_key"); clientKey != "" {
		jobConfig.ClientKey = &clientKey
	}

	if clientPassphrase := ctx.PostForm("client_passphrase"); clientPassphrase != "" {
		jobConfig.ClientPassphrase = &clientPassphrase
	}

	if ignorePattern := ctx.PostForm("ignore_pattern"); ignorePattern != "" {
		_, err = filter.NewPatternMatcher([]string{ignorePattern})
		if err != nil {
			return nil, err
		}
		jobConfig.IgnorePattern = &ignorePattern
	}

	header := ctx.PostFormArray("header")
	if header == nil {
		header = []string{}
	}
	jobConfig.Headers = models.ParseHeaders(header)

	var file *multipart.FileHeader
	if file, err = ctx.FormFile("client_cert"); err == nil {
		open, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer open.Close()

		certFile := bytes.NewBuffer(nil)
		if _, err := io.Copy(certFile, open); err != nil {
			return nil, err
		}
		jobConfig.ClientCerts = certFile.Bytes()
	}

	if startRange, ok := ctx.GetPostForm("start_range"); ok {
		startRangeTime, err := time.Parse("2006-01-02", startRange)
		if err != nil {
			return nil, err
		}
		jobConfig.StartRange = &startRangeTime
	}

	if endRange, ok := ctx.GetPostForm("end_range"); ok {
		endRangeTime, err := time.Parse("2006-01-02", endRange)
		if err != nil {
			return nil, err
		}
		jobConfig.EndRange = &endRangeTime
	}

	if (jobConfig.StartRange == nil) != (jobConfig.EndRange == nil) {
		return nil, errors.New("specify both ranges'")
	}

	if rate, ok := ctx.GetPostForm("rate"); ok {
		rateFloat, err := strconv.ParseFloat(rate, 32)
		if err != nil {
			return nil, err
		}
		jobConfig.Rate = &rateFloat
	}

	if worker, ok := ctx.GetPostForm("worker"); ok {
		workerInt, err := strconv.ParseInt(worker, 10, 64)
		if err != nil {
			return nil, err
		}
		jobConfig.Worker = int(workerInt)
	}
	return &jobConfig, nil
}

// addJob creates a new job configuration.
func (c *Controller) addJob(ctx *gin.Context) {
	jobConfig, err := parseJobConfig(ctx, false)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var id *int64
	if err = c.db.Run(ctx, func(ctx context.Context, conn *pgxpool.Conn) error {
		id, err = models.InsertJob(ctx, conn, *jobConfig)
		return err
	}, 0); err != nil {
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

// updateJob updates a job configuration.
func (c *Controller) updateJob(ctx *gin.Context) {
	jobConfig, err := parseJobConfig(ctx, true)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var id *int64
	if err = c.db.Run(ctx, func(ctx context.Context, conn *pgxpool.Conn) error {
		id, err = models.UpdateJob(ctx, conn, *jobConfig)
		return err
	}, 0); err != nil {
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id": id,
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
	var id *int64
	if err = c.db.Run(ctx, func(ctx context.Context, conn *pgxpool.Conn) error {
		id, err = models.DeleteJob(ctx, conn, jobID)
		return err
	}, 0); err != nil {
		slog.Error("database error", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

// viewJobs returns all configured jobs.
func (c *Controller) viewJobs(ctx *gin.Context) {
	var jobs []models.JobConfig

	var err error
	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			jobs, err = models.GetJobs(ctx, conn)
			return err
		}, 0); err != nil {
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
