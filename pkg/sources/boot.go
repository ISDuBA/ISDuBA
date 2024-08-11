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
		sourcesSQL = `SELECT id, name, url, rate, slots, active FROM sources`
		feedsSQL   = `SELECT id, label, sources_id, url, rolie, log_lvl::text FROM feeds`
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
				return &s, row.Scan(
					&s.id,
					&s.name,
					&s.url,
					&s.rate,
					&s.slots,
					&s.active)
			})
			if err != nil {
				return fmt.Errorf("collecting sources failed: %w", err)
			}
			// Collect feeds.
			frows, err := tx.Query(rctx, feedsSQL)
			if err != nil {
				return fmt.Errorf("querying sources failed: %w", err)
			}
			defer frows.Close()
			for frows.Next() {
				var (
					f   feed
					sid int64
					raw string
				)
				if err := frows.Scan(
					&f.id,
					&f.label,
					&sid,
					&raw,
					&f.rolie,
					&f.logLevel,
				); err != nil {
					return err
				}
				parsed, err := url.Parse(raw)
				if err != nil {
					return fmt.Errorf("invalid URL: %w", err)
				}
				f.url = parsed
				// Add to list of active feeds.
				idx := slices.IndexFunc(m.sources, func(s *source) bool { return s.id == sid })
				if idx == -1 {
					// Should really not happen! Considering a panic.
					return fmt.Errorf("cannot find source id %d", sid)
				}
				s := m.sources[idx]
				s.feeds = append(s.feeds, &f)
				f.source = s
			}
			if err := frows.Err(); err != nil {
				return fmt.Errorf("collecting feeds failed: %w", err)
			}
			return tx.Commit(rctx)
		}, 0,
	); err != nil {
		return err
	}

	activeFeeds := m.numActiveFeeds()

	slog.Info("number of sources", "num", len(m.sources))
	slog.Info("number of active feeds", "num", activeFeeds)

	// Trigger a refresh of the loaded feeds.
	if activeFeeds > 0 {
		m.backgroundPing()
	}
	return nil
}
