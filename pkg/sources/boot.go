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
	"net/url"
	"slices"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Boot loads the sources from database.
func (m *Manager) Boot(ctx context.Context) error {
	const (
		sourcesSQL = `SELECT id, rate, slots, active FROM sources`
		feedsSQL   = `SELECT id, sources_id, url, rolie, log_lvl::text FROM feeds`
	)
	if err := m.db.Run(
		ctx,
		func(rctx context.Context, con *pgxpool.Conn) error {
			tx, err := con.BeginTx(rctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})
			if err != nil {
				return fmt.Errorf("starting transaction failed: %w", err)
			}
			defer tx.Rollback(rctx)
			// Collect sources.
			srows, err := tx.Query(rctx, sourcesSQL)
			if err != nil {
				return fmt.Errorf("querying sources failed: %w", err)
			}
			m.sources, err = pgx.CollectRows(srows, func(row pgx.CollectableRow) (*source, error) {
				var s source
				return &s, row.Scan(&s.id, &s.rate, &s.slots, &s.active)
			})
			if err != nil {
				return fmt.Errorf("collecting sources failed: %w", err)
			}
			// Collect feeds.
			frows, err := tx.Query(rctx, feedsSQL)
			if err != nil {
				return fmt.Errorf("querying sources failed: %w", err)
			}
			m.feeds, err = pgx.CollectRows(frows, func(row pgx.CollectableRow) (*feed, error) {
				var (
					f   feed
					sid int64
					raw string
				)
				if err := row.Scan(&f.id, &sid, &raw, &f.rolie, &f.logLevel); err != nil {
					return nil, err
				}
				parsed, err := url.Parse(raw)
				if err != nil {
					return nil, fmt.Errorf("invalid URL: %w", err)
				}
				f.url = parsed
				// Add to list of active feeds.
				idx := slices.IndexFunc(m.sources, func(s *source) bool { return s.id == sid })
				if idx == -1 {
					// Should really not happen! Considering a panic.
					return nil, fmt.Errorf("cannot find source id %d", sid)
				}
				src := m.sources[idx]
				src.feeds = append(src.feeds, &f)
				f.source = src
				return &f, nil
			})
			if err != nil {
				return fmt.Errorf("collecting feeds failed: %w", err)
			}
			return tx.Commit(rctx)
		}, 0,
	); err != nil {
		return err
	}

	slog.Info("number of sources", "num", len(m.sources))
	slog.Info("number of feeds", "num", len(m.feeds))

	// Trigger a refresh of the loaded feeds.
	if len(m.feeds) > 0 {
		go func() { m.fns <- (*Manager).ping }()
	}

	return nil
}
