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
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ISDuBA/ISDuBA/internal/config"
	"github.com/ISDuBA/ISDuBA/internal/database"
	"github.com/ISDuBA/ISDuBA/internal/version"
	"github.com/gocsaf/csaf/v3/util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/time/rate"
)

const deactivatedDueToClientCertIssue = `Deactivated due to client cert issue.`

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
	logLevel atomic.Int32

	invalid atomic.Bool

	nextCheck time.Time
	queue     []location
	source    *source

	refreshBlocked bool
	lastETag       string
	lastModified   time.Time
}

type ignorePatterns []*regexp.Regexp

type source struct {
	id        int64
	name      string
	url       string
	active    bool
	feeds     []*feed
	usedSlots int
	status    []string

	rate           *float64
	limiter        *rate.Limiter
	slots          *int
	headers        []string
	strictMode     *bool
	secure         *bool
	signatureCheck *bool
	age            *time.Duration
	ignorePatterns ignorePatterns

	clientCertPublic     []byte
	clientCertPrivate    []byte
	clientCertPassphrase []byte
	tlsCertificates      []tls.Certificate

	checksum        []byte
	checksumAck     time.Time
	checksumUpdated time.Time
}

// ignore returns true if the given url should be ignored.
func (ip ignorePatterns) ignore(u *url.URL) bool {
	if len(ip) == 0 {
		return false
	}
	p := u.String()
	for _, pattern := range ip {
		if pattern.MatchString(p) {
			return true
		}
	}
	return false
}

// refresh fetches the feed index and accordingly updates
// the list of locations if needed.
func (f *feed) refresh(m *Manager) {
	f.log(m, config.InfoFeedLogLevel, "refreshing feed")

	// Fetching the index is too expensive for the manager main loop.
	// So we do it async and call back when its is done.
	f.fetchIndex(m, func(candidates []location, err error) {
		if err != nil {
			f.log(m, config.ErrorFeedLogLevel, "fetching feed index failed: %v", err)
			return
		}
		if candidates == nil {
			slog.Debug("feed has not changed", "feed", f.id)
			f.log(m, config.InfoFeedLogLevel, "feed %d has not changed", f.id)
			f.log(m, config.InfoFeedLogLevel, "entries to download: %d", len(f.queue))
			return
		}

		slog.Debug("feed has new candidates", "feed", f.id, "candidates", len(candidates))

		// The manager is the owner of the feed so let it do the changes.
		m.fns <- func(m *Manager, ctx context.Context) {
			// Filter out candidates which are already in the database with same or newer.
			if candidates, err = f.removeOlder(ctx, m.db, candidates); err != nil {
				f.log(m, config.ErrorFeedLogLevel,
					"feed refresh failed with database error: %v", err)
				return
			}

			if len(candidates) == 0 { // Nothing to do.
				slog.Debug("feed has no candidates left", "feed", f.id)
				return
			}

			// Candidates may pile up on same urls so only keep
			// the latest ones.
			f.removeOutdatedWaiting(candidates)

			// Merge candidates into list of locations.
			f.queue = append(f.queue, candidates...)
			slices.SortFunc(f.queue, func(a, b location) int {
				return a.updated.Compare(b.updated)
			})

			slog.Debug("feed entries to download", "feed", f.id, "queue", len(f.queue))
			f.log(m, config.InfoFeedLogLevel, "entries to download: %d", len(f.queue))
		}
	})
}

// removeOutdatedWaiting removes locations with urls from queue which
// have newer update candidates.
func (f *feed) removeOutdatedWaiting(candidates []location) {
	if len(f.queue) == 0 {
		return
	}
	urls := make(map[string]time.Time, len(candidates))
	for i := range candidates {
		cand := &candidates[i]
		urls[cand.doc.String()] = cand.updated
	}
	f.queue = slices.DeleteFunc(f.queue, func(l location) bool {
		if l.state == waiting {
			updated, ok := urls[l.doc.String()]
			return ok && updated.After(l.updated)
		}
		return false
	})
}

// resetIndexTags resets the tags used to signal
// that we know the feed index.
func (f *feed) resetIndexTags() {
	f.lastETag = ""
	f.lastModified = time.Time{}
}

