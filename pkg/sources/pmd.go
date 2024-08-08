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
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/csaf-poc/csaf_distribution/v3/csaf"
	"github.com/csaf-poc/csaf_distribution/v3/util"
)

// holdingPMDsDuration is the duration how long PMDs are cached.
const holdingPMDsDuration = time.Minute * 15

type pmdCacheEntry struct {
	expires time.Time
	lpmd    *csaf.LoadedProviderMetadata
}

type pmdCache struct {
	mu      sync.Mutex
	entries map[string]*pmdCacheEntry
}

func newPMDCache() *pmdCache {
	return &pmdCache{
		entries: map[string]*pmdCacheEntry{},
	}
}

func (pc *pmdCache) cleanup() {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	now := time.Now()
	for url, entry := range pc.entries {
		if entry.expires.Before(now) {
			delete(pc.entries, url)
		}
	}
}

func (pc *pmdCache) pmd(url string) *csaf.LoadedProviderMetadata {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	if e := pc.entries[url]; e != nil {
		if e.expires.After(time.Now()) {
			return e.lpmd
		}
		delete(pc.entries, url)
	}
	header := http.Header{}
	header.Add("User-Agent", UserAgent)
	client := util.Client(&util.HeaderClient{
		Client: &http.Client{},
		Header: header,
	})
	if slog.Default().Enabled(context.Background(), slog.LevelDebug) {
		client = &util.LoggingClient{
			Client: client,
			Log: func(method, url string) {
				slog.Debug("looking up PMD", "method", method, "url", url)
			},
		}
	}
	pmdLoader := csaf.NewProviderMetadataLoader(client)
	lpmd := pmdLoader.Load(url)
	pc.entries[url] = &pmdCacheEntry{
		expires: time.Now().Add(holdingPMDsDuration),
		lpmd:    lpmd,
	}
	return lpmd
}
