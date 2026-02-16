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
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path/filepath"
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

	// Start the automatic forwarders.
	for _, forwarder := range fm.forwarders {
		if forwarder.cfg.Automatic {
			go forwarder.run(ctx)
			defer forwarder.kill()
		}
	}

	poller := newPoller(fm)
	go poller.run(ctx)
	defer poller.kill()

	ticker := time.NewTicker(fm.cfg.UpdateInterval / 2)
	defer ticker.Stop()
	for !fm.done {
		select {
		case fn := <-fm.fns:
			fn(fm)
		case <-ctx.Done():
			return
		case <-ticker.C:
			fm.runTargets(ctx)
		}
	}
}

// changesDetected tries to deliver detected advisory changes to
// the manager. If the manager is ready a fresh changedAdvisories map
// is returned. If the delivery would block the given map is
// returned so that the poller can go on detecting avoiding duplicates.
func (m *ForwardManager) changesDetected(changes changedAdvisories) changedAdvisories {
	select {
	case m.fns <- func(m *ForwardManager) { m.changes = changes }:
		return changedAdvisories{}
	default:
		return changes
	}
}

// runTargets fetches and sends all new documents to the configured targets.
func (fm *ForwardManager) runTargets(ctx context.Context) {

	/*
		for _, target := range fm.targets {
			if !target.automatic || !target.running.TryLock() {
				continue
			}
			documentIDs, err := fm.fetchNewDocuments(ctx, target.url, target.publisher)
			if err != nil {
				slog.Error("could not fetch documents to forward", "err", err)
				target.running.Unlock()
				continue
			}
			if len(documentIDs) > 0 {
				go fm.uploadDocuments(ctx, target, documentIDs)
			} else {
				target.running.Unlock()
			}
		}
	*/
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