// fetchIndex fetches the content of the feed index.
func (f *feed) fetchIndex(m *Manager, fn func([]location, error)) {
	// Prevent stacked calling
	f.refreshBlocked = true

	indexURL := f.url.String()
	if !f.rolie {
		var err error
		if indexURL, err = url.JoinPath(indexURL, "changes.csv"); err != nil {
			fn(nil, err)
			return
		}
	}
	slog.Debug("fetching index", "url", indexURL, "rolie", f.rolie)
	req, err := http.NewRequest(http.MethodGet, indexURL, nil)
	if err != nil {
		fn(nil, err)
		return
	}
	if f.lastETag != "" {
		req.Header.Add("If-None-Match", f.lastETag)
	}
	if !f.lastModified.IsZero() {
		req.Header.Add("If-Modified-Since", f.lastModified.Format(http.TimeFormat))
	}
	client := f.source.httpClient(m)
	// Copy relevant data to avoid races.
	fi := feedIndex{
		base:           f.url,
		age:            f.source.age,
		ignorePatterns: f.source.ignorePatterns,
		sameOrNewer:    f.sameOrNewer(),
	}
	// Do the actual fetching async.
	go func() {
		defer func() {
			client.CloseIdleConnections()
			// Re-enable refreshing
			m.fns <- func(*Manager, context.Context) { f.refreshBlocked = false }
		}()
		resp, err := f.source.doRequest(client, m, req)
		if err != nil {
			fn(nil, err)
			return
		}
		defer resp.Body.Close()
		// Nothing changed since last call.
		if resp.StatusCode == http.StatusNotModified {
			fn(nil, nil)
			return
		}
		if resp.StatusCode != http.StatusOK {
			fn(nil, fmt.Errorf("status code %d", resp.StatusCode))
			return
		}
		var locations []location
		if f.rolie {
			locations, err = fi.rolieLocations(resp.Body)
		} else {
			locations, err = fi.directoryLocations(resp.Body)
		}
		if err != nil {
			fn(nil, err)
			return
		}
		fn(locations, nil)
		m.fns <- func(*Manager, context.Context) {
			f.lastETag = resp.Header.Get("Etag")
			if m := resp.Header.Get("Last-Modified"); m != "" {
				f.lastModified, _ = time.Parse(http.TimeFormat, m)
			}
		}
	}()
}

// removeOlder takes a list of locations and removes the items which are already
// in the database with a same or newer update time.
func (f *feed) removeOlder(
	ctx context.Context, db *database.DB,
	candidates []location,
) ([]location, error) {
	var remove [][2]int

	exists := func(idx int) func(pgx.Row) error {
		return func(row pgx.Row) error {
			var have bool
			if err := row.Scan(&have); err != nil {
				return fmt.Errorf("looking for same or newer in db failed: %w", err)
			}
			if have {
				if n := len(remove); n > 0 && remove[n-1][1] == idx-1 {
					remove[n-1][1] = idx
				} else {
					remove = append(remove, [2]int{idx, idx})
				}
			}
			return nil
		}
	}

	const sql = `SELECT EXISTS(SELECT 1 FROM changes ` +
		`WHERE url = $1 AND feeds_id = $2 AND time >= $3)`

	batch := pgx.Batch{}

	for i := range candidates {
		cand := &candidates[i]
		batch.Queue(sql, cand.doc.String(), f.id, cand.updated).QueryRow(exists(i))
	}

	if err := db.Run(
		ctx,
		func(ctx context.Context, conn *pgxpool.Conn) error {
			return conn.SendBatch(ctx, &batch).Close()
		}, 0,
	); err != nil {
		return nil, fmt.Errorf("sending same or newer batch failed: %w", err)
	}

	for i := len(remove) - 1; i >= 0; i-- {
		candidates = slices.Delete(candidates, remove[i][0], remove[i][1]+1)
	}

	return candidates, nil
}

