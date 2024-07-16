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
	"maps"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type parseError string

type exprType int

const (
	cnst exprType = iota
	cast
	and
	or
	not
	eq
	ne
	gt
	lt
	ge
	le
	access
	search
	mentioned
	involved
	ilike
	ilikePID
	now
	add
	sub
	mul
	div
)

type valueType int

const (
	intType valueType = iota
	floatType
	boolType
	stringType
	timeType
	workflowType
	durationType
	eventsType
)

// ParserMode represents the operation mode of the parser.
type ParserMode int

const (
	DocumentMode ParserMode = iota // DocumentMode operates on documents.
	AdvisoryMode                   // AdvisoryMode operates on advisories.
	EventMode                      // EventMode operates on events.
)

// Parser helps parsing database queries,
type Parser struct {
	// Mode indicates that only advisories should be considered.
	Mode ParserMode
	// MinSearchLength enforces a minimal lengths of search phrases.
	MinSearchLength int
	// Me is a replacement text for the "me" keyword.
	Me string

	// aliases defined by 'as'.
	aliases map[string]struct{}
}

// Expr encapsulates a parsed expression to be converted to an SQL WHERE clause.
type Expr struct {
	exprType  exprType
	valueType valueType

	stringValue   string
	intValue      int64
	floatValue    float64
	timeValue     time.Time
	boolValue     bool
	durationValue time.Duration
	alias         string

	children []*Expr
}

type documentColumn struct {
	name           string
	valueType      valueType
	modes          []ParserMode
	projectionOnly bool
}

type binaryCompat struct {
	left     valueType
	operator exprType
	right    valueType
}

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

// String implements [fmt.Stringer].
func (vt valueType) String() string {
	switch vt {
	case intType:
		return "integer"
	case floatType:
		return "float"
	case boolType:
		return "bool"
	case stringType:
		return "string"
	case timeType:
		return "timestamp"
	case workflowType:
		return "workflow"
	case durationType:
		return "duration"
	case eventsType:
		return "events"
	default:
		return fmt.Sprintf("unknown value type %d", vt)
	}
}

// String implements [fmt.Stringer].
func (pm ParserMode) String() string {
	switch pm {
	case DocumentMode:
		return "document"
	case AdvisoryMode:
		return "advisory"
	case EventMode:
		return "event"
	default:
		return fmt.Sprintf("Unknown parser mode: %d", pm)
	}
}

// FieldEqInt is a shortcut mainly for building expressions
// accessing an integer column like 'id's.
func FieldEqInt(field string, value int64) *Expr {
	return &Expr{
		valueType: boolType,
		exprType:  eq,
		children: []*Expr{
			{valueType: intType, exprType: access, stringValue: field},
			{valueType: intType, exprType: cnst, intValue: value},
		},
	}
}

// FieldEqString is a shortcut mainly for building expressions
// accessing a string column.
func FieldEqString(field, value string) *Expr {
	return &Expr{
		valueType: boolType,
		exprType:  eq,
		children: []*Expr{
			{valueType: stringType, exprType: access, stringValue: field},
			{valueType: stringType, exprType: cnst, stringValue: value},
		},
	}
}

// BoolField returns an access term that returns a bool value.
func BoolField(field string) *Expr {
	return &Expr{
		valueType:   boolType,
		exprType:    access,
		stringValue: field,
	}
}

// String implements [fmt.Stringer].
func (et exprType) String() string {
	switch et {
	case cnst:
		return "constant"
	case cast:
		return "cast"
	case and:
		return "and"
	case or:
		return "or"
	case not:
		return "not"
	case eq:
		return "eq"
	case ne:
		return "ne"
	case gt:
		return "gt"
	case lt:
		return "lt"
	case ge:
		return "ge"
	case le:
		return "le"
	case access:
		return "access"
	case search:
		return "search"
	case mentioned:
		return "mentioned"
	case involved:
		return "involved"
	case ilike:
		return "ilike"
	case ilikePID:
		return "ilikepid"
	case now:
		return "now"
	case add:
		return "+"
	case sub:
		return "-"
	case mul:
		return "*"
	case div:
		return "/"
	default:
		return fmt.Sprintf("unknown expression type %d", et)
	}
}

func (pe parseError) Error() string {
	return string(pe)
}

var (
	docAdvEvtModes = []ParserMode{DocumentMode, AdvisoryMode, EventMode}
	advModes       = []ParserMode{AdvisoryMode}
	evtsModes      = []ParserMode{EventMode}
)

