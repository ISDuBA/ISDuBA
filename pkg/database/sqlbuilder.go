// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package database

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

// SQLBuilder helps constructing a SQL query.
type SQLBuilder struct {
	WhereClause  string
	Replacements []any
	replToIdx    map[string]int
	Aliases      map[string]string
	Advisory     bool
	TextTables   bool
}

// CreateWhere construct a WHERE clause for a given expression.
func (sb *SQLBuilder) CreateWhere(e *Expr) string {
	var b strings.Builder
	sb.whereRecurse(e, &b)
	sb.WhereClause = b.String()
	return b.String()
}

func (sb *SQLBuilder) searchWhere(e *Expr, b *strings.Builder) {
	const tsquery = `websearch_to_tsquery`

	b.WriteString(`ts @@ ` + tsquery + `('`)
	b.WriteString(e.langValue)
	b.WriteString("',$")
	idx := sb.replacementIndex(e.stringValue)
	b.WriteString(strconv.Itoa(idx + 1))
	b.WriteByte(')')
	// Handle alias
	if e.alias == "" {
		return
	}
	repl := fmt.Sprintf(
		"ts_headline('%[1]s',txt,"+tsquery+"('%[1]s', $%[2]d))",
		e.langValue, idx+1)
	if sb.Aliases == nil {
		sb.Aliases = map[string]string{}
	}
	// We need the text tables to be joined.
	sb.TextTables = true
	sb.Aliases[e.alias] = repl
}

func (sb *SQLBuilder) csearchWhere(_ *Expr, _ *strings.Builder) {
	// TODO: Implement me!
	slog.Debug("csearch is not implemented, yet!")
}

func (sb *SQLBuilder) castWhere(e *Expr, b *strings.Builder) {
	b.WriteString("CAST(")
	sb.whereRecurse(e.children[0], b)
	b.WriteString(" AS ")
	switch e.valueType {
	case stringType:
		b.WriteString("text")
	case intType:
		b.WriteString("int")
	case floatType:
		b.WriteString("float")
	case timeType:
		b.WriteString("timestamptz")
	case boolType:
		b.WriteString("boolean")
	case workflowType:
		b.WriteString("workflow")
	case durationType:
		b.WriteString("interval")
	}
	b.WriteByte(')')
}

func (sb *SQLBuilder) cnstWhere(e *Expr, b *strings.Builder) {

	switch e.valueType {
	case stringType:
		b.WriteByte('$')
		idx := sb.replacementIndex(e.stringValue)
		b.WriteString(strconv.Itoa(idx + 1))
	case intType:
		b.WriteString(strconv.FormatInt(e.intValue, 10))
	case floatType:
		b.WriteString(strconv.FormatFloat(e.floatValue, 'f', -1, 64))
	case timeType:
		b.WriteByte('\'')
		utc := e.timeValue.UTC()
		b.WriteString(utc.Format("2006-01-02T15:04:05-0700"))
		b.WriteString("'::timestamptz")
	case boolType:
		if e.boolValue {
			b.WriteString("TRUE")
		} else {
			b.WriteString("FALSE")
		}
	case workflowType:
		b.WriteByte('\'')
		b.WriteString(e.stringValue)
		b.WriteString("'::workflow")
	case durationType:
		fmt.Fprintf(b, "'%.2f seconds'::interval", e.durationValue.Seconds())
	}
}

func (sb *SQLBuilder) binaryWhere(e *Expr, b *strings.Builder, op string) {
	b.WriteByte('(')
	sb.whereRecurse(e.children[0], b)
	b.WriteString(op)
	sb.whereRecurse(e.children[1], b)
	b.WriteByte(')')
}

func (sb *SQLBuilder) notWhere(e *Expr, b *strings.Builder) {
	b.WriteString("(NOT ")
	sb.whereRecurse(e.children[0], b)
	b.WriteByte(')')
}

const versionsCount = `(SELECT count(*) FROM documents WHERE ` +
	`documents.publisher = advisories.publisher AND ` +
	`documents.tracking_id = advisories.tracking_id)`

