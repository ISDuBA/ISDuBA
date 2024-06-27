// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package worker

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/ISDuBA/ISDuBA/pkg/config"
	"github.com/ISDuBA/ISDuBA/pkg/database"
	"github.com/ISDuBA/ISDuBA/pkg/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	ctx        context.Context
	db         *database.DB
	cfg        *config.Config
	downloader *DownloadWorker
	cron       *cron.Cron
	notify     chan bool
}

func NewScheduler(ctx context.Context, db *database.DB, cfg *config.Config) *Scheduler {
	c := cron.New()
	c.Start()

	downloadWorker := NewDownloadWorker()

	return &Scheduler{
		ctx:        ctx,
		db:         db,
		cfg:        cfg,
		downloader: downloadWorker,
		cron:       c,
		notify:     make(chan bool),
	}
}

func (s *Scheduler) AddTask(jobID int64) (*int64, error) {
	const insertSQL = `INSERT INTO tasks (` +
		`job_id,` +
		`status` +
		`) VALUES ($1, $2)` +
		`RETURNING id`

	var taskID int64

	if err := s.db.Run(
		s.ctx,
		func(rctx context.Context, conn *pgxpool.Conn) error {
			return conn.QueryRow(rctx, insertSQL,
				jobID,
				models.QUEUED,
			).Scan(&taskID)
		}, 0,
	); err != nil {
		slog.Error("database error", "err", err)
		return nil, err
	}

	s.notify <- true
	return &taskID, nil
}

func (s *Scheduler) AddCron(cron models.Cron) (*int64, error) {
	const insertSQL = `INSERT INTO cron (` +
		`name,` +
		`job_id,` +
		`cront_timing` +
		`) VALUES ($1, $2, $3)` +
		`RETURNING id`

	var cronID int64

	if err := s.db.Run(
		s.ctx,
		func(rctx context.Context, conn *pgxpool.Conn) error {
			return conn.QueryRow(rctx, insertSQL,
				cron.Name,
				cron.JobId,
				cron.CronTiming,
			).Scan(&cronID)
		}, 0,
	); err != nil {
		var pgErr *pgconn.PgError
		// Unique constraint violation
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, errors.New("already in database")
		}
		slog.Error("database error", "err", err)
		return nil, err
	}
	s.notify <- true
	return &cronID, nil
}

func (s *Scheduler) runCron() {
	var crons []models.Cron

	if err := s.db.Run(
		s.ctx,
		func(rctx context.Context, conn *pgxpool.Conn) error {
			const fetchSQL = `SELECT id, name, job_id, cron_timing FROM cron`
			rows, _ := conn.Query(rctx, fetchSQL)
			var err error
			crons, err = pgx.CollectRows(
				rows,
				func(row pgx.CollectableRow) (models.Cron, error) {
					var c models.Cron
					err := row.Scan(&c.Id, &c.Name, &c.JobId, &c.CronTiming)
					return c, err
				})
			return err
		}, 0,
	); err != nil {
		slog.Error("database error", "err", err)
		return
	}

	s.cron.Stop()
	// TODO: use cron.Remove(id EntryID) to delete cron job
	s.cron = cron.New()
	for _, c := range crons {
		cronJob := func() {
			_, err := s.AddTask(c.JobId)
			if err != nil {
				slog.Error("could not add task", "err", err)
				return
			}
		}
		_, err := s.cron.AddFunc(c.CronTiming, cronJob)
		if err != nil {
			slog.Error("could not add cron job", "err", err)
			return
		}
	}
}

func (s *Scheduler) runTasks() {
	var tasks []models.Task

	if err := s.db.Run(
		s.ctx,
		func(rctx context.Context, conn *pgxpool.Conn) error {
			const fetchSQL = `UPDATE tasks SET status = 'RUNNING' WHERE status = 'QUEUED' RETURNING id, created, job_id, status`
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
		return
	}

	for _, task := range tasks {
		expr := database.FieldEqInt("id", task.JobId)
		builder := database.SQLBuilder{}
		builder.CreateWhere(expr)

		var jobConf models.JobConfig

		if err := s.db.Run(
			s.ctx,
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
			return
		}

		t := time.Now().UTC()

		go func() {
			var status models.Status
			if err := s.downloader.Run(s.ctx, DownloadJob{
				Config:         jobConf,
				ForwardQueue:   0,
				Presets:        s.cfg.Importer.RemoteValidatorPresets,
				ValidationMode: s.cfg.Importer.ValidationMode,
				Db:             s.db,
				LogFile:        s.cfg.Importer.LogPath + "-" + jobConf.Name + "-" + t.Format(time.RFC3339) + ".log",
				LogLevel:       s.cfg.Importer.LogLevel,
			}); err != nil {
				slog.Error("download error", "err", err)
				status = models.FAILED
			} else {
				status = models.COMPLETED

				_, err := s.setTaskState(task.Id, status)
				if err != nil {
					slog.Error("setTaskState error", "err", err)
				}
			}
		}()
	}
}

func (s *Scheduler) setTaskState(taskID int64, status models.Status) (*int64, error) {
	expr := database.FieldEqInt("id", taskID)
	builder := database.SQLBuilder{}
	builder.CreateWhere(expr)

	var updateTaskId int64

	updateSql := `UPDATE tasks SET status = $1  WHERE ` +
		builder.WhereClause +
		`RETURNING id`

	if err := s.db.Run(
		s.ctx,
		func(rctx context.Context, conn *pgxpool.Conn) error {
			return conn.QueryRow(rctx, updateSql, status).Scan(&updateTaskId)
		}, 0,
	); err != nil {
		slog.Error("database error", "err", err)
		return nil, err
	}

	return &updateTaskId, nil
}

func (s *Scheduler) Run() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-s.notify:
			func() {
				s.runTasks()
				s.runCron()
			}()
		}
	}
}

func (s *Scheduler) Close() {
	s.cron.Stop()
	close(s.notify)
}
