// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

// Package models implements the handling with advisories.
package models

//go:generate go run ./internal/generators/generate_workflow_diagram.go -o ../../docs/workflow.svg
//go:generate go run ./internal/generators/generate_workflow_ts.go -o ../../client/src/lib/workflow.ts
