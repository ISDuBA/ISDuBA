// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package query

import (
	"errors"
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gocsaf/csaf/v3/csaf"
)

// ParserMode represents the operation mode of the parser.
type ParserMode int

const (
	// DocumentMode operates on documents.
	DocumentMode ParserMode = iota
	// AdvisoryMode operates on advisories.
	AdvisoryMode
	// EventMode operates on events.
	EventMode
)

// Parser helps parsing database queries,
type Parser struct {
	// Mode indicates that only advisories should be considered.
	Mode ParserMode
	// MinSearchLength enforces a minimal lengths of search phrases.
	MinSearchLength int
	// Me is a replacement text for the "me" keyword.
	Me string

	// UsedSources are the sources found during parsing.
	UsedSources columnSource

	// aliases defined by 'as'.
	aliases map[string]*Expr
}

// columnSource is a type to accumulate the tables needed to build an SQL query.
type columnSource int

const (
	// noTable means that no table is needed.
	noTable columnSource = 0
	// documentsTable means that the documents table is needed.
	documentsTable columnSource = 1 << iota
	// advisoriesTable means that the advisories table is needed.
	advisoriesTable
	// textTable means that the tables for fulltext search are needed.
	textTable
	// ssvcHistoryTable means that the ssvc_history table is needed.
	ssvcHistoryTable
	// eventsLogTable means that the events_log table is needed.
	eventsLogTable
	// commentsTable means that the comments table is needed.
	commentsTable
)

var columnSourceTables = [...]string{
	noTable:          "",
	documentsTable:   "documents",
	advisoriesTable:  "advisories",
	textTable:        "text",
	ssvcHistoryTable: "ssvc_history",
	eventsLogTable:   "events_log",
	commentsTable:    "comments",
}

type documentColumn struct {
	name           string
	valueType      valueType
	modes          []ParserMode
	projectionOnly bool
	sources        columnSource
}

type binaryCompat struct {
	left     valueType
	operator exprType
	right    valueType
}

type stack []*Expr

// documentColumns are the documentColumns which can be accessed.
var documentColumns = []documentColumn{
	{"id", intType, docAdvEvtModes, false, documentsTable},
	{"latest", boolType, docAdvEvtModes, false, documentsTable},
	{"tracking_id", stringType, docAdvEvtModes, false, advisoriesTable},
	{"version", stringType, docAdvEvtModes, false, documentsTable},
	{"publisher", stringType, docAdvEvtModes, false, advisoriesTable},
	{"current_release_date", timeType, docAdvEvtModes, false, documentsTable},
	{"initial_release_date", timeType, docAdvEvtModes, false, documentsTable},
	{"rev_history_length", intType, docAdvEvtModes, false, documentsTable},
	{"title", stringType, docAdvEvtModes, false, documentsTable},
	{"tlp", stringType, docAdvEvtModes, false, documentsTable},
	{"ssvc", stringType, docAdvEvtModes, false, ssvcHistoryTable},
	{"cvss_v2_score", floatType, docAdvEvtModes, false, documentsTable},
	{"cvss_v3_score", floatType, docAdvEvtModes, false, documentsTable},
	{"critical", floatType, docAdvEvtModes, false, documentsTable},
	{"four_cves", stringType, docAdvEvtModes, true, documentsTable},
	{"comments", intType, docAdvEvtModes, false, documentsTable},
	{"tracking_status", statusType, docAdvEvtModes, false, documentsTable},
	// Advisories only
	{"state", workflowType, advModes, false, advisoriesTable},
	{"recent", timeType, advModes, false, advisoriesTable},
	// ToDo: Column "versions" does not exist, but table versions does?
	{"versions", intType, advModes, false, documentsTable | advisoriesTable},
	// Events only
	{"event", eventsType, evtsModes, false, eventsLogTable},
	{"event_state", workflowType, evtsModes, false, eventsLogTable},
	{"time", timeType, evtsModes, false, eventsLogTable},
	{"actor", stringType, evtsModes, false, eventsLogTable},
	{"comments_id", intType, evtsModes, false, commentsTable},
	{"message", stringType, evtsModes, false, commentsTable},
}

