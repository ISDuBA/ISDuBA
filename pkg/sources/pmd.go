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
	"crypto/sha1"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/ISDuBA/ISDuBA/pkg/cache"
	"github.com/ISDuBA/ISDuBA/pkg/config"
	"github.com/gocsaf/csaf/v3/csaf"
	"github.com/gocsaf/csaf/v3/util"
)

// holdingPMDsDuration is the duration how long PMDs are cached.
const holdingPMDsDuration = time.Minute * 15

// CachedProviderMetadata holds a loaded PMD and enables access to
// the respective model.
type CachedProviderMetadata struct {
	Loaded  *csaf.LoadedProviderMetadata
	modelMu sync.Mutex
	model   *csaf.ProviderMetadata
}

type pmdCache struct {
	*cache.ExpirationCache[string, *CachedProviderMetadata]
}

type resolvedPMD struct {
	url string
	pmd *csaf.ProviderMetadata
}

type resolvedPMDs []resolvedPMD

func newPMDCache() *pmdCache {
	return &pmdCache{
		ExpirationCache: cache.NewExpirationCache[string, *CachedProviderMetadata](holdingPMDsDuration),
	}
}

func (pc *pmdCache) pmd(url string, cfg *config.Config) *CachedProviderMetadata {

	if cpmd, ok := pc.Get(url); ok {
		return cpmd
	}

	header := http.Header{}
	header.Add("User-Agent", UserAgent)

	baseClient := &http.Client{
		Transport: cfg.General.Transport(),
	}
	if timeout := cfg.Sources.Timeout; timeout > 0 {
		baseClient.Timeout = timeout
	}

	client := util.Client(&util.HeaderClient{
		Client: baseClient,
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
	cpmd := &CachedProviderMetadata{Loaded: lpmd}
	pc.Set(url, cpmd)
	return cpmd
}

// Valid returns true if the loaded PMD is valid.
func (cpmd *CachedProviderMetadata) Valid() bool {
	return cpmd != nil && cpmd.Loaded.Valid()
}

// Model returns the model for the loaded PMD.
func (cpmd *CachedProviderMetadata) Model() (*csaf.ProviderMetadata, error) {
	if !cpmd.Valid() {
		return nil, InvalidArgumentError("PMD is invalid")
	}
	cpmd.modelMu.Lock()
	defer cpmd.modelMu.Unlock()
	if cpmd.model != nil {
		return cpmd.model, nil
	}
	model := new(csaf.ProviderMetadata)
	// XXX: This is ugly! We should better keep the original data when loading the PMD.
	if err := util.ReMarshalJSON(model, cpmd.Loaded.Document); err != nil {
		return nil, InvalidArgumentError(
			fmt.Sprintf("re-marshaling of PDM failed: %v", err.Error()))
	}
	cpmd.model = model
	return model, nil
}

// availableFeeds returns a list of the feeds available for the given provider.
func availableFeeds(pmd *csaf.ProviderMetadata) []string {
	var feeds []string
	add := func(feed string) {
		if !slices.Contains(feeds, feed) {
			feeds = append(feeds, feed)
		}
	}
	// ROLIE feeds
	for i := range pmd.Distributions {
		d := pmd.Distributions[i]
		if d.Rolie == nil {
			continue
		}
		feeds := d.Rolie.Feeds
		for j := range feeds {
			if f := &feeds[j]; f.URL != nil {
				add(string(*f.URL))
			}
		}
	}
	// Directory feeds
	for i := range pmd.Distributions {
		if d := pmd.Distributions[i]; d.Rolie == nil && d.DirectoryURL != "" {
			add(d.DirectoryURL)
		}
	}
	return feeds
}

// checksumPMD calculates a checksum over the relevant fields in a PMD.
// Currently only the feed paths are used.
func checksumPMD(pmd *csaf.ProviderMetadata) []byte {
	feeds := availableFeeds(pmd)
	hash := sha1.New()
	for _, feed := range feeds {
		hash.Write([]byte(feed))
	}
	return hash.Sum(nil)
}

// isROLIEFeed checks if the given url leads to a ROLIE feed.
func isROLIEFeed(pmd *csaf.ProviderMetadata, url string) bool {
	for i := range pmd.Distributions {
		d := pmd.Distributions[i]
		if d.Rolie == nil {
			continue
		}
		feeds := d.Rolie.Feeds
		for j := range feeds {
			if f := &feeds[j]; f.URL != nil && string(*f.URL) == url {
				return true
			}
		}
	}
	return false
}

// isDirectoryFeed checks if the given url leads to a directory based feed.
func isDirectoryFeed(pmd *csaf.ProviderMetadata, url string) bool {
	for i := range pmd.Distributions {
		if d := pmd.Distributions[i]; d.Rolie == nil && d.DirectoryURL == url {
			return true
		}
	}
	return false
}

// add deduplicates urls as each lookup is expensive.
func (rps *resolvedPMDs) add(urls ...string) {
	for _, url := range urls {
		if !slices.ContainsFunc(*rps, func(rp resolvedPMD) bool {
			return rp.url == url
		}) {
			*rps = append(*rps, resolvedPMD{url: url})
		}
	}
}

const numURLResolvers = 5

// resolve resolves all urls added with add to PMDs.
func (rps resolvedPMDs) resolve(cache *pmdCache, cfg *config.Config) {
	var (
		wg        sync.WaitGroup
		toResolve = make(chan *resolvedPMD)
	)
	worker := func() {
		defer wg.Done()
		for tr := range toResolve {
			cpmd := cache.pmd(tr.url, cfg)
			if !cpmd.Valid() {
				slog.Debug("Invalid PMD", "url", tr.url)
				continue
			}
			pmd, err := cpmd.Model()
			if err != nil {
				slog.Debug("Invalid PMD model", "url", tr.url, "err", err)
				continue
			}
			tr.pmd = pmd
		}
	}
	for range max(1, min(len(rps), numURLResolvers)) {
		wg.Add(1)
		go worker()
	}
	for i := range rps {
		toResolve <- &rps[i]
	}
	close(toResolve)
	wg.Wait()
}

func (rps resolvedPMDs) pmd(url string) *csaf.ProviderMetadata {
	if idx := slices.IndexFunc(rps, func(rp resolvedPMD) bool { return rp.url == url }); idx >= 0 {
		return rps[idx].pmd
	}
	return nil
}
