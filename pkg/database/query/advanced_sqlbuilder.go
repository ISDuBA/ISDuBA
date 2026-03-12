// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package query

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

// AdvancedSQLBuilder helps to construct a SQL query.
type AdvancedSQLBuilder struct {
	expr         *Expr
	parser       *AdvancedParser
	orderFields  []string
	fields       []string
	whereClause  string
	order        string
	Replacements []any
	replToIdx    map[string]int
	aliases      map[string]string
	ignoreFields map[string]struct{}
	usedSources  columnSource
}

// AdvancedSQLBuilderOption is an option to create an advanced SQL builder.
type AdvancedSQLBuilderOption func(*AdvancedSQLBuilder)

// AdvancedSQLBuilderExpr creates an option to create an advanced SQL builder
// with an expression.
func AdvancedSQLBuilderExpr(e *Expr) AdvancedSQLBuilderOption {
	return func(ab *AdvancedSQLBuilder) {
		ab.expr = e
	}
}

// AdvancedSQLBuilderOrderFields creates an option to create an advanced SQL builder
// with a order fields.
func AdvancedSQLBuilderOrderFields(orderFields []string) AdvancedSQLBuilderOption {
	return func(ab *AdvancedSQLBuilder) {
		ab.orderFields = orderFields
	}
}

// AdvancedSQLBuilderFields creates an option to create an advanced SQL builder
// with projection fields.
func AdvancedSQLBuilderFields(fields []string) AdvancedSQLBuilderOption {
	return func(ab *AdvancedSQLBuilder) {
		ab.fields = fields
	}
}

// AdvancedSQLBuilderParser creates an option to create an advanced SQL builder
// with a given parser.
func AdvancedSQLBuilderParser(parser *AdvancedParser) AdvancedSQLBuilderOption {
	return func(ab *AdvancedSQLBuilder) {
		ab.parser = parser
	}
}

// NewAdvancedSQLBuilder creates a new advanced builder with a list of options.
func NewAdvancedSQLBuilder(options ...AdvancedSQLBuilderOption) (*AdvancedSQLBuilder, error) {
	ab := new(AdvancedSQLBuilder)
	for _, option := range options {
		option(ab)
	}
	// If given the parser has already a list of used source tables.
	if ab.parser != nil {
		ab.usedSources = ab.parser.UsedSources
	}
	if err := ab.check(); err != nil {
		return nil, fmt.Errorf("creating advanced SQL builder failed: %w", err)
	}
	return ab, nil
}

// HasFields returns true if the builder has projection fields.
func (sb *AdvancedSQLBuilder) HasFields() bool {
	return len(sb.fields) > 0
}

// createWhere construct a WHERE clause for a given expression.
func (sb *AdvancedSQLBuilder) createWhere() string {
	if sb.whereClause != "" {
		return sb.whereClause
	}
	var b strings.Builder
	sb.whereRecurse(sb.expr, &b)
	sb.whereClause = b.String()
	return sb.whereClause
}

func (sb *AdvancedSQLBuilder) mode() ParserMode {
	if sb.parser != nil {
		return sb.parser.Mode
	}
	return DocumentMode
}

// RemoveIgnoredFields removes fields that should be ignored.
func (sb *AdvancedSQLBuilder) RemoveIgnoredFields() []string {
	filtered := make([]string, 0, len(sb.fields))
	for _, f := range sb.fields {
		if _, found := sb.ignoreFields[f]; !found {
			filtered = append(filtered, f)
		}
	}
	return filtered
}

func (sb *AdvancedSQLBuilder) searchWhere(e *Expr, b *strings.Builder) {
	fmt.Fprintf(b, "txt ILIKE $%d",
		sb.replacementIndex(LikeEscape(e.stringValue))+1)

	// Handle alias
	if e.alias == "" {
		return
	}
	switch sb.mode() {
	case AdvisoryMode, DocumentMode:
		fmt.Fprintf(b, "EXISTS(SELECT 1 FROM documents_texts "+
			"JOIN unique_texts ON unique_texts.id = documents_texts.txt_id "+
			"WHERE txt ILIKE $%d "+
			"AND documents_texts.documents_id = documents.id)",
			sb.replacementIndex(LikeEscape(e.stringValue))+1)
	case EventMode:
		// TODO clarify how to handle event search
	}
	if sb.ignoreFields == nil {
		sb.ignoreFields = map[string]struct{}{}
	}
	sb.ignoreFields[e.alias] = struct{}{}
}

