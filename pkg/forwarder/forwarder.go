// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

// Package forwarder implements the document forwarder.
package forwarder

import (
	"bytes"
	"context"
	"crypto/tls"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path/filepath"
	"strings"
	"sync"
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

type target struct {
	url       string
	publisher *string
	header    http.Header
	client    http.Client
	running   *sync.Mutex
	enabled   bool
}

// ForwardManager forwards documents to specified targets.
type ForwardManager struct {
	cfg     *config.Forwarder
	targets []target
	db      *database.DB
	fns     chan func(*ForwardManager)
	done    bool
}

// NewForwardManager creates a new forward manager.
func NewForwardManager(cfg *config.Forwarder, db *database.DB) *ForwardManager {
	var targets []target

	for _, targetCfg := range cfg.Targets {

		// Init http clients
		client := http.Client{
			Timeout: targetCfg.Timeout,
		}

		var tlsConfig tls.Config
		if targetCfg.ClientPrivateCert != "" && targetCfg.ClientPublicCert != "" {
			clientCert, err := tls.LoadX509KeyPair(targetCfg.ClientPublicCert, targetCfg.ClientPrivateCert)
			if err != nil {
				slog.Error("could not configure forward target client cert", "err", err)
			} else {
				tlsConfig.Certificates = []tls.Certificate{clientCert}
			}
		}

		client.Transport = &http.Transport{
			TLSClientConfig: &tlsConfig,
		}

		headers := http.Header{}
		for _, header := range targetCfg.Header {
			h := strings.Split(header, ":")
			if len(h) != 2 {
				slog.Error("forwarder init: could not set invalid header key:pair value", "header", header)
				continue
			}
			headers.Add(h[0], h[1])
		}

		t := target{
			client:    client,
			url:       targetCfg.URL,
			publisher: targetCfg.Publisher,
			header:    headers,
			running:   &sync.Mutex{},
			enabled:   targetCfg.Enabled,
		}

		targets = append(targets, t)
	}
	return &ForwardManager{
		cfg:     cfg,
		targets: targets,
		db:      db,
		fns:     make(chan func(manager *ForwardManager)),
	}
}

// Run runs the forward manager. To be used in a Go routine.
func (fm *ForwardManager) Run(ctx context.Context) {
	ticker := time.NewTicker(fm.cfg.UpdateInterval)
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

// runTargets fetches and sends all new documents to the configured targets.
func (fm *ForwardManager) runTargets(ctx context.Context) {
	for index := range fm.targets {
		target := &fm.targets[index]
		if !target.enabled || !target.running.TryLock() {
			continue
		}
		documentIDs, err := fm.fetchNewDocuments(ctx, target.url, target.publisher)
		if err != nil {
			slog.Error("could not fetch documents to forward", "err", err)
		}
		go fm.uploadDocuments(ctx, target, documentIDs)
	}
}

func (fm *ForwardManager) uploadDocuments(ctx context.Context, target *target, documentIDs []int64) {
	defer target.running.Unlock()
	for _, documentID := range documentIDs {
		document, filename, err := fm.loadDocument(ctx, documentID)
		if err != nil {
			slog.Error("could not load document to forward", "err", err)
			continue
		}
		documentString := string(document)
		if err := fm.uploadDocument(ctx, documentString, filename, documentID, target); err != nil {
			slog.Error("could not forward document", "err", err)
		}
	}
}

func (fm *ForwardManager) uploadDocument(ctx context.Context, doc string, filename string, documentID int64, target *target) error {
	valStatus := fm.fetchValidationStatus(ctx, documentID)
	req, err := fm.buildRequest(doc, filename, valStatus, target)
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

func (fm *ForwardManager) buildRequest(doc string, filename string, status validationStatus, target *target) (*http.Request, error) {
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

	if err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, target.url, body)
	if err != nil {
		return nil, err
	}
	contentType := writer.FormDataContentType()
	if target.header != nil {
		req.Header = target.header.Clone()
	}
	req.Header.Set("Content-Type", contentType)
	return req, nil
}

func (fm *ForwardManager) loadDocument(ctx context.Context, documentID int64) ([]byte, string, error) {
	expr := query.FieldEqInt("id", documentID)

	fields := []string{"original", "filename"}
	builder := query.SQLBuilder{}
	builder.CreateWhere(expr)
	sql := builder.CreateQuery(fields, "", -1, -1)

	var original []byte
	var filename string

	if err := fm.db.Run(
		ctx,
		func(rctx context.Context, conn *pgxpool.Conn) error {
			return conn.QueryRow(rctx, sql, builder.Replacements...).Scan(&original, &filename)
		}, 0,
	); err != nil {
		return nil, "", err
	}
	return original, filename, nil
}

func (fm *ForwardManager) fetchNewDocuments(ctx context.Context, url string, publisher *string) ([]int64, error) {
	var wherePublisher string
	args := []any{url}
	if publisher != nil {
		wherePublisher = "publisher = $2 AND "
		args = append(args, *publisher)
	}
	var documentIDs []int64

	if err := fm.db.Run(
		ctx,
		func(rctx context.Context, conn *pgxpool.Conn) error {
			fetchSQL := `SELECT documents_id FROM events_log ` +
				`WHERE documents_id in ` +
				`(SELECT id FROM documents WHERE ` + wherePublisher +
				`id NOT IN (SELECT documents_id FROM forwarded_documents WHERE url = $1)) ORDER BY time DESC`
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
		return documentIDs, err
	}

	return documentIDs, nil
}

func (fm *ForwardManager) fetchValidationStatus(ctx context.Context, documentID int64) validationStatus {
	status := notValidatedValidationStatus

	if err := fm.db.Run(ctx, func(rctx context.Context, conn *pgxpool.Conn) error {
		fetchSQL := `SELECT bool_or(filename_failed OR remote_failed OR checksum_failed OR signature_failed) FROM downloads ` +
			`WHERE documents_id = $1`
		var failedValidation bool
		if err := conn.QueryRow(rctx, fetchSQL, documentID).Scan(&failedValidation); err != nil {
			// Check if there was no download record
			if err == sql.ErrNoRows {
				return nil
			}
			return err
		}
		if failedValidation {
			status = invalidValidationStatus
		} else {
			status = validValidationStatus
		}
		return nil
	}, 0); err != nil {
		slog.Error("could not fetch validation status", "err", err)
	}
	return status
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
	URL string `json:"url"`
	ID  int    `json:"id"`
}

// Targets returns a list of forward targets.
func (fm *ForwardManager) Targets() []ForwardTarget {
	result := make(chan []ForwardTarget)
	fm.fns <- func(fm *ForwardManager) {
		targets := make([]ForwardTarget, len(fm.targets))
		for i := range fm.targets {
			target := &fm.targets[i]
			targets[i] = ForwardTarget{ID: i, URL: target.url}
		}
		result <- targets
	}
	return <-result
}

// ForwardDocument sends the document to the specified target.
func (fm *ForwardManager) ForwardDocument(ctx context.Context, targetID int, documentID int64) error {
	result := make(chan error)
	fm.fns <- func(fm *ForwardManager) {
		if len(fm.targets) <= targetID {
			result <- errors.New("could not find target with specified id")
			return
		}
		document, filename, err := fm.loadDocument(ctx, documentID)
		if err != nil {
			slog.Error("could not load document to forward", "err", err)
			result <- err
			return
		}
		documentString := string(document)
		result <- fm.uploadDocument(ctx, documentString, filename, documentID, &fm.targets[targetID])
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
