// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package models

import "time"

type Status string

const (
	QUEUED    Status = "QUEUED"
	RUNNING   Status = "RUNNING"
	ABORTED   Status = "ABORTED"
	FAILED    Status = "FAILED"
	COMPLETED Status = "COMPLETED"
)

type Task struct {
	Id      int64     `json:"task_id"`
	Created time.Time `json:"created"`
	JobId   int64     `json:"job_id"`
	Status  Status    `json:"status"`
}

type Cron struct {
	Id         int64  `json:"cron_id"`
	Name       string `json:"name"`
	JobId      int64  `json:"job_id"`
	CronTiming string `json:"cron_timing"`
}

// JobConfig represents a job configuration.
type JobConfig struct {
	ID                   int64  `json:"id"`
	Name                 string `json:"name"`
	Insecure             bool   `json:"insecure"`
	IgnoreSignatureCheck bool   `json:"ignore_signature_check"`
	// TODO: Implement client certs
	// ClientCerts []string `json:"client_certs"`
	ClientKey        *string    `json:"client_key,omitempty"`
	ClientPassphrase *string    `json:"client_passphrase,omitempty"`
	Rate             *float64   `json:"rate,omitempty"`
	Worker           int        `json:"worker"`
	StartRange       *time.Time `json:"start_range,omitempty"`
	EndRange         *time.Time `json:"end_range,omitempty"`
	IgnorePattern    *string    `json:"ignore_pattern,omitempty"`
	Domains          []string   `json:"domains"`
}
