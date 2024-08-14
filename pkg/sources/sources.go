// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

// Package sources implements the download from sources.
package sources

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ISDuBA/ISDuBA/pkg/config"
	"github.com/ISDuBA/ISDuBA/pkg/database"
	"github.com/ISDuBA/ISDuBA/pkg/version"
	"github.com/csaf-poc/csaf_distribution/v3/util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/time/rate"
)

// UserAgent is the name of the http client
var UserAgent = "isduba/" + version.SemVersion

type state int

const (
	waiting state = iota
	running
	done
)

type location struct {
	updated   time.Time
	doc       *url.URL
	hash      *url.URL
	signature *url.URL
	state     state
	id        int64
}

type feed struct {
	id       int64
	label    string
	url      *url.URL
	rolie    bool
	logLevel config.FeedLogLevel

	invalid atomic.Bool

	nextCheck time.Time
	locations []location
	source    *source

	lastETag     string
	lastModified time.Time
}

type source struct {
	id     int64
	name   string
	url    string
	rate   *float64
	slots  *int
	active bool

	feeds     []*feed
	usedSlots int
	limiterMu sync.Mutex
	limiter   *rate.Limiter
	headers   []string
}

// refresh fetches the feed index and accordingly updates
// the list of locations if needed.
func (f *feed) refresh(m *Manager) error {

	f.log(m, config.InfoFeedLogLevel, "refreshing feed")

	candidates, err := f.fetchIndex()
	if err != nil {
		return fmt.Errorf("fetching feed index failed: %w", err)
	}
	if candidates == nil {
		f.log(m, config.InfoFeedLogLevel, "feed %d has not changed", f.id)
		f.log(m, config.InfoFeedLogLevel, "entries to download: %d", len(f.locations))
		return nil
	}

	// Filter out candidates which are already in the database with same or newer.
	if candidates, err = f.removeOlder(m.db, candidates); err != nil {
		return fmt.Errorf("removing candidates by looking at database failed: %w", err)
	}

	if len(candidates) == 0 { // Nothing to do.
		return nil
	}

	// Candidates may pile up on same urls so only keep
	// the latest ones.
	f.removeOutdatedWaiting(candidates)

	// Merge candidates into list of locations.
	f.locations = append(f.locations, candidates...)
	slices.SortFunc(f.locations, func(a, b location) int {
		return a.updated.Compare(b.updated)
	})

	f.log(m, config.InfoFeedLogLevel, "entries to download: %d", len(f.locations))
	return nil
}

// removeOutdatedWaiting removes locations with urls from queue which
// have newer update candidates.
func (f *feed) removeOutdatedWaiting(candidates []location) {
	if len(f.locations) == 0 {
		return
	}
	urls := make(map[string]time.Time, len(candidates))
	for i := range candidates {
		cand := &candidates[i]
		urls[cand.doc.String()] = cand.updated
	}
	f.locations = slices.DeleteFunc(f.locations, func(l location) bool {
		if l.state == waiting {
			updated, ok := urls[l.doc.String()]
			return ok && updated.After(l.updated)
		}
		return false
	})
}

