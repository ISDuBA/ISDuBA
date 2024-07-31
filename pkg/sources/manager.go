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
	"slices"
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

	feeds     []*activeFeed
	sources   []*activeSource
	usedSlots int
	uniqueID  int64
}

// NewManager creates a new downloader.
func NewManager(cfg *config.Sources, db *database.DB) *Manager {
	return &Manager{
		cfg: cfg,
		db:  db,
		fns: make(chan func(*Manager)),
	}
}

// refreshFeeds checks if there are feeds that need reloading
// and does so in that case.
func (m *Manager) refreshFeeds() {
	now := time.Now()
	for _, feed := range m.feeds {
		// Does the feed need a refresh?
		if feed.nextCheck.IsZero() || !now.Before(feed.nextCheck) {
			if err := feed.refresh(m.db); err != nil {
				slog.Error("feed refresh failed", "feed", feed.id, "err", err)
				// TODO: Log to database
			}
			// Even if there was an error try again later.
			feed.nextCheck = time.Now().Add(m.cfg.FeedRefresh)
		}
	}
}

// startDownloads starts downloads if there are enough slots and
// there are things to download.
func (m *Manager) startDownloads() {
	for m.usedSlots < m.cfg.DownloadSlots {
		started := false
		for _, feed := range m.feeds {
			// Has this feed a free slot?
			maxSlots := min(m.cfg.SlotsPerSource, m.cfg.DownloadSlots)
			if feed.source.slots != nil {
				maxSlots = min(maxSlots, *feed.source.slots)
			}
			if feed.source.usedSlots >= maxSlots {
				continue
			}
			// Find a candidate to download.
			location := feed.findWaiting()
			if location == nil {
				continue
			}
			m.usedSlots++
			feed.source.usedSlots++
			location.state = running
			location.id = m.generateID()
			started = true
			// Calling reciever by value is indented here!
			go (*location).download(m, feed)
			if m.usedSlots >= m.cfg.SlotsPerSource {
				return
			}
		}
		if !started {
			return
		}
	}
}

// compactDone removes the locations the feeds which are downloaded.
func (m *Manager) compactDone() {
	for _, feed := range m.feeds {
		feed.locations = slices.DeleteFunc(feed.locations, func(l location) bool {
			return l.state == done
		})
	}
}

func (m *Manager) generateID() (id int64) {
	id = m.uniqueID
	m.uniqueID++
	return
}

func (af *activeFeed) findLocationByID(id int64) *location {
	for i := len(af.locations) - 1; i >= 0; i-- {
		if location := &af.locations[i]; location.id == id {
			return location
		}
	}
	return nil
}

func (af *activeFeed) findWaiting() *location {
	for i := len(af.locations) - 1; i >= 0; i-- {
		if location := &af.locations[i]; location.state == waiting {
			return location
		}
	}
	return nil
}

func (l location) download(m *Manager, f *activeFeed) {
	defer func() {
		m.fns <- func(m *Manager) {
			f.source.usedSlots--
			m.usedSlots--
			f.findLocationByID(l.id).state = done
			slog.Debug("download done", "id", l.id)
		}
	}()
	// TODO: Implement me!
}

// Run runs the manager. To be used in a Go routine.
func (m *Manager) Run(ctx context.Context) {
	refreshTicker := time.NewTicker(refreshDuration)
	defer refreshTicker.Stop()
	for !m.done {
		m.compactDone()
		m.refreshFeeds()
		m.startDownloads()
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
