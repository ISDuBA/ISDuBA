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
	"log/slog"
	"time"

	"github.com/ISDuBA/ISDuBA/pkg/config"
	"github.com/ISDuBA/ISDuBA/pkg/database"
)

// refreshDuration is the fallback duration for feeds to be checked for refresh.
const refreshDuration = time.Minute

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

func (m *Manager) refreshFeeds() {
	now := time.Now()
	for _, feed := range m.feeds {
		// Does the feed need a refresh?
		if feed.nextCheck.IsZero() || !feed.nextCheck.Before(now) {
			if err := feed.refresh(); err != nil {
				slog.Error("feed refresh failed", "feed", feed.id, "err", err)
				// TODO: Log to database
			}
			// Even if there was an error try again later.
			feed.nextCheck = time.Now().Add(m.cfg.FeedRefresh)
		}
	}
}

// Run runs the manager. To be used in a Go routine.
func (m *Manager) Run(ctx context.Context) {
	refreshTicker := time.NewTicker(refreshDuration)
	defer refreshTicker.Stop()
	for !m.done {
		m.refreshFeeds()
		select {
		case fn := <-m.fns:
			fn(m)
		case <-ctx.Done():
			return
		case <-refreshTicker.C:
		}
	}
}

// Kill stops the manager.
func (m *Manager) Kill() {
	m.fns <- func(m *Manager) { m.done = true }
}

// AddSource registers a new source.
func (*Manager) AddSource(int64) error {
	// There is currently no code needed here as new sources
	// do not have feeds in particular no active ones.
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
