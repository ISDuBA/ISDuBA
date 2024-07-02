// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

// Package importer implements the import of documents.
package importer

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/ISDuBA/ISDuBA/pkg/database/query"

	"github.com/ISDuBA/ISDuBA/pkg/config"
	"github.com/ISDuBA/ISDuBA/pkg/database"
	"github.com/ISDuBA/ISDuBA/pkg/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
)

// Scheduler schedules download jobs.
type Scheduler struct {
	ctx             context.Context
	db              *database.DB
	cfg             *config.Config
	downloader      *downloadWorker
	cron            *cron.Cron
	cronLock        sync.Mutex
	runningTasks    map[int64]context.CancelFunc
	runningTaskLock sync.Mutex
	notify          chan bool
	wg              sync.WaitGroup
}

// NewScheduler returns a new scheduler.
func NewScheduler(ctx context.Context, db *database.DB, cfg *config.Config) *Scheduler {
	c := cron.New()
	c.Start()

	downloadWorker := newDownloadWorker()

	s := &Scheduler{
		ctx:          ctx,
		db:           db,
		cfg:          cfg,
		downloader:   downloadWorker,
		cron:         c,
		runningTasks: make(map[int64]context.CancelFunc),
		notify:       make(chan bool),
	}
	s.init()
	return s
}

func (s *Scheduler) init() {
	// Mark all running tasks as aborted
	var tasks []models.Task
	if err := s.db.Run(
		s.ctx,
		func(rctx context.Context, conn *pgxpool.Conn) error {
			const fetchSQL = `UPDATE tasks SET status = 'ABORTED' WHERE status = 'RUNNING' RETURNING id, created, job_id, status`
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
		return
	}

	for _, task := range tasks {
		slog.Info("aborted task", "id", task.ID)
	}
}

// AddTask adds a task into the task queue.
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

// AddCron creates a new cron job.
func (s *Scheduler) AddCron(cron models.Cron) (*int64, error) {
	const insertSQL = `INSERT INTO cron (` +
		`name,` +
		`job_id,` +
		`cron_timing` +
		`) VALUES ($1, $2, $3)` +
		`RETURNING id`

	var cronID int64

	if err := s.db.Run(
		s.ctx,
		func(rctx context.Context, conn *pgxpool.Conn) error {
			return conn.QueryRow(rctx, insertSQL,
				cron.Name,
				cron.JobID,
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

// DeleteCron deletes a cron job.
func (s *Scheduler) DeleteCron(cronID int64) (*int64, error) {
	expr := query.FieldEqInt("id", cronID)
	builder := query.SQLBuilder{}
	builder.CreateWhere(expr)

	deleteSQL := `DELETE FROM cron WHERE ` +
		builder.WhereClause

	if err := s.db.Run(
		s.ctx,
		func(rctx context.Context, conn *pgxpool.Conn) error {
			if _, err := conn.Exec(rctx, deleteSQL); err != nil {
				return err
			}
			return nil
		}, 0,
	); err != nil {
		slog.Error("database error", "err", err)
		return nil, err
	}
	s.notify <- true
	return &cronID, nil
}

// AbortTask aborts the specified task.
func (s *Scheduler) AbortTask(taskID int64) error {
	s.runningTaskLock.Lock()
	cancel := s.runningTasks[taskID]
	if cancel == nil {
		return errors.New("task not found")
	}
	delete(s.runningTasks, taskID)
	s.runningTaskLock.Unlock()
	cancel()

	return nil
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
					err := row.Scan(&c.ID, &c.Name, &c.JobID, &c.CronTiming)
					return c, err
				})
			return err
		}, 0,
	); err != nil {
		slog.Error("database error", "err", err)
		return
	}

	s.cronLock.Lock()
	s.cron.Stop()
	// TODO: use cron.Remove(id EntryID) to delete cron job
	s.cron = cron.New()
	s.cronLock.Unlock()
	for _, c := range crons {
		cronJob := func() {
			_, err := s.AddTask(c.JobID)
			if err != nil {
				slog.Error("could not add task", "err", err)
				return
			}
		}
		s.cronLock.Lock()
		_, err := s.cron.AddFunc(c.CronTiming, cronJob)
		s.cronLock.Unlock()
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
					err := row.Scan(&task.ID, &task.Created, &task.JobID, &task.Status)
					return task, err
				})
			return err
		}, 0,
	); err != nil {
		slog.Error("database error", "err", err)
		return
	}

	for _, task := range tasks {
		expr := query.FieldEqInt("id", task.JobID)
		builder := query.SQLBuilder{}
		builder.CreateWhere(expr)

		var jobConf models.JobConfig
		downloadCtx, cancel := context.WithCancel(context.Background())
		s.runningTaskLock.Lock()
		s.runningTasks[task.ID] = cancel
		s.runningTaskLock.Unlock()

		if err := s.db.Run(
			s.ctx,
			func(rctx context.Context, conn *pgxpool.Conn) error {
				fetchSQL := `SELECT name, insecure, ignore_signature_check, rate, worker, start_range, end_range, ignore_pattern, domains, http_headers FROM jobs WHERE ` +
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

						var headers []string
						err := row.Scan(&j.Name, &j.Insecure, &j.IgnoreSignatureCheck, &j.Rate, &j.Worker, &startRange, &endRange, &ignorePattern, &j.Domains, &headers)
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
							jobConf.Headers[s[0]] = s[1]
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
			logFileLocation := s.cfg.Importer.LogPath + "-" + jobConf.Name + "-" + t.Format(time.RFC3339) + ".log"
			s.wg.Add(1)
			if err := s.db.Run(
				s.ctx,
				func(rctx context.Context, conn *pgxpool.Conn) error {
					const fetchSQL = `UPDATE tasks SET log_file = $2 WHERE id = $1 RETURNING id`
					return conn.QueryRow(rctx, fetchSQL, task.ID, logFileLocation).Scan(&task.ID)
				}, 0,
			); err != nil {
				slog.Error("database error", "err", err)
			}
			s.wg.Done()
			var status models.Status
			if err := s.downloader.run(downloadCtx, DownloadJob{
				Config:         jobConf,
				ForwardQueue:   0,
				Presets:        s.cfg.Importer.RemoteValidatorPresets,
				ValidationMode: s.cfg.Importer.ValidationMode,
				Db:             s.db,
				LogFile:        logFileLocation,
				LogLevel:       s.cfg.Importer.LogLevel,
			}); err != nil {
				slog.Error("download error", "err", err)
				if errors.Is(err, context.Canceled) {
					status = models.ABORTED
				} else {
					status = models.FAILED
				}
			} else {
				status = models.COMPLETED
			}

			s.wg.Add(1)
			_, err := s.setTaskState(task.ID, status)
			if err != nil {
				slog.Error("setTaskState error", "err", err)
			}
			s.wg.Done()
		}()
	}
}

