// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package sources

import (
	"context"
	"fmt"
	"log/slog"
	"slices"

	"github.com/ISDuBA/ISDuBA/pkg/config"
	"github.com/ISDuBA/ISDuBA/pkg/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Manager fetches advisories from sources.
type Manager struct {
	cfg  *config.Sources
	db   *database.DB
	fns  chan func(*Manager)
	done bool

	feeds   []*activeFeed
	sources []*activeSource
}

// NewManager creates a new downloader.
func NewManager(cfg *config.Sources, db *database.DB) *Manager {
	return &Manager{
		cfg: cfg,
		db:  db,
		fns: make(chan func(*Manager)),
	}
}

// Boot loads the sources from database.
func (m *Manager) Boot(ctx context.Context) error {

	const (
		sourcesSQL = `SELECT id, rate, slots FROM sources WHERE active`
		feedsSQL   = `SELECT f.id, sources_id, url, rolie ` +
			`FROM feeds f JOIN sources s ON f.sources_id = s.id ` +
			`WHERE active`
	)

	if err := m.db.Run(
		ctx,
		func(rctx context.Context, con *pgxpool.Conn) error {
			tx, err := con.BeginTx(rctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})
			if err != nil {
				return fmt.Errorf("starting transaction failed: %w", err)
			}
			defer tx.Rollback(rctx)

			// Collect active sources.
			srows, err := tx.Query(rctx, sourcesSQL)
			if err != nil {
				return fmt.Errorf("querying sources failed: %w", err)
			}
			m.sources, err = pgx.CollectRows(srows, func(row pgx.CollectableRow) (*activeSource, error) {
				var s activeSource
				return &s, row.Scan(&s.id, &s.rate, &s.slots)
			})
			if err != nil {
				return fmt.Errorf("collecting active sources failed: %w", err)
			}

			// Collect active feeds
			frows, err := tx.Query(rctx, feedsSQL)
			if err != nil {
				return fmt.Errorf("querying sources failed: %w", err)
			}
			m.feeds, err = pgx.CollectRows(frows, func(row pgx.CollectableRow) (*activeFeed, error) {
				var f activeFeed
				var sid int64
				if err := row.Scan(&f.id, &sid, &f.url, &f.rolie); err != nil {
					return nil, err
				}
				// Add to list of active feeds.
				idx := slices.IndexFunc(m.sources, func(s *activeSource) bool { return s.id == sid })
				if idx == -1 {
					return nil, fmt.Errorf("cannot find source id %d", sid)
				}
				m.sources[idx].feeds = append(m.sources[idx].feeds, &f)
				return &f, nil
			})
			if err != nil {
				return fmt.Errorf("collecting active feeds failed: %w", err)
			}

			return tx.Commit(rctx)
		}, 0,
	); err != nil {
		return fmt.Errorf("fetching active feeds failed: %w", err)
	}

	slog.Info("number of active sources", "num", len(m.sources))
	slog.Info("number of active feeds", "num", len(m.feeds))

	return nil
}

// Run runs the manager. To be used in a Go routine.
func (m *Manager) Run(ctx context.Context) {
	for !m.done {
		select {
		case fn := <-m.fns:
			fn(m)
		case <-ctx.Done():
			return
		}
	}
}

// Kill stops the manager.
func (m *Manager) Kill() {
	m.fns <- func(m *Manager) { m.done = true }
}

// AddSource registers a new source.
func (*Manager) AddSource(int64) error {
	return nil
}

// RemoveSource removes a sources from manager.
func (m *Manager) RemoveSource(sourceID int64) error {
	result := make(chan error)
	m.fns <- func(_ *Manager) {
		_ = sourceID
		// TODO: Implement me!
		result <- nil
	}
	return <-result
}

// AddFeed adds a new feed to a source.
func (m *Manager) AddFeed(feedID int64) error {
	result := make(chan error)
	m.fns <- func(_ *Manager) {
		_ = feedID
		// TODO: Implement me!
		result <- nil
	}
	return <-result
}

// RemoveFeed removes a feed from a source.
func (m *Manager) RemoveFeed(feedID int64) error {
	result := make(chan error)
	m.fns <- func(_ *Manager) {
		_ = feedID
		// TODO: Implement me!
		result <- nil
	}
	return <-result
}
