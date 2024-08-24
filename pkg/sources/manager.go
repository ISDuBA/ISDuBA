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
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/ISDuBA/ISDuBA/pkg/config"
	"github.com/ISDuBA/ISDuBA/pkg/database"
	"github.com/csaf-poc/csaf_distribution/v3/csaf"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	// NoSuchEntryError is returned if a given feed or source does not exists.
	NoSuchEntryError string
	// InvalidArgumentError is returned if a given argument is unsuited.
	InvalidArgumentError string
)

// Error implements [builtin.error].
func (nsee NoSuchEntryError) Error() string { return string(nsee) }

// Error implements [builtin.error].
func (iae InvalidArgumentError) Error() string { return string(iae) }

// Is supports [errors.Is].
func (NoSuchEntryError) Is(target error) bool {
	_, ok := target.(NoSuchEntryError)
	return ok
}

// Is supports [errors.Is].
func (InvalidArgumentError) Is(target error) bool {
	_, ok := target.(InvalidArgumentError)
	return ok
}

// refreshDuration is the fallback duration for feeds to be checked for refresh.
const refreshDuration = time.Minute

// Manager fetches advisories from sources.
type Manager struct {
	cfg  *config.Config
	db   *database.DB
	fns  chan func(*Manager)
	done bool

	cipherKey []byte

	sources []*source

	pmdCache  *pmdCache
	keysCache *keysCache

	val csaf.RemoteValidator

	usedSlots int
	uniqueID  int64
}

// SourceUpdateResult is return by UpdateSource.
type SourceUpdateResult int

const (
	// SourceUnchanged is returned if there was no change in the source.
	SourceUnchanged SourceUpdateResult = iota
	// SourceUpdated is returned if the source was updated.
	SourceUpdated
	// SourceDeactivated is returned if the source was deactivated during the update.
	SourceDeactivated
)

func (sur SourceUpdateResult) String() string {
	switch sur {
	case SourceUnchanged:
		return "unchanged"
	case SourceUpdated:
		return "updated"
	case SourceDeactivated:
		return "deactivated"
	default:
		return fmt.Sprintf("unknown update result: %d", sur)
	}
}

