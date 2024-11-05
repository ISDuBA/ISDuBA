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
	"bufio"
	"bytes"
	"context"
	"crypto/sha1"
	"log/slog"
	"sync"
	"time"

	"github.com/ISDuBA/ISDuBA/pkg/config"
	"github.com/ISDuBA/ISDuBA/pkg/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const maxPMDWorkers = 10

// Manager handles the refreshing of the aggregators.
type Manager struct {
	Cache *Cache

	done bool
	fns  chan func(*Manager)
	cfg  *config.Aggregators
	db   *database.DB
}

// NewManager creates a new aggregators manager.
func NewManager(cfg *config.Aggregators, db *database.DB) *Manager {
	return &Manager{
		Cache: newCache(cfg.Timeout),
		fns:   make(chan func(*Manager)),
		cfg:   cfg,
		db:    db,
	}
}

// Run runs the aggregators manager.
func (m *Manager) Run(ctx context.Context) {
	ticker := time.NewTicker(m.cfg.UpdateInterval)
	defer ticker.Stop()
	cacheTicker := time.NewTicker(holdingDuration)
	defer cacheTicker.Stop()
	for !m.done {
		select {
		case fn := <-m.fns:
			fn(m)
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.refresh(ctx)
		case <-cacheTicker.C:
			m.Cache.Cleanup()
		}
	}
}

func aggregatorChecksum(cagg *CachedAggregator) []byte {
	hash := sha1.New()
	w := bufio.NewWriter(hash)
	for _, url := range cagg.SourceURLs() {
		w.WriteString(url)
	}
	return hash.Sum(nil)
}

func (m *Manager) refresh(ctx context.Context) {
	type aggregator struct {
		id          int64
		url         string
		checksum    []byte
		newChecksum []byte
	}
	const (
		selectSQL = `SELECT id, url, checksum FROM aggregators`
		updateSQL = `UPDATE aggregators ` +
			`SET (checksum, checksum_updated) = ($1, $2) ` +
			`WHERE id = $3`
	)
	var aggregators []aggregator
	if err := m.db.Run(
		ctx,
		func(ctx context.Context, conn *pgxpool.Conn) error {
			rows, _ := conn.Query(ctx, selectSQL)
			var err error
			aggregators, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (aggregator, error) {
				var agg aggregator
				err := row.Scan(&agg.id, &agg.url, &agg.checksum)
				return agg, err
			})
			return err
		}, 0,
	); err != nil {
		slog.Error("fetching aggregators failed", "error", err)
		return
	}
	if len(aggregators) == 0 {
		return
	}
	var (
		toFetch    = make(chan int)
		numWorkers = min(maxPMDWorkers, len(aggregators))
		wg         sync.WaitGroup
	)
	fetch := func() {
		wg.Done()
		for index := range toFetch {
			cagg, err := m.Cache.GetAggregator(aggregators[index].url, false)
			if err != nil {
				slog.Warn("fetching aggregator failed", "err", err)
				continue
			}
			aggregators[index].newChecksum = aggregatorChecksum(cagg)
		}
	}
	for range numWorkers {
		wg.Add(1)
		go fetch()
	}
	for index := range aggregators {
		toFetch <- index
	}
	close(toFetch)
	wg.Wait()
	var (
		batch pgx.Batch
		now   = time.Now()
	)
	for i := range aggregators {
		agg := &aggregators[i]
		if !bytes.Equal(agg.checksum, agg.newChecksum) {
			batch.Queue(updateSQL, agg.newChecksum, now, agg.id)
		}
	}
	if batch.Len() == 0 {
		return
	}
	if err := m.db.Run(
		ctx,
		func(ctx context.Context, conn *pgxpool.Conn) error {
			tx, err := conn.Begin(ctx)
			if err != nil {
				return err
			}
			defer tx.Rollback(ctx)
			if err := tx.SendBatch(ctx, &batch).Close(); err != nil {
				return err
			}
			return tx.Commit(ctx)
		}, 0,
	); err != nil {
		slog.Error("fetching aggregators failed", "error", err)
	}
}

func (m *Manager) kill() { m.done = true }

// Kill shuts down the aggregators manager.
func (m *Manager) Kill() {
	m.fns <- (*Manager).kill
}
