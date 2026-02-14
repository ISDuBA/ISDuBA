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
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/ISDuBA/ISDuBA/pkg/config"
	"github.com/ISDuBA/ISDuBA/pkg/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// forwarderWakeupInterval is a saftey net wakeup interval for each
// forwarder if a ping from the manager is missed some how.
const forwarderWakeupInterval = 2 * time.Minute

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
		clientCert, err := tls.LoadX509KeyPair(
			cfg.ClientPublicCert,
			cfg.ClientPrivateCert)
		if err != nil {
			return nil, fmt.Errorf(
				"cannot load client cert for forward target %q: %w",
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
	ticker := time.NewTicker(forwarderWakeupInterval)
	defer ticker.Stop()
	// TODO: Implement me!
	for !f.done {
		if err := f.forward(ctx); err != nil {
			slog.Error("forwarder has issues", "error", err, "forwarder", f.cfg.URL)
		}
		select {
		case fn := <-f.fns:
			fn(f)
		case <-ctx.Done():
			return
		}
	}
}

func (f *forwarder) forward(ctx context.Context) error {
	const docIDsSQL = `` +
		`SELECT` +
		` documents_id ` +
		`FROM forwarders_queue fwq ` +
		`JOIN forwarders fw ON fwq.forwarders_id = fw.id ` +
		`WHERE` +
		` fwq.state IN ('pending', 'failed')` +
		` AND fw.url = $1 ` +
		`ORDER BY documents_id DESC ` + // Assuming older docs have lower ids.
		`LIMIT 20` // Poll in smaller batches.

	for {
		var docIDs []int64
		if err := f.db.Run(
			ctx,
			func(rctx context.Context, conn *pgxpool.Conn) error {
				rows, _ := conn.Query(rctx, docIDsSQL, f.cfg.URL)
				var err error
				docIDs, err = pgx.CollectRows(
					rows,
					func(row pgx.CollectableRow) (int64, error) {
						var docID int64
						err := row.Scan(&docID)
						return docID, err
					})
				return err
			}, 0,
		); err != nil {
			return fmt.Errorf(
				"scanning for pending/failed documents failed: %w", err)
		}
		if len(docIDs) == 0 {
			break
		}
		// TODO: Load docs, forward them and store the result.
	}
	return nil
}

func (f *forwarder) kill() {
	f.fns <- func(f *forwarder) { f.done = true }
}

// ping should be called by the manager to signal the fowarder
// that there are documents to forward.
func (f *forwarder) ping() {
	f.fns <- func(f *forwarder) {}
}

func (f *forwarder) acceptsPublisher(publisher string) bool {
	return f.cfg.Publisher == nil || publisher == *f.cfg.Publisher
}
