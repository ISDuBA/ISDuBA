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
	"bytes"
	"cmp"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/ISDuBA/ISDuBA/pkg/database/query"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ISDuBA/ISDuBA/pkg/config"
	"github.com/ISDuBA/ISDuBA/pkg/database"
)

// validationStatus represents the validation status
// known to the HTTP endpoint.
type validationStatus string

const (
	validValidationStatus        = validationStatus("valid")
	invalidValidationStatus      = validationStatus("invalid")
	notValidatedValidationStatus = validationStatus("not_validated")
)

type document struct {
	data      string
	filename  string
	valStatus validationStatus
}

// ForwardManager forwards documents to specified targets.
type ForwardManager struct {
	cfg        *config.Forwarder
	db         *database.DB
	fns        chan func(*ForwardManager)
	done       bool
	forwarders []*forwarder
	changes    changedAdvisories
}

type (
	trackingStatus int

	versionInfo struct {
		id            int64
		version       string
		status        trackingStatus
		historyLength int
		current       *time.Time
		initial       *time.Time
	}
)

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

// NewForwardManager creates a new forward manager.
func NewForwardManager(
	cfg *config.Forwarder,
	db *database.DB,
) (*ForwardManager, error) {
	forwarders := make([]*forwarder, 0, len(cfg.Targets))
	for i := range cfg.Targets {
		tcfg := &cfg.Targets[i]
		forwarder, err := newForwarder(tcfg, db)
		if err != nil {
			return nil,
				fmt.Errorf("create automatic forwarder for %q failed: %w",
					tcfg.URL, err)
		}
		forwarders = append(forwarders, forwarder)
	}
	return &ForwardManager{
		cfg:        cfg,
		db:         db,
		fns:        make(chan func(manager *ForwardManager)),
		forwarders: forwarders,
	}, nil
}