// fetchIndex fetches the content of the feed index.
func (f *feed) fetchIndex() ([]location, error) {
	req, err := http.NewRequest(http.MethodGet, f.url.String(), nil)
	if err != nil {
		return nil, err
	}
	if f.lastETag != "" {
		req.Header.Add("If-None-Match", f.lastETag)
	}
	if !f.lastModified.IsZero() {
		req.Header.Add("If-Modified-Since", f.lastModified.Format(http.TimeFormat))
	}
	resp, err := f.source.doRequest(req)
	if err != nil {
		return nil, err
	}
	// Nothing changed since last call.
	if resp.StatusCode == http.StatusNotModified {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	var locations []location
	if f.rolie {
		locations, err = f.rolieLocations(resp.Body)
	} else {
		locations, err = f.directoryLocations(resp.Body)
	}
	if err != nil {
		return nil, err
	}
	f.lastETag = resp.Header.Get("Etag")
	if m := resp.Header.Get("Last-Modified"); m != "" {
		f.lastModified, _ = time.Parse(http.TimeFormat, m)
	}
	return locations, nil
}

// removeOlder takes a list of locations and removes the items which are already
// in the database with a same or newer update time.
func (f *feed) removeOlder(db *database.DB, candidates []location) ([]location, error) {

	var remove []int

	batch := pgx.Batch{}

	const sql = `SELECT EXISTS(SELECT 1 FROM changes ` +
		`WHERE url = $1 AND feeds_id = $2 AND time >= $3)`

	exists := func(idx int) func(pgx.Row) error {
		return func(row pgx.Row) error {
			var have bool
			if err := row.Scan(&have); err != nil {
				return fmt.Errorf("looking for same or newer in db failed: %w", err)
			}
			if have {
				remove = append(remove, idx)
			}
			return nil
		}
	}

	for i := range candidates {
		cand := &candidates[i]
		batch.Queue(sql, cand.doc.String(), f.id, cand.updated).QueryRow(exists(i))
	}

	if err := db.Run(
		context.Background(),
		func(ctx context.Context, conn *pgxpool.Conn) error {
			return conn.SendBatch(ctx, &batch).Close()
		}, 0,
	); err != nil {
		return nil, fmt.Errorf("sending same or newer batch failed: %w", err)
	}

	// XXX: This could be optimized by passing ranges to Delete.
	for i := len(remove) - 1; i >= 0; i-- {
		candidates = slices.Delete(candidates, remove[i], remove[i])
	}

	return candidates, nil
}

// sameOrNewer returns a function which checks if a given location is
// already in the current list of locations with the same or newer update time.
func (f *feed) sameOrNewer() func(*location) bool {
	// XXX: Maybe this extra indexing could be replaced by something
	// which uses the fact that the locations are already sorted by updated?!
	have := make(map[string]time.Time, len(f.locations))
	for _, location := range f.locations {
		url := location.doc.String()
		if t, ok := have[url]; !ok || location.updated.After(t) {
			have[url] = location.updated
		}
	}
	return func(location *location) bool {
		t, ok := have[location.doc.String()]
		return ok && !t.Before(location.updated)
	}
}

// findLocationByID looks for location with a given id.
func (f *feed) findLocationByID(id int64) *location {
	for i := len(f.locations) - 1; i >= 0; i-- {
		if location := &f.locations[i]; location.id == id {
			return location
		}
	}
	return nil
}

// findWaiting looks for a location ready to download.
func (f *feed) findWaiting() *location {
	// Backwards because the new ones are at the end.
	for i := len(f.locations) - 1; i >= 0; i-- {
		if location := &f.locations[i]; location.state == waiting {
			return location
		}
	}
	return nil
}

func (s *source) setRate(rate *float64) {
	s.limiterMu.Lock()
	s.rate = rate
	s.limiter = nil
	s.limiterMu.Unlock()
}

// wait establishes the request rate per source.
func (s *source) wait(ctx context.Context) {
	s.limiterMu.Lock()
	defer s.limiterMu.Unlock()
	if s.rate != nil {
		if s.limiter == nil {
			s.limiter = rate.NewLimiter(rate.Limit(*s.rate), 1)
		}
		s.limiter.Wait(ctx)
	}
}

func (s *source) httpClient() *http.Client {
	client := http.Client{}
	// TODO: Implement me!
	return &client
}

// doRequest executes an HTTP request with the source specific parameters.
func (s *source) doRequest(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", UserAgent)

	// TODO: Implement me!
	client := s.httpClient()
	s.wait(context.Background())
	return client.Do(req)
}

func (s *source) httpGet(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return s.doRequest(req)
}

// loadHash fetches text form of a hash from remote location.
func (s *source) loadHash(url string) ([]byte, error) {
	resp, err := s.httpGet(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s (%d)", http.StatusText(resp.StatusCode), resp.StatusCode)
	}
	return util.HashFromReader(resp.Body)
}

// storeLastChanges is intented to be called in the transaction storing the
// importing the document after is was successful. It helps to remember the
// last changes per location so we don't need to download them all again and again.
func (f *feed) storeLastChanges(l *location) func(context.Context, pgx.Tx, int64) error {
	return func(ctx context.Context, tx pgx.Tx, _ int64) error {
		if f.invalid.Load() {
			return nil
		}
		const updatedSQL = `INSERT INTO changes (url, feeds_id, time) ` +
			`VALUES ($1, $2, $3) ` +
			`ON CONFLICT (url, feeds_id) DO ` +
			`UPDATE SET time = $3`
		_, err := tx.Exec(ctx, updatedSQL, l.doc.String(), f.id, l.updated)
		return err
	}
}