func (sb *SQLBuilder) accessWhere(e *Expr, b *strings.Builder) {
	switch column := e.stringValue; column {
	case "tracking_id", "publisher":
		b.WriteString("documents.")
		b.WriteString(column)
	case "versions":
		b.WriteString(versionsCount)
	default:
		b.WriteString(column)
	}
}

func (sb *SQLBuilder) nowWhere(_ *Expr, b *strings.Builder) {
	b.WriteString("current_timestamp")
}

func (sb *SQLBuilder) ilikeWhere(e *Expr, b *strings.Builder) {
	b.WriteByte('(')
	sb.whereRecurse(e.children[0], b)
	b.WriteString(" ILIKE ")
	sb.whereRecurse(e.children[1], b)
	b.WriteByte(')')
}

func (sb *SQLBuilder) ilikePIDWhere(e *Expr, b *strings.Builder) {

	b.WriteString(`EXISTS (` +
		`WITH product_ids AS (SELECT jsonb_path_query(` +
		`document, '$.product_tree.**.product.product_id')::int num ` +
		`FROM documents ds WHERE ds.id = documents.id)` +
		`SELECT * FROM documents_texts dts JOIN product_ids ` +
		`ON product_ids.num = dts.num JOIN unique_texts ON dts.txt_id = unique_texts.id ` +
		`WHERE dts.documents_id = documents.id AND ` +
		`unique_texts.txt ILIKE `)
	sb.whereRecurse(e.children[0], b)
	b.WriteByte(')')
	/*
		b.WriteString(`EXISTS (` +
			`SELECT jsonb_path_query(` +
			`document, '$.product_tree.**.product.product_id')::int ` +
			`FROM documents ds WHERE ds.id = documents.id ` +
			`INTERSECT ` +
			`SELECT num FROM documents_texts ` +
			`WHERE documents_id = documents.id AND ` +
			`txt ILIKE `)
		recurse(e.children[0])
		b.WriteByte(')')
	*/
	/*
		b.WriteString(`EXISTS (` +
			`SELECT num FROM documents_texts ` +
			`WHERE documents_id = documents.id AND ` +
			`txt ILIKE `)
		recurse(e.children[0])
		b.WriteString(` INTERSECT ` +
			`SELECT jsonb_path_query(` +
			`document, '$.product_tree.**.product.product_id')::int ` +
			`FROM documents ds WHERE ds.id = documents.id)`)
	*/
}

func (sb *SQLBuilder) whereRecurse(e *Expr, b *strings.Builder) {
	b.WriteByte('(')
	switch e.exprType {
	case access:
		sb.accessWhere(e, b)
	case cnst:
		sb.cnstWhere(e, b)
	case cast:
		sb.castWhere(e, b)
	case eq:
		sb.binaryWhere(e, b, "=")
	case ne:
		sb.binaryWhere(e, b, "<>")
	case lt:
		sb.binaryWhere(e, b, "<")
	case gt:
		sb.binaryWhere(e, b, ">")
	case le:
		sb.binaryWhere(e, b, "<=")
	case ge:
		sb.binaryWhere(e, b, ">=")
	case not:
		sb.notWhere(e, b)
	case and:
		sb.binaryWhere(e, b, "AND")
	case or:
		sb.binaryWhere(e, b, "OR")
	case search:
		sb.searchWhere(e, b)
	case csearch:
		sb.csearchWhere(e, b)
	case ilike:
		sb.ilikeWhere(e, b)
	case ilikePID:
		sb.ilikePIDWhere(e, b)
	case now:
		sb.nowWhere(e, b)
	case add:
		sb.binaryWhere(e, b, "+")
	case sub:
		sb.binaryWhere(e, b, "-")
	case mul:
		sb.binaryWhere(e, b, "*")
	case div:
		sb.binaryWhere(e, b, "/")
	}
	b.WriteByte(')')
}

