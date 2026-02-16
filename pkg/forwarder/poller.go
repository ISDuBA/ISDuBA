// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024, 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024, 2026 Intevation GmbH <https://intevation.de>

package forwarder

import (
	"cmp"
	"context"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	changedAdvisory struct {
		first     time.Time
		publisher string
	}
	changedAdvisories map[int64]*changedAdvisory
	changedIDAdvisory struct {
		*changedAdvisory
		id int64
	}
)

type poller struct {
	manager *ForwardManager
	done    bool
	fns     chan func(*poller)
	changes changedAdvisories
	latest  time.Time
}

func newPoller(manager *ForwardManager) *poller {
	return &poller{
		manager: manager,
		fns:     make(chan func(*poller)),
		changes: changedAdvisories{},
	}
}

func (p *poller) run(ctx context.Context) {
	ticker := time.NewTicker(p.manager.cfg.UpdateInterval)
	defer ticker.Stop()
	for !p.done {
		select {
		case fn := <-p.fns:
			fn(p)
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.poll(ctx)
		}
	}
}

func (p *poller) kill() {
	p.fns <- func(p *poller) { p.done = true }
}

func (p *poller) poll(ctx context.Context) {
	const recentSQL = `` +
		`SELECT` +
		` ads.id AS id,` +
		` publisher,` +
		` recent ` +
		`FROM advisories ads` +
		` JOIN documents docs ON docs.advisories_id = ads.id ` +
		` JOIN events_log el  ON docs.id            = el.documents_id ` +
		`WHERE` +
		` ads.recent >= $1 AND el.event = 'import_document`

	if err := p.manager.db.Run(
		ctx,
		func(rctx context.Context, conn *pgxpool.Conn) error {
			rows, err := conn.Query(rctx, recentSQL, p.latest)
			if err != nil {
				return fmt.Errorf("poll: loading latest failed: %w", err)
			}
			defer rows.Close()
			for rows.Next() {
				var (
					id int64
					ca changedAdvisory
				)
				if err := rows.Scan(
					&id,
					&ca.publisher,
					&ca.first,
				); err != nil {
					return fmt.Errorf("poll: scaning rows failed: %w", err)
				}
				ca.first = ca.first.UTC()
				// To avoid poll all the advisories in the round again
				// we remember where we where the last time.
				// This should only deliver many results the first time.
				// Subsequent polls should only deliver what have
				// changed since.
				if ca.first.After(p.latest) {
					p.latest = ca.first
				}
				if old := p.changes[id]; old == nil {
					p.changes[id] = &ca
				} else if ca.first.Before(old.first) {
					old.first = ca.first
				}
			}
			return nil
		}, 0,
	); err != nil {
		// If there is a DB error it might be no a good idea
		// to directly inform the manager about found changes
		// as this would lead to more db operations likely to fail.
		// Instead do so the next time we wake up.
		// The db will be up then again, hopefully.
		slog.Error("forwarder", "error", err)
		return
	}

	// Try to deliver changes to manager.
	if len(p.changes) > 0 {
		// This a none blocking call.
		// If the manager is currently filling
		// the forwarder queues we will send the
		// the changes when it is not busy.
		p.changes = p.manager.changesDetected(p.changes)
	}
}

func (cas changedAdvisories) order() []changedIDAdvisory {
	caids := make([]changedIDAdvisory, 0, len(cas))
	for id, advisory := range cas {
		caids = append(caids, changedIDAdvisory{changedAdvisory: advisory, id: id})
	}
	slices.SortFunc(caids, func(a, b changedIDAdvisory) int {
		return cmp.Or(
			b.first.Compare(a.first), // descending in age
			cmp.Compare(a.id, b.id),  // ascending in id
		)
	})
	return caids
}