var (
	// baseAction are the action available in every parser.
	baseAction = map[string]func(*Parser, *stack){
		"true":       (*Parser).pushTrue,
		"false":      (*Parser).pushFalse,
		"not":        (*Parser).pushNot,
		"and":        curry3((*Parser).pushBinary, and),
		"or":         curry3((*Parser).pushBinary, or),
		"float":      (*Parser).pushFloat,
		"integer":    (*Parser).pushInteger,
		"timestamp":  (*Parser).pushTimestamp,
		"workflow":   pushEnum(workflowType, parseWorkflow),
		"events":     pushEnum(eventsType, parseEvents),
		"status":     pushEnum(statusType, parseStatus),
		"=":          curry3((*Parser).pushCmp, eq),
		"!=":         curry3((*Parser).pushCmp, ne),
		"<":          curry3((*Parser).pushCmp, lt),
		"<=":         curry3((*Parser).pushCmp, le),
		">":          curry3((*Parser).pushCmp, gt),
		">=":         curry3((*Parser).pushCmp, ge),
		"ilike":      (*Parser).pushILike,
		"ilikepname": (*Parser).pushILikePName,
		"ilikepid":   (*Parser).pushILikePID,
		"now":        (*Parser).pushNow,
		"duration":   (*Parser).pushDuration,
		"+":          curry3((*Parser).pushBinary, add),
		"-":          curry3((*Parser).pushBinary, sub),
		"/":          curry3((*Parser).pushBinary, div),
		"*":          curry3((*Parser).pushBinary, mul),
		"me":         (*Parser).pushMe,
		"mentioned":  (*Parser).pushMentioned,
		"involved":   (*Parser).pushInvolved,
		"search":     (*Parser).pushSearch,
		"as":         (*Parser).pushAs,
	}
	// action is for fast looking up actions along the parser mode.
	action = map[ParserMode]map[string]func(*Parser, *stack){
		DocumentMode: buildActions(DocumentMode),
		AdvisoryMode: buildActions(AdvisoryMode),
		EventMode:    buildActions(EventMode),
	}
)

var binaryCompatMatrix = map[binaryCompat]valueType{
	{boolType, and, boolType}:         boolType,
	{boolType, or, boolType}:          boolType,
	{intType, add, intType}:           intType,
	{intType, sub, intType}:           intType,
	{intType, mul, intType}:           intType,
	{intType, div, intType}:           intType,
	{intType, add, floatType}:         floatType,
	{intType, sub, floatType}:         floatType,
	{intType, mul, floatType}:         floatType,
	{intType, div, floatType}:         floatType,
	{floatType, add, intType}:         floatType,
	{floatType, sub, intType}:         floatType,
	{floatType, mul, intType}:         floatType,
	{floatType, div, intType}:         floatType,
	{timeType, add, durationType}:     timeType,
	{timeType, sub, durationType}:     timeType,
	{durationType, add, timeType}:     timeType,
	{durationType, sub, timeType}:     timeType,
	{durationType, add, durationType}: durationType,
	{durationType, sub, durationType}: durationType,
	{durationType, mul, intType}:      durationType,
	{durationType, div, intType}:      durationType,
	{intType, mul, durationType}:      durationType,
	{durationType, mul, floatType}:    durationType,
	{durationType, div, floatType}:    durationType,
}

var (
	docAdvEvtModes = []ParserMode{DocumentMode, AdvisoryMode, EventMode}
	advModes       = []ParserMode{AdvisoryMode}
	evtsModes      = []ParserMode{EventMode}
)

func curry3[A, B, C any](fn func(A, B, C), c C) func(A, B) {
	return func(a A, b B) { fn(a, b, c) }
}