// NewManager creates a new downloader.
func NewManager(
	cfg *config.Config,
	db *database.DB,
	val csaf.RemoteValidator,
) (*Manager, error) {
	cipherKey, err := createCipherKey(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating cipher failed: %w", err)
	}
	return &Manager{
		cfg:       cfg,
		db:        db,
		fns:       make(chan func(*Manager)),
		cipherKey: cipherKey,
		pmdCache:  newPMDCache(),
		keysCache: newKeysCache(cfg.Sources.OpenPGPCaching),
		val:       val,
	}, nil
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

func (m *Manager) findFeedByID(feedID int64) *feed {
	for _, s := range m.sources {
		if idx := slices.IndexFunc(s.feeds, func(f *feed) bool { return f.id == feedID }); idx >= 0 {
			return s.feeds[idx]
		}
	}
	return nil
}

func (m *Manager) findSourceByID(sourceID int64) *source {
	if idx := slices.IndexFunc(m.sources, func(s *source) bool { return s.id == sourceID }); idx >= 0 {
		return m.sources[idx]
	}
	return nil
}

func (m *Manager) findSourceByName(name string) *source {
	if idx := slices.IndexFunc(m.sources, func(s *source) bool { return s.name == name }); idx >= 0 {
		return m.sources[idx]
	}
	return nil
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
			// Funtion to be called when the download is finished.
			done := m.downloadDone(f, location.id)
			// Calling reciever by value is intended here!
			go (*location).download(m, f, done)
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
		f.queue = slices.DeleteFunc(f.queue, func(l location) bool {
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
		m.pmdCache.Cleanup()
		m.keysCache.Cleanup()
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

// AllSources iterates over all sources.
func (m *Manager) AllSources(fn func(
	id int64,
	name string,
	url string,
	active bool,
	rate *float64,
	slots *int,
	headers []string,
	strictmode *bool,
	insecure *bool,
	signatureCheck *bool,
	age *time.Duration,
	ignorepatterns []*regexp.Regexp,
	hasClientCertPublic bool,
	hasClientCertPrivate bool,
	hasClientCertPassphrase bool,
)) {
	done := make(chan struct{})
	m.fns <- func(m *Manager) {
		defer close(done)
		for _, s := range m.sources {
			fn(s.id, s.name, s.url, s.active, s.rate, s.slots, s.headers,
				s.strictMode, s.insecure, s.signatureCheck, s.age, s.ignorePatterns,
				s.clientCertPublic != nil,
				s.clientCertPrivate != nil,
				s.clientCertPassphrase != nil,
			)
		}
	}
	<-done
}

// Feeds passes the fields of the feeds of a given source to a given function.
func (m *Manager) Feeds(sourceID int64, fn func(
	id int64,
	label string,
	url *url.URL,
	rolie bool,
	lvl config.FeedLogLevel,
)) error {
	errCh := make(chan error)
	m.fns <- func(m *Manager) {
		s := m.findSourceByID(sourceID)
		if s == nil {
			errCh <- NoSuchEntryError("no such source")
			return
		}
		for _, f := range s.feeds {
			fn(f.id, f.label, f.url, f.rolie, config.FeedLogLevel(f.logLevel.Load()))
		}
		errCh <- nil
	}
	return <-errCh
}

// Feed passes the fields of feed to a given function.
func (m *Manager) Feed(feedID int64, fn func(
	label string,
	url *url.URL,
	rolie bool,
	lvl config.FeedLogLevel,
)) error {
	errCh := make(chan error)
	m.fns <- func(m *Manager) {
		f := m.findFeedByID(feedID)
		if f == nil {
			errCh <- NoSuchEntryError("no such feed")
			return
		}
		fn(f.label, f.url, f.rolie, config.FeedLogLevel(f.logLevel.Load()))
		errCh <- nil
	}
	return <-errCh
}

// FeedLog sends the log of the feed with the given id to the given function.
func (m *Manager) FeedLog(feedID int64, fn func(
	t time.Time,
	lvl config.FeedLogLevel,
	msg string,
)) error {
	errCh := make(chan error)
	m.fns <- func(m *Manager) {
		const sql = `SELECT time, lvl::text, msg FROM feed_logs WHERE feeds_id = $1 ` +
			`ORDER by time DESC`
		errCh <- m.db.Run(
			context.Background(),
			func(ctx context.Context, con *pgxpool.Conn) error {
				rows, err := con.Query(ctx, sql, feedID)
				if err != nil {
					return fmt.Errorf("querying feed logs failed: %w", err)
				}
				defer rows.Close()
				var (
					t   time.Time
					lvl config.FeedLogLevel
					msg string
				)
				for rows.Next() {
					if err := rows.Scan(&t, &lvl, &msg); err != nil {
						return fmt.Errorf("scanning log failed: %w", err)
					}
					fn(t, lvl, msg)
				}
				return rows.Err()
			}, 0,
		)
	}
	return <-errCh
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
	if m.findSourceByID(sourceID) == nil {
		return NoSuchEntryError("no such source")
	}
	const sql = `DELETE FROM sources WHERE id = $1`
	notFound := false
	if err := m.db.Run(
		context.Background(),
		func(rctx context.Context, con *pgxpool.Conn) error {
			tags, err := con.Exec(rctx, sql, sourceID)
			if err != nil {
				return fmt.Errorf("removing source failed: %w", err)
			}
			notFound = tags.RowsAffected() == 0
			return nil
		}, 0,
	); err != nil {
		return fmt.Errorf("deleting source from db failed: %w", err)
	}
	m.sources = slices.DeleteFunc(m.sources, func(s *source) bool {
		if s.id == sourceID {
			s.active = false
			s.feeds = nil
			return true
		}
		return false
	})
	// XXX: Should not happen!
	if notFound {
		return NoSuchEntryError("no such source")
	}
	return nil
}

func (m *Manager) removeFeed(feedID int64) error {
	f := m.findFeedByID(feedID)
	if f == nil {
		return NoSuchEntryError("no such feed")
	}
	f.invalid.Store(true)
	const sql = `DELETE FROM feeds WHERE id = $1`
	if err := m.db.Run(
		context.Background(),
		func(ctx context.Context, con *pgxpool.Conn) error {
			_, err := con.Exec(ctx, sql, feedID)
			return err
		}, 0,
	); err != nil {
		return fmt.Errorf("deleting feed failed: %w", err)
	}
	s := f.source
	s.feeds = slices.DeleteFunc(s.feeds, func(g *feed) bool { return f == g })
	return nil
}

func (m *Manager) inManager(fn func(*Manager)) {
	done := make(chan struct{})
	m.fns <- func(m *Manager) {
		defer close(done)
		fn(m)
	}
	<-done
}

func (m *Manager) asManager(fn func(*Manager, int64) error, id int64) error {
	err := make(chan error)
	m.fns <- func(m *Manager) { err <- fn(m, id) }
	return <-err
}

// AddSource registers a new source.
func (m *Manager) AddSource(
	name string,
	url string,
	rate *float64,
	slots *int,
	headers []string,
	strictMode *bool,
	insecure *bool,
	signatureCheck *bool,
	age *time.Duration,
	ignorePatterns []*regexp.Regexp,
	clientCertPublic []byte,
	clientCertPrivate []byte,
	clientCertPassphrase []byte,
) (int64, error) {
	lpmd := m.PMD(url)
	if !lpmd.Valid() {
		return 0, InvalidArgumentError("PMD is invalid")
	}
	errCh := make(chan error)
	s := &source{
		name:                 name,
		url:                  url,
		rate:                 rate,
		slots:                slots,
		headers:              headers,
		strictMode:           strictMode,
		insecure:             insecure,
		signatureCheck:       signatureCheck,
		age:                  age,
		ignorePatterns:       ignorePatterns,
		clientCertPublic:     clientCertPublic,
		clientCertPrivate:    clientCertPrivate,
		clientCertPassphrase: clientCertPassphrase,
	}
	if clientCertPrivate != nil {
		var err error
		if clientCertPrivate, err = m.encrypt(clientCertPrivate); err != nil {
			return 0, err
		}
	}
	if clientCertPassphrase != nil {
		var err error
		if clientCertPassphrase, err = m.encrypt(clientCertPassphrase); err != nil {
			return 0, err
		}
	}
	m.fns <- func(m *Manager) {
		if m.findSourceByName(name) != nil {
			errCh <- InvalidArgumentError("source already exists")
			return
		}
		const sql = `INSERT INTO sources (` +
			`name, url, rate, slots, headers, ` +
			`strict_mode, insecure, signature_check, age, ignore_patterns, ` +
			`client_cert_public, client_cert_private, client_cert_passphrase) ` +
			`VALUES (` +
			`$1, $2, $3, $4, $5, ` +
			`$6, $7, $8, $9, $10, ` +
			`$11, $12, $13) ` +
			`RETURNING id`
		if err := m.db.Run(
			context.Background(),
			func(rctx context.Context, con *pgxpool.Conn) error {
				return con.QueryRow(rctx, sql,
					name, url, rate, slots, headers,
					strictMode, insecure, signatureCheck, age, ignorePatterns,
					clientCertPublic, clientCertPrivate, clientCertPassphrase,
				).Scan(&s.id)
			}, 0,
		); err != nil {
			errCh <- fmt.Errorf("adding source to database failed: %w", err)
			return
		}
		m.sources = append(m.sources, s)
		errCh <- nil
	}
	return s.id, <-errCh
}

// AddFeed adds a new feed to a source.
func (m *Manager) AddFeed(
	sourceID int64,
	label string,
	url *url.URL,
	logLevel config.FeedLogLevel,
) (int64, error) {
	var feedID int64
	errCh := make(chan error)
	m.fns <- func(m *Manager) {
		s := m.findSourceByID(sourceID)
		if s == nil {
			errCh <- NoSuchEntryError("no such source")
			return
		}
		if slices.ContainsFunc(s.feeds, func(f *feed) bool { return f.label == label }) {
			errCh <- InvalidArgumentError("label already exists")
			return
		}
		pmd, err := asProviderMetaData(m.PMD(s.url))
		if err != nil {
			errCh <- err
			return
		}
		rolie := isROLIEFeed(pmd, url.String())
		if !rolie && !isDirectoryFeed(pmd, url.String()) {
			errCh <- InvalidArgumentError("feed is neither ROLIE nor directory based")
			return
		}
		const sql = `INSERT INTO feeds (label, sources_id, url, rolie, log_lvl) ` +
			`VALUES ($1, $2, $3, $4, $5::feed_logs_level) ` +
			`RETURNING id`
		if err := m.db.Run(
			context.Background(),
			func(ctx context.Context, conn *pgxpool.Conn) error {
				return conn.QueryRow(ctx, sql,
					label,
					sourceID,
					url.String(),
					rolie,
					logLevel,
				).Scan(&feedID)
			}, 0,
		); err != nil {
			errCh <- fmt.Errorf("inserting feed failed: %w", err)
			return
		}
		f := &feed{
			id:     feedID,
			label:  label,
			url:    url,
			rolie:  rolie,
			source: s,
		}
		f.logLevel.Store(int32(logLevel))
		s.feeds = append(s.feeds, f)
		if s.active {
			m.backgroundPing()
		}
		errCh <- nil
	}
	if err := <-errCh; err != nil {
		return 0, err
	}
	return feedID, nil
}

// RemoveSource removes a sources from manager.
func (m *Manager) RemoveSource(sourceID int64) error {
	return m.asManager((*Manager).removeSource, sourceID)
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

// PMD returns the provider metadata from the given url.
func (m *Manager) PMD(url string) *csaf.LoadedProviderMetadata {
	return m.pmdCache.pmd(m, url)
}

// SourceUpdater offers a protocol to update a source. Call the UpdateX
// (with X in Name, Rate, ...) methods to update specfic fields.
type SourceUpdater struct {
	manager           *Manager
	source            *source
	clientCertUpdated bool
	changes           []func(*source)
	fields            []string
	values            []any
}

func (su *SourceUpdater) addChange(ch func(*source), field string, value any) {
	if !slices.Contains(su.fields, field) {
		su.changes = append(su.changes, ch)
		su.fields = append(su.fields, field)
		su.values = append(su.values, value)
	}
}

func (su *SourceUpdater) applyChanges() bool {
	for _, ch := range su.changes {
		if ch != nil {
			ch(su.source)
		}
	}
	return len(su.changes) > 0
}

// UpdateName requests a name update.
func (su *SourceUpdater) UpdateName(name string) error {
	if name == su.source.name {
		return nil
	}
	if name == "" || su.manager.findSourceByName(name) != nil {
		return InvalidArgumentError("invalid name")
	}
	su.addChange(func(s *source) { s.name = name }, "name", name)
	return nil
}

// UpdateRate requests a rate update.
func (su *SourceUpdater) UpdateRate(rate *float64) error {
	if rate == nil && su.source.rate == nil {
		return nil
	}
	if rate != nil && su.source.rate != nil && *rate == *su.source.rate {
		return nil
	}
	if rate != nil && (*rate <= 0 || *rate > su.manager.cfg.Sources.MaxRatePerSource) {
		return InvalidArgumentError("invalid rate value")
	}
	su.addChange(func(s *source) { s.setRate(rate) }, "rate", rate)
	return nil
}

// UpdateSlots requests a slots update.
func (su *SourceUpdater) UpdateSlots(slots *int) error {
	if slots == nil && su.source.slots == nil {
		return nil
	}
	if slots != nil && su.source.slots != nil && *slots == *su.source.slots {
		return nil
	}
	if slots != nil && (*slots < 1 || *slots > su.manager.cfg.Sources.MaxSlotsPerSource) {
		return InvalidArgumentError("invalid slot value")
	}
	su.addChange(func(s *source) { s.slots = slots }, "slots", slots)
	return nil
}

// UpdateActive requests an active update.
func (su *SourceUpdater) UpdateActive(active bool) error {
	if active == su.source.active {
		return nil
	}
	su.addChange(func(s *source) {
		s.active = active
		if active {
			su.manager.backgroundPing()
		}
	}, "active", active)
	return nil
}

// clone as slices.Clone sadly does not work that way.
func clone[S ~[]E, E any](s S) S {
	if len(s) == 0 {
		return nil
	}
	return append(make(S, 0, len(s)), s...)
}

// UpdateHeaders requests a headers update.
func (su *SourceUpdater) UpdateHeaders(headers []string) error {
	if slices.Equal(headers, su.source.headers) {
		return nil
	}
	headers = clone(headers)
	su.addChange(func(s *source) { s.headers = headers }, "headers", headers)
	return nil
}

// UpdateStrictMode requests an update on strictMode.
func (su *SourceUpdater) UpdateStrictMode(strictMode *bool) error {
	if su.source.strictMode == nil && strictMode == nil {
		return nil
	}
	if su.source.strictMode != nil && strictMode != nil && *su.source.strictMode == *strictMode {
		return nil
	}
	su.addChange(func(s *source) { s.strictMode = strictMode }, "strict_mode", strictMode)
	return nil
}

// UpdateInsecure requests an update on insecure.
func (su *SourceUpdater) UpdateInsecure(insecure *bool) error {
	if su.source.insecure == nil && insecure == nil {
		return nil
	}
	if su.source.insecure != nil && insecure != nil && *su.source.insecure == *insecure {
		return nil
	}
	su.addChange(func(s *source) { s.insecure = insecure }, "insecure", insecure)
	return nil
}

// UpdateSignatureCheck requests an update on signatureCheck.
func (su *SourceUpdater) UpdateSignatureCheck(signatureCheck *bool) error {
	if su.source.signatureCheck == nil && signatureCheck == nil {
		return nil
	}
	if su.source.signatureCheck != nil && signatureCheck != nil && *su.source.signatureCheck == *signatureCheck {
		return nil
	}
	su.addChange(func(s *source) { s.signatureCheck = signatureCheck }, "signature_check", signatureCheck)
	return nil
}

// UpdateAge requests an update on age.
func (su *SourceUpdater) UpdateAge(age *time.Duration) error {
	if su.source.age == nil && age == nil {
		return nil
	}
	if su.source.age != nil && age != nil && *su.source.age == *age {
		return nil
	}
	su.addChange(func(s *source) { s.setAge(age) }, "age", age)
	return nil
}

// UpdateIgnorePatterns requests an update on ignorepatterns.
func (su *SourceUpdater) UpdateIgnorePatterns(ignorePatterns []*regexp.Regexp) error {
	if slices.EqualFunc(su.source.ignorePatterns, ignorePatterns,
		func(a, b *regexp.Regexp) bool { return a != nil && b != nil && a.String() == b.String() }) {
		return nil
	}
	ignorePatterns = clone(ignorePatterns)
	su.addChange(func(s *source) { s.setIgnorePatterns(ignorePatterns) }, "ignore_patterns", ignorePatterns)
	return nil
}

// UpdateClientCertPublic requests an update ob client cert public part.
func (su *SourceUpdater) UpdateClientCertPublic(data []byte) error {
	if data == nil && su.source.clientCertPublic == nil {
		return nil
	}
	if data != nil && su.source.clientCertPublic != nil && slices.Equal(data, su.source.clientCertPublic) {
		return nil
	}
	data = clone(data)
	su.addChange(func(s *source) {
		su.clientCertUpdated = true
		s.clientCertPublic = data
	}, "client_cert_public", data)
	return nil
}

// UpdateClientCertPrivate requests an update ob client cert private part.
func (su *SourceUpdater) UpdateClientCertPrivate(data []byte) error {
	orig := su.source.clientCertPrivate
	if data == nil && orig == nil {
		return nil
	}
	if data != nil && orig != nil && slices.Equal(data, orig) {
		return nil
	}
	encrypted, err := su.manager.encrypt(data)
	if err != nil {
		return err
	}
	data = clone(data)
	su.addChange(func(s *source) {
		su.clientCertUpdated = true
		s.clientCertPrivate = data
	}, "client_cert_private", encrypted)
	return nil
}

// UpdateClientCertPassphrase requests an update ob client cert private part.
func (su *SourceUpdater) UpdateClientCertPassphrase(data []byte) error {
	orig := su.source.clientCertPassphrase
	if data == nil && orig == nil {
		return nil
	}
	if data != nil && orig != nil && slices.Equal(data, orig) {
		return nil
	}
	encrypted, err := su.manager.encrypt(data)
	if err != nil {
		return err
	}
	data = clone(data)
	su.addChange(func(s *source) {
		su.clientCertUpdated = true
		s.clientCertPassphrase = data
	}, "client_cert_passphrase", encrypted)
	return nil
}

func (su *SourceUpdater) updateDB() error {
	if len(su.fields) == 0 {
		return nil
	}
	var ob, cb string
	if len(su.fields) > 1 {
		ob, cb = "(", ")"
	}
	sql := fmt.Sprintf(
		"UPDATE sources SET %[1]s%[3]s%[2]s = %[1]s%[4]s%[2]s WHERE id = %[5]d",
		ob, cb,
		strings.Join(su.fields, ","),
		placeholders(len(su.values)),
		su.source.id)
	return su.manager.db.Run(
		context.Background(),
		func(ctx context.Context, conn *pgxpool.Conn) error {
			_, err := conn.Exec(ctx, sql, su.values...)
			return err
		}, 0)
}

func placeholders(n int) string {
	var b strings.Builder
	for i := range n {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteByte('$')
		b.WriteString(strconv.Itoa(i + 1))
	}
	return b.String()
}

// UpdateSource passes an updater to manipulate a source for a given id to a given callback.
func (m *Manager) UpdateSource(
	sourceID int64,
	updates func(*SourceUpdater) error,
) (SourceUpdateResult, error) {
	type result struct {
		v   SourceUpdateResult
		err error
	}
	resCh := make(chan result)
	m.fns <- func(m *Manager) {
		s := m.findSourceByID(sourceID)
		if s == nil {
			resCh <- result{err: NoSuchEntryError("no such source")}
			return
		}
		su := SourceUpdater{manager: m, source: s}
		if err := updates(&su); err != nil {
			resCh <- result{err: fmt.Errorf("updates failed: %w", err)}
			return
		}
		if err := su.updateDB(); err != nil {
			resCh <- result{err: fmt.Errorf("updating database failed: %w", err)}
			return
		}
		// Only apply changes if database updates went through.
		if !su.applyChanges() {
			resCh <- result{v: SourceUnchanged}
			return
		}
		if su.clientCertUpdated {
			if err := s.updateCertificate(); err != nil {
				slog.Warn("updating client cert failed", "warn", err)
				if s.active {
					s.active = false
					x := SourceUpdater{manager: m, source: s}
					x.addChange(nil, "active", false)
					if err := x.updateDB(); err != nil {
						slog.Error("deactivating source failed", "err", err)
					}
					resCh <- result{v: SourceDeactivated}
					return
				}
			}
		}
		resCh <- result{v: SourceUpdated}
	}
	res := <-resCh
	return res.v, res.err
}
