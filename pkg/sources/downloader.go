// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package sources

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/ISDuBA/ISDuBA/pkg/version"
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

func (af *activeFeed) refresh() error {

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
		slog.Info("Feed does not change", "feed", af.id)
		return nil
	}

	// TODO: Virtualize the None-ROLIE case
	rolie, err := rolieFromData(content)
	if err != nil {
		return fmt.Errorf("de-serializing rolie failed: %w", err)
	}
	locations, err := rolie.toLocations(af.url)
	slog.Info("Entries in feed", "num", len(locations), "feed", af.id)
	return nil
}
