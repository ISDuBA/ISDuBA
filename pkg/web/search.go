// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

package web

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"iter"
	"log/slog"
	"maps"
	"net/http"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/ISDuBA/ISDuBA/pkg/database/query"
	"github.com/ISDuBA/ISDuBA/pkg/itertools"
	"github.com/ISDuBA/ISDuBA/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	aggregatedDocument struct {
		id   int64
		rows []map[string]any
	}
	aggregatedDocuments []aggregatedDocument
)

func (c *Controller) aggregatedResults(
	ctx *gin.Context,
	calcCount bool,
	limit, offset int64,
	builder *query.AdvancedSQLBuilder,
) {
	var (
		ads  aggregatedDocuments
		sql  = builder.CreateQuery(-1, -1)
		rctx = ctx.Request.Context()
	)

	if slog.Default().Enabled(rctx, slog.LevelDebug) {
		slog.Debug("documents", "SQL", query.InterpolateSQLqnd(sql, builder.Replacements))
	}
	if err := c.db.Run(
		rctx,
		func(rctx context.Context, conn *pgxpool.Conn) error {
			rows, err := conn.Query(rctx, sql, builder.Replacements...)
			if err != nil {
				return fmt.Errorf("cannot fetch results: %w", err)
			}
			defer rows.Close()
			if ads, err = scanAggregatedDocuments(rows, builder); err != nil {
				return fmt.Errorf("loading data failed: %w", err)
			}
			return nil
		},
		c.cfg.Database.MaxQueryDuration, // In case the user provided a very expensive query.
	); err != nil {
		slog.Error("database error", "err", err)
		models.SendError(ctx, http.StatusInternalServerError, err)
		return
	}

	var ilikes query.ILikeExpr
	if expr := builder.Expr(); expr != nil {
		if searches := slices.Collect(itertools.Apply(
			expr.Searches(),
			func(search string) string {
				return query.LikeEscape(search)
			},
		)); len(searches) > 0 {
			var err error
			ilikes, err = query.CompileILike(searches...)
			if err != nil {
				slog.Error("compiling ilikes failed", "err", err)
				models.SendError(ctx, http.StatusInternalServerError, err)
				return
			}
		}
	}

	const (
		buffer = 20    // Reading context
		fill   = "..." // Gap filler
	)
	delims := [2]string{`[!<`, `>!]`} // Used to mark the sections.

	// We load texts from the database lazily so we render in another connection.
	// XXX: Think about moving it to the DB stuff above.
	if err := c.db.Run(
		rctx,
		func(rctx context.Context, conn *pgxpool.Conn) error {
			// Load a text only if we don't already have it.
			const uniqueSQL = `SELECT txt FROM unique_texts WHERE id = $1`
			// Cache the shortend texts.
			txts := map[int64]string{}
			fetchText := func(txtID int64) (string, error) {
				txt, ok := txts[txtID]
				if ok {
					return txt, nil
				}
				if err := conn.QueryRow(rctx, uniqueSQL, txtID).Scan(&txt); err != nil {
					return "", fmt.Errorf("fetching unique text failed: %w", err)
				}
				if ilikes.Regexp != nil {
					sections := ilikes.Search(txt)
					txt = sections.Shorten(txt, buffer, fill, delims)
				}
				txts[txtID] = txt
				return txt, nil
			}
			// Catch errors occuring while rendering to logged outside.
			var trackedErr error
			ctx.Render(http.StatusOK, jsonStream(trackError(
				&trackedErr,
				func(w http.ResponseWriter) error {
					fmt.Fprint(w, "{")
					if calcCount {
						fmt.Fprintf(w, `"count":%d`, len(ads))
					}
					// Only produce documents when we have them.
					if len(ads) == 0 {
						_, err := fmt.Fprint(w, "}")
						return err
					}
					if calcCount {
						fmt.Fprint(w, ",")
					}
					fmt.Fprint(w, `"documents":[`)
					firstDocument := true
					enc := json.NewEncoder(w)

					encoded := make(map[string]any)

					for ad := range ads.window(limit, offset) {
						if firstDocument {
							firstDocument = false
						} else {
							fmt.Fprint(w, ",")
						}
						fmt.Fprintf(w, `{"id":%d,"data":[`, ad.id)

						for i, row := range ad.rows {
							if i > 0 {
								fmt.Fprint(w, ",")
							}
							clear(encoded)
							for k, v := range row {
								if !builder.HasAlias(k) {
									encoded[k] = v
									continue
								}
								// If we have an alias we need to load the text from the database.
								id, ok := asInt64(v)
								if !ok {
									return fmt.Errorf("alias %q has not an int value", k)
								}
								txt, err := fetchText(id)
								if err != nil {
									return err
								}
								encoded[k] = txt
							}
							if err := enc.Encode(encoded); err != nil {
								return fmt.Errorf("serializing row failed: %w err", err)
							}
						}
						if _, err := fmt.Fprint(w, "]}"); err != nil {
							return fmt.Errorf("writing window failed %w", err)
						}
					}
					_, err := fmt.Fprint(w, "]}")
					return err
				})))
			return trackedErr
		},
		c.cfg.Database.MaxQueryDuration, // In case the user provided a very expensive query.
	); err != nil {
		slog.Error("database error", "err", err)
		// Too late to send an error to the client.
	}
}

