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
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/gocsaf/csaf/v3/csaf"
)

type feedIndex struct {
	base           *url.URL
	age            *time.Duration
	ignorePatterns ignorePatterns
	sameOrNewer    func(*location) bool
}

// rolieLocations assumes that the feed index is ROLIE.
func (fi *feedIndex) rolieLocations(r io.Reader) ([]location, error) {
	rolie, err := csaf.LoadROLIEFeed(r)
	if err != nil {
		return nil, fmt.Errorf("loading rolie feed from data failed: %w", err)
	}
	resolve := func(href string, store **url.URL) error {
		u, err := url.Parse(href)
		if err != nil {
			return fmt.Errorf("invalid href: %v", href)
		}
		if !u.IsAbs() {
			u = joinURL(fi.base, u)
		}
		*store = u
		return nil
	}
	// If we have a max age set calculate the cut time.
	var cut time.Time
	if fi.age != nil {
		cut = time.Now().Add(-*fi.age)
	}
	// Extract the locations
	entries := rolie.Feed.Entry
	dls := make([]location, 0, len(entries))
nextEntry:
	for _, entry := range entries {
		links := entry.Link
		updated := time.Time(entry.Updated)
		// Apply age filter
		if fi.age != nil && updated.Before(cut) {
			continue
		}
		dl := location{updated: updated}
		sha512 := false
	nextLink:
		for j := range links {
			link := &links[j]
			switch link.Rel {
			case "self":
				if err := resolve(link.HRef, &dl.doc); err != nil {
					return nil, err
				}
				// Apply ignore patterns
				if fi.ignorePatterns.ignore(dl.doc) {
					continue nextEntry
				}
			case "signature":
				if err := resolve(link.HRef, &dl.signature); err != nil {
					return nil, err
				}
			case "hash":
				if sha512 {
					// If we already have SHA512 don't bother with others.
					continue nextLink
				}
				switch href := strings.ToLower(link.HRef); {
				case strings.HasSuffix(href, ".sha512"):
					if err := resolve(link.HRef, &dl.hash); err != nil {
						return nil, err
					}
					sha512 = true
				case strings.HasSuffix(href, ".sha256"):
					if err := resolve(link.HRef, &dl.hash); err != nil {
						return nil, err
					}
				default:
					slog.Warn("unknown hash format", "href", link.HRef)
				}
			}
		}
		// Only append if we don't have already the same or we are
		// waiting to request a new one.
		if dl.doc != nil {
			if fi.sameOrNewer != nil && fi.sameOrNewer(&dl) {
				continue
			}
			dls = append(dls, dl)
		}
	}
	return dls, nil
}

// directoryLocations assumes that the feed index is changes.csv
func (fi *feedIndex) directoryLocations(r io.Reader) ([]location, error) {
	c := csv.NewReader(r)
	c.FieldsPerRecord = 2
	c.ReuseRecord = true

	// If we have a max age set calculate the cut time.
	var cut time.Time
	if fi.age != nil {
		cut = time.Now().Add(-*fi.age)
	}

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
			doc = joinURL(fi.base, doc)
		}
		// Apply age filter
		if fi.age != nil && updated.Before(cut) {
			continue
		}
		// Apply ignore patterns
		if fi.ignorePatterns.ignore(doc) {
			continue
		}
		dl := location{
			updated: updated,
			doc:     doc,
		}
		if fi.sameOrNewer != nil && fi.sameOrNewer(&dl) {
			continue
		}
		dls = append(dls, dl)
	}

	return dls, nil
}
