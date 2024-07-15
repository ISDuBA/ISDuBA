// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package models

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ISDuBA/ISDuBA/pkg/database/query"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
	"time"
)

// Status of the task.
type Status string

// Task states.
const (
	Queued    Status = "QUEUED"
	Running   Status = "RUNNING"
	Aborted   Status = "ABORTED"
	Failed    Status = "FAILED"
	Completed Status = "COMPLETED"
)

// Task represents a job that can be run.
type Task struct {
	ID      int64     `json:"task_id"`
	Created time.Time `json:"created"`
	JobID   int64     `json:"job_id"`
	Status  Status    `json:"status"`
}

// Cron represents a cron job.
type Cron struct {
	ID         int64  `json:"cron_id"`
	Name       string `json:"name"`
	JobID      int64  `json:"job_id"`
	CronTiming string `json:"cron_timing"`
}

// JobConfig represents a job configuration.
type JobConfig struct {
	ID                   int64             `json:"id"`
	Name                 string            `json:"name"`
	Insecure             bool              `json:"insecure"`
	IgnoreSignatureCheck bool              `json:"ignore_signature_check"`
	ClientCerts          []byte            `json:"client_certs,omitempty"`
	ClientKey            *string           `json:"client_key,omitempty"`
	ClientPassphrase     *string           `json:"client_passphrase,omitempty"`
	Rate                 *float64          `json:"rate,omitempty"`
	Worker               int               `json:"worker"`
	StartRange           *time.Time        `json:"start_range,omitempty"`
	EndRange             *time.Time        `json:"end_range,omitempty"`
	IgnorePattern        *string           `json:"ignore_pattern,omitempty"`
	Domains              []string          `json:"domains"`
	Headers              map[string]string `json:"headers,omitempty"`
}

func ParseHeaders(headers []string) map[string]string {
	var parsed map[string]string
	for _, h := range headers {
		s := strings.Split(h, ":")
		if len(s) < 2 {
			continue
		}
		parsed[s[0]] = s[1]
	}
	return parsed
}

func SerializeHeaders(headers map[string]string) []string {
	serialized := make([]string, 0)
	for k, v := range headers {
		serialized = append(serialized, fmt.Sprintf("%s:%s", k, v))
	}
	return serialized
}

func InsertJob(ctx context.Context, conn *pgxpool.Conn, jobConfig JobConfig) (*int64, error) {
	headers := SerializeHeaders(jobConfig.Headers)
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

	if err := conn.QueryRow(ctx, insertSQL,
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
		headers,
		jobConfig.ClientCerts,
	).Scan(&jobID); err != nil {
		return nil, err
	}

	return &jobID, nil
}

func UpdateJob(ctx context.Context, conn *pgxpool.Conn, jobConfig JobConfig) (*int64, error) {
	headers := SerializeHeaders(jobConfig.Headers)

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

	if err := conn.QueryRow(ctx, updateSQL,
		jobConfig.Name,
		jobConfig.Insecure,
		jobConfig.IgnoreSignatureCheck,
		jobConfig.Rate,
		jobConfig.Worker,
		jobConfig.StartRange,
		jobConfig.EndRange,
		jobConfig.IgnorePattern,
		jobConfig.Domains,
		headers,
		jobConfig.ClientCerts,
	).Scan(&jobID); err != nil {
		return nil, err
	}

	return &jobID, nil
}

func DeleteJob(ctx context.Context, conn *pgxpool.Conn, jobID int64) (*int64, error) {
	expr := query.FieldEqInt("id", jobID)
	builder := query.SQLBuilder{}
	builder.CreateWhere(expr)

	deleteSQL := `DELETE FROM jobs WHERE ` +
		builder.WhereClause

	if _, err := conn.Exec(ctx, deleteSQL); err != nil {
		return nil, err
	}

	return &jobID, nil
}

func GetJobs(ctx context.Context, conn *pgxpool.Conn) ([]JobConfig, error) {
	var jobs []JobConfig

	const fetchSQL = `SELECT id, name, insecure, ignore_signature_check, rate, worker, start_range, end_range, ignore_pattern, domains, http_headers FROM jobs`
	rows, _ := conn.Query(ctx, fetchSQL)
	var err error
	jobs, err = pgx.CollectRows(
		rows,
		func(row pgx.CollectableRow) (JobConfig, error) {
			var j JobConfig
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
	return jobs, err
}
