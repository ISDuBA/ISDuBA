// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

// Package cache implements an in-memory cache with expiration for its items.
package cache

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

// ExpirationCache is a cache with a expiration duration for its items.
type ExpirationCache[K comparable, V any] struct {
	expiration time.Duration
	mu         sync.Mutex
	items      map[K]*item[V]
}

// NewExpirationCache creates a new cache with a given expiration duration.
func NewExpirationCache[K comparable, V any](expiration time.Duration) *ExpirationCache[K, V] {
	return &ExpirationCache[K, V]{
		expiration: expiration,
		items:      map[K]*item[V]{},
	}
}

// Cleanup removes all expired items from the cache.
func (c *ExpirationCache[K, V]) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	for k, v := range c.items {
		if !v.expires.IsZero() && v.expires.Before(now) {
			delete(c.items, k)
		}
	}
}

// Get fetches a value for a given key.
func (c *ExpirationCache[K, V]) Get(k K) (V, bool) {
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

// Set stores a value for a given key.
func (c *ExpirationCache[K, V]) Set(k K, v V) {
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

// SetWithExpiration stores a value for a given key with an explicit expiration.
func (c *ExpirationCache[K, V]) SetWithExpiration(k K, v V, expiration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[k] = &item[V]{
		expires: time.Now().Add(expiration),
		value:   v,
	}
}
