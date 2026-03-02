// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024, 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024, 2026 Intevation GmbH <https://intevation.de>

// Package forwarder implements the document forwarder.
package forwarder

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"slices"
	"strconv"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ISDuBA/ISDuBA/pkg/config"
	"github.com/ISDuBA/ISDuBA/pkg/database"
)

// Manager forwards documents to specified targets.
type Manager struct {
	cfg        *config.Forwarder
	db         *database.DB
	fns        chan func(*Manager)
	done       bool
	forwarders []*forwarder
	changes    changedAdvisories
}

type (
	versionInfo struct {
		id            int64
		version       string
		status        trackingStatus
		historyLength int
		current       *time.Time
		initial       *time.Time
	}
	versionInfos []versionInfo
)

type trackingStatus int

const (
	unknownStatus = trackingStatus(iota)
	draftStatus
	interimStatus
	finalStatus
)

func parseTrackingStatus(s *string) trackingStatus {
	if s == nil {
		return unknownStatus
	}
	switch *s {
	case "final":
		return finalStatus
	case "interim":
		return interimStatus
	case "draft":
		return draftStatus
	default:
		return unknownStatus
	}
}

// NewManager creates a new forward manager.
func NewManager(
	cfg *config.Config,
	db *database.DB,
) (*Manager, error) {
	// TODO: Move this parsing to config.
	var extURL *url.URL
	if cfg.Web.ExternalURL != "" {
		eu, err := url.Parse(cfg.Web.ExternalURL)
		if err != nil {
			return nil, fmt.Errorf("external URL is invalid: %w", err)
		}
		extURL = eu
	}
	fwdCfg := &cfg.Forwarder
	forwarders := make([]*forwarder, 0, len(fwdCfg.Targets))
	for i := range fwdCfg.Targets {
		tcfg := &fwdCfg.Targets[i]
		forwarder, err := newForwarder(tcfg, extURL, db)
		if err != nil {
			return nil,
				fmt.Errorf("create automatic forwarder for %q failed: %w",
					tcfg.URL, err)
		}
		forwarders = append(forwarders, forwarder)
	}
	return &Manager{
		cfg:        fwdCfg,
		db:         db,
		fns:        make(chan func(manager *Manager)),
		forwarders: forwarders,
	}, nil
}