// UnmarshalText implements [encoding.TextUnmarshaler].
func (pm *ParserMode) UnmarshalText(text []byte) error {
	switch s := string(text); s {
	case "advisories":
		*pm = AdvisoryMode
	case "documents":
		*pm = DocumentMode
	case "events":
		*pm = EventMode
	default:
		return fmt.Errorf("unknown parser mode %q", s)
	}
	return nil
}

// MarshalText implements [encoding.TextMarshaler].
func (pm ParserMode) MarshalText() ([]byte, error) {
	switch pm {
	case AdvisoryMode:
		return []byte("advisories"), nil
	case DocumentMode:
		return []byte("documents"), nil
	case EventMode:
		return []byte("events"), nil
	default:
		return nil, fmt.Errorf("unknown parser mode %d", pm)
	}
}

// String implements [fmt.Stringer].
func (pm ParserMode) String() string {
	switch pm {
	case DocumentMode:
		return "documents"
	case AdvisoryMode:
		return "advisories"
	case EventMode:
		return "events"
	default:
		return fmt.Sprintf("Unknown parser mode: %d", pm)
	}
}

// Scan implements [sql.Scanner].
func (pm *ParserMode) Scan(src any) error {
	if s, ok := src.(string); ok {
		return pm.UnmarshalText([]byte(s))
	}
	return errors.New("unsupported type")
}

func (cs columnSource) contains(other columnSource) bool {
	return cs&other == other
}

func (cs *columnSource) add(other columnSource) {
	*cs |= other
}

// String implements [fmt.Stringer].
func (cs columnSource) String() string {
	var tables []string
	for i := range len(columnSourceTables) - 1 {
		mask := columnSource(1) << i
		if cs.contains(mask) {
			tables = append(tables, columnSourceTables[i+1])
			cs &= ^mask
		}
	}
	if cs != noTable {
		tables = append(tables, fmt.Sprintf("unknown table(s): %b", cs))
	}
	return strings.Join(tables, "|")
}

func existsDocumentColumn(name string, mode ParserMode) bool {
	return findDocumentColumn(name, mode) != nil
}

// findDocumentColumn returns a column if it exists.
func findDocumentColumn(name string, mode ParserMode) *documentColumn {
	for i := range documentColumns {
		if col := &documentColumns[i]; col.name == name {
			if !slices.Contains(col.modes, mode) {
				return nil
			}
			return col
		}
	}
	return nil
}

func buildActions(mode ParserMode) map[string]func(*Parser, *stack) {
	actions := maps.Clone(baseAction)
	for i := range documentColumns {
		col := &documentColumns[i]
		if !col.projectionOnly && slices.Contains(col.modes, mode) {
			actions["$"+col.name] = func(p *Parser, st *stack) {
				p.pushAccess(st, col)
			}
		}
	}
	return actions
}

func (*Parser) pushTrue(st *stack)  { st.push(True()) }
func (*Parser) pushFalse(st *stack) { st.push(False()) }

func (*Parser) pushNot(st *stack) {
	e := st.pop()
	e.checkValueType(boolType)
	st.push(&Expr{
		exprType:  not,
		valueType: boolType,
		children:  []*Expr{e},
	})
}

func (*Parser) pushBinary(st *stack, et exprType) {
	right := st.pop()
	left := st.pop()
	resultValueType, ok := binaryCompatMatrix[binaryCompat{
		left:     left.valueType,
		operator: et,
		right:    right.valueType,
	}]
	if !ok {
		panic(parseError(
			fmt.Sprintf("invalid binary operation: %q %q %q",
				left.valueType, et, right.valueType)))
	}
	st.push(&Expr{
		exprType:  et,
		valueType: resultValueType,
		children:  []*Expr{left, right},
	})
}

func (p *Parser) pushAccess(st *stack, column *documentColumn) {
	p.UsedSources.add(column.sources)
	st.push(&Expr{
		exprType:    access,
		valueType:   column.valueType,
		stringValue: column.name,
	})
}

