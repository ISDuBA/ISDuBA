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
	"net/http"
	"sync"
	"time"

	"github.com/ISDuBA/ISDuBA/internal/cache"
	"github.com/csaf-poc/csaf_distribution/v3/csaf"
	"github.com/csaf-poc/csaf_distribution/v3/util"
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

func newPMDCache() *pmdCache {
	return &pmdCache{
		ExpirationCache: cache.NewExpirationCache[string, *CachedProviderMetadata](holdingPMDsDuration),
	}
}

func (pc *pmdCache) pmd(m *Manager, url string) *CachedProviderMetadata {

	if cpmd, ok := pc.Get(url); ok {
		return cpmd
	}

	header := http.Header{}
	header.Add("User-Agent", UserAgent)

	baseClient := &http.Client{}
	if m.cfg.Sources.Timeout > 0 {
		baseClient.Timeout = m.cfg.Sources.Timeout
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
