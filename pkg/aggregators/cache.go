// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package aggregators

import (
	"time"

	"github.com/ISDuBA/ISDuBA/internal/cache"
)

// holdingDuration is the duration how long PMDs are cached.
const holdingDuration = time.Minute * 15

// CachedAggregator are cached aggregators.
type CachedAggregator struct {
}

// Cache is cache of aggregators.
type Cache struct {
	*cache.ExpirationCache[string, *CachedAggregator]
}

func newCache() *Cache {
	return &Cache{
		ExpirationCache: cache.NewExpirationCache[string, *CachedAggregator](holdingDuration),
	}
}
