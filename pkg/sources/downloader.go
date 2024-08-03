// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package sources

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/ISDuBA/ISDuBA/pkg/config"
	"github.com/ISDuBA/ISDuBA/pkg/database"
	"github.com/ISDuBA/ISDuBA/pkg/models"
	"github.com/ISDuBA/ISDuBA/pkg/version"
	"github.com/csaf-poc/csaf_distribution/v3/csaf"
	"github.com/csaf-poc/csaf_distribution/v3/util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/time/rate"
)

var userAgent = "isduba/" + version.SemVersion

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
	url      *url.URL
	rolie    bool
	logLevel config.FeedLogLevel

	nextCheck time.Time
	locations []location
	source    *source

	lastETag     string
	lastModified time.Time
}

type source struct {
	id     int64
	rate   *float64
	slots  *int
	active bool

	feeds     []*feed
	usedSlots int
	limiterMu sync.Mutex
	limiter   *rate.Limiter
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
	if candidates, err = removeOlder(m.db, candidates); err != nil {
		return fmt.Errorf("removing candidates by looking at database failed: %w", err)
	}

	if len(candidates) == 0 { // Nothing to do.
		return nil
	}

	// Merge candidates into list of locations.
	f.locations = append(f.locations, candidates...)
	slices.SortFunc(f.locations, func(a, b location) int {
		return a.updated.Compare(b.updated)
	})

	f.log(m, config.InfoFeedLogLevel, "entries to download: %d", len(f.locations))
	return nil
}

// wait establishes the request rate per source.
func (s *source) wait(ctx context.Context) {
	if s.rate != nil {
		s.limiterMu.Lock()
		defer s.limiterMu.Unlock()
		if s.limiter == nil {
			s.limiter = rate.NewLimiter(rate.Limit(*s.rate), 1)
		}
		s.limiter.Wait(ctx)
	}
}

// doRequest executes an HTTP request with the source specific parameters.
func (s *source) doRequest(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", userAgent)

	client := http.Client{}
	s.wait(context.Background())

	// TODO: Implement me!

	return client.Do(req)
}

func (s *source) httpGet(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return s.doRequest(req)
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
func removeOlder(db *database.DB, candidates []location) ([]location, error) {

	var remove []int

	batch := pgx.Batch{}

	const sql = `SELECT EXISTS(SELECT 1 FROM changes WHERE url = $1 AND time >= $2)`

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
		candidate := &candidates[i]
		batch.Queue(sql, candidate.doc.String(), candidate.updated).QueryRow(exists(i))
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

// download fetches the files of a document and stores
// them into the database.
func (l location) download(m *Manager, f *feed, done func()) {
	defer done()

	var (
		writers        []io.Writer
		checksum       hash.Hash
		remoteChecksum []byte
	)

	// Loading the hash
	if l.hash != nil { // ROLIE gave us an URL to hash file.
		hashFile := l.hash.String()
		switch lc := strings.ToLower(hashFile); {
		case strings.HasSuffix(lc, ".sha256"):
			checksum = sha256.New()
		case strings.HasSuffix(lc, ".sha512"):
			checksum = sha512.New()
		}
		if checksum != nil {
			var err error
			if remoteChecksum, err = f.source.loadHash(hashFile); err != nil {
				f.log(m, config.WarnFeedLogLevel, "fetching hash %q failed: %v", hashFile, err)
			} else {
				writers = append(writers, checksum)
			}
		}
	} else if !f.rolie { // If we are directory based, do some guessing:
		// TODO: Implement me!
		slog.Warn("Hash loading for none.ROLIE feeds is not implement, yet.")
	}

	// Download the CSAF document.
	resp, err := f.source.httpGet(l.doc.String())
	if err != nil {
		f.log(m, config.ErrorFeedLogLevel, "downloading %q failed: %v", l.doc, err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		f.log(m, config.ErrorFeedLogLevel, "downloading %q failed: %s (%d)",
			l.doc, http.StatusText(resp.StatusCode), resp.StatusCode)
		return
	}
	defer resp.Body.Close()

	var data bytes.Buffer
	writers = append(writers, &data)

	// Prevent over-sized downloads.
	limited := io.LimitReader(resp.Body, int64(m.cfg.General.AdvisoryUploadLimit))

	tee := io.TeeReader(limited, io.MultiWriter(writers...))

	// Decode document into JSON.
	var doc any
	if err := json.NewDecoder(tee).Decode(&doc); err != nil {
		f.log(m, config.ErrorFeedLogLevel, "decoding document %q failed: %v", l.doc, err)
		return
	}

	// Compare checksums.
	if remoteChecksum != nil {
		if !bytes.Equal(checksum.Sum(nil), remoteChecksum) {
			f.log(m, config.ErrorFeedLogLevel, "Checksum mismatch for document %q", l.doc)
			return
		}
	}

	// Check document against schema.
	if errors, err := csaf.ValidateCSAF(doc); err != nil || len(errors) > 0 {
		if err != nil {
			f.log(m, config.ErrorFeedLogLevel,
				"Schema validation of document %q failed: %v", l.doc, err)
		} else {
			f.log(m, config.ErrorFeedLogLevel,
				"Schema validation of document %q has %d errors", l.doc, len(errors))
		}
		return
	}

	// TODO: Check against remote validator
	// TODO: Filename check. (???)
	// TODO: Check signatures
	// TODO: Statistics

	// Remember last changes.
	inTx := func(ctx context.Context, tx pgx.Tx, _ int64) error {
		const updatedSQL = `INSERT INTO changes (url, feeds_id, time) ` +
			`VALUES ($1, $2, $3) ` +
			`ON CONFLICT (url, feeds_id) DO ` +
			`UPDATE SET time = $3`
		_, err := tx.Exec(ctx, updatedSQL, l.doc.String(), f.id, l.updated)
		return err
	}

	ctx := context.Background()
	if err := m.db.Run(ctx, func(ctx context.Context, conn *pgxpool.Conn) error {
		_, err := models.ImportDocumentData(
			ctx, conn,
			doc, data.Bytes(),
			m.cfg.Sources.FeedImporter,
			m.cfg.Sources.PublishersTLPs,
			inTx,
			false)
		return err
	}, 0); err != nil {
		f.log(m, config.ErrorFeedLogLevel, "storing %q failed: %v", l.doc, err)
		return
	}

	f.log(m, config.InfoFeedLogLevel, "downloading %q done", l.doc)
}

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