func (*Parser) pushFloat(st *stack) {
	if st.top().valueType == floatType {
		return
	}
	switch e := st.pop(); e.exprType {
	case cnst:
		switch e.valueType {
		case stringType:
			st.push(&Expr{
				exprType:   cnst,
				valueType:  floatType,
				floatValue: parseFloat(e.stringValue),
			})
		case intType:
			st.push(&Expr{
				exprType:   cnst,
				valueType:  intType,
				floatValue: float64(e.intValue),
			})
		}
	default:
		switch e.valueType {
		case stringType, intType:
			st.push(&Expr{
				exprType:  cast,
				valueType: floatType,
				children:  []*Expr{e},
			})
		default:
			panic(parseError("unsupported cast"))
		}
	}
}

func (*Parser) pushInteger(st *stack) {
	if st.top().valueType == intType {
		return
	}
	switch e := st.pop(); e.exprType {
	case cnst:
		switch e.valueType {
		case stringType:
			st.push(&Expr{
				exprType:  cnst,
				valueType: intType,
				intValue:  parseInt(e.stringValue),
			})
		case floatType:
			st.push(&Expr{
				exprType:  cnst,
				valueType: intType,
				intValue:  int64(e.floatValue),
			})
		}
	default:
		switch e.valueType {
		case stringType, floatType:
			st.push(&Expr{
				exprType:  cast,
				valueType: intType,
				children:  []*Expr{e},
			})
		default:
			panic(parseError("unsupported cast"))
		}
	}
}

func (*Parser) pushTimestamp(st *stack) {
	if st.top().valueType == timeType {
		return
	}
	switch e := st.pop(); e.exprType {
	case cnst:
		switch e.valueType {
		case stringType:
			st.push(&Expr{
				exprType:  cnst,
				valueType: timeType,
				timeValue: parseTime(e.stringValue),
			})
		}
	default:
		switch e.valueType {
		case stringType:
			st.push(&Expr{
				exprType:  cast,
				valueType: timeType,
				children:  []*Expr{e},
			})
		default:
			panic(parseError("unsupported cast"))
		}
	}
}

func (*Parser) pushCmp(st *stack, et exprType) {
	right := st.pop()
	left := st.pop()
	if right.valueType != left.valueType {
		panic(parseError(
			fmt.Sprintf("incompatible types: left %q right %q",
				left.valueType, right.valueType)))
	}
	st.push(&Expr{
		exprType:  et,
		valueType: boolType,
		children:  []*Expr{left, right},
	})
}

func pushEnum(vtype valueType, parse func(string) string) func(*Parser, *stack) {
	return func(_ *Parser, st *stack) {
		if st.top().valueType == vtype {
			return
		}
		switch e := st.pop(); e.exprType {
		case cnst:
			switch e.valueType {
			case stringType:
				st.push(&Expr{
					exprType:    cnst,
					valueType:   vtype,
					stringValue: parse(e.stringValue),
				})
			}
		default:
			switch e.valueType {
			case stringType:
				st.push(&Expr{
					exprType:  cast,
					valueType: vtype,
					children:  []*Expr{e},
				})
			default:
				panic(parseError(
					fmt.Sprintf("unsupported cast from %q to %q", e.valueType, vtype)))
			}
		}
	}
}

func (p *Parser) checkSearchLength(term string) {
	if p.MinSearchLength > 0 && len(term) < p.MinSearchLength {
		panic(parseError(
			fmt.Sprintf("search term too short (must be at least %d chars long)",
				p.MinSearchLength)))
	}
}

func (p *Parser) pushSearch(st *stack) {
	term := st.pop()
	term.checkValueType(stringType)
	p.checkSearchLength(term.stringValue)
	p.UsedSources.add(documentsTable | textTable)
	st.push(&Expr{
		exprType:    search,
		valueType:   boolType,
		stringValue: term.stringValue,
	})
}

func (p *Parser) pushMentioned(st *stack) {
	term := st.pop()
	term.checkValueType(stringType)
	p.checkSearchLength(term.stringValue)
	p.UsedSources.add(eventsLogTable)
	st.push(&Expr{
		exprType:    mentioned,
		valueType:   boolType,
		stringValue: term.stringValue,
	})
}

