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
	"slices"
	"strconv"
	"strings"
)

// AdvancedSQLBuilder helps to construct a SQL query.
type AdvancedSQLBuilder struct {
	expr         *Expr
	parser       *AdvancedParser
	orderFields  []string
	fields       []string
	Replacements []any
	replToIdx    map[string]int
	usedSources  columnSource
}

type statementMode interface {
	projection(sb *AdvancedSQLBuilder, b *strings.Builder, name string)
	from(sb *AdvancedSQLBuilder, b *strings.Builder)
	accessWhere(sb *AdvancedSQLBuilder, e *Expr, b *strings.Builder)
	searchWhere(sb *AdvancedSQLBuilder, e *Expr, b *strings.Builder)
	mentionedWhere(sb *AdvancedSQLBuilder, e *Expr, b *strings.Builder)
	involvedWhere(sb *AdvancedSQLBuilder, e *Expr, b *strings.Builder)
	ilikePNameWhere(sb *AdvancedSQLBuilder, e *Expr, b *strings.Builder)
	ilikePIDWhere(sb *AdvancedSQLBuilder, e *Expr, b *strings.Builder)
	order(sb *AdvancedSQLBuilder, b *strings.Builder, name string)
}

type (
	classicMode struct{}
	cteMode     struct{ classicMode }
)

func (classicMode) projectionCommon(
	sb *AdvancedSQLBuilder, b *strings.Builder,
	name string,
	versionsCount, commentsCountDocuments string,
) {
	switch name {
	case "state", "event", "tracking_status":
		b.WriteString(name)
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
			b.WriteString(name)
		case DocumentMode:
			b.WriteString(commentsCountDocuments + `AS comments`)
		case EventMode:
			b.WriteString(commentsCountEvents + `AS comments`)
		}
	default:
		b.WriteString(name)
	}
}

func (cm classicMode) projection(sb *AdvancedSQLBuilder, b *strings.Builder, name string) {
	switch name {
	case "tracking_id", "publisher":
		b.WriteString("advisories.")
		b.WriteString(name)
		b.WriteString(` AS `)
		b.WriteString(name)
	case "id":
		b.WriteString("documents.")
		b.WriteString(name)
		b.WriteString(` AS `)
		b.WriteString(name)
	default:
		cm.projectionCommon(sb, b, name,
			versionsCountClassic, commentsCountDocumentsClassic)
	}
}

func (cm cteMode) projection(sb *AdvancedSQLBuilder, b *strings.Builder, name string) {
	switch name {
	case "tracking_id", "publisher":
		b.WriteString("docads.")
		b.WriteString(name)
		b.WriteString(` AS `)
		b.WriteString(name)
	case "id":
		b.WriteString("docads.")
		b.WriteString(name)
		b.WriteString(` AS `)
		b.WriteString(name)
	default:
		cm.projectionCommon(sb, b, name,
			versionsCountCTE, commentsCountDocumentsCTE)
	}
}

func (classicMode) from(sb *AdvancedSQLBuilder, b *strings.Builder) {
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
	if sb.usedSources.contains(ssvcHistoryTable) {
		b.WriteString(` LEFT JOIN LATERAL ( ` +
			`SELECT ssvc FROM ssvc_history ` +
			`WHERE documents_id = documents.id ` +
			`ORDER BY changedate DESC, change_number DESC LIMIT 1 ` +
			`) AS ssvc_current ON TRUE`)
	}
	if sb.usedSources.contains(textTable) {
		b.WriteString(` JOIN documents_texts ON documents.id = documents_texts.documents_id ` +
			`JOIN unique_texts ON documents_texts.txt_id = unique_texts.id`)
	}
}

func (cteMode) from(sb *AdvancedSQLBuilder, b *strings.Builder) {
	switch sb.mode() {
	case AdvisoryMode, DocumentMode:
		b.WriteString(`docads`)
	case EventMode:
		b.WriteString(`events_log JOIN docads ON events_log.documents_id = docads.id ` +
			`LEFT JOIN (SELECT message, id FROM comments) AS comment ` +
			`ON events_log.comments_id = comment.id`)
	}
	// Add SSVC if exists
	if sb.usedSources.contains(ssvcHistoryTable) {
		b.WriteString(` LEFT JOIN LATERAL ( ` +
			`SELECT ssvc FROM ssvc_history ` +
			`WHERE documents_id = docads.id ` +
			`ORDER BY changedate DESC, change_number DESC LIMIT 1 ` +
			`) AS ssvc_current ON TRUE`)
	}
	if sb.usedSources.contains(textTable) {
		b.WriteString(` JOIN documents_texts ON docads.id = documents_texts.documents_id ` +
			`JOIN unique_texts ON documents_texts.txt_id = unique_texts.id`)
	}
}