func (sb *AdvancedSQLBuilder) mentionedWhere(e *Expr, b *strings.Builder) {
	switch sb.mode() {
	case AdvisoryMode:
		fmt.Fprintf(b, "EXISTS(SELECT 1 FROM comments "+
			"JOIN documents docs ON comments.documents_id = docs.id "+
			"WHERE message ILIKE $%d "+
			"AND docs.advisories_id = documents.advisories_id)",
			sb.replacementIndex(LikeEscape(e.stringValue))+1)
	case DocumentMode:
		fmt.Fprintf(b, "EXISTS(SELECT 1 FROM comments WHERE message ILIKE $%d "+
			"AND comments.documents_id = documents.id)",
			sb.replacementIndex(LikeEscape(e.stringValue))+1)
	case EventMode:
		fmt.Fprintf(b, "EXISTS(SELECT 1 FROM comments WHERE message ILIKE $%d "+
			"AND comments.id = events_log.comments_id)",
			sb.replacementIndex(LikeEscape(e.stringValue))+1)
	}
}

func (sb *AdvancedSQLBuilder) involvedWhere(e *Expr, b *strings.Builder) {
	switch sb.mode() {
	case AdvisoryMode, EventMode:
		fmt.Fprintf(b, "EXISTS(SELECT 1 FROM events_log JOIN documents docs "+
			"ON events_log.documents_id = docs.id "+
			"WHERE actor = $%d "+
			"AND docs.advisories_id = documents.advisories_id)",
			sb.replacementIndex(e.stringValue)+1)
	case DocumentMode:
		fmt.Fprintf(b, "EXISTS(SELECT 1 FROM events_log WHERE actor = $%d "+
			"AND events_log.documents_id = documents.id)",
			sb.replacementIndex(e.stringValue)+1)
	}
}

func (sb *AdvancedSQLBuilder) castWhere(e *Expr, b *strings.Builder) {
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
	case eventsType:
		b.WriteString("events")
	case statusType:
		b.WriteString("status")
	case durationType:
		b.WriteString("interval")
	}
	b.WriteByte(')')
}

func (sb *AdvancedSQLBuilder) cnstWhere(e *Expr, b *strings.Builder) {
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
	case eventsType:
		b.WriteByte('\'')
		b.WriteString(e.stringValue)
		b.WriteString("'::events")
	case statusType:
		b.WriteByte('\'')
		b.WriteString(e.stringValue)
		b.WriteString("'::status")
	case durationType:
		fmt.Fprintf(b, "'%.2f seconds'::interval", e.durationValue.Seconds())
	}
}

func (sb *AdvancedSQLBuilder) binaryWhere(e *Expr, b *strings.Builder, op string) {
	b.WriteByte('(')
	sb.whereRecurse(e.children[0], b)
	b.WriteString(op)
	sb.whereRecurse(e.children[1], b)
	b.WriteByte(')')
}

func (sb *AdvancedSQLBuilder) notWhere(e *Expr, b *strings.Builder) {
	b.WriteString("(NOT ")
	sb.whereRecurse(e.children[0], b)
	b.WriteByte(')')
}

func (sb *AdvancedSQLBuilder) accessWhere(e *Expr, b *strings.Builder) {
	switch column := e.stringValue; column {
	case "id":
		b.WriteString("documents.")
		b.WriteString(column)
	case "tracking_id", "publisher":
		b.WriteString("advisories.")
		b.WriteString(column)
	case "versions":
		b.WriteString(versionsCount)
	case "comments":
		switch sb.mode() {
		case AdvisoryMode:
			b.WriteString(column)
		case DocumentMode:
			b.WriteString(commentsCountDocuments)
		case EventMode:
			b.WriteString(commentsCountEvents)
		}
	case "event_state":
		b.WriteString("events_log.state")
	case "ssvc":
		b.WriteString("ssvc_current.ssvc")
	default:
		b.WriteString(column)
	}
}

func (sb *AdvancedSQLBuilder) nowWhere(_ *Expr, b *strings.Builder) {
	b.WriteString("current_timestamp")
}

func (sb *AdvancedSQLBuilder) ilikeWhere(e *Expr, b *strings.Builder) {
	b.WriteByte('(')
	sb.whereRecurse(e.children[0], b)
	b.WriteString(` ILIKE ` + ilikePrefix)
	sb.whereRecurse(e.children[1], b)
	b.WriteString(ilikeSuffix + `)`)
}

