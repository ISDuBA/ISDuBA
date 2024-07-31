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
		sourcesSQL = `SELECT id, rate, slots FROM sources WHERE active`
		feedsSQL   = `SELECT f.id, sources_id, url, rolie ` +
			`FROM feeds f JOIN sources s ON f.sources_id = s.id ` +
			`WHERE active`
	)

	if err := m.db.Run(
		ctx,
		func(rctx context.Context, con *pgxpool.Conn) error {
			tx, err := con.BeginTx(rctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})
			if err != nil {
				return fmt.Errorf("starting transaction failed: %w", err)
			}
			defer tx.Rollback(rctx)

			// Collect active sources.
			srows, err := tx.Query(rctx, sourcesSQL)
			if err != nil {
				return fmt.Errorf("querying sources failed: %w", err)
			}
			m.sources, err = pgx.CollectRows(srows, func(row pgx.CollectableRow) (*activeSource, error) {
				var s activeSource
				return &s, row.Scan(&s.id, &s.rate, &s.slots)
			})
			if err != nil {
				return fmt.Errorf("collecting active sources failed: %w", err)
			}

			// Collect active feeds
			frows, err := tx.Query(rctx, feedsSQL)
			if err != nil {
				return fmt.Errorf("querying sources failed: %w", err)
			}
			m.feeds, err = pgx.CollectRows(frows, func(row pgx.CollectableRow) (*activeFeed, error) {
				var f activeFeed
				var sid int64
				var raw string
				if err := row.Scan(&f.id, &sid, &raw, &f.rolie); err != nil {
					return nil, err
				}
				parsed, err := url.Parse(raw)
				if err != nil {
					return nil, fmt.Errorf("invalid URL: %w", err)
				}
				f.url = parsed

				// Add to list of active feeds.
				idx := slices.IndexFunc(m.sources, func(s *activeSource) bool { return s.id == sid })
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
				return fmt.Errorf("collecting active feeds failed: %w", err)
			}

			return tx.Commit(rctx)
		}, 0,
	); err != nil {
		return fmt.Errorf("fetching active feeds failed: %w", err)
	}

	slog.Info("number of active sources", "num", len(m.sources))
	slog.Info("number of active feeds", "num", len(m.feeds))

	// Trigger a refresh of the loaded feeds.
	if len(m.feeds) > 0 {
		go func() { m.fns <- func(*Manager) {} }()
	}

	return nil
}
