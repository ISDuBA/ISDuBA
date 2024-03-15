// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package ginkeycloak

import (
	"sync"
	"time"
)

type item[V any] struct {
	expires time.Time
	value   V
}

func (i *item[V]) expired() bool {
	if i.expires.IsZero() {
		return false
	}
	return i.expires.Before(time.Now())
}

type cache[K comparable, V any] struct {
	expiration time.Duration
	mu         sync.Mutex
	items      map[K]*item[V]
}

func newCache[K comparable, V any](expiration time.Duration) *cache[K, V] {
	return &cache[K, V]{
		expiration: expiration,
		items:      map[K]*item[V]{},
	}
}

func (c *cache[K, V]) get(k K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	it := c.items[k]
	if it == nil {
		var zero V
		return zero, false
	}
	if it.expired() {
		delete(c.items, k)
		var zero V
		return zero, false
	}
	return it.value, true
}

func (c *cache[K, V]) set(k K, v V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	var expires time.Time
	if c.expiration > 0 {
		expires = time.Now().Add(c.expiration)
	}
	c.items[k] = &item[V]{
		expires: expires,
		value:   v,
	}
}
