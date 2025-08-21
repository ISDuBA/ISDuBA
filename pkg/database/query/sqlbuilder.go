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
	"regexp"
	"strconv"
	"strings"
)

// SQLBuilder helps to construct a SQL query.
type SQLBuilder struct {
	WhereClause         string
	Replacements        []any
	replToIdx           map[string]int
	Aliases             map[string]string
	IgnoreFields        map[string]struct{}
	Mode                ParserMode
	TextTables          bool
	ReturnSearchResults bool
}

// CreateWhere construct a WHERE clause for a given expression.
func (sb *SQLBuilder) CreateWhere(e *Expr) string {
	var b strings.Builder
	sb.whereRecurse(e, &b)
	sb.WhereClause = b.String()
	return sb.WhereClause
}

var (
	escapeLike = strings.NewReplacer(
		`%`, `\%`,
		`_`, `\_`).Replace
	whiteSpaces = regexp.MustCompile(`\s+`)
)

// LikeEscape quotes a query string to be more convenient
// to use with LIKE filters.
func LikeEscape(query string) string {
	query = strings.TrimSpace(query)
	query = escapeLike(query)
	query = whiteSpaces.ReplaceAllString(query, `%`)
	return `%` + query + `%`
}

func (sb *SQLBuilder) searchWhere(e *Expr, b *strings.Builder) {
	if sb.ReturnSearchResults {
		fmt.Fprintf(b, "txt ILIKE $%d",
			sb.replacementIndex(LikeEscape(e.stringValue))+1)

		// We need the text tables to be joined.
		sb.TextTables = true

		// Handle alias
		if e.alias == "" {
			return
		}
		if sb.Aliases == nil {
			sb.Aliases = map[string]string{}
		}
		sb.Aliases[e.alias] = `txt`
	} else {
		switch sb.Mode {
		case AdvisoryMode, DocumentMode:
			fmt.Fprintf(b, "EXISTS(SELECT 1 FROM unique_texts "+
				"JOIN documents_texts ON unique_texts.id = documents_texts.txt_id "+
				"WHERE txt ILIKE $%d "+
				"AND documents_texts.documents_id = documents.id)", sb.replacementIndex(LikeEscape(e.stringValue))+1)
		case EventMode:
			// TODO clarify how to handle event search
		}

		// Ignore alias for now to avoid breaking change
		if e.alias == "" {
			return
		}
		if sb.IgnoreFields == nil {
			sb.IgnoreFields = map[string]struct{}{}
		}
		sb.IgnoreFields[e.alias] = struct{}{}
	}

}

