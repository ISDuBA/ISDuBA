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
	"time"

	"github.com/ISDuBA/ISDuBA/internal/cache"
	"github.com/csaf-poc/csaf_distribution/v3/csaf"
	"github.com/csaf-poc/csaf_distribution/v3/util"
)

// holdingPMDsDuration is the duration how long PMDs are cached.
const holdingPMDsDuration = time.Minute * 15

type pmdCache struct {
	*cache.ExpirationCache[string, *csaf.LoadedProviderMetadata]
}

func newPMDCache() *pmdCache {
	return &pmdCache{
		ExpirationCache: cache.NewExpirationCache[string, *csaf.LoadedProviderMetadata](holdingPMDsDuration),
	}
}

func (pc *pmdCache) pmd(m *Manager, url string) *csaf.LoadedProviderMetadata {

	if lpmd, ok := pc.Get(url); ok {
		return lpmd
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
	pc.Set(url, lpmd)
	return lpmd
}

func asProviderMetaData(lpmd *csaf.LoadedProviderMetadata) (*csaf.ProviderMetadata, error) {
	if !lpmd.Valid() {
		return nil, InvalidArgumentError("PMD is invalid")
	}
	pmd := new(csaf.ProviderMetadata)
	// XXX: This is ugly! We should better keep the original data when loading the PMD.
	if err := util.ReMarshalJSON(pmd, lpmd.Document); err != nil {
		return nil, InvalidArgumentError(
			fmt.Sprintf("re-marshaling of PDM failed: %v", err.Error()))
	}
	return pmd, nil
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