func (sb *AdvancedSQLBuilder) ilikePNameWhere(e *Expr, b *strings.Builder) {
	b.WriteString(`EXISTS (` +
		`WITH product_names AS (SELECT jsonb_path_query(` +
		`document, '$.product_tree.**.product.name')::int num ` +
		`FROM documents ds WHERE ds.id = documents.id)` +
		`SELECT * FROM documents_texts dts JOIN product_names ` +
		`ON product_names.num = dts.num JOIN unique_texts ON dts.txt_id = unique_texts.id ` +
		`WHERE dts.documents_id = documents.id AND ` +
		`unique_texts.txt ILIKE ` + ilikePrefix)
	sb.whereRecurse(e.children[0], b)
	b.WriteString(ilikeSuffix + `)`)
}
func (sb *AdvancedSQLBuilder) ilikePIDWhere(e *Expr, b *strings.Builder) {
	b.WriteString(`EXISTS (` +
		`WITH product_ids AS (SELECT jsonb_path_query(` +
		`document, '$.product_tree.**.product.product_id')::int num ` +
		`FROM documents ds WHERE ds.id = documents.id)` +
		`SELECT * FROM documents_texts dts JOIN product_ids ` +
		`ON product_ids.num = dts.num JOIN unique_texts ON dts.txt_id = unique_texts.id ` +
		`WHERE dts.documents_id = documents.id AND ` +
		`unique_texts.txt ILIKE ` + ilikePrefix)
	sb.whereRecurse(e.children[0], b)
	b.WriteString(ilikeSuffix + `)`)
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

func (sb *AdvancedSQLBuilder) whereRecurse(e *Expr, b *strings.Builder) {
	if e == nil {
		return
	}
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
	case mentioned:
		sb.mentionedWhere(e, b)
	case involved:
		sb.involvedWhere(e, b)
	case ilike:
		sb.ilikeWhere(e, b)
	case ilikePName:
		sb.ilikePNameWhere(e, b)
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

func (sb *AdvancedSQLBuilder) replacementIndex(s string) int {
	if idx, ok := sb.replToIdx[s]; ok {
		return idx
	}
	if sb.replToIdx == nil {
		sb.replToIdx = map[string]int{}
	}
	sb.Replacements = append(sb.Replacements, s)
	idx := len(sb.replToIdx)
	sb.replToIdx[s] = idx
	return idx
}

func (sb *AdvancedSQLBuilder) createFrom(b *strings.Builder) {
	switch sb.mode() {
	case AdvisoryMode, DocumentMode:
		b.WriteString(`documents ` +
			`JOIN advisories ON ` +
			`advisories.id = documents.advisories_id`)
	case EventMode:
		b.WriteString(`events_log JOIN documents ON events_log.documents_id = documents.id ` +
			`JOIN advisories ON advisories.id = documents.advisories_id ` +
			`LEFT JOIN (SELECT message, id FROM comments) AS comment ON events_log.comments_id = comment.id`)
	}

	// Add SSVC if exists
	b.WriteString(` LEFT JOIN LATERAL ( ` +
		`SELECT ssvc FROM ssvc_history ` +
		`WHERE documents_id = documents.id ` +
		`ORDER BY changedate DESC, change_number DESC LIMIT 1 ` +
		`) AS ssvc_current ON TRUE`)

	if sb.usedSources.contains(textTable) {
		b.WriteString(` JOIN documents_texts ON documents.id = documents_texts.documents_id ` +
			`JOIN unique_texts ON documents_texts.txt_id = unique_texts.id`)
	}
}

// CreateCountSQL returns an SQL count statement to count
// the number of rows which are possible to fetch by the
// given filter.
func (sb *AdvancedSQLBuilder) CreateCountSQL() string {
	var b strings.Builder
	b.WriteString("SELECT count(*) FROM ")
	sb.createFrom(&b)
	b.WriteString(" WHERE ")
	b.WriteString(sb.createWhere())
	return b.String()
}

// createOrder returns a ORDER BY clause for given columns.
func (sb *AdvancedSQLBuilder) createOrder() string {
	if sb.order != "" {
		return sb.order
	}
	var b strings.Builder
	for _, field := range sb.orderFields {
		desc := strings.HasPrefix(field, "-")
		if desc {
			field = field[1:]
		}
		if b.Len() > 0 {
			b.WriteByte(',')
		}
		switch field {
		case "tracking_id", "publisher", "id":
			b.WriteString("advisories.")
			b.WriteString(field)
		case "cvss_v2_score", "cvss_v3_score", "critical":
			b.WriteString("COALESCE(")
			b.WriteString(field)
			b.WriteString(",0)")
		case "ssvc":
			b.WriteString("ssvc_current.ssvc")
		case "version":
			// TODO: This is not optimal (SemVer).
			b.WriteString(
				`CASE WHEN version ~ '^[[:digit:]]+$' THEN version::int END`)
		default:
			b.WriteString(field)
		}

		if desc {
			b.WriteString(" DESC")
		} else {
			b.WriteString(" ASC")
		}
	}
	sb.order = b.String()
	return sb.order
}

// CreateQuery creates an SQL statement to query the documents
// table and the associated texts if needed.
// WARN: Make sure that the input is vetted against injections.
func (sb *AdvancedSQLBuilder) CreateQuery(
	limit, offset int64,
) string {
	var b strings.Builder

	b.WriteString("SELECT ")
	sb.projectionsWithCasts(&b, sb.fields)
	b.WriteString(" FROM ")
	sb.createFrom(&b)
	b.WriteString(" WHERE ")
	b.WriteString(sb.createWhere())

	if order := sb.createOrder(); order != "" {
		b.WriteString(" ORDER BY ")
		b.WriteString(order)
	}

	if limit >= 0 {
		b.WriteString(" LIMIT ")
		b.WriteString(strconv.FormatInt(limit, 10))
	}
	if offset > 0 {
		b.WriteString(" OFFSET ")
		b.WriteString(strconv.FormatInt(offset, 10))
	}

	return b.String()
}

// projectionsWithCasts joins given projection adding casts if needed.
func (sb *AdvancedSQLBuilder) projectionsWithCasts(b *strings.Builder, proj []string) {
	for i, p := range proj {
		if _, found := sb.ignoreFields[p]; found {
			continue
		}
		if i > 0 {
			b.WriteByte(',')
		}
		if alias, found := sb.aliases[p]; found {
			b.WriteString(`CASE WHEN length(`)
			b.WriteString(alias)
			b.WriteString(`)<= 200 THEN `)
			b.WriteString(alias)
			b.WriteString(` ELSE substring(`)
			b.WriteString(alias)
			b.WriteString(`, 0, 197)END||'...'AS `)
			b.WriteString(p)
			continue
		}
		switch p {
		case "tracking_id", "publisher":
			b.WriteString("advisories.")
			b.WriteString(p)
			b.WriteString(` AS `)
			b.WriteString(p)
		case "id":
			b.WriteString("documents.")
			b.WriteString(p)
			b.WriteString(` AS `)
			b.WriteString(p)
		case "state", "event", "tracking_status":
			b.WriteString(p)
			b.WriteString("::text")
		case "event_state":
			b.WriteString("events_log.state::text AS event_state")
		case "versions":
			b.WriteString(versionsCount + `AS versions`)
		case "ssvc":
			b.WriteString("ssvc_current.ssvc AS ssvc")
		case "comments":
			switch sb.mode() {
			case AdvisoryMode:
				b.WriteString(p)
			case DocumentMode:
				b.WriteString(commentsCountDocuments + `AS comments`)
			case EventMode:
				b.WriteString(commentsCountEvents + `AS comments`)
			}
		default:
			b.WriteString(p)
		}
	}
}

// check tests for the existence of used columns.
func (sb *AdvancedSQLBuilder) check() error {
	// check projections.
	for _, p := range sb.fields {
		if _, found := sb.aliases[p]; found {
			continue
		}
		if _, found := sb.ignoreFields[p]; found {
			continue
		}
		col := findDocumentColumn(p, sb.mode())
		if col == nil {
			return fmt.Errorf("column %q does not exists", p)
		}
		sb.usedSources.add(col.sources)
	}
	// check order
	for _, field := range sb.orderFields {
		if desc := strings.HasPrefix(field, "-"); desc {
			field = field[1:]
		}
		if _, found := sb.aliases[field]; !found {
			col := findDocumentColumn(field, sb.mode())
			if col == nil {
				return fmt.Errorf("order field %q does not exists", field)
			}
			sb.usedSources.add(col.sources)
		}
	}
	slog.Debug("advanced sqlbuilder", "used sources", sb.usedSources)
	return nil
}
