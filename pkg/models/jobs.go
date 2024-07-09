// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package models

import "time"

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