// sameOrNewer returns a function which checks if a given location is
// already in the current list of locations with the same or newer update time.
func (f *feed) sameOrNewer() func(*location) bool {
	// XXX: Maybe this extra indexing could be replaced by something
	// which uses the fact that the locations are already sorted by updated?!
	have := make(map[string]time.Time, len(f.queue))
	for _, location := range f.queue {
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
	for i := len(f.queue) - 1; i >= 0; i-- {
		if location := &f.queue[i]; location.id == id {
			return location
		}
	}
	return nil
}

// findWaiting looks for a location ready to download.
func (f *feed) findWaiting() *location {
	// Backwards because the new ones are at the end.
	for i := len(f.queue) - 1; i >= 0; i-- {
		if location := &f.queue[i]; location.state == waiting {
			return location
		}
	}
	return nil
}

func (f *feed) addStats(st *Stats) {
	for i := range f.queue {
		switch f.queue[i].state {
		case waiting:
			st.Waiting++
		case running:
			st.Downloading++
		}
	}
}

func (s *source) addStats(st *Stats) {
	for _, f := range s.feeds {
		if !f.invalid.Load() {
			f.addStats(st)
		}
	}
}

// forceIndexRefresh forces an index refresh on all feeds of a source.
func (s *source) forceIndexRefresh() {
	past := time.Now().Add(-time.Minute)
	for _, f := range s.feeds {
		if !f.invalid.Load() {
			f.nextCheck = past
			f.resetIndexTags()
		}
	}
}

// deleteTooOld removes locations from the feeds of the source
// which are before the accepted age.
func (s *source) deleteTooOld() {
	if s.age == nil {
		return
	}
	cut := time.Now().Add(-*s.age)
	for _, f := range s.feeds {
		if f.invalid.Load() {
			continue
		}
		f.queue = slices.DeleteFunc(f.queue, func(l location) bool {
			return l.state == waiting && l.updated.Before(cut)
		})
	}
}

func (s *source) setAge(age *time.Duration) {
	s.age = age
	s.deleteTooOld()
	s.forceIndexRefresh()
}

// deleteIgnore remove the location from the feeds of this source
// which should be ignored.
func (s *source) deleteIgnore() {
	if len(s.ignorePatterns) == 0 {
		return
	}
	for _, f := range s.feeds {
		if f.invalid.Load() {
			continue
		}
		f.queue = slices.DeleteFunc(f.queue, func(l location) bool {
			return l.state == waiting && s.ignorePatterns.ignore(l.doc)
		})
	}
}

func (s *source) setIgnorePatterns(ignorePatterns []*regexp.Regexp) {
	s.ignorePatterns = ignorePatterns
	s.deleteIgnore()
	s.forceIndexRefresh()
}

func (s *source) setRate(rate *float64) {
	s.rate = rate
	s.limiter = nil
}

// wait establishes the request rate per source.
func (s *source) wait() *rate.Limiter {
	if s.rate == nil {
		return nil
	}
	if s.limiter == nil {
		s.limiter = rate.NewLimiter(rate.Limit(*s.rate), 1)
	}
	return s.limiter
}

func (s *source) httpClient(m *Manager) *http.Client {
	var tlsConfig tls.Config

	if s.secure != nil {
		tlsConfig.InsecureSkipVerify = !*s.secure
	} else {
		tlsConfig.InsecureSkipVerify = !m.cfg.Sources.Secure
	}

	if len(s.tlsCertificates) > 0 {
		tlsConfig.Certificates = s.tlsCertificates
	}

	transport := m.cfg.General.Transport()
	transport.TLSClientConfig = &tlsConfig

	client := http.Client{Transport: transport}
	if m.cfg.Sources.Timeout > 0 {
		client.Timeout = m.cfg.Sources.Timeout
	}
	return &client
}

func (s *source) applyHeaders(req *http.Request) {
	for _, header := range s.headers {
		if k, v, ok := strings.Cut(header, ":"); ok {
			req.Header.Add(k, v)
		}
	}
	if req.Header.Get("User-Agent") == "" {
		req.Header.Add("User-Agent", UserAgent)
	}
}

// doRequest executes an HTTP request with the source specific parameters.
func (s *source) doRequest(client *http.Client, m *Manager, req *http.Request) (*http.Response, error) {
	// The manager owns the configuration.
	// So we let the manager do the adjustment of the request.

	var limiter *rate.Limiter

	m.inManager(func(m *Manager, _ context.Context) {
		s.applyHeaders(req)
		if client == nil {
			client = s.httpClient(m)
		}
		limiter = s.wait()
	})

	if limiter != nil {
		limiter.Wait(context.Background())
	}
	return client.Do(req)
}

func (s *source) httpGet(client *http.Client, m *Manager, url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return s.doRequest(client, m, req)
}

// loadHash fetches text form of a hash from remote location.
func (s *source) loadHash(client *http.Client, m *Manager, url string) ([]byte, error) {
	resp, err := s.httpGet(client, m, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s (%d)",
			http.StatusText(resp.StatusCode), resp.StatusCode)
	}
	return util.HashFromReader(resp.Body)
}

// checkSignature tells if the signature check should be taken seriously.
func (s *source) checkSignature(m *Manager) bool {
	if s.signatureCheck != nil {
		return *s.signatureCheck
	}
	return m.cfg.Sources.SignatureCheck
}

// useStrictMode tells if the check results should be taken seriously.
func (s *source) useStrictMode(m *Manager) bool {
	if s.strictMode != nil {
		return *s.strictMode
	}
	return m.cfg.Sources.StrictMode
}

// storeLastChanges is intended to be called in the transaction storing the
// imported document after was successful. It helps to remember the
// last changes per location so we don't need to download them all again and again.
func (f *feed) storeLastChanges(l *location) func(context.Context, pgx.Tx, int64, bool) error {
	return func(ctx context.Context, tx pgx.Tx, _ int64, _ bool) error {
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
