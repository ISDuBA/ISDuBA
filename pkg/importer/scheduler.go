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
	db           *database.DB
	cfg          *config.Config
	downloader   *downloadWorker
	cron         *cron.Cron
	runningTasks map[int64]context.CancelFunc
	fns          chan func(scheduler *Scheduler, ctx context.Context)
	done         bool
}

// NewScheduler returns a new scheduler.
func NewScheduler(db *database.DB, cfg *config.Config) *Scheduler {
	c := cron.New()
	c.Start()

	downloadWorker := newDownloadWorker()

	s := &Scheduler{
		db:           db,
		cfg:          cfg,
		downloader:   downloadWorker,
		cron:         c,
		runningTasks: make(map[int64]context.CancelFunc),
		fns:          make(chan func(scheduler *Scheduler, ctx context.Context)),
	}
	return s
}

func (s *Scheduler) init(ctx context.Context) {
	// Mark all running tasks as aborted
	var tasks []models.Task
	if err := s.db.Run(
		ctx,
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
func (s *Scheduler) AddTask(ctx context.Context, jobID int64) (*int64, error) {
	const insertSQL = `INSERT INTO tasks (` +
		`job_id,` +
		`status` +
		`) VALUES ($1, $2)` +
		`RETURNING id`

	var taskID int64

	if err := s.db.Run(
		ctx,
		func(rctx context.Context, conn *pgxpool.Conn) error {
			return conn.QueryRow(rctx, insertSQL,
				jobID,
				models.Queued,
			).Scan(&taskID)
		}, 0,
	); err != nil {
		slog.Error("database error", "err", err)
		return nil, err
	}

	s.fns <- func(s *Scheduler, ctx context.Context) {
		s.runTasks(ctx)
	}

	return &taskID, nil
}

// AddCron creates a new cron job.
func (s *Scheduler) AddCron(ctx context.Context, cron models.Cron) (*int64, error) {
	if err := checkCronTiming(cron.CronTiming); err != nil {
		return nil, err
	}
	const insertSQL = `INSERT INTO cron (` +
		`name,` +
		`job_id,` +
		`cron_timing` +
		`) VALUES ($1, $2, $3)` +
		`RETURNING id`

	var cronID int64

	if err := s.db.Run(
		ctx,
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

	s.fns <- func(s *Scheduler, ctx context.Context) {
		s.runCron(ctx)
	}

	return &cronID, nil
}

func checkCronTiming(cronTiming string) error {
	var options cron.ParseOption
	parser := cron.NewParser(options)

	_, err := parser.Parse(cronTiming)
	return err
}

// UpdateCron updates the specified cron job.
func (s *Scheduler) UpdateCron(ctx context.Context, cron models.Cron) (*int64, error) {
	if err := checkCronTiming(cron.CronTiming); err != nil {
		return nil, err
	}

	expr := query.FieldEqInt("id", cron.ID)
	builder := query.SQLBuilder{}
	builder.CreateWhere(expr)

	var updateCronID int64

	updateSQL := `UPDATE cron SET ` +
		`name = $1,` +
		`job_id = $2,` +
		`cron_timing = $3 ` +
		`WHERE ` +
		builder.WhereClause +
		`RETURNING id`

	if err := s.db.Run(
		ctx,
		func(rctx context.Context, conn *pgxpool.Conn) error {
			return conn.QueryRow(rctx, updateSQL, cron.Name, cron.JobID, cron.CronTiming).Scan(&updateCronID)
		}, 0,
	); err != nil {
		slog.Error("database error", "err", err)
		return nil, err
	}
	s.fns <- func(s *Scheduler, ctx context.Context) {
		s.runCron(ctx)
	}

	return &updateCronID, nil
}

// DeleteCron deletes a cron job.
func (s *Scheduler) DeleteCron(ctx context.Context, cronID int64) (*int64, error) {
	expr := query.FieldEqInt("id", cronID)
	builder := query.SQLBuilder{}
	builder.CreateWhere(expr)

	deleteSQL := `DELETE FROM cron WHERE ` +
		builder.WhereClause

	if err := s.db.Run(
		ctx,
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

	s.fns <- func(s *Scheduler, ctx context.Context) {
		s.runCron(ctx)
	}
	return &cronID, nil
}

// abortTask aborts the specified task.
func (s *Scheduler) abortTask(taskID int64) error {
	cancel := s.runningTasks[taskID]
	if cancel == nil {
		return errors.New("task not found")
	}
	delete(s.runningTasks, taskID)
	cancel()

	return nil
}

// AbortTask aborts the specified task.
func (s *Scheduler) AbortTask(taskID int64) error {
	errChan := make(chan error)
	s.fns <- func(s *Scheduler, ctx context.Context) {
		errChan <- (*Scheduler).abortTask(s, taskID)
	}
	err := <-errChan
	return err
}

func (s *Scheduler) runCron(ctx context.Context) {
	var crons []models.Cron

	if err := s.db.Run(
		ctx,
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

	s.cron.Stop()
	s.cron = cron.New()
	for _, c := range crons {
		cronJob := func() {
			_, err := s.AddTask(ctx, c.JobID)
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

func (s *Scheduler) runTasks(ctx context.Context) {
	var tasks []models.Task

	if err := s.db.Run(
		ctx,
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
		s.runningTasks[task.ID] = cancel

		if err := s.db.Run(
			ctx,
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
						jobConf.Headers = models.ParseHeaders(headers)
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
			if err := s.db.Run(
				ctx,
				func(rctx context.Context, conn *pgxpool.Conn) error {
					const fetchSQL = `UPDATE tasks SET log_file = $2 WHERE id = $1 RETURNING id`
					return conn.QueryRow(rctx, fetchSQL, task.ID, logFileLocation).Scan(&task.ID)
				}, 0,
			); err != nil {
				slog.Error("database error", "err", err)
			}
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
					status = models.Aborted
				} else {
					status = models.Failed
				}
			} else {
				status = models.Completed
			}

			_, err := s.setTaskState(ctx, task.ID, status)
			if err != nil {
				slog.Error("setTaskState error", "err", err)
			}
		}()
	}
}

func (s *Scheduler) setTaskState(ctx context.Context, taskID int64, status models.Status) (*int64, error) {
	expr := query.FieldEqInt("id", taskID)
	builder := query.SQLBuilder{}
	builder.CreateWhere(expr)

	var updateTaskID int64

	updateSQL := `UPDATE tasks SET status = $1  WHERE ` +
		builder.WhereClause +
		`RETURNING id`

	if err := s.db.Run(
		ctx,
		func(rctx context.Context, conn *pgxpool.Conn) error {
			return conn.QueryRow(rctx, updateSQL, status).Scan(&updateTaskID)
		}, 0,
	); err != nil {
		slog.Error("database error", "err", err)
		return nil, err
	}

	return &updateTaskID, nil
}

// Run starts the scheduler.
func (s *Scheduler) Run(ctx context.Context) {
	s.init(ctx)
	defer func() {
		s.cron.Stop()
		for t := range s.runningTasks {
			if err := s.abortTask(t); err != nil {
				slog.Error("could not abort task", "err", err)
			}
		}
	}()

	for !s.done {
		select {
		case <-ctx.Done():
			return
		case fn := <-s.fns:
			fn(s, ctx)
		}
	}
}

func (s *Scheduler) kill(_ context.Context) {
	s.done = true
}

// Kill kills the scheduler.
func (s *Scheduler) Kill() {
	s.fns <- (*Scheduler).kill
}
