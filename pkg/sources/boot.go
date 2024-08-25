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

	"github.com/ISDuBA/ISDuBA/pkg/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Boot loads the sources from database.
func (m *Manager) Boot(ctx context.Context) error {
	const (
		sourcesSQL = `SELECT id, name, url, rate, slots, active, headers, ` +
			`strict_mode, insecure, signature_check, age, ignore_patterns, ` +
			`client_cert_public, client_cert_private, client_cert_passphrase ` +
			`FROM sources`
		feedsSQL = `SELECT id, label, sources_id, url, rolie, log_lvl::text FROM feeds`
	)
	if err := m.db.Run(
		ctx,
		func(rctx context.Context, con *pgxpool.Conn) error {
			tx, err := con.Begin(rctx)
			if err != nil {
				return fmt.Errorf("starting transaction failed: %w", err)
			}
			defer tx.Rollback(rctx)
			// Collect sources.
			srows, err := tx.Query(rctx, sourcesSQL)
			if err != nil {
				return fmt.Errorf("querying sources failed: %w", err)
			}
			var bads []int64
			m.sources, err = pgx.CollectRows(srows, func(row pgx.CollectableRow) (*source, error) {
				var (
					s                                       source
					patterns                                []string
					clientCertPrivate, clientCertPassphrase []byte
				)
				if err := row.Scan(
					&s.id, &s.name, &s.url, &s.rate, &s.slots, &s.active, &s.headers,
					&s.strictMode, &s.insecure, &s.signatureCheck, &s.age, &patterns,
					&s.clientCertPublic, &clientCertPrivate, &clientCertPassphrase,
				); err != nil {
					return nil, err
				}
				regexps, err := AsRegexps(patterns)
				if err != nil {
					return nil, err
				}
				s.ignorePatterns = regexps

				var bad bool
				if s.clientCertPrivate, err = m.decrypt(clientCertPrivate); err != nil {
					bad = true
				}
				if s.clientCertPassphrase, err = m.decrypt(clientCertPassphrase); err != nil {
					bad = true
				}
				if !bad {
					if err := s.updateCertificate(); err != nil {
						bad = true
					}
				}
				if bad && s.active {
					s.active = false
					bads = append(bads, s.id)
				}
				return &s, nil
			})
			if err != nil {
				return fmt.Errorf("collecting sources failed: %w", err)
			}
			// If we have sources with bad crypto deactivate these.
			if len(bads) > 0 {
				const deactivateSQL = `UPDATE sources SET active = FALSE WHERE id = $1`
				batch := &pgx.Batch{}
				for _, id := range bads {
					batch.Queue(deactivateSQL, id)
				}
				if err := tx.SendBatch(ctx, batch).Close(); err != nil {
					return fmt.Errorf("deactivating bad sources failed: %w", err)
				}
			}

			// Collect feeds.
			frows, err := tx.Query(rctx, feedsSQL)
			if err != nil {
				return fmt.Errorf("querying sources failed: %w", err)
			}
			defer frows.Close()
			for frows.Next() {
				var (
					f        feed
					sid      int64
					raw      string
					logLevel config.FeedLogLevel
				)
				if err := frows.Scan(
					&f.id,
					&f.label,
					&sid,
					&raw,
					&f.rolie,
					&logLevel,
				); err != nil {
					return err
				}
				parsed, err := url.Parse(raw)
				if err != nil {
					return fmt.Errorf("invalid URL: %w", err)
				}
				f.url = parsed
				f.logLevel.Store(int32(logLevel))
				// Add to list of active feeds.
				s := m.findSourceByID(sid)
				if s == nil {
					// Should really not happen! Considering a panic.
					return fmt.Errorf("cannot find source id %d", sid)
				}
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