func (classicMode) accessWhereCommon(
	sb *AdvancedSQLBuilder, e *Expr, b *strings.Builder,
	versionsCount, commentsCountDocuments string,
) {
	switch column := e.stringValue; column {
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

func (cm classicMode) accessWhere(sb *AdvancedSQLBuilder, e *Expr, b *strings.Builder) {
	switch column := e.stringValue; column {
	case "id":
		b.WriteString("documents.")
		b.WriteString(column)
	case "tracking_id", "publisher":
		b.WriteString("advisories.")
		b.WriteString(column)
	default:
		cm.accessWhereCommon(sb, e, b,
			versionsCountClassic, commentsCountDocumentsClassic)
	}
}

func (cm cteMode) accessWhere(sb *AdvancedSQLBuilder, e *Expr, b *strings.Builder) {
	switch column := e.stringValue; column {
	case "id":
		b.WriteString("docads.")
		b.WriteString(column)
	case "tracking_id", "publisher":
		b.WriteString("docads.")
		b.WriteString(column)
	default:
		cm.classicMode.accessWhereCommon(sb, e, b,
			versionsCountCTE, commentsCountDocumentsCTE)
	}
}

func (classicMode) searchWhere(sb *AdvancedSQLBuilder, e *Expr, b *strings.Builder) {
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
}

func (cteMode) searchWhere(sb *AdvancedSQLBuilder, e *Expr, b *strings.Builder) {
	switch sb.mode() {
	case AdvisoryMode, DocumentMode:
		fmt.Fprintf(b, "EXISTS(SELECT 1 FROM documents_texts "+
			"JOIN unique_texts ON unique_texts.id = documents_texts.txt_id "+
			"WHERE txt ILIKE $%d "+
			"AND documents_texts.documents_id = docads.id)",
			sb.replacementIndex(LikeEscape(e.stringValue))+1)
	case EventMode:
		// TODO clarify how to handle event search
	}
}

func (classicMode) mentionedWhereCommon(sb *AdvancedSQLBuilder, e *Expr, b *strings.Builder) {
	switch sb.mode() {
	case EventMode:
		fmt.Fprintf(b, "EXISTS(SELECT 1 FROM comments WHERE message ILIKE $%d "+
			"AND comments.id = events_log.comments_id)",
			sb.replacementIndex(LikeEscape(e.stringValue))+1)
	}
}

func (cm classicMode) mentionedWhere(sb *AdvancedSQLBuilder, e *Expr, b *strings.Builder) {
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
	default:
		cm.mentionedWhereCommon(sb, e, b)
	}
}

func (cm cteMode) mentionedWhere(sb *AdvancedSQLBuilder, e *Expr, b *strings.Builder) {
	switch sb.mode() {
	case AdvisoryMode:
		fmt.Fprintf(b, "EXISTS(SELECT 1 FROM comments "+
			"JOIN docads ON comments.documents_id = docads.id "+
			"WHERE message ILIKE $%d "+
			"AND docads.advisories_id = docads.advisories_id)",
			sb.replacementIndex(LikeEscape(e.stringValue))+1)
	case DocumentMode:
		fmt.Fprintf(b, "EXISTS(SELECT 1 FROM comments WHERE message ILIKE $%d "+
			"AND comments.documents_id = docads.id)",
			sb.replacementIndex(LikeEscape(e.stringValue))+1)
	default:
		cm.mentionedWhereCommon(sb, e, b)
	}
}

func (classicMode) involvedWhere(sb *AdvancedSQLBuilder, e *Expr, b *strings.Builder) {
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

func (cteMode) involvedWhere(sb *AdvancedSQLBuilder, e *Expr, b *strings.Builder) {
	switch sb.mode() {
	case AdvisoryMode, EventMode:
		fmt.Fprintf(b, "EXISTS(SELECT 1 FROM events_log JOIN documents docads "+
			"ON events_log.documents_id = docads.id "+
			"WHERE actor = $%d)",
			sb.replacementIndex(e.stringValue)+1)
	case DocumentMode:
		fmt.Fprintf(b, "EXISTS(SELECT 1 FROM events_log WHERE actor = $%d "+
			"AND events_log.documents_id = docads.id)",
			sb.replacementIndex(e.stringValue)+1)
	}
}

func (cm classicMode) ilikePNameWhere(sb *AdvancedSQLBuilder, e *Expr, b *strings.Builder) {
	b.WriteString(`EXISTS (` +
		`WITH product_names AS (SELECT jsonb_path_query(` +
		`document, '$.product_tree.**.product.name')::int num ` +
		`FROM documents ds WHERE ds.id = documents.id)` +
		`SELECT * FROM documents_texts dts JOIN product_names ` +
		`ON product_names.num = dts.num JOIN unique_texts ON dts.txt_id = unique_texts.id ` +
		`WHERE dts.documents_id = documents.id AND ` +
		`unique_texts.txt ILIKE ` + ilikePrefix)
	sb.whereRecurse(e.children[0], b, cm)
	b.WriteString(ilikeSuffix + `)`)
}

func (cm cteMode) ilikePNameWhere(sb *AdvancedSQLBuilder, e *Expr, b *strings.Builder) {
	b.WriteString(`EXISTS (` +
		`WITH product_names AS (SELECT jsonb_path_query(` +
		`document, '$.product_tree.**.product.name')::int num ` +
		`FROM docads ds WHERE ds.id = docads.id)` +
		`SELECT * FROM documents_texts dts JOIN product_names ` +
		`ON product_names.num = dts.num JOIN unique_texts ON dts.txt_id = unique_texts.id ` +
		`WHERE dts.documents_id = docads.id AND ` +
		`unique_texts.txt ILIKE ` + ilikePrefix)
	sb.whereRecurse(e.children[0], b, cm)
	b.WriteString(ilikeSuffix + `)`)
}

func (cm classicMode) ilikePIDWhere(sb *AdvancedSQLBuilder, e *Expr, b *strings.Builder) {
	b.WriteString(`EXISTS (` +
		`WITH product_ids AS (SELECT jsonb_path_query(` +
		`document, '$.product_tree.**.product.product_id')::int num ` +
		`FROM documents ds WHERE ds.id = documents.id)` +
		`SELECT * FROM documents_texts dts JOIN product_ids ` +
		`ON product_ids.num = dts.num JOIN unique_texts ON dts.txt_id = unique_texts.id ` +
		`WHERE dts.documents_id = documents.id AND ` +
		`unique_texts.txt ILIKE ` + ilikePrefix)
	sb.whereRecurse(e.children[0], b, cm)
	b.WriteString(ilikeSuffix + `)`)
}

func (cm cteMode) ilikePIDWhere(sb *AdvancedSQLBuilder, e *Expr, b *strings.Builder) {
	b.WriteString(`EXISTS (` +
		`WITH product_ids AS (SELECT jsonb_path_query(` +
		`document, '$.product_tree.**.product.product_id')::int num ` +
		`FROM docads ds WHERE ds.id = docads.id)` +
		`SELECT * FROM documents_texts dts JOIN product_ids ` +
		`ON product_ids.num = dts.num JOIN unique_texts ON dts.txt_id = unique_texts.id ` +
		`WHERE dts.documents_id = docads.id AND ` +
		`unique_texts.txt ILIKE ` + ilikePrefix)
	sb.whereRecurse(e.children[0], b, cm)
	b.WriteString(ilikeSuffix + `)`)
}

func (classicMode) orderCommon(b *strings.Builder, name string) {
	switch name {
	case "cvss_v2_score", "cvss_v3_score", "critical":
		b.WriteString("COALESCE(")
		b.WriteString(name)
		b.WriteString(",0)")
	case "ssvc":
		b.WriteString("ssvc_current.ssvc")
	case "version":
		// TODO: This is not optimal (SemVer).
		b.WriteString(
			`CASE WHEN version ~ '^[[:digit:]]+$' THEN version::int END`)
	default:
		b.WriteString(name)
	}
}

func (cm classicMode) order(_ *AdvancedSQLBuilder, b *strings.Builder, name string) {
	switch name {
	case "tracking_id", "publisher", "id":
		b.WriteString("advisories.")
		b.WriteString(name)
	default:
		cm.orderCommon(b, name)
	}
}

func (cm cteMode) order(_ *AdvancedSQLBuilder, b *strings.Builder, name string) {
	switch name {
	case "tracking_id", "publisher", "id":
		b.WriteString("docads.")
		b.WriteString(name)
	default:
		cm.orderCommon(b, name)
	}
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

// Fields returns the projection fields of the query.
func (sb *AdvancedSQLBuilder) Fields() []string {
	return sb.fields
}

func (sb *AdvancedSQLBuilder) alias(name string) *Expr {
	if sb.parser == nil {
		return nil
	}
	return sb.parser.aliases[name]
}

// createWhere construct a WHERE clause for a given expression.
func (sb *AdvancedSQLBuilder) createWhere(b *strings.Builder, sm statementMode) {
	sb.whereRecurse(sb.expr, b, sm)
}

func (sb *AdvancedSQLBuilder) mode() ParserMode {
	if sb.parser != nil {
		return sb.parser.Mode
	}
	return DocumentMode
}

func (sb *AdvancedSQLBuilder) castWhere(e *Expr, b *strings.Builder, sm statementMode) {
	b.WriteString("CAST(")
	sb.whereRecurse(e.children[0], b, sm)
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

func (sb *AdvancedSQLBuilder) binaryWhere(e *Expr, b *strings.Builder, op string, sm statementMode) {
	b.WriteByte('(')
	sb.whereRecurse(e.children[0], b, sm)
	b.WriteString(op)
	sb.whereRecurse(e.children[1], b, sm)
	b.WriteByte(')')
}

func (sb *AdvancedSQLBuilder) notWhere(e *Expr, b *strings.Builder, sm statementMode) {
	b.WriteString("(NOT ")
	sb.whereRecurse(e.children[0], b, sm)
	b.WriteByte(')')
}

func (sb *AdvancedSQLBuilder) nowWhere(b *strings.Builder) {
	b.WriteString("current_timestamp")
}

func (sb *AdvancedSQLBuilder) ilikeWhere(e *Expr, b *strings.Builder, sm statementMode) {
	b.WriteByte('(')
	sb.whereRecurse(e.children[0], b, sm)
	b.WriteString(` ILIKE ` + ilikePrefix)
	sb.whereRecurse(e.children[1], b, sm)
	b.WriteString(ilikeSuffix + `)`)
}

func (sb *AdvancedSQLBuilder) whereRecurse(e *Expr, b *strings.Builder, sm statementMode) {
	if e == nil {
		return
	}
	b.WriteByte('(')
	switch e.exprType {
	case access:
		sm.accessWhere(sb, e, b)
	case cnst:
		sb.cnstWhere(e, b)
	case cast:
		sb.castWhere(e, b, sm)
	case eq:
		sb.binaryWhere(e, b, "=", sm)
	case ne:
		sb.binaryWhere(e, b, "<>", sm)
	case lt:
		sb.binaryWhere(e, b, "<", sm)
	case gt:
		sb.binaryWhere(e, b, ">", sm)
	case le:
		sb.binaryWhere(e, b, "<=", sm)
	case ge:
		sb.binaryWhere(e, b, ">=", sm)
	case not:
		sb.notWhere(e, b, sm)
	case and:
		sb.binaryWhere(e, b, "AND", sm)
	case or:
		sb.binaryWhere(e, b, "OR", sm)
	case search:
		sm.searchWhere(sb, e, b)
	case mentioned:
		sm.mentionedWhere(sb, e, b)
	case involved:
		sm.involvedWhere(sb, e, b)
	case ilike:
		sb.ilikeWhere(e, b, sm)
	case ilikePName:
		sm.ilikePNameWhere(sb, e, b)
	case ilikePID:
		sm.ilikePIDWhere(sb, e, b)
	case now:
		sb.nowWhere(b)
	case add:
		sb.binaryWhere(e, b, "+", sm)
	case sub:
		sb.binaryWhere(e, b, "-", sm)
	case mul:
		sb.binaryWhere(e, b, "*", sm)
	case div:
		sb.binaryWhere(e, b, "/", sm)
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

// CreateCountSQL returns an SQL count statement to count
// the number of rows which are possible to fetch by the
// given filter.
func (sb *AdvancedSQLBuilder) CreateCountSQL() string {
	var b strings.Builder
	sm := statementMode(classicMode{})
	if sb.usedSources.contains(documentsTable | advisoriesTable) {
		sm = cteMode{}
		sb.prefixCTE(&b)
	}
	b.WriteString("SELECT count(*) FROM ")
	sm.from(sb, &b)
	b.WriteString(" WHERE ")
	sb.createWhere(&b, sm)
	return b.String()
}

func (sb *AdvancedSQLBuilder) prefixCTE(b *strings.Builder) {
	b.WriteString(`WITH docads AS (` +
		`SELECT `)
	for i, field := range enumerate(
		unique(
			slices.Values(sb.fields),
			apply(
				slices.Values(sb.orderFields),
				func(s string) string { return strings.TrimPrefix(s, "-") }),
			sb.expr.accessedColumns(),
		)) {
		if i > 0 {
			b.WriteByte(',')
		}
		switch field {
		case "id":
			b.WriteString("documents.id AS id")
		case "versions":
			b.WriteString(versionsCountClassic + ` AS versions`)
		case "ssvc":
			b.WriteString(`(` +
				`SELECT ssvc FROM ssvc_history ` +
				`WHERE documents_id = documents.id ` +
				`ORDER BY changedate DESC, change_number DESC LIMIT 1)`)
		default:
			b.WriteString(field)
		}
	}
	b.WriteString(` FROM documents JOIN advisories` +
		` ON documents.advisories_id = advisories.id)`)
}

// CreateQuery creates an SQL statement to query the documents
// table and the associated texts if needed.
// WARN: Make sure that the input is vetted against injections.
func (sb *AdvancedSQLBuilder) CreateQuery(
	limit, offset int64,
) string {
	var b strings.Builder
	sm := statementMode(classicMode{})
	if sb.usedSources.contains(documentsTable | advisoriesTable) {
		sm = cteMode{}
		sb.prefixCTE(&b)
	}

	b.WriteString("SELECT ")
	sb.createProjectionsWithCasts(&b, sm)
	b.WriteString(" FROM ")
	sm.from(sb, &b)
	b.WriteString(" WHERE ")
	sb.createWhere(&b, sm)

	if len(sb.orderFields) > 0 {
		b.WriteString(" ORDER BY ")
		sb.createOrder(&b, sm)
	}

	if limit >= 0 {
		b.WriteString(" LIMIT ")
		b.WriteString(strconv.FormatInt(limit, 10))
	}
	if offset > 0 {
		b.WriteString(" OFFSET ")
		b.WriteString(strconv.FormatInt(offset, 10))
	}

	query := b.String()
	slog.Debug("sql builder", "query", query)
	return query
}

// createOrder returns a ORDER BY clause for given columns.
func (sb *AdvancedSQLBuilder) createOrder(b *strings.Builder, sm statementMode) {
	for i, field := range sb.orderFields {
		desc := strings.HasPrefix(field, "-")
		if desc {
			field = field[1:]
		}
		if i > 0 {
			b.WriteByte(',')
		}
		sm.order(sb, b, field)
		if desc {
			b.WriteString(" DESC")
		} else {
			b.WriteString(" ASC")
		}
	}
}

// createProjectionsWithCasts joins given projection adding casts if needed.
func (sb *AdvancedSQLBuilder) createProjectionsWithCasts(b *strings.Builder, sm statementMode) {
	for i, name := range sb.fields {
		if i > 0 {
			b.WriteByte(',')
		}
		if alias := sb.alias(name); alias != nil {
			txt := fmt.Sprintf("txt ILIKE $%d",
				sb.replacementIndex(alias.stringValue)+1)
			b.WriteString(`CASE WHEN length(`)
			b.WriteString(txt)
			b.WriteString(`)<= 200 THEN `)
			b.WriteString(txt)
			b.WriteString(` ELSE substring(`)
			b.WriteString(txt)
			b.WriteString(`, 0, 197)END||'...'AS `)
			b.WriteString(name)
			continue
		}
		sm.projection(sb, b, name)
	}
}

// check tests for the existence of used columns.
func (sb *AdvancedSQLBuilder) check() error {
	// A field is valid if there is a named column or
	// there exists an alias for it.
	checkField := func(field string) bool {
		col := findDocumentColumn(field, sb.mode())
		if col != nil {
			sb.usedSources.add(col.sources)
			return true
		}
		if sb.alias(field) != nil {
			sb.usedSources.add(documentsTable | textTable)
			return true
		}
		return false
	}
	// check projections.
	for _, f := range sb.fields {
		if !checkField(f) {
			return fmt.Errorf("column %q does not exists", f)
		}
	}
	// check order
	for _, f := range sb.orderFields {
		if f := strings.TrimPrefix(f, "-"); !checkField(f) {
			return fmt.Errorf("order field %q does not exists", f)
		}
	}
	slog.Debug("advanced sqlbuilder", "used sources", sb.usedSources)
	return nil
}
