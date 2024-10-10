// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

// Package tempstore implements a temporary store for documents.
package forwarder

import (
	"context"
	"time"

	"github.com/ISDuBA/ISDuBA/pkg/config"
	"github.com/ISDuBA/ISDuBA/pkg/database"
	"github.com/ISDuBA/ISDuBA/pkg/database/query"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ForwardManager struct {
	cfg  *config.Forwarder
	db   *database.DB
	fns  chan func(*ForwardManager)
	done bool
}

func NewForwardManager(cfg *config.Forwarder) *ForwardManager {
	return &ForwardManager{
		cfg: cfg,
	}
}

// Run runs the forward manager. To be used in a Go routine.
func (fm *ForwardManager) Run(ctx context.Context) {
	ticker := time.NewTicker(fm.cfg.UpdateInterval)
	defer ticker.Stop()
	for !fm.done {
		select {
		case fn := <-fm.fns:
			fn(fm)
		case <-ctx.Done():
			return
		case <-ticker.C:
			fm.runTargets()
		}
	}
}

// runTargets fetches and sends all new documents to the configured targets.
func (fm *ForwardManager) runTargets() {
}

func (fm *ForwardManager) fetchNewDocuments(ctx context.Context, publisher *string) ([]int64, error) {
	builder := query.SQLBuilder{}

	if publisher != nil {
		builder.CreateWhere(query.FieldEqString("publisher", *publisher))
	}
	var documentIDs []int64

	if err := fm.db.Run(
		ctx,
		func(rctx context.Context, conn *pgxpool.Conn) error {
			fetchSQL := `SELECT documents_id FROM events_log ` +
				`WHERE documents_id in (SELECT id FROM documents WHERE ` + builder.WhereClause + `) ORDER BY time DESC`
			rows, _ := conn.Query(rctx, fetchSQL, builder.Replacements...)
			var err error
			documentIDs, err = pgx.CollectRows(
				rows,
				func(row pgx.CollectableRow) (int64, error) {
					var documentID int64
					err := row.Scan(&documentID)
					return documentID, err
				})
			return err
		}, 0,
	); err != nil {
		return documentIDs, err
	}

	return documentIDs, nil
}

func (fm *ForwardManager) ForwardDocument(documentID int64) error {
	result := make(chan error)
	fm.fns <- func(fm *ForwardManager) {
		result <- nil
	}
	return <-result
}

func (fm *ForwardManager) kill() { fm.done = true }

// Kill shuts down the forward manager.
func (fm *ForwardManager) Kill() { fm.fns <- (*ForwardManager).kill }