type jsonStream func(http.ResponseWriter) error

func (js jsonStream) Render(w http.ResponseWriter) error { return js(w) }

func (jsonStream) WriteContentType(w http.ResponseWriter) {
	render.JSON{}.WriteContentType(w)
}

func trackError(
	err *error,
	fn func(http.ResponseWriter) error,
) func(http.ResponseWriter) error {
	return func(rw http.ResponseWriter) error {
		*err = fn(rw)
		return *err
	}
}

// scanAggregatedDocuments returns a result set into a slice of aggregated documents.
func scanAggregatedDocuments(
	rows pgx.Rows,
	builder *query.AdvancedSQLBuilder,
) (aggregatedDocuments, error) {
	fields := builder.Fields()
	idIdx := slices.Index(fields, "id")
	if idIdx == -1 {
		return nil, errors.New("missing id column to aggregate")
	}
	values := make([]any, len(fields))
	ptrs := make([]any, len(fields))
	for i := range ptrs {
		ptrs[i] = &values[i]
	}
	ads := make(aggregatedDocuments, 0, 512)
	lastID := int64(-1)
	have := make(map[string]any)
	for rows.Next() {
		if err := rows.Scan(ptrs...); err != nil {
			return nil, fmt.Errorf("scanning row failed: %w", err)
		}
		id, ok := asInt64(values[idIdx])
		if !ok {
			// XXX: Should we panic here!?
			return nil, fmt.Errorf("id column is not an int: %T", values[idIdx])
		}
		if id != lastID { // A new document bundle.
			lastID = id
			clear(have)
			for j, v := range values {
				// Ignore id as it will be already stored in the aggregatedDocument.
				if name := fields[j]; name != "id" {
					have[name] = v
				}
			}
			ads = append(ads, aggregatedDocument{
				id:   id,
				rows: []map[string]any{maps.Clone(have)},
			})
		} else { // We already have documents for this bundle.
			row := make(map[string]any, 1) // The diff should be small.
			for j, v := range values {
				name := fields[j]
				if name == "id" {
					continue
				}
				if x, ok := have[name]; !ok ||
					builder.HasAlias(name) || !reflect.DeepEqual(x, v) {
					row[name] = v
					have[name] = v
				}
			}
			if len(row) > 0 {
				last := &ads[len(ads)-1].rows
				*last = append(*last, row)
			}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scanning failed: %w", err)
	}
	return ads, nil
}

func asInt64(x any) (int64, bool) {
	switch v := x.(type) {
	case int64:
		return v, true
	case int8:
		return int64(v), true
	case int16:
		return int64(v), true
	case int32:
		return int64(v), true
	case uint64:
		return int64(v), true
	case uint8:
		return int64(v), true
	case uint16:
		return int64(v), true
	case uint32:
		return int64(v), true
	case int:
		return int64(v), true
	case uint:
		return int64(v), true
	}
	return 0, false
}

func (ads aggregatedDocuments) window(limit, offset int64) iter.Seq[*aggregatedDocument] {
	var start, end int
	if offset < 0 {
		start = 0
	} else {
		start = min(int(offset), len(ads))
	}
	if limit < 0 {
		end = len(ads)
	} else {
		end = min(int(offset+limit), len(ads))
	}
	return func(yield func(*aggregatedDocument) bool) {
		for i := range ads[start:end] {
			if !yield(&ads[i]) {
				return
			}
		}
	}
}

// documentTexts is an end point to return JSON paths and highlighting positions.
//
//	@Summary		Returns List of JSON paths with highlighting positions.
//	@Description	Returns List of JSON paths with highlighting positions inside text matches found by the search query.
//	@Param			id		path	int	true	"Document ID"
//	@Param			query	query	string	false	"Document query"
//	@Param			include	query	bool	false	"Include the texts in the result"
//	@Produce		json
//	@Success		200	{object}	models.TextPaths
//	@Failure		400	{object}	models.Error
//	@Failure		401
//	@Failure		500	{object}	models.Error
//	@Router			/documents/texts/{id} [get]
func (c *Controller) documentTexts(ctx *gin.Context) {
	// Get an ID from context
	docID, ok := parse(ctx, toInt64, ctx.Param("id"))
	if !ok {
		return
	}

	parser := query.Parser{
		Mode:            query.DocumentMode,
		MinSearchLength: MinSearchLength,
	}

	// The query to filter the documents.
	expr, ok := parse(ctx, parser.Parse, ctx.DefaultQuery("query", "true"))
	if !ok {
		return
	}

	// Should the texts be included in the result?
	include, ok := parse(ctx, strconv.ParseBool, ctx.DefaultQuery("include", "false"))
	if !ok {
		return
	}

	searches := slices.Collect(itertools.Apply(
		expr.Searches(),
		func(search string) string {
			return query.LikeEscape(search)
		},
	))
	if len(searches) == 0 {
		// No search texts found in query -> no need to proceed.
		ctx.JSON(http.StatusOK, models.TextPaths{})
		return
	}
	ilikes, err := query.CompileILike(searches...)
	if err != nil {
		slog.Error("compiling ilikes failed", "err", err)
		models.SendError(ctx, http.StatusInternalServerError, err)
		return
	}

	searchTerm, replacements := query.
		CreateTextSearchWhereClause("ut.txt", expr)
	replacements = append(replacements, docID)

	tlps := c.tlps(ctx)

	const (
		publisherTLPSQL = `` +
			`SELECT` +
			` ads.publisher,` +
			` docs.tlp ` +
			`FROM` +
			` documents docs JOIN` +
			` advisories ads ON docs.advisories_id = ads.id ` +
			`WHERE docs.id = $1`
		uniqueTextsSQL = `` +
			`SELECT` +
			` dt.num,` +
			` ut.txt ` +
			`FROM` +
			` documents_texts dt JOIN` +
			` unique_texts ut ON ut.id = dt.txt_id ` +
			`WHERE` +
			` dt.documents_id = $%d AND %s`
		documentSQL = `` +
			`SELECT document ` +
			`FROM documents ` +
			`WHERE id = $1`
	)

	//  Data to be loaded from the database.
	var (
		uniqueTexts  = make(map[int64]string)
		documentData []byte
	)

	var forbidden bool
	switch err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			// Check permissions.
			var publisher, tlp string
			if err := conn.QueryRow(rctx, publisherTLPSQL, docID).Scan(
				&publisher,
				&tlp,
			); err != nil {
				return fmt.Errorf("failed to extract publisher and tlp: %w", err)
			}
			if len(tlps) > 0 && !tlps.Allowed(publisher, models.TLP(tlp)) {
				forbidden = true
				return nil
			}
			// Load unique texts.
			query := fmt.Sprintf(uniqueTextsSQL, len(replacements), searchTerm)
			rows, err := conn.Query(rctx, query, replacements...)
			if err != nil {
				return fmt.Errorf("loading unique texts failed: %w", err)
			}
			defer rows.Close()
			for rows.Next() {
				var (
					textID int64
					text   string
				)
				if err := rows.Scan(&textID, &text); err != nil {
					return fmt.Errorf("scanning unique texts failed: %w", err)
				}
				uniqueTexts[textID] = text
			}
			if err := rows.Err(); err != nil {
				return fmt.Errorf("scanning unique texts failed: %w", err)
			}
			if len(uniqueTexts) == 0 {
				// No matches.
				return nil
			}
			rows.Close()
			// Load the document itself.
			if err := conn.QueryRow(rctx, documentSQL, docID).Scan(&documentData); err != nil {
				return fmt.Errorf("loading document failed: %w", err)
			}
			return nil
		},
		c.cfg.Database.MaxQueryDuration, // In case the user provided a very expensive query.
	); {
	case errors.Is(err, pgx.ErrNoRows):
		models.SendErrorMessage(ctx, http.StatusNotFound, "document not found")
		return
	case err != nil:
		slog.Error("database error", "err", err)
		models.SendError(ctx, http.StatusInternalServerError, err)
		return
	case forbidden:
		models.SendErrorMessage(ctx, http.StatusForbidden, "access denied")
		return
	}

	if len(uniqueTexts) == 0 {
		ctx.JSON(http.StatusOK, models.TextPaths{})
		return
	}

	var document any
	dec := json.NewDecoder(bytes.NewReader(documentData))
	// We use Numbers here to make it easier to
	// filter out strings that apparently are not IDs like 3.1415
	// which would be otherwise parse as float64s. 3.000 would
	// also be parsed into float64s as would 3. Not using Numbers
	// would make hard to tell them apart.
	dec.UseNumber()
	if err := dec.Decode(&document); err != nil {
		models.SendError(ctx, http.StatusInternalServerError, err)
		return
	}
	paths := buildTextPaths(document, uniqueTexts, ilikes, include)
	ctx.JSON(http.StatusOK, paths)
}