func (sb *SQLBuilder) mentionedWhere(e *Expr, b *strings.Builder) {
	switch sb.Mode {
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

func (sb *SQLBuilder) involvedWhere(e *Expr, b *strings.Builder) {
	switch sb.Mode {
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
	case eventsType:
		b.WriteString("events")
	case statusType:
		b.WriteString("status")
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

const (
	versionsCount = `(SELECT count(*) FROM documents WHERE ` +
		`documents.advisories_id = advisories.id)`
	commentsCountDocuments = `(SELECT count(*) FROM comments WHERE ` +
		`comments.documents_id = documents.id)`
	commentsCountEvents = `(SELECT count(*) FROM comments WHERE ` +
		`comments.documents_id = documents_id)`
)

func (sb *SQLBuilder) accessWhere(e *Expr, b *strings.Builder) {
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
		switch sb.Mode {
		case AdvisoryMode:
			b.WriteString(column)
		case DocumentMode:
			b.WriteString(commentsCountDocuments)
		case EventMode:
			b.WriteString(commentsCountEvents)
		}
	case "event_state":
		b.WriteString("events_log.state")
	default:
		b.WriteString(column)
	}
}

func (sb *SQLBuilder) nowWhere(_ *Expr, b *strings.Builder) {
	b.WriteString("current_timestamp")
}

const (
	ilikePrefix = `'%'||regexp_replace(regexp_replace(`
	ilikeSuffix = `,'(%|_)','\\\1','g'),'(\s+)','%','g')||'%'`
)

func (sb *SQLBuilder) ilikeWhere(e *Expr, b *strings.Builder) {
	b.WriteByte('(')
	sb.whereRecurse(e.children[0], b)
	b.WriteString(` ILIKE ` + ilikePrefix)
	sb.whereRecurse(e.children[1], b)
	b.WriteString(ilikeSuffix + `)`)
}

func (sb *SQLBuilder) ilikePNameWhere(e *Expr, b *strings.Builder) {
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
func (sb *SQLBuilder) ilikePIDWhere(e *Expr, b *strings.Builder) {
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

func (sb *SQLBuilder) replacementIndex(s string) int {
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

func (sb *SQLBuilder) createFrom(b *strings.Builder) {
	switch sb.Mode {
	case AdvisoryMode, DocumentMode:
		b.WriteString(`documents ` +
			`JOIN advisories ON ` +
			`advisories.id = documents.advisories_id`)
	case EventMode:
		b.WriteString(`events_log JOIN documents ON events_log.documents_id = documents.id ` +
			`JOIN advisories ON advisories.id = documents.advisories_id ` +
			`LEFT JOIN (SELECT message, id FROM comments) AS comment ON events_log.comments_id = comment.id`)
	}

	if sb.TextTables {
		b.WriteString(` JOIN documents_texts ON documents.id = documents_texts.documents_id ` +
			`JOIN unique_texts ON documents_texts.txt_id = unique_texts.id`)
	}
}

// CreateCountSQL returns an SQL count statement to count
// the number of rows which are possible to fetch by the
// given filter.
func (sb *SQLBuilder) CreateCountSQL() string {
	var b strings.Builder
	b.WriteString("SELECT count(*) FROM ")
	sb.createFrom(&b)
	b.WriteString(" WHERE ")
	b.WriteString(sb.WhereClause)
	return b.String()
}

// CreateOrder returns a ORDER BY clause for given columns.
func (sb *SQLBuilder) CreateOrder(fields []string) (string, error) {
	var b strings.Builder
	for _, field := range fields {
		desc := strings.HasPrefix(field, "-")
		if desc {
			field = field[1:]
		}
		if _, found := sb.Aliases[field]; !found && !ExistsDocumentColumn(field, sb.Mode) {
			return "", fmt.Errorf("order field %q does not exists", field)
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
	return b.String(), nil
}

// CreateQuery creates an SQL statement to query the documents
// table and the associated texts if needed.
// WARN: Make sure that the input is vetted against injections.
func (sb *SQLBuilder) CreateQuery(
	fields []string,
	order string,
	limit, offset int64,
) string {
	var b strings.Builder

	b.WriteString("SELECT ")
	sb.projectionsWithCasts(&b, fields)
	b.WriteString(" FROM ")
	sb.createFrom(&b)
	b.WriteString(" WHERE ")
	b.WriteString(sb.WhereClause)

	if order != "" {
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
func (sb *SQLBuilder) projectionsWithCasts(b *strings.Builder, proj []string) {
	for i, p := range proj {
		if _, found := sb.IgnoreFields[p]; found {
			continue
		}
		if i > 0 {
			b.WriteByte(',')
		}
		if alias, found := sb.Aliases[p]; found {
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
		case "comments":
			switch sb.Mode {
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

// CheckProjections checks if the requested projections are valid.
func (sb *SQLBuilder) CheckProjections(proj []string) error {
	for _, p := range proj {
		if _, found := sb.Aliases[p]; found {
			continue
		}
		if _, found := sb.IgnoreFields[p]; found {
			continue
		}
		if !ExistsDocumentColumn(p, sb.Mode) {
			return fmt.Errorf("column %q does not exists", p)
		}
	}
	return nil
}