func (s *Scheduler) setTaskState(taskID int64, status models.Status) (*int64, error) {
	expr := query.FieldEqInt("id", taskID)
	builder := query.SQLBuilder{}
	builder.CreateWhere(expr)

	var updateTaskID int64

	updateSQL := `UPDATE tasks SET status = $1  WHERE ` +
		builder.WhereClause +
		`RETURNING id`

	if err := s.db.Run(
		s.ctx,
		func(rctx context.Context, conn *pgxpool.Conn) error {
			return conn.QueryRow(rctx, updateSQL, status).Scan(&updateTaskID)
		}, 0,
	); err != nil {
		slog.Error("database error", "err", err)
		return nil, err
	}

	return &updateTaskID, nil
}

// Start starts the scheduler.
func (s *Scheduler) Start() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
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
	}()
}

// Close closes the scheduler
func (s *Scheduler) Close() {
	s.wg.Wait()
	close(s.notify)
	s.cronLock.Lock()
	s.cron.Stop()
	s.cronLock.Unlock()

	// Stop all running tasks
	s.runningTaskLock.Lock()
	tasks := make([]int64, len(s.runningTasks))
	i := 0
	for t := range s.runningTasks {
		tasks[i] = t
		i++
	}
	s.runningTaskLock.Unlock()
	for _, t := range tasks {
		if err := s.AbortTask(t); err != nil {
			slog.Error("could not abort task", "err", err)
		}
	}
}
