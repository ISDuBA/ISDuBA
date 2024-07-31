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
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/ISDuBA/ISDuBA/pkg/database"
	"github.com/ISDuBA/ISDuBA/pkg/version"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var userAgent = "isduba/" + version.SemVersion

type location struct {
	updated   time.Time
	doc       *url.URL
	hash      *url.URL
	signature *url.URL
}

type activeFeed struct {
	id        int64
	url       *url.URL
	rolie     bool
	nextCheck time.Time
	locations []location
	source    *activeSource

	lastETag      string
	modifiedSince time.Time
}

type activeSource struct {
	id    int64
	rate  *float64
	slots *int
	feeds []*activeFeed
}

func (af *activeFeed) fetchIndex() ([]byte, error) {

	req, err := http.NewRequest(http.MethodGet, af.url.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", userAgent)
	if af.lastETag != "" {
		req.Header.Add("If-None-Match", af.lastETag)
	}
	if !af.modifiedSince.IsZero() {
		req.Header.Add("If-Modified-Since", af.modifiedSince.Format(http.TimeFormat))
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusNotModified {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	af.lastETag = resp.Header.Get("Etag")
	if m := resp.Header.Get("Last-Modified"); m != "" {
		af.modifiedSince, _ = time.Parse(http.TimeFormat, m)
	}
	return content, nil
}

func (af *activeFeed) refresh(db *database.DB) error {

	if !af.rolie {
		// TODO: Implement me!
		slog.Warn("None-ROLIE feeds are not implemented, yet", "feed", af.id)
		return nil
	}
	content, err := af.fetchIndex()
	if err != nil {
		return fmt.Errorf("fetching feed index failed: %w", err)
	}
	if content == nil {
		slog.Info("Feed has not changed", "feed", af.id)
		return nil
	}

	// TODO: Virtualize the None-ROLIE case
	rolie, err := rolieFromData(content)
	if err != nil {
		return fmt.Errorf("de-serializing rolie failed: %w", err)
	}

	// Extract candidates from feed leaving out location where we already
	// have requests in memory which are more recent.
	candidates, err := rolie.toLocations(af.url, af.sameOrNewer())
	if err != nil {
		return fmt.Errorf("extracting locations from feed failed: %w", err)
	}

	// Filter out candidates which are already in the database with same or newer.
	if candidates, err = removeOlder(db, candidates); err != nil {
		return fmt.Errorf("removing candidates by looking at database failed: %w", err)
	}

	if len(candidates) == 0 { // Nothing to do.
		return nil
	}

	// Merge candidates into list of locations.
	af.locations = append(af.locations, candidates...)
	slices.SortFunc(af.locations, func(a, b location) int {
		return a.updated.Compare(b.updated)
	})

	slog.Info("Entries in feed", "num", len(af.locations), "feed", af.id)
	return nil
}

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

func (af *activeFeed) sameOrNewer() func(*location) bool {
	// XXX: Maybe this extra indexing could be replaced by something
	// which uses the fact that the locations are already sorted by updated?!
	have := make(map[string]time.Time, len(af.locations))
	for _, location := range af.locations {
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
