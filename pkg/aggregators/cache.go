// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package aggregators

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/ISDuBA/ISDuBA/internal/cache"
	"github.com/ISDuBA/ISDuBA/pkg/sources"
	"github.com/gocsaf/csaf/v3/csaf"
)

// holdingDuration is the duration how long PMDs are cached.
const holdingDuration = time.Minute * 15

// CachedAggregator are cached aggregators.
type CachedAggregator struct {
	Raw        []byte
	Aggregator *csaf.Aggregator
}

// Cache is cache of aggregators.
type Cache struct {
	*cache.ExpirationCache[string, *CachedAggregator]
	timeout time.Duration
}

func newCache(timeout time.Duration) *Cache {
	return &Cache{
		ExpirationCache: cache.NewExpirationCache[string, *CachedAggregator](holdingDuration),
		timeout:         timeout,
	}
}

// GetAggregator fetches a cached aggregator.
func (c *Cache) GetAggregator(url string) (*CachedAggregator, error) {
	if !strings.HasSuffix(url, "/aggregator.json") {
		return nil, errors.New("invalid aggregator url")
	}
	if ca, ok := c.Get(url); ok {
		return ca, nil
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", sources.UserAgent)
	client := &http.Client{}
	if c.timeout > 0 {
		client.Timeout = c.timeout
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code %s (%d)", resp.Status, resp.StatusCode)
	}
	var data bytes.Buffer
	var doc any
	r := io.TeeReader(resp.Body, &data)
	if err := json.NewDecoder(r).Decode(&doc); err != nil {
		return nil, fmt.Errorf("invalid json: %w", err)
	}
	msgs, err := csaf.ValidateAggregator(doc)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	if len(msgs) > 0 {
		return nil, fmt.Errorf("validation failed: %s", strings.Join(msgs, ", "))
	}
	raw := bytes.Clone(data.Bytes())
	agg := new(csaf.Aggregator)
	if err := json.Unmarshal(raw, agg); err != nil {
		return nil, fmt.Errorf("cannot unmarshal aggregator: %w", err)
	}
	ca := &CachedAggregator{
		Raw:        raw,
		Aggregator: agg,
	}
	c.Set(url, ca)
	return ca, nil
}

// SourceURLs extracts the source URLs from the cached aggregator.
func (ca *CachedAggregator) SourceURLs() []string {
	var urls []string
	unique := func(url string) {
		if !slices.Contains(urls, url) {
			urls = append(urls, url)
		}
	}
	add := func(metadata *csaf.AggregatorCSAFProviderMetadata) {
		if metadata != nil {
			if url := metadata.URL; url != nil {
				unique(string(*url))
			}
		}
	}
	for _, provider := range ca.Aggregator.CSAFProviders {
		if provider != nil {
			add(provider.Metadata)
		}
	}
	for _, publisher := range ca.Aggregator.CSAFPublishers {
		if publisher != nil {
			add(publisher.Metadata)
			for _, m := range publisher.Mirrors {
				unique(string(m))
			}
		}
	}
	return urls
}
