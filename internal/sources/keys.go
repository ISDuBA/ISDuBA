// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package sources

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ISDuBA/ISDuBA/internal/cache"
	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

type keysCache struct {
	*cache.ExpirationCache[int64, *crypto.KeyRing]
}

func newKeysCache(expiration time.Duration) *keysCache {
	return &keysCache{
		ExpirationCache: cache.NewExpirationCache[int64, *crypto.KeyRing](expiration),
	}
}

// openPGPKeys extracts the OpenPGP key from them PMD of a source if not already
// in cache.
func (m *Manager) openPGPKeys(source *source) (*crypto.KeyRing, error) {
	if keys, ok := m.keysCache.Get(source.id); ok {
		return keys, nil
	}
	keys, _ := crypto.NewKeyRing(nil)
	cpmd := m.pmdCache.pmd(source.url, m.cfg)
	if !cpmd.Valid() {
		// Try again soon.
		m.keysCache.SetWithExpiration(source.id, keys, holdingPMDsDuration)
		return nil, fmt.Errorf("PMD of %q is invalid", source.url)
	}
	pmd, err := cpmd.Model()
	if err != nil {
		// Try again soon.
		m.keysCache.SetWithExpiration(source.id, keys, holdingPMDsDuration)
		return nil, fmt.Errorf("re-marshaling failed: %w", err)
	}
	base, err := url.Parse(source.url)
	if err != nil {
		// XXX: This should not happen.
		m.keysCache.SetWithExpiration(source.id, keys, holdingPMDsDuration)
		return nil, fmt.Errorf("invalid PMD url: %q", source.url)
	}
	client := source.httpClient(m)
	defer client.CloseIdleConnections()
	for i := range pmd.PGPKeys {
		key := &pmd.PGPKeys[i]
		if key.URL == nil {
			continue
		}
		u, err := url.Parse(*key.URL)
		if err != nil {
			slog.Warn("Invalid OpenPGP url", "url", *key.URL, "err", err)
			continue
		}
		if !u.IsAbs() {
			u = joinURL(base, u)
		}
		res, err := source.httpGet(client, m, u.String())
		if err != nil {
			slog.Warn(
				"Fetching public OpenPGP key failed",
				"url", u,
				"error", err)
			continue
		}
		if res.StatusCode != http.StatusOK {
			res.Body.Close()
			slog.Warn(
				"Fetching public OpenPGP key failed",
				"url", u,
				"status_code", res.StatusCode,
				"status", res.Status)
			continue
		}
		ckey, err := func() (*crypto.Key, error) {
			defer res.Body.Close()
			return crypto.NewKeyFromArmoredReader(res.Body)
		}()
		if err != nil {
			slog.Warn(
				"Reading public OpenPGP key failed",
				"url", u,
				"error", err)
			continue
		}
		if key.Fingerprint != "" &&
			!strings.EqualFold(ckey.GetFingerprint(), string(key.Fingerprint)) {
			slog.Warn(
				"Fingerprint of public OpenPGP key does not match remotely loaded",
				"url", u)
			continue
		}
		if err := keys.AddKey(ckey); err != nil {
			slog.Warn(
				"Could not add public OpenPGP key to key ring",
				"url", u)
		}
	}
	m.keysCache.Set(source.id, keys)
	return keys, nil
}

// loadSignature loads an ascii armored OpenPGP signature file from a given url.
func (s *source) loadSignature(client *http.Client, m *Manager, u *url.URL) (*crypto.PGPSignature, []byte, error) {
	resp, err := s.httpGet(client, m, u.String())
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, nil, fmt.Errorf(
			"fetching signature from %q failed: %s (%d)", u, resp.Status, resp.StatusCode)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	sign, err := crypto.NewPGPSignatureFromArmored(string(data))
	if err != nil {
		return nil, nil, err
	}
	return sign, data, nil
}