// Run runs the forward manager. To be used in a Go routine.
func (fm *Manager) Run(ctx context.Context) {
	hasAutomatic := false
	// Start the automatic forwarders.
	for _, forwarder := range fm.forwarders {
		if forwarder.cfg.Automatic {
			if err := fm.createForwarder(ctx, forwarder.cfg.URL); err != nil {
				slog.Error("forwarder", "error", err)
			} else {
				hasAutomatic = true
				go forwarder.run(ctx)
				defer forwarder.kill()
			}
		}
	}
	// No need to poll if there are no automatic forwarders.
	if !hasAutomatic {
		for !fm.done {
			select {
			case fn := <-fm.fns:
				fn(fm)
			case <-ctx.Done():
				return
			}
		}
		return
	}
	// Start the poller
	poller := newPoller(fm)
	go poller.run(ctx)
	defer poller.kill()
	// The poller should wake us up but in case wake up
	// on our own timer based.
	ticker := time.NewTicker(fm.cfg.UpdateInterval / 2)
	defer ticker.Stop()
	for !fm.done {
		fm.fillForwarderQueues(ctx)
		select {
		case fn := <-fm.fns:
			fn(fm)
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

// createForwarder ensure the existence of a forwarder in the forwarder
// lookup table.
func (fm *Manager) createForwarder(ctx context.Context, url string) error {
	const insertForwarderSQL = `` +
		`INSERT INTO forwarders (url) VALUES ($1) ` +
		`ON CONFLICT (url) DO NOTHING`
	if err := fm.db.Run(
		ctx,
		func(rctx context.Context, conn *pgxpool.Conn) error {
			_, err := conn.Exec(rctx, insertForwarderSQL, url)
			return err
		}, 0,
	); err != nil {
		return fmt.Errorf("inserting forwarder %q failed: %w", url, err)
	}
	return nil
}

// changesDetected tries to deliver detected advisory changes to
// the manager. If the manager is ready a fresh changedAdvisories map
// is returned. If the delivery would block the given map is
// returned so that the poller can go on detecting avoiding duplicates.
func (fm *Manager) changesDetected(changes changedAdvisories) changedAdvisories {
	select {
	case fm.fns <- func(fm *Manager) { fm.changes = changes }:
		return changedAdvisories{}
	default:
		return changes
	}
}

var (
	filterIndex = map[config.ForwarderStrategy]int{
		config.ForwarderStrategyAll:         0,
		config.ForwarderStrategyNewAndMajor: 1,
	}
	filters = [2]func(versionInfos) []int{
		versionInfos.filterAll,
		versionInfos.filterNewAndMajor,
	}
)

func (vis versionInfos) filterAll() []int {
	// Index all.
	indices := make([]int, len(vis))
	for i := range indices {
		indices[i] = i
	}
	return indices
}

func (vis versionInfos) filterNewAndMajor() []int {
	// TODO: Implement more cases.
	var (
		lastSemVer *semver.Version
		indices    = make([]int, 0, len(vis))
	)
	if len(vis) == 0 {
		return indices
	}
	// The first one is always important.
	if _, err := strconv.Atoi(vis[0].version); err == nil {
		// Even if its a draft include it.
		indices = append(indices, 0)
	} else if sv, err := semver.NewVersion(vis[0].version); err == nil {
		lastSemVer = sv
		indices = append(indices, 0)
	}

	// Handle the rest.
	for i := range vis[1:] {
		vi := &vis[i]
		// All versions with numbers should be forwarded
		// if they are not in draft status.
		if _, err := strconv.Atoi(vi.version); err == nil {
			if vi.status != draftStatus {
				indices = append(indices, i)
			}
			continue
		}
		sv, err := semver.NewVersion(vi.version)
		if err != nil {
			// It could only be integer or SemVer.
			// Ignore version which are none of these.
			continue
		}
		// If the major part increases its important.
		if lastSemVer == nil || sv.Major() > lastSemVer.Major() {
			indices = append(indices, i)
		}
		lastSemVer = sv
	}
	return indices
}

// fillForwarderQueues takes the advisory changes aggregated by the poller
func (fm *Manager) fillForwarderQueues(ctx context.Context) {
	if len(fm.changes) == 0 {
		return
	}
	ordered := fm.changes.order()
	fm.changes = nil

	pings := map[*forwarder]struct{}{}
	defer func() {
		// Notify forwarders that have new jobs.
		for fw := range pings {
			fw.ping()
		}
	}()
	// Do not recalculate indices when having more than forwarder
	// with same strategy.
	indicesCache := make([][]int, len(filterIndex))
	if err := fm.db.Run(
		ctx,
		func(rctx context.Context, conn *pgxpool.Conn) error {
			for ; len(ordered) > 0; ordered = ordered[1:] {
				adv := &ordered[0]
				// Ignore advisories no forwarder is interested in.
				if !slices.ContainsFunc(fm.forwarders, func(fw *forwarder) bool {
					return fw.cfg.Automatic && fw.acceptsPublisher(adv.publisher)
				}) {
					continue
				}
				vis, err := loadVersionInfos(rctx, conn, adv.id)
				if err != nil {
					return err
				}
				clear(indicesCache)
				for _, fw := range fm.forwarders {
					if !fw.cfg.Automatic || !fw.acceptsPublisher(adv.publisher) {
						continue
					}
					strategy := fm.cfg.Strategy
					if fw.cfg.Strategy != nil {
						strategy = *fw.cfg.Strategy
					}
					fi := filterIndex[strategy]
					cachedIndices := indicesCache[fi]
					if indicesCache[fi] == nil {
						cachedIndices = filters[fi](vis)
						indicesCache[fi] = cachedIndices
					}
					// Nothing to do.
					if len(cachedIndices) == 0 {
						continue
					}
					if err := storeIndicesInQueue(
						ctx, conn,
						vis, cachedIndices,
						fw.cfg.URL,
					); err != nil {
						return err
					}
					// Forwarder needs a ping afterwards.
					pings[fw] = struct{}{}
				}
			}
			return nil
		}, 0,
	); err != nil {
		// Store the remaining unhandled changes back for later.
		fm.changes = ordered.changes()
		slog.Error("forwarder", "error", err)
	}
}

func storeIndicesInQueue(
	ctx context.Context,
	conn *pgxpool.Conn,
	vis []versionInfo,
	indices []int,
	url string,
) error {
	const upsertSQL = `` +
		`INSERT INTO forwarders_queue` +
		` (forwarders_id, documents_id) ` +
		`SELECT id, $1` +
		` FROM forwarders` +
		` WHERE url = $2` +
		` ON CONFLICT (forwarders_id, documents_id) DO NOTHING`
	// XXX: Maybe using batches here is a bit to aggressive?!
	batch := &pgx.Batch{}
	for _, idx := range indices {
		batch.Queue(upsertSQL, vis[idx].id, url)
	}
	if err := conn.SendBatch(ctx, batch).Close(); err != nil {
		return fmt.Errorf(
			"sending documents to queue failed: %w", err)
	}
	return nil
}

func loadVersionInfos(
	ctx context.Context,
	conn *pgxpool.Conn,
	advisoryID int64,
) (versionInfos, error) {
	const versionSQL = `` +
		`SELECT` +
		` id,` +
		` version,` +
		` tracking_status::text,` +
		` coalesce(rev_history_length, 0),` +
		` current_release_date,` +
		` initial_release_date ` +
		`FROM documents ` +
		`WHERE` +
		` advisories_id = $1`
	rows, _ := conn.Query(ctx, versionSQL, advisoryID)
	vs, err := pgx.CollectRows(
		rows,
		func(row pgx.CollectableRow) (versionInfo, error) {
			var vi versionInfo
			var status *string
			err := row.Scan(
				&vi.id,
				&vi.version,
				&status,
				&vi.historyLength,
				&vi.current,
				&vi.initial)
			vi.status = parseTrackingStatus(status)
			if vi.current != nil {
				*vi.current = vi.current.UTC()
			}
			if vi.initial != nil {
				*vi.initial = vi.initial.UTC()
			}
			return vi, err
		})
	if err != nil {
		return nil, fmt.Errorf("loading version infos failed: %w", err)
	}
	vis := versionInfos(vs)
	vis.order()
	return vis, nil
}

func compare[T interface{ Compare(T) int }](a, b *T) int {
	switch {
	case a == nil && b == nil:
		return 0
	case a == nil:
		return +1
	case b == nil:
		return -1
	default:
		return (*a).Compare(*b)
	}
}

func (vis versionInfos) order() {
	slices.SortFunc(vis, func(a, b versionInfo) int {
		return cmp.Or(
			compare(b.initial, a.initial),
			compare(b.current, a.current),
			cmp.Compare(b.status, a.status),
			cmp.Compare(b.historyLength, a.historyLength),
		)
	})
}

// ForwardTarget contains information about the available target.
type ForwardTarget struct {
	URL  string `json:"url"`
	Name string `json:"name,omitempty"`
	ID   int    `json:"id"`
}

// Targets returns a list of forward targets.
func (fm *Manager) Targets() []ForwardTarget {
	result := make(chan []ForwardTarget)
	fm.fns <- func(fm *Manager) {
		forwarders := make([]ForwardTarget, 0, len(fm.forwarders))
		for i, forwarder := range fm.forwarders {
			if !forwarder.cfg.Automatic {
				forwarders = append(
					forwarders, ForwardTarget{
						ID:   i,
						URL:  forwarder.cfg.URL,
						Name: forwarder.cfg.Name,
					})
			}
		}
		result <- forwarders
	}
	return <-result
}

// ForwardDocument sends the document to the specified target.
func (fm *Manager) ForwardDocument(ctx context.Context, targetID int, docID int64) error {
	result := make(chan error)
	fm.fns <- func(fm *Manager) {
		if targetID < 0 || targetID >= len(fm.forwarders) || fm.forwarders[targetID].cfg.Automatic {
			result <- errors.New("could not find target with specified id")
			return
		}
		result <- fm.forwarders[targetID].forwardDocument(ctx, docID)
	}
	return <-result
}

// Kill shuts down the forward manager.
func (fm *Manager) Kill() {
	fm.fns <- func(fm *Manager) { fm.done = true }
}
