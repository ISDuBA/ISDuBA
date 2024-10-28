// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

// Package aggregators handles the refreshing of the managed aggregators.
package aggregators

import (
	"context"
	"time"

	"github.com/ISDuBA/ISDuBA/pkg/config"
)

// Manager handles the refreshing of the aggregators.
type Manager struct {
	Cache *Cache

	done bool
	fns  chan func(*Manager)
	cfg  *config.Aggregators
}

// NewManager creates a new aggregators manager.
func NewManager(cfg *config.Aggregators) *Manager {
	return &Manager{
		Cache: newCache(),
		fns:   make(chan func(*Manager)),
		cfg:   cfg,
	}
}

// Run runs the aggregators manager.
func (m *Manager) Run(ctx context.Context) {
	ticker := time.NewTicker(m.cfg.UpdateInterval)
	defer ticker.Stop()
	for !m.done {
		select {
		case fn := <-m.fns:
			fn(m)
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.refresh(ctx)
		}
	}
}

func (m *Manager) refresh(ctx context.Context) {
	// TODO: Implement me!
	_ = ctx
}

func (m *Manager) kill() { m.done = true }

// Kill shuts down the aggregators manager.
func (m *Manager) Kill() {
	m.fns <- (*Manager).kill
}