// documentColumns are the documentColumns which can be accessed.
var documentColumns = []documentColumn{
	{"id", intType, docAdvEvtModes, false},
	{"latest", boolType, docAdvEvtModes, false},
	{"tracking_id", stringType, docAdvEvtModes, false},
	{"version", stringType, docAdvEvtModes, false},
	{"publisher", stringType, docAdvEvtModes, false},
	{"current_release_date", timeType, docAdvEvtModes, false},
	{"initial_release_date", timeType, docAdvEvtModes, false},
	{"rev_history_length", intType, docAdvEvtModes, false},
	{"title", stringType, docAdvEvtModes, false},
	{"tlp", stringType, docAdvEvtModes, false},
	{"ssvc", stringType, docAdvEvtModes, false},
	{"cvss_v2_score", floatType, docAdvEvtModes, false},
	{"cvss_v3_score", floatType, docAdvEvtModes, false},
	{"critical", floatType, docAdvEvtModes, false},
	{"four_cves", stringType, docAdvEvtModes, true},
	{"comments", intType, docAdvEvtModes, false},
	// Advisories only
	{"state", workflowType, advModes, false},
	{"recent", timeType, advModes, false},
	{"versions", intType, advModes, false},
	// Events only
	{"event", eventsType, evtsModes, false},
	{"event_state", workflowType, evtsModes, false},
	{"time", timeType, evtsModes, false},
	{"actor", stringType, evtsModes, false},
	{"comments_id", intType, evtsModes, false},
}

var (
	// baseActions are the action available in every parser.
	baseActions = map[string]func(*Parser, *stack){
		"true":      (*Parser).pushTrue,
		"false":     (*Parser).pushFalse,
		"not":       (*Parser).pushNot,
		"and":       curry3((*Parser).pushBinary, and),
		"or":        curry3((*Parser).pushBinary, or),
		"float":     (*Parser).pushFloat,
		"integer":   (*Parser).pushInteger,
		"timestamp": (*Parser).pushTimestamp,
		"workflow":  (*Parser).pushWorkflow,
		"events":    (*Parser).pushEvents,
		"=":         curry3((*Parser).pushCmp, eq),
		"!=":        curry3((*Parser).pushCmp, ne),
		"<":         curry3((*Parser).pushCmp, lt),
		"<=":        curry3((*Parser).pushCmp, le),
		">":         curry3((*Parser).pushCmp, gt),
		">=":        curry3((*Parser).pushCmp, ge),
		"ilike":     (*Parser).pushILike,
		"ilikepid":  (*Parser).pushILikePID,
		"now":       (*Parser).pushNow,
		"duration":  (*Parser).pushDuration,
		"+":         curry3((*Parser).pushBinary, add),
		"-":         curry3((*Parser).pushBinary, sub),
		"/":         curry3((*Parser).pushBinary, div),
		"*":         curry3((*Parser).pushBinary, mul),
		"me":        (*Parser).pushMe,
		"mentioned": (*Parser).pushMentioned,
		"involved":  (*Parser).pushInvolved,
	}
	// advancedActions are action only available is documents and advisories.
	advancedActions = map[string]func(*Parser, *stack){
		"search": (*Parser).pushSearch,
		"as":     (*Parser).pushAs,
	}
	// actions is for fast looking up actions along the parser mode.
	actions = map[ParserMode]map[string]func(*Parser, *stack){
		DocumentMode: buildActions(DocumentMode),
		AdvisoryMode: buildActions(AdvisoryMode),
		EventMode:    buildActions(EventMode),
	}
)

// ExistsDocumentColumn returns true if a column in document exists.
func ExistsDocumentColumn(name string, mode ParserMode) bool {
	return findDocumentColumn(name, mode) != nil
}

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

// And concats two expressions and-wise.
func (e *Expr) And(o *Expr) *Expr {
	if e.valueType != boolType || o.valueType != boolType {
		return False()
	}
	if e.exprType == cnst {
		if !e.boolValue {
			return False()
		}
		return o
	}
	if o.exprType == cnst {
		if !o.boolValue {
			return False()
		}
		return e
	}
	return &Expr{
		exprType:  and,
		valueType: boolType,
		children:  []*Expr{e, o},
	}
}

// Or concats two expressions or-wise.
func (e *Expr) Or(o *Expr) *Expr {
	if e.valueType != boolType || o.valueType != boolType {
		return False()
	}
	if e.exprType == cnst {
		if e.boolValue {
			return True()
		}
		return o
	}
	if o.exprType == cnst {
		if o.boolValue {
			return True()
		}
		return e
	}
	return &Expr{
		exprType:  or,
		valueType: boolType,
		children:  []*Expr{e, o},
	}
}

// Not negates an expresssion.
func (e *Expr) Not() *Expr {
	if e.valueType != boolType {
		return False()
	}
	if e.exprType == cnst {
		if e.boolValue {
			return False()
		}
		return True()
	}
	return &Expr{
		exprType:  not,
		valueType: boolType,
		children:  []*Expr{e},
	}
}

type stack []*Expr

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

// False returns a false expression.
func False() *Expr {
	return &Expr{
		exprType:  cnst,
		valueType: boolType,
		boolValue: false,
	}
}