// Run runs the forward manager. To be used in a Go routine.
func (fm *ForwardManager) Run(ctx context.Context) {
	hasAutomatic := false
	// Start the automatic forwarders.
	for _, forwarder := range fm.forwarders {
		if forwarder.cfg.Automatic {
			hasAutomatic = true
			go forwarder.run(ctx)
			defer forwarder.kill()
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

// changesDetected tries to deliver detected advisory changes to
// the manager. If the manager is ready a fresh changedAdvisories map
// is returned. If the delivery would block the given map is
// returned so that the poller can go on detecting avoiding duplicates.
func (fm *ForwardManager) changesDetected(changes changedAdvisories) changedAdvisories {
	select {
	case fm.fns <- func(fm *ForwardManager) { fm.changes = changes }:
		return changedAdvisories{}
	default:
		return changes
	}
}

var (
	filterIndex = map[config.ForwarderStrategy]int{
		config.ForwarderStrategyAll:       0,
		config.ForwarderStrategyImportant: 1,
	}
	filters = [2]func([]versionInfo) []int{
		filterAll,
		filterImportant,
	}
)

func filterAll(vis []versionInfo) []int {
	indices := make([]int, len(vis))
	for i := range indices {
		indices[i] = i
	}
	return indices
}

func filterImportant(vis []versionInfo) []int {
	// TODO: Implement me!
	return filterAll(vis)
}

// fillForwarderQueues takes the advisory changes aggregated by the poller
func (fm *ForwardManager) fillForwarderQueues(ctx context.Context) {
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
					if indicesCache[fi] == nil {
						indicesCache[fi] = filters[fi](vis)
					}
					if err := storeIndicesInQueue(
						ctx, conn,
						vis, indicesCache[fi],
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
) ([]versionInfo, error) {
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
	vis, err := pgx.CollectRows(
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
	orderVersionInfos(vis)
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

func orderVersionInfos(vis []versionInfo) {
	slices.SortFunc(vis, func(a, b versionInfo) int {
		return cmp.Or(
			compare(a.initial, b.initial),
			compare(a.current, b.current),
			cmp.Compare(a.status, b.status),
			cmp.Compare(a.historyLength, b.historyLength),
		)
	})
}

func (fm *ForwardManager) uploadDocuments(ctx context.Context, target *forwarder, documentIDs []int64) {
	/*
		defer target.running.Unlock()
		for _, documentID := range documentIDs {
			document, err := fm.loadDocument(ctx, documentID)
			if err != nil {
				slog.Error("could not load document to forward", "err", err)
				continue
			}
			if err := fm.uploadDocument(ctx, document, documentID, target); err != nil {
				slog.Error("could not forward document", "err", err)
			}
		}
	*/
}

func (fm *ForwardManager) uploadDocument(ctx context.Context, doc *document, documentID int64, target *forwarder) error {
	/*
		req, err := fm.buildRequest(doc.data, doc.filename, doc.valStatus, target)
		if err != nil {
			slog.Error("building forward request failed", "err", err)
			return err
		}
		res, err := target.client.Do(req)
		if err != nil {
			slog.Error("sending forward request failed", "err", err)
			return err
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusCreated {
			slog.Error("forwarding failed", "status", res.StatusCode)
			body, err := io.ReadAll(res.Body)
			if err == nil {
				slog.Error("forward request failed", "response", string(body))
			}
			return errors.New("forwarding failed " + res.Status)
		}
		fm.logDocument(ctx, target.url, documentID)
	*/
	return nil
}

var escapeQuotes = strings.NewReplacer("\\", "\\\\", `"`, "\\\"").Replace

// CreateFormFile creates an [io.Writer] like [mime/multipart.Writer.CreateFromFile].
// This version allows to set the mime type, too.
func createFormFile(w *multipart.Writer, fieldname, filename, mimeType string) (io.Writer, error) {
	// Source: https://cs.opensource.google/go/go/+/refs/tags/go1.20:src/mime/multipart/writer.go;l=140
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldname), escapeQuotes(filename)))
	h.Set("Content-Type", mimeType)
	return w.CreatePart(h)
}

func buildRequest(
	doc, filename string,
	status validationStatus,
	url string,
	headers http.Header,
) (*http.Request, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	var err error
	part := func(name, fname, mimeType, content string) {
		if err != nil {
			return
		}
		if fname == "" {
			err = writer.WriteField(name, content)
			return
		}
		var w io.Writer
		if w, err = createFormFile(writer, name, fname, mimeType); err == nil {
			_, err = w.Write([]byte(content))
		}
	}

	base := filepath.Base(filename)
	part("advisory", base, "application/json", doc)
	part("validation_status", "", "text/plain", string(status))

	if err := errors.Join(err, writer.Close()); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	for k, vs := range headers {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}
	contentType := writer.FormDataContentType()
	req.Header.Set("Content-Type", contentType)
	return req, nil
}

func (fm *ForwardManager) loadDocument(ctx context.Context, documentID int64) (*document, error) {

	builder := query.SQLBuilder{}
	builder.CreateWhere(query.FieldEqInt("id", documentID))

	var (
		original         []byte
		filename         string
		failedValidation *bool
	)

	const selectSQL = `SELECT ` +
		`original,` +
		`filename,` +
		`(filename_failed OR remote_failed OR checksum_failed OR signature_failed)` +
		`FROM documents JOIN downloads ` +
		`ON documents.id = downloads.documents_id ` +
		`WHERE `

	if err := fm.db.Run(
		ctx,
		func(rctx context.Context, conn *pgxpool.Conn) error {
			fetchSQL := selectSQL + builder.WhereClause
			return conn.QueryRow(rctx, fetchSQL, builder.Replacements...).Scan(
				&original,
				&filename,
				&failedValidation)
		}, 0,
	); err != nil {
		return nil, err
	}

	validationStatus := notValidatedValidationStatus
	if failedValidation != nil {
		if *failedValidation {
			validationStatus = invalidValidationStatus
		} else {
			validationStatus = validValidationStatus
		}
	}
	return &document{
		data:      string(original),
		filename:  filename,
		valStatus: validationStatus,
	}, nil
}

func (fm *ForwardManager) fetchNewDocuments(
	ctx context.Context,
	url string,
	publisher *string,
) ([]int64, error) {

	var wherePublisher string
	args := []any{url}
	if publisher != nil {
		wherePublisher = "publisher = $2 AND "
		args = append(args, *publisher)
	}

	const selectSQL = `SELECT ` +
		`documents_id ` +
		`FROM events_log ` +
		`WHERE documents_id IN ` +
		`(SELECT id FROM documents` +
		` WHERE %s` +
		` id NOT IN (SELECT documents_id FROM forwarded_documents WHERE url = $1)) ` +
		`ORDER BY time DESC`
	fetchSQL := fmt.Sprintf(selectSQL, wherePublisher)

	var documentIDs []int64
	if err := fm.db.Run(
		ctx,
		func(rctx context.Context, conn *pgxpool.Conn) error {
			rows, _ := conn.Query(rctx, fetchSQL, args...)
			var err error
			documentIDs, err = pgx.CollectRows(
				rows,
				func(row pgx.CollectableRow) (int64, error) {
					var documentID int64
					err := row.Scan(&documentID)
					return documentID, err
				})
			return err
		}, 0,
	); err != nil {
		return nil, err
	}
	return documentIDs, nil
}

func (fm *ForwardManager) logDocument(ctx context.Context, url string, documentID int64) {
	const sql = `INSERT INTO forwarded_documents (url, documents_id) VALUES ($1, $2)`
	if err := fm.db.Run(
		ctx,
		func(ctx context.Context, con *pgxpool.Conn) error {
			_, err := con.Exec(ctx, sql, url, documentID)
			return err
		}, 0,
	); err != nil {
		slog.Error("database error", "err", err)
	}
}

// ForwardTarget contains information about the available target.
type ForwardTarget struct {
	URL  string `json:"url"`
	Name string `json:"name,omitempty"`
	ID   int    `json:"id"`
}

// Targets returns a list of forward targets.
func (fm *ForwardManager) Targets() []ForwardTarget {
	result := make(chan []ForwardTarget)
	fm.fns <- func(fm *ForwardManager) {
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
func (fm *ForwardManager) ForwardDocument(ctx context.Context, targetID int, documentID int64) error {
	result := make(chan error)
	fm.fns <- func(fm *ForwardManager) {
		if targetID < 0 || targetID >= len(fm.forwarders) || fm.forwarders[targetID].cfg.Automatic {
			result <- errors.New("could not find target with specified id")
			return
		}
		document, err := fm.loadDocument(ctx, documentID)
		if err != nil {
			slog.Error("could not load document to forward", "err", err)
			result <- err
			return
		}
		result <- fm.uploadDocument(ctx, document, documentID, fm.forwarders[targetID])
	}
	return <-result
}

func (fm *ForwardManager) kill() {
	fm.done = true
}

// Kill shuts down the forward manager.
func (fm *ForwardManager) Kill() {
	fm.fns <- (*ForwardManager).kill
}
