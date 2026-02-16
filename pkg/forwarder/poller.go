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
	"slices"
	"time"
)

type (
	changedAdvisory struct {
		when      time.Time
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
	// TODO: Implement me!
	_ = ctx
	if len(p.changes) > 0 {
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
			b.when.Compare(a.when),  // descending in age
			cmp.Compare(a.id, b.id), // ascending in id
		)
	})
	return caids
}