// True returns a true expression.
func True() *Expr {
	return &Expr{
		exprType:  cnst,
		valueType: boolType,
		boolValue: true,
	}
}

func (*Parser) pushTrue(st *stack)  { st.push(True()) }
func (*Parser) pushFalse(st *stack) { st.push(False()) }

func (st *stack) pushString(s string) {
	st.push(&Expr{
		exprType:    cnst,
		valueType:   stringType,
		stringValue: s,
	})
}

func (e *Expr) checkValueType(vt valueType) {
	if e.valueType != vt {
		panic(parseError(
			fmt.Sprintf("value type mismatch: %q %q", e.valueType, vt)))
	}
}

func (e *Expr) checkExprType(eTypes ...exprType) {
	for _, et := range eTypes {
		if e.exprType == et {
			return
		}
	}
	panic(parseError(
		fmt.Sprintf("expression type mismatch: %q %v", e.exprType, eTypes)))
}

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

func (*Parser) pushAccess(st *stack, column *documentColumn) {
	st.push(&Expr{
		exprType:    access,
		valueType:   column.valueType,
		stringValue: column.name,
	})
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

var validWorkflows = []string{
	"new", "read", "assessing",
	"review", "archive", "delete",
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

func (*Parser) pushWorkflow(st *stack) {
	if st.top().valueType == workflowType {
		return
	}
	switch e := st.pop(); e.exprType {
	case cnst:
		switch e.valueType {
		case stringType:
			st.push(&Expr{
				exprType:    cnst,
				valueType:   workflowType,
				stringValue: parseWorkflow(e.stringValue),
			})
		}
	default:
		switch e.valueType {
		case stringType:
			st.push(&Expr{
				exprType:  cast,
				valueType: workflowType,
				children:  []*Expr{e},
			})
		default:
			panic(parseError("unsupported cast"))
		}
	}
}

func (*Parser) pushEvents(st *stack) {
	if st.top().valueType == eventsType {
		return
	}
	switch e := st.pop(); e.exprType {
	case cnst:
		switch e.valueType {
		case stringType:
			st.push(&Expr{
				exprType:    cnst,
				valueType:   eventsType,
				stringValue: parseEvents(e.stringValue),
			})
		}
	default:
		switch e.valueType {
		case stringType:
			st.push(&Expr{
				exprType:  cast,
				valueType: eventsType,
				children:  []*Expr{e},
			})
		default:
			panic(parseError("unsupported cast"))
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
	st.push(&Expr{
		exprType:    mentioned,
		valueType:   boolType,
		stringValue: term.stringValue,
	})
}

func (*Parser) pushInvolved(st *stack) {
	term := st.pop()
	term.checkValueType(stringType)
	st.push(&Expr{
		exprType:    involved,
		valueType:   boolType,
		stringValue: term.stringValue,
	})
}

func (*Parser) pushILike(st *stack) {
	needle := st.pop()
	haystack := st.pop()
	needle.checkValueType(stringType)
	haystack.checkValueType(stringType)
	st.push(&Expr{
		exprType:  ilike,
		valueType: boolType,
		children:  []*Expr{haystack, needle},
	})
}

func (*Parser) pushILikePID(st *stack) {
	needle := st.pop()
	needle.checkValueType(stringType)
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

var aliasRe = regexp.MustCompile(`[a-zA-Z_0-9]+`)

func validAlias(s string) {
	if !aliasRe.MatchString(s) {
		panic(parseError(fmt.Sprintf("invalid alias %q", s)))
	}
}

func (p *Parser) pushAs(st *stack) {
	alias := st.pop()
	srch := st.top()
	alias.checkValueType(stringType)
	srch.checkExprType(search) // TODO: Add csearch?
	validAlias(alias.stringValue)
	if p.aliases == nil {
		p.aliases = map[string]struct{}{}
	}
	if _, already := p.aliases[alias.stringValue]; already {
		panic(parseError(fmt.Sprintf("duplicate alias %q", alias.stringValue)))
	}
	p.aliases[alias.stringValue] = struct{}{}
	srch.alias = alias.stringValue
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

func buildActions(mode ParserMode) map[string]func(*Parser, *stack) {
	actions := maps.Clone(baseActions)
	for i := range documentColumns {
		col := &documentColumns[i]
		if !col.projectionOnly && slices.Contains(col.modes, mode) {
			actions["$"+col.name] = func(p *Parser, st *stack) { p.pushAccess(st, col) }
		}
	}
	// Fill in extra actions
	switch mode {
	case DocumentMode, AdvisoryMode:
		maps.Copy(actions, advancedActions)
	}
	return actions
}

func curry3[A, B, C any](fn func(A, B, C), c C) func(A, B) {
	return func(a A, b B) { fn(a, b, c) }
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

func (p *Parser) parse(input string) (*Expr, error) {

	p.aliases = nil

	st := stack{}
	acts := actions[p.Mode]

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