func (p *Parser) pushInvolved(st *stack) {
	term := st.pop()
	term.checkValueType(stringType)
	p.UsedSources.add(eventsLogTable)
	st.push(&Expr{
		exprType:    involved,
		valueType:   boolType,
		stringValue: term.stringValue,
	})
}

func (p *Parser) pushILike(st *stack) {
	needle := st.pop()
	haystack := st.pop()
	needle.checkValueType(stringType)
	haystack.checkValueType(stringType)
	p.UsedSources.add(documentsTable | textTable)
	st.push(&Expr{
		exprType:  ilike,
		valueType: boolType,
		children:  []*Expr{haystack, needle},
	})
}

func (p *Parser) pushILikePName(st *stack) {
	needle := st.pop()
	needle.checkValueType(stringType)
	p.UsedSources.add(documentsTable | textTable)
	st.push(&Expr{
		exprType:  ilikePName,
		valueType: boolType,
		children:  []*Expr{needle},
	})
}

func (p *Parser) pushILikePID(st *stack) {
	needle := st.pop()
	needle.checkValueType(stringType)
	p.UsedSources.add(documentsTable | textTable)
	st.push(&Expr{
		exprType:  ilikePID,
		valueType: boolType,
		children:  []*Expr{needle},
	})
}

func (*Parser) pushNow(st *stack) {
	st.push(&Expr{
		exprType:  now,
		valueType: timeType,
	})
}

func (p *Parser) pushMe(st *stack) {
	st.pushString(p.Me)
}

func (*Parser) pushDuration(st *stack) {
	if st.top().valueType == durationType {
		return
	}
	switch e := st.pop(); e.exprType {
	case cnst:
		switch e.valueType {
		case stringType:
			st.push(&Expr{
				exprType:      cnst,
				valueType:     durationType,
				durationValue: parseDuration(e.stringValue),
			})
		default:
			panic(parseError(
				fmt.Sprintf("value type %q is not supported as duration",
					e.valueType)))
		}
	default:
		panic(parseError(
			fmt.Sprintf("type %q is not supported as duration", e.exprType)))
	}
}

func (p *Parser) pushAs(st *stack) {
	alias := st.pop()
	srch := st.top()
	alias.checkValueType(stringType)
	srch.checkExprType(search) // TODO: Add csearch?
	validAlias(alias.stringValue)
	if p.aliases == nil {
		p.aliases = map[string]*Expr{}
	}
	if p.aliases[alias.stringValue] != nil {
		panic(parseError(fmt.Sprintf("duplicate alias %q", alias.stringValue)))
	}
	p.UsedSources.add(documentsTable | textTable)
	srch.alias = alias.stringValue
	srch.intValue = int64(len(p.aliases))
	p.aliases[alias.stringValue] = srch
}

func (p *Parser) parse(input string) (*Expr, error) {
	p.aliases = nil

	st := stack{}
	acts := action[p.Mode]

	split(input, func(field string, isString bool) {
		if !isString {
			if act := acts[field]; act != nil {
				act(p, &st)
				return
			}
		}
		st.pushString(field)
	})

	// If there are more than 2 open bool valued expressions on
	// the stack automatically and them together.
	st.andReduce()

	if len(st) != 1 {
		return nil, parseError(fmt.Sprintf(
			"invalid number of expression roots: expected 1 have %d", len(st)))
	}
	e := st.top()
	e.checkValueType(boolType)
	return e, nil
}

// Parse returns an expression.
func (p *Parser) Parse(input string) (expr *Expr, err error) {
	defer func() {
		if x := recover(); x != nil {
			if pe, ok := x.(parseError); ok {
				err = pe
			} else {
				panic(x)
			}
		}
	}()
	return p.parse(input)
}

func (st *stack) push(v *Expr) {
	*st = append(*st, v)
}

func (st *stack) pop() *Expr {
	if l := len(*st); l > 0 {
		x := (*st)[l-1]
		(*st)[l-1] = nil
		*st = (*st)[:l-1]
		return x
	}
	panic(parseError("stack empty"))
}