// There are integers in CSAF documents which are not text ids.
var knownNoneTexts = map[string]struct{}{
	"/document/tracking/version": {},
	// TODO: Fill me!
}

// buildTextPaths traverses the document tree and attemps to resolve
// the unique texts producing a list of matches.
func buildTextPaths(
	document any,
	uniqueTexts map[int64]string,
	ilikes query.ILikeExpr,
	include bool,
) models.TextPaths {
	paths := make(models.TextPaths, 0, len(uniqueTexts))

	// XXX: Leaving this currently in for further debugging efforts.
	// usedTexts := make(map[int64]struct{}, len(uniqueTexts))

	store := func(path []string, id int64) {
		joined := strings.Join(path, "")
		// Ignore the ids known to be regular texts.
		if _, ok := knownNoneTexts[joined]; ok {
			return
		}
		// Only include this text if we have text
		// with this id.
		if text, ok := uniqueTexts[id]; ok {
			// XXX: Leaving this currently in for further debugging efforts.
			// usedTexts[id] = struct{}{}
			var t *string
			if include {
				t = &text
			}
			paths = append(paths, models.TextPath{
				Path:      joined,
				Text:      t,
				Positions: ilikes.Search(text),
			})
		}
	}

	/* XXX: Leaving this currently in for further debugging efforts.
	defer func() {
		if tl, ul := len(uniqueTexts), len(usedTexts); tl > ul {
			fmt.Printf("len unique texts: %d\n", tl)
			fmt.Printf("len used texts: %d\n", ul)
			for id, txt := range uniqueTexts {
				if _, ok := usedTexts[id]; !ok {
					fmt.Printf("\tmissing: %d %q\n", id, txt)
				}
			}

		}
	}()
	*/

	// Recursively traverse the document and find the text ids.
	var recurse func(any, []string)
	recurse = func(curr any, path []string) {
		switch x := curr.(type) {
		case map[string]any:
			// Sort to make output deterministic.
			keys := slices.Sorted(maps.Keys(x))
			for _, key := range keys {
				path = append(path, "/", key)
				recurse(x[key], path)
				path = path[:len(path)-2]
			}
		case []any:
			for i, y := range x {
				path = append(path, "["+strconv.Itoa(i)+"]")
				recurse(y, path)
				path = path[:len(path)-1]
			}
		case json.Number:
			if id, err := x.Int64(); err == nil { // Only if its an integer.
				store(path, id)
			}
		case string:
			if id, err := strconv.ParseInt(x, 10, 64); err == nil {
				store(path, id)
			} else if positions := ilikes.Search(x); len(positions) != 0 {
				// Its not indexed but matches the search.
				joined := strings.Join(path, "")
				paths = append(paths, models.TextPath{
					Path:      joined,
					Text:      &x,
					Positions: positions,
				})
			}
		}
	}
	recurse(document, nil)
	return paths
}
