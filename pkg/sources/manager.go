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
	cfg  *config.Config
	db   *database.DB
	fns  chan func(*Manager)
	done bool

	sources   []*source
	usedSlots int
	uniqueID  int64
}

// NewManager creates a new downloader.
func NewManager(cfg *config.Config, db *database.DB) *Manager {
	return &Manager{
		cfg: cfg,
		db:  db,
		fns: make(chan func(*Manager)),
	}
}

func (m *Manager) numActiveFeeds() int {
	sum := 0
	for _, s := range m.sources {
		if s.active {
			sum += len(s.feeds)
		}
	}
	return sum
}

func (m *Manager) activeFeeds(fn func(*feed) bool) {
	for _, s := range m.sources {
		if s.active {
			for _, f := range s.feeds {
				if !fn(f) {
					return
				}
			}
		}
	}
}

func (m *Manager) allFeeds(fn func(*feed) bool) {
	for _, s := range m.sources {
		for _, f := range s.feeds {
			if !fn(f) {
				return
			}
		}
	}
}

// refreshFeeds checks if there are feeds that need reloading
// and does so in that case.
func (m *Manager) refreshFeeds() {
	now := time.Now()
	m.activeFeeds(func(f *feed) bool {
		// Does the feed need a refresh?
		if f.nextCheck.IsZero() || !now.Before(f.nextCheck) {
			if err := f.refresh(m); err != nil {
				f.log(m, config.ErrorFeedLogLevel, "feed refresh failed: %v", err.Error())
			}
			// Even if there was an error try again later.
			f.nextCheck = time.Now().Add(m.cfg.Sources.FeedRefresh)
		}
		return true
	})
}

// startDownloads starts downloads if there are enough slots and
// there are things to download.
func (m *Manager) startDownloads() {
	for m.usedSlots < m.cfg.Sources.DownloadSlots {
		started := false
		m.activeFeeds(func(f *feed) bool {
			// Has this feed a free slot?
			maxSlots := min(m.cfg.Sources.MaxSlotsPerSource, m.cfg.Sources.DownloadSlots)
			if f.source.slots != nil {
				maxSlots = min(maxSlots, *f.source.slots)
			}
			if f.source.usedSlots >= maxSlots {
				return true
			}
			// Find a candidate to download.
			location := f.findWaiting()
			if location == nil {
				return true
			}
			m.usedSlots++
			f.source.usedSlots++
			location.state = running
			location.id = m.generateID()
			started = true
			// Calling reciever by value is intended here!
			go (*location).download(m, f, m.downloadDone(f, location.id))
			return m.usedSlots < m.cfg.Sources.DownloadSlots
		})
		if !started {
			return
		}
	}
}

// compactDone removes the locations the feeds which are downloaded.
func (m *Manager) compactDone() {
	m.allFeeds(func(f *feed) bool {
		f.locations = slices.DeleteFunc(f.locations, func(l location) bool {
			return l.state == done
		})
		return true
	})
}

func (m *Manager) generateID() int64 {
	// Start with 1 to avoid clashes with zeroed locations.
	m.uniqueID++
	return m.uniqueID
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

func (m *Manager) backgroundPing() {
	go func() { m.fns <- (*Manager).ping }()
}

// Kill stops the manager.
func (m *Manager) Kill() {
	m.fns <- func(m *Manager) { m.done = true }
}

func (m *Manager) removeSource(sourceID int64) error {
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
	for _, s := range m.sources {
		before := len(s.feeds)
		s.feeds = slices.DeleteFunc(s.feeds, func(f *feed) bool {
			return f.id == feedID
		})
		if before > len(s.feeds) {
			return nil
		}
	}
	return nil
}

func (m *Manager) addSource(sourceID int64) error {
	// Ignore it if we already have it.
	if slices.ContainsFunc(m.sources, func(s *source) bool { return s.id == sourceID }) {
		return nil
	}
	const (
		sourceSQL = `SELECT rate, slots, active FROM sources WHERE id = $1`
		feedsSQL  = `SELECT id, url, rolie, log_lvl::text FROM feeds WHERE sources_id = $1`
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
			return tx.Commit(ctx)
		}, 0,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("fetching source failed: %w", err)
	}
	s.feeds = fs
	m.sources = append(m.sources, s)

	if s.active && len(fs) > 0 {
		m.backgroundPing()
	}
	return nil
}

func (m *Manager) addFeed(feedID int64) error {
	// Ignore it if we already have it.
	for _, s := range m.sources {
		if slices.ContainsFunc(s.feeds, func(f *feed) bool { return f.id == feedID }) {
			return nil
		}
	}
	const feedSQL = `SELECT sources_id, url, rolie, log_lvl::text FROM feeds WHERE sources_id = $1`
	var (
		f *feed
		s *source
	)
	if err := m.db.Run(
		context.Background(),
		func(ctx context.Context, con *pgxpool.Conn) error {
			tx, err := con.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})
			if err != nil {
				return fmt.Errorf("starting transaction failed: %w", err)
			}
			defer tx.Rollback(ctx)
			// Collect feed.
			frows, err := tx.Query(ctx, feedSQL, feedID)
			if err != nil {
				return fmt.Errorf("querying feed failed: %w", err)
			}
			var sid int64
			f, err = pgx.CollectOneRow(frows, func(row pgx.CollectableRow) (*feed, error) {
				var (
					f   feed
					raw string
				)
				if err := row.Scan(&sid, &raw, &f.rolie, &f.logLevel); err != nil {
					return nil, err
				}
				parsed, err := url.Parse(raw)
				if err != nil {
					return nil, err
				}
				f.url = parsed
				return &f, nil
			})
			if err != nil {
				return err
			}
			// Do we have the source already?
			idx := slices.IndexFunc(m.sources, func(s *source) bool { return s.id == sid })
			if idx == -1 {
				// XXX: Maybe we should load the source and all the other feeds of this source?
				return errors.New("source is missing")
			}
			s = m.sources[idx]
			return tx.Commit(ctx)
		}, 0,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("fetching feed failed: %w", err)
	}
	f.source = s
	s.feeds = append(s.feeds, f)
	if s.active {
		m.backgroundPing()
	}
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

// downloadDone returns a function which does the needed
// book keeping when a download is finished. To be used
// as a defer function in the download.
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
