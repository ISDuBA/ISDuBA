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

	"github.com/ISDuBA/ISDuBA/pkg/config"
	"github.com/ISDuBA/ISDuBA/pkg/database"
)

// Manager fetches advisories from sources.
type Manager struct {
	cfg  *config.Sources
	db   *database.DB
	fns  chan func(*Manager)
	done bool
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
func (m *Manager) Boot() error {
	// TODO: Implement me!:
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
func (m *Manager) AddFeed(sourceID, feedID int64) error {
	result := make(chan error)
	m.fns <- func(_ *Manager) {
		_ = sourceID
		_ = feedID
		// TODO: Implement me!
		result <- nil
	}
	return <-result
}