func (st stack) top() *Expr { return st.topN(0) }

func (st stack) topN(n int) *Expr {
	if l := len(st); l > n {
		return st[l-n-1]
	}
	panic(parseError("stack empty"))
}

func (st *stack) pushString(s string) {
	st.push(&Expr{
		exprType:    cnst,
		valueType:   stringType,
		stringValue: s,
	})
}

func (st *stack) andReduce() {
	for len(*st) > 1 {
		a, b := st.topN(0), st.topN(1)
		if a.valueType != boolType || b.valueType != boolType {
			return
		}
		st.pop()
		st.pop()
		st.push(a.And(b))
	}
}

func parseFloat(s string) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(parseError(fmt.Sprintf("%q is not a float: %v", s, err)))
	}
	return v
}

func parseInt(s string) int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(parseError(fmt.Sprintf("%q is not an int: %v", s, err)))
	}
	return v
}

func parseDuration(s string) time.Duration {
	duration, err := time.ParseDuration(s)
	if err != nil {
		panic(parseError(fmt.Sprintf("cannot parse %q as duration", s)))
	}
	return duration
}

func parseTime(s string) time.Time {
	for _, format := range []string{
		time.RFC3339,
		"2006-01-02",
		"2006-01-02T15:04:05-0700",
		"2006-01-02 15:04:05-0700",
	} {
		t, err := time.Parse(format, s)
		if err == nil {
			return t
		}
	}
	panic(parseError(fmt.Sprintf("cannot parse %q as time", s)))
}

var validWorkflows = []string{
	"new", "read", "assessing",
	"review", "archived", "delete",
}

func parseWorkflow(s string) string {
	if !slices.Contains(validWorkflows, s) {
		panic(parseError(fmt.Sprintf("%q is not a valid workflow", s)))
	}
	return s
}

var validEvents = []string{
	"import_document", "delete_document",
	"state_change",
	"add_sscv", "change_sscv", "delete_sscv",
	"add_comment", "change_comment", "delete_comment",
}

func parseEvents(s string) string {
	if !slices.Contains(validEvents, s) {
		panic(parseError(fmt.Sprintf("%q is not a valid event", s)))
	}
	return s
}

func parseStatus(s string) string {
	switch st := csaf.TrackingStatus(s); st {
	case csaf.CSAFTrackingStatusDraft, csaf.CSAFTrackingStatusFinal, csaf.CSAFTrackingStatusInterim:
		return s
	default:
		panic(parseError(fmt.Sprintf("%q is not a valid status", s)))
	}
}

var aliasRe = regexp.MustCompile(`[a-zA-Z_0-9]+`)

func validAlias(s string) {
	if !aliasRe.MatchString(s) {
		panic(parseError(fmt.Sprintf("invalid alias %q", s)))
	}
}

func split(input string, fn func(string, bool)) {
	var b strings.Builder
	state := 0
	for _, r := range input {
		switch state {
		case 0: // white space
			switch r {
			case '"':
				state = 1
			case '\\':
				state = 2
			default:
				if !unicode.IsSpace(r) {
					b.WriteRune(r)
					state = 3
				}
			}
		case 1: // quoted string
			switch r {
			case '\\':
				state = 5
			case '"':
				fn(b.String(), true)
				b.Reset()
				state = 0
			default:
				b.WriteRune(r)
			}
		case 2: // \ in white space
			b.WriteRune(r)
			state = 3
		case 3: // unquoted string
			if r == '\\' {
				state = 4
			} else if unicode.IsSpace(r) {
				fn(b.String(), false)
				b.Reset()
				state = 0
			} else {
				b.WriteRune(r)
			}
		case 4: // \ in unquoted string
			b.WriteRune(r)
			state = 3
		case 5: // \ in quoted string
			b.WriteRune(r)
			state = 1
		}
	}
	if state != 0 {
		fn(b.String(), state == 1 || state == 5)
	}
}
