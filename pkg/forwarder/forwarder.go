// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

package forwarder

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"github.com/ISDuBA/ISDuBA/pkg/config"
	"github.com/ISDuBA/ISDuBA/pkg/database"
)

type forwarder struct {
	cfg     *config.ForwardTarget
	db      *database.DB
	fns     chan (func(*forwarder))
	done    bool
	client  *http.Client
	headers http.Header
}

func newForwarder(cfg *config.ForwardTarget, db *database.DB) (*forwarder, error) {
	// Init http clients
	var tlsConfig tls.Config
	if cfg.ClientPrivateCert != "" && cfg.ClientPublicCert != "" {
		clientCert, err := tls.LoadX509KeyPair(cfg.ClientPublicCert, cfg.ClientPrivateCert)
		if err != nil {
			return nil, fmt.Errorf("cannot load client cert for forward target %q: %w",
				cfg.URL, err)
		}
		tlsConfig.Certificates = []tls.Certificate{clientCert}
	}
	client := &http.Client{
		Timeout: cfg.Timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tlsConfig,
		},
	}
	headers := make(http.Header, len(cfg.Header))
	for _, header := range cfg.Header {
		if k, v, ok := strings.Cut(header, ":"); ok {
			headers.Add(k, v)
			continue
		}
		return nil, fmt.Errorf(
			"header %q of forwarder target %q is missing ':'",
			header, cfg.URL)
	}
	return &forwarder{
		cfg:     cfg,
		db:      db,
		fns:     make(chan func(*forwarder)),
		client:  client,
		headers: headers,
	}, nil
}

func (f *forwarder) run(ctx context.Context) {
	// TODO: Implement me!
	for !f.done {
		select {
		case fn := <-f.fns:
			fn(f)
		case <-ctx.Done():
			return
		}
	}
}

func (f *forwarder) kill() {
	f.fns <- func(f *forwarder) { f.done = true }
}
