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
	"errors"
	"fmt"
	"net/url"
	"slices"
	"time"

	"github.com/ISDuBA/ISDuBA/pkg/config"
	"github.com/ISDuBA/ISDuBA/pkg/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// refreshDuration is the fallback duration for feeds to be checked for refresh.
const refreshDuration = time.Minute

// Manager fetches advisories from sources.
type Manager struct {
	cfg  *config.Sources
	db   *database.DB
	fns  chan func(*Manager)
	done bool

	feeds     []*feed
	sources   []*source
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
		// Ignore inactive sources.
		if !feed.source.active {
			continue
		}
		// Does the feed need a refresh?
		if feed.nextCheck.IsZero() || !now.Before(feed.nextCheck) {
			if err := feed.refresh(m); err != nil {
				feed.log(m, config.ErrorFeedLogLevel, "feed refresh failed: %v", err.Error())
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
			// Ignore inactive sources.
			if !feed.source.active {
				continue
			}
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
			// Calling reciever by value is intended here!
			go (*location).download(m, feed, m.downloadDone(feed, location.id))
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

func (m *Manager) generateID() int64 {
	// Start with 1 to avoid clashes with zeroed locations.
	m.uniqueID++
	return m.uniqueID
}

func (f *feed) findLocationByID(id int64) *location {
	for i := len(f.locations) - 1; i >= 0; i-- {
		if location := &f.locations[i]; location.id == id {
			return location
		}
	}
	return nil
}

func (f *feed) findWaiting() *location {
	for i := len(f.locations) - 1; i >= 0; i-- {
		if location := &f.locations[i]; location.state == waiting {
			return location
		}
	}
	return nil
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

// ping wakes up the manager.
func (m *Manager) ping() {}

// Kill stops the manager.
func (m *Manager) Kill() {
	m.fns <- func(m *Manager) { m.done = true }
}

func (m *Manager) removeSource(sourceID int64) error {
	m.feeds = slices.DeleteFunc(m.feeds, func(f *feed) bool {
		return f.source.id == sourceID
	})
	m.sources = slices.DeleteFunc(m.sources, func(s *source) bool {
		if s.id == sourceID {
			s.feeds = nil
			return true
		}
		return false
	})
	return nil
}

func (m *Manager) removeFeed(feedID int64) error {
	m.feeds = slices.DeleteFunc(m.feeds, func(f *feed) bool {
		if f.id == feedID {
			f.source.feeds = slices.DeleteFunc(f.source.feeds, func(f *feed) bool {
				return f.id == feedID
			})
			return true
		}
		return false
	})
	return nil
}

func (m *Manager) addSource(sourceID int64) error {
	// Ignore it if we already have it.
	if slices.ContainsFunc(m.sources, func(s *source) bool { return s.id == sourceID }) {
		return nil
	}
	const (
		sourceSQL = `SELECT rate, slots, active FROM sources WHERE id = $1`
		feedsSQL  = `SELECT id, url, rolie, log_lvl::text FROM feeds WHERE source_id = $1`
	)
	var (
		s  *source
		fs []*feed
	)
	if err := m.db.Run(
		context.Background(),
		func(ctx context.Context, con *pgxpool.Conn) error {
			tx, err := con.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})
			if err != nil {
				return fmt.Errorf("starting transaction failed: %w", err)
			}
			defer tx.Rollback(ctx)
			// Collect source.
			srows, err := tx.Query(ctx, sourceSQL, sourceID)
			if err != nil {
				return fmt.Errorf("querying source failed: %w", err)
			}
			s, err = pgx.CollectOneRow(srows, func(row pgx.CollectableRow) (*source, error) {
				var s source
				return &s, row.Scan(&s.rate, &s.slots, &s.active)
			})
			if err != nil {
				return err
			}
			// Collect feeds.
			frows, err := tx.Query(ctx, feedsSQL, sourceID)
			if err != nil {
				return fmt.Errorf("querying feeds failed: %w", err)
			}
			fs, err = pgx.CollectRows(frows, func(row pgx.CollectableRow) (*feed, error) {
				var (
					f   feed
					raw string
				)
				if err := row.Scan(&f.id, &raw, &f.rolie, &f.logLevel); err != nil {
					return nil, err
				}
				parsed, err := url.Parse(raw)
				if err != nil {
					return nil, fmt.Errorf("invalid URL: %w", err)
				}
				f.url = parsed
				f.source = s
				return &f, nil
			})
			if err != nil {
				return fmt.Errorf("collecting feeds failed: %w", err)
			}
			s.feeds = fs
			return tx.Commit(ctx)
		}, 0,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("fetching source failed: %w", err)
	}
	m.sources = append(m.sources, s)
	m.feeds = append(m.feeds, fs...)

	if s.active && len(fs) > 0 {
		go func() { m.fns <- (*Manager).ping }()
	}
	return nil
}

func (m *Manager) addFeed(feedID int64) error {
	// Ignore it if we already have it.
	if slices.ContainsFunc(m.feeds, func(f *feed) bool { return f.id == feedID }) {
		return nil
	}
	// TODO: Implement me!
	return nil
}

func (m *Manager) asManager(fn func(*Manager, int64) error, id int64) error {
	err := make(chan error)
	m.fns <- func(m *Manager) { err <- fn(m, id) }
	return <-err
}

// AddSource registers a new source.
func (m *Manager) AddSource(sourceID int64) error {
	return m.asManager((*Manager).addSource, sourceID)
}

// RemoveSource removes a sources from manager.
func (m *Manager) RemoveSource(sourceID int64) error {
	return m.asManager((*Manager).removeSource, sourceID)
}

// AddFeed adds a new feed to a source.
func (m *Manager) AddFeed(feedID int64) error {
	return m.asManager((*Manager).addFeed, feedID)
}

// RemoveFeed removes a feed from a source.
func (m *Manager) RemoveFeed(feedID int64) error {
	return m.asManager((*Manager).removeFeed, feedID)
}

func (m *Manager) downloadDone(f *feed, id int64) func() {
	return func() {
		m.fns <- func(m *Manager) {
			f.source.usedSlots = max(0, f.source.usedSlots-1)
			m.usedSlots = max(0, m.usedSlots-1)
			if l := f.findLocationByID(id); l != nil {
				l.state = done
			}
		}
	}
}
