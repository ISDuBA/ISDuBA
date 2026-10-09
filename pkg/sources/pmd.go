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
	v20 "github.com/gocsaf/csaf/v3/csaf/v20"
	v21 "github.com/gocsaf/csaf/v3/csaf/v21"
	"github.com/gocsaf/csaf/v3/util"
)

// holdingPMDsDuration is the duration how long PMDs are cached.
const holdingPMDsDuration = time.Minute * 15

type pmdModel struct {
	version csaf.SpecVersion
	v20     *v20.Provider
	v21     *v21.Provider
}

// CachedProviderMetadata holds a loaded PMD and enables access to
// the respective model.
type CachedProviderMetadata struct {
	Loaded  *csaf.LoadedProviderMetadata
	modelMu sync.Mutex
	model   *pmdModel
}

type pmdCache struct {
	*cache.ExpirationCache[string, *CachedProviderMetadata]
}

type resolvedPMD struct {
	url string
	pmd *pmdModel
}

type resolvedPMDs []resolvedPMD

func newPMDCache() *pmdCache {
	return &pmdCache{
		ExpirationCache: cache.NewExpirationCache[string, *CachedProviderMetadata](holdingPMDsDuration),
	}
}

// CanonicalURL returns the provider canonical url string
func (pm *pmdModel) CanonicalURL() string {
	switch {
	case pm.v20 != nil:
		return string(pm.v20.CanonicalURL)
	case pm.v21 != nil:
		return string(pm.v21.CanonicalURL)
	default:
		return ""
	}
}

type pgpKeyElem struct {
	URL         string
	Fingerprint string
}

// PgpKeyElems returns the list of strings of the provider openpgpkey urls
func (pm *pmdModel) PgpKeyElems() []pgpKeyElem {
	var keyElems []pgpKeyElem
	switch {
	case pm.v20 != nil:
		for _, elem := range pm.v20.PublicOpenpgpKeys {
			keyElems = append(keyElems, pgpKeyElem{URL: string(elem.URL), Fingerprint: *elem.Fingerprint})
		}
	case pm.v21 != nil:
		for _, elem := range pm.v21.PGPKeys {
			keyElems = append(keyElems, pgpKeyElem{URL: string(elem.URL), Fingerprint: elem.Fingerprint})
		}
	}
	return keyElems
}

// ROLIEFeedURLs returns the list of strings of provider rolie feed urls
func (pm *pmdModel) ROLIEFeedURLs() []string {
	var rolieFeedURLs []string
	switch {
	case pm.v20 != nil:
		for _, distElem := range pm.v20.Distributions {
			for _, feedsElem := range distElem.Rolie.Feeds {
				rolieFeedURLs = append(rolieFeedURLs, string(feedsElem.URL))
			}
		}
	case pm.v21 != nil:
		for _, distElem := range pm.v21.Distributions {
			for _, feedsElem := range distElem.Rolie.Feeds {
				rolieFeedURLs = append(rolieFeedURLs, string(feedsElem.URL))
			}
		}
	}
	return rolieFeedURLs
}

// DirectoryURLs returns the list of strings of provider directory urls
func (pm *pmdModel) DirectoryURLs() []string {
	var directoryURLs []string
	switch {
	case pm.v20 != nil:
		for _, distElem := range pm.v20.Distributions {
			if distElem.Rolie == nil && distElem.DirectoryURL != nil {
				directoryURLs = append(directoryURLs, string(*distElem.DirectoryURL))
			}
		}
	case pm.v21 != nil:
		for _, distElem := range pm.v21.Distributions {
			if distElem.Rolie == nil && distElem.Directory.URL != "" {
				directoryURLs = append(directoryURLs, string(distElem.Directory.URL))
			}
		}
	}
	return directoryURLs
}

func (pc *pmdCache) pmd(ctx context.Context, url string, cfg *config.Config) *CachedProviderMetadata {

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

	if slog.Default().Enabled(ctx, slog.LevelDebug) {
		client = &util.LoggingClient{
			Client: client,
			Log: func(method, url string) {
				slog.Debug("looking up PMD", "method", method, "url", url)
			},
		}
	}
	pmdLoader := csaf.NewProviderMetadataLoader(client)
	lpmd := pmdLoader.LoadWithContext(ctx, url)
	cpmd := &CachedProviderMetadata{Loaded: lpmd}
	pc.Set(url, cpmd)
	return cpmd
}

// Valid returns true if the loaded PMD is valid.
func (cpmd *CachedProviderMetadata) Valid() bool {
	return cpmd != nil && cpmd.Loaded.Valid()
}

// Model returns the model for the loaded PMD.
func (cpmd *CachedProviderMetadata) Model() (*pmdModel, error) {
	if !cpmd.Valid() {
		return nil, InvalidArgumentError("PMD is invalid")
	}
	cpmd.modelMu.Lock()
	defer cpmd.modelMu.Unlock()
	if cpmd.model != nil {
		return cpmd.model, nil
	}
	doc, ok := cpmd.Loaded.Document.(map[string]any)
	if !ok {
		return nil, InvalidArgumentError("PMD is not valid JSON")
	}
	version, _ := doc["metadata_version"].(csaf.SpecVersion)
	model := &pmdModel{version: version}
	var err error
	switch version {
	case csaf.Version20:
		model.v20 = new(v20.Provider)
		// XXX: This is ugly! We should better keep the original data when loading the PMD.
		err = util.ReMarshalJSON(model.v20, doc)
	case csaf.Version21:
		model.v21 = new(v21.Provider)
		// XXX: This is ugly! We should better keep the original data when loading the PMD.
		err = util.ReMarshalJSON(model.v21, doc)
	default:
		return nil, InvalidArgumentError(fmt.Sprintf("unsupported metadata_version %q", version))
	}
	if err != nil {
		return nil, err
	}
	cpmd.model = model
	return model, nil
}

// availableFeeds returns a list of the feeds available for the given provider.
func availableFeeds(pmd *pmdModel) []string {
	var feeds []string
	for _, feed := range append(pmd.ROLIEFeedURLs(), pmd.DirectoryURLs()...) {
		if !slices.Contains(feeds, feed) {
			feeds = append(feeds, feed)
		}
	}
	return feeds
}

// checksumPMD calculates a checksum over the relevant fields in a PMD.
// Currently only the feed paths are used.
func checksumPMD(pmd *pmdModel) []byte {
	feeds := availableFeeds(pmd)
	hash := sha1.New()
	for _, feed := range feeds {
		hash.Write([]byte(feed))
	}
	return hash.Sum(nil)
}

// isROLIEFeed checks if the given url leads to a ROLIE feed.
func isROLIEFeed(pmd *pmdModel, url string) bool {
	return slices.Contains(pmd.ROLIEFeedURLs(), url)
}

// isDirectoryFeed checks if the given url leads to a directory based feed.
func isDirectoryFeed(pmd *pmdModel, url string) bool {
	return slices.Contains(pmd.DirectoryURLs(), url)
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
func (rps resolvedPMDs) resolve(ctx context.Context, cache *pmdCache, cfg *config.Config) {
	var (
		wg        sync.WaitGroup
		toResolve = make(chan *resolvedPMD)
	)
	worker := func() {
		defer wg.Done()
		for tr := range toResolve {
			cpmd := cache.pmd(ctx, tr.url, cfg)
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

func (rps resolvedPMDs) pmd(url string) *pmdModel {
	if idx := slices.IndexFunc(rps, func(rp resolvedPMD) bool { return rp.url == url }); idx >= 0 {
		return rps[idx].pmd
	}
	return nil
}
