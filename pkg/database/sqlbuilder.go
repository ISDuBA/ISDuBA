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

// ConstructWhere construct a WHERE clause for a given expression.
func (sb *SQLBuilder) ConstructWhere(e *Expr) string {
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
