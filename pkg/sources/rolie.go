// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package sources

import (
	"encoding/json"
	"fmt"
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

func (r *rolie) toLocations(base *url.URL) ([]location, error) {
	resolve := func(href string, store **url.URL) error {
		u, err := url.Parse(href)
		if err != nil {
			return fmt.Errorf("invalid href: %v", href)
		}
		if u.IsAbs() {
			*store = u
		} else {
			*store = base.ResolveReference(u)
		}
		return nil
	}
	entries := r.Feed.Entries
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
		if dl.doc != nil {
			dls = append(dls, dl)
		}
	}
	return dls, nil
}

func rolieFromData(data []byte) (*rolie, error) {
	var rolie rolie
	if err := json.Unmarshal(data, &rolie); err != nil {
		return nil, fmt.Errorf("rolie from data failed: %w", err)
	}
	return &rolie, nil
}