func (sb *SQLBuilder) replacementIndex(s string) int {
	if idx, ok := sb.replToIdx[s]; ok {
		return idx
	}
	if sb.replToIdx == nil {
		sb.replToIdx = map[string]int{}
	}
	idx := len(sb.replToIdx)
	sb.replToIdx[s] = idx
	return idx
}

func (sb *SQLBuilder) createFrom() string {
	var from string
	if sb.Advisory {
		from = `documents ` +
			`JOIN advisories ON ` +
			`advisories.tracking_id = documents.tracking_id AND ` +
			`advisories.publisher = documents.publisher`
	} else {
		from = `documents`
	}

	if len(sb.Aliases) > 0 { // XXX: This is wrong
		from += ` JOIN documents_texts ON id = documents_texts.documents_id ` +
			`JOIN unique_texts ON documents_texts.txt_id = unique_texts.id`
	}
	return from
}

// CreateCountSQL returns an SQL count statement to count
// the number of rows which are possible to fetch by the
// given filter.
func (sb *SQLBuilder) CreateCountSQL() string {
	from := sb.createFrom()
	return "SELECT count(*) FROM " + from + " WHERE " + sb.WhereClause
}

// CreateOrder returns a ORDER BY clause for given columns.
func (sb *SQLBuilder) CreateOrder(fields []string) (string, error) {
	var b strings.Builder
	for _, field := range fields {
		desc := strings.HasPrefix(field, "-")
		if desc {
			field = field[1:]
		}
		if _, found := sb.Aliases[field]; !found && !ExistsDocumentColumn(field, sb.Advisory) {
			return "", fmt.Errorf("order field %q does not exists", field)
		}
		if b.Len() > 0 {
			b.WriteByte(',')
		}
		switch field {
		case "tracking_id", "publisher":
			b.WriteString("documents.")
			b.WriteString(field)
		case "cvss_v2_score", "cvss_v3_score", "critical":
			b.WriteString("COALESCE(")
			b.WriteString(field)
			b.WriteString(",0)")
		case "version":
			// TODO: This is not optimal (SemVer).
			b.WriteString(
				`CASE WHEN pg_input_is_valid(version, 'integer') THEN version::int END`)
		default:
			b.WriteString(field)
		}

		if desc {
			b.WriteString(" DESC")
		} else {
			b.WriteString(" ASC")
		}
	}
	return b.String(), nil
}

// CreateQuery creates an SQL statement to query the documents
// table and the associated texts if needed.
// WARN: Make sure that the iput is vetted against injections.
func (sb *SQLBuilder) CreateQuery(
	fields []string,
	order string,
	limit, offset int64,
) string {
	projs := sb.projectionsWithCasts(fields)

	from := sb.createFrom()

	sql := "SELECT " + projs + " FROM " + from + " WHERE " + sb.WhereClause

	if order != "" {
		sql += " ORDER BY " + order
	}

	if limit >= 0 {
		sql += " LIMIT " + strconv.FormatInt(limit, 10)
	}
	if offset > 0 {
		sql += " OFFSET " + strconv.FormatInt(offset, 10)
	}

	return sql
}

// projectionsWithCasts joins given projection adding casts if needed.
func (sb *SQLBuilder) projectionsWithCasts(proj []string) string {
	var b strings.Builder
	for i, p := range proj {
		if i > 0 {
			b.WriteByte(',')
		}
		if alias, found := sb.Aliases[p]; found {
			b.WriteString(alias)
			continue
		}
		switch p {
		case "id", "tracking_id", "publisher":
			b.WriteString("documents.")
			b.WriteString(p)
		case "state":
			b.WriteString("state::text")
		case "versions":
			b.WriteString(versionsCount + `AS versions`)
		default:
			b.WriteString(p)
		}
	}
	return b.String()
}

// CheckProjections checks if the requested projections are valid.
func (sb *SQLBuilder) CheckProjections(proj []string) error {
	for _, p := range proj {
		if _, found := sb.Aliases[p]; found {
			continue
		}
		if !ExistsDocumentColumn(p, sb.Advisory) {
			return fmt.Errorf("column %q does not exists", p)
		}
	}
	return nil
}
