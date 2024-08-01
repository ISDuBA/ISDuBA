// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package sources

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"strings"
	"time"
)

type rolieLink struct {
	Rel  string `json:"rel"`
	HRef string `json:"href"`
}

type rolieEntry struct {
	Links   []rolieLink `json:"link"`
	Updated time.Time   `json:"updated"`
}

type rolieFeed struct {
	Updated time.Time    `json:"updated"`
	Entries []rolieEntry `json:"entry"`
}

type rolie struct {
	Feed rolieFeed `json:"feed"`
}

// rolieLocations assumes that the feed index is ROLIE.
func (af *activeFeed) rolieLocations(r io.Reader) ([]location, error) {
	// De-serialize JSON
	var rolie rolie
	if err := json.NewDecoder(r).Decode(&rolie); err != nil {
		return nil, fmt.Errorf("rolie from data failed: %w", err)
	}
	resolve := func(href string, store **url.URL) error {
		u, err := url.Parse(href)
		if err != nil {
			return fmt.Errorf("invalid href: %v", href)
		}
		if u.IsAbs() {
			*store = u
		} else {
			*store = af.url.ResolveReference(u)
		}
		return nil
	}
	sameOrNewer := af.sameOrNewer()
	// Extract the locations
	entries := rolie.Feed.Entries
	dls := make([]location, 0, len(entries))
	for i := range entries {
		entry := &entries[i]
		links := entry.Links
		dl := location{
			updated: entry.Updated,
		}
		for j := range links {
			link := &links[j]
			switch link.Rel {
			case "self":
				if err := resolve(link.HRef, &dl.doc); err != nil {
					return nil, err
				}
			case "signature":
				if err := resolve(link.HRef, &dl.signature); err != nil {
					return nil, err
				}
			case "hash":
				if h := strings.ToLower(link.HRef); strings.HasSuffix(h, ".sha256") ||
					strings.HasSuffix(h, ".sha512") {
					if err := resolve(link.HRef, &dl.hash); err != nil {
						return nil, err
					}
				} else {
					slog.Warn("unknown hash format", "href", link.HRef)
				}
			}
		}
		// Only append if we don't have already the same or we are
		// waiting to request a new one.
		if dl.doc != nil && !sameOrNewer(&dl) {
			dls = append(dls, dl)
		}
	}
	return dls, nil
}

// directoryLocations assumes that the feed index is changes.csv
func (af *activeFeed) directoryLocations(r io.Reader) ([]location, error) {
	c := csv.NewReader(r)
	c.FieldsPerRecord = 2
	c.ReuseRecord = true

	sameOrNewer := af.sameOrNewer()

	var dls []location

	for lineNo := 1; ; lineNo++ {
		record, err := c.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("CSV line %d is invalid: %w", lineNo, err)
		}
		doc, err := url.Parse(record[0])
		if err != nil {
			return nil, fmt.Errorf("column 1 in line %d is not a valid URL: %w", lineNo, err)
		}
		updated, err := time.Parse(time.RFC3339, record[1])
		if err != nil {
			return nil, fmt.Errorf("column 2 in line %d is not a valid RFC3339 time: %w", lineNo, err)
		}
		if !doc.IsAbs() {
			doc = af.url.ResolveReference(doc)
		}
		dl := location{
			updated: updated,
			doc:     doc,
		}
		if !sameOrNewer(&dl) {
			dls = append(dls, dl)
		}
	}

	return dls, nil
}
