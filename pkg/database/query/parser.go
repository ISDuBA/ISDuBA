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
	csearch
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
)

// Parser helps parsing database queries,
type Parser struct {
	// Advisory indicates that only advisories should be considered.
	Advisory bool
	// Languages are the languages supported by full-text search.
	Languages []string
	// MinSearchLength enforces a minimal lengths of search phrases.
	MinSearchLength int
	// Me is a replacement text for the "me" keyword.
	Me string
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
	langValue     string
	durationValue time.Duration
	alias         string

	children []*Expr
}

type documentColumn struct {
	name           string
	valueType      valueType
	advisoryOnly   bool
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
	{durationType, mul, floatType}:    durationType,
	{durationType, div, floatType}:    durationType,
}

// String implements [fmt.Stringer].
func (vt valueType) String() string {
	switch vt {
	case intType:
		return "int"
	case floatType:
		return "float"
	case boolType:
		return "bool"
	case stringType:
		return "string"
	case timeType:
		return "time"
	case workflowType:
		return "workflow"
	case durationType:
		return "duration"
	default:
		return fmt.Sprintf("unknown value type %d", vt)
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
	case csearch:
		return "csearch"
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

var columns = []documentColumn{
	{"id", intType, false, false},
	{"latest", boolType, false, false},
	{"tracking_id", stringType, false, false},
	{"version", stringType, false, false},
	{"publisher", stringType, false, false},
	{"current_release_date", timeType, false, false},
	{"initial_release_date", timeType, false, false},
	{"rev_history_length", intType, false, false},
	{"title", stringType, false, false},
	{"tlp", stringType, false, false},
	{"ssvc", stringType, false, false},
	{"cvss_v2_score", floatType, false, false},
	{"cvss_v3_score", floatType, false, false},
	{"critical", floatType, false, false},
	{"four_cves", stringType, false, true},
	{"comments", intType, false, false},
	// Advisories only
	{"state", workflowType, true, false},
	{"recent", timeType, true, false},
	{"versions", intType, true, false},
}

// supportedLangs are the default languages.
// Can be overwritten in Parser.
var supportedLangs = []string{
	"english",
	"german",
}

// ExistsDocumentColumn returns true if a column in document exists.
func ExistsDocumentColumn(name string, advisory bool) bool {
	return findDocumentColumn(name, advisory) != nil
}

func findDocumentColumn(name string, advisory bool) *documentColumn {
	for i := range columns {
		if col := &columns[i]; col.name == name {
			if col.advisoryOnly && !advisory {
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

func (st stack) top() *Expr {
	if l := len(st); l > 0 {
		return st[l-1]
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

func (st *stack) pushTrue()  { st.push(True()) }
func (st *stack) pushFalse() { st.push(False()) }

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

func (st *stack) not() {
	e := st.pop()
	e.checkValueType(boolType)
	st.push(&Expr{
		exprType:  not,
		valueType: boolType,
		children:  []*Expr{e},
	})
}

func (st *stack) binary(et exprType) {
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

func (st *stack) access(field string, advisory bool) {
	col := findDocumentColumn(field, advisory)
	if col == nil {
		panic(parseError(fmt.Sprintf("unknown column %q", field)))
	}
	if col.projectionOnly {
		panic(parseError(fmt.Sprintf("column %q is for projection only", field)))
	}
	st.push(&Expr{
		exprType:    access,
		valueType:   col.valueType,
		stringValue: field,
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

func (st *stack) float() {
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

func (st *stack) integer() {
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

func (st *stack) time() {
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

func (st *stack) cmp(et exprType) {
	right := st.pop()
	left := st.pop()
	if right.valueType != left.valueType {
		panic(parseError("incompatible types"))
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

func (st *stack) workflow() {
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

func (p *Parser) checkLanguage(lang string) {
	var langs []string
	if p.Languages != nil {
		langs = p.Languages
	} else {
		langs = supportedLangs
	}
	if !slices.Contains(langs, lang) {
		panic(parseError(
			fmt.Sprintf("unsupported search language %q", lang)))
	}
}

func (p *Parser) checkSearchLength(term string) {
	if p.MinSearchLength > 0 && len(term) < p.MinSearchLength {
		panic(parseError(
			fmt.Sprintf("search term too short (must be at least %d chars long)", p.MinSearchLength)))
	}
}

func (st *stack) search(p *Parser) {
	lang := st.pop()
	term := st.pop()
	lang.checkValueType(stringType)
	term.checkValueType(stringType)
	p.checkLanguage(lang.stringValue)
	p.checkSearchLength(term.stringValue)
	st.push(&Expr{
		exprType:    search,
		valueType:   boolType,
		langValue:   lang.stringValue,
		stringValue: term.stringValue,
	})
}

func (st *stack) csearch(p *Parser) {
	lang := st.pop()
	term := st.pop()
	lang.checkValueType(stringType)
	term.checkValueType(stringType)
	p.checkLanguage(lang.stringValue)
	p.checkSearchLength(term.stringValue)
	st.push(&Expr{
		exprType:    csearch,
		valueType:   boolType,
		langValue:   lang.stringValue,
		stringValue: term.stringValue,
	})
}

func (st *stack) mentioned(p *Parser) {
	term := st.pop()
	term.checkValueType(stringType)
	p.checkSearchLength(term.stringValue)
	st.push(&Expr{
		exprType:    mentioned,
		valueType:   boolType,
		stringValue: term.stringValue,
	})
}

func (st *stack) involved() {
	term := st.pop()
	term.checkValueType(stringType)
	st.push(&Expr{
		exprType:    involved,
		valueType:   boolType,
		stringValue: term.stringValue,
	})
}

func (st *stack) ilike() {
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

func (st *stack) ilikePID() {
	needle := st.pop()
	needle.checkValueType(stringType)
	st.push(&Expr{
		exprType:  ilikePID,
		valueType: boolType,
		children:  []*Expr{needle},
	})
}

func (st *stack) now() {
	st.push(&Expr{
		exprType:  now,
		valueType: timeType,
	})
}

func (st *stack) duration() {
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

var aliasRe = regexp.MustCompile(`[a-zA-Z][a-zA-Z_0-9]*`)

func validAlias(s string) {
	if !aliasRe.MatchString(s) {
		panic(parseError(fmt.Sprintf("invalid alias %q", s)))
	}
}

func (st *stack) as(aliases map[string]struct{}) {
	alias := st.pop()
	srch := st.top()
	alias.checkValueType(stringType)
	srch.checkExprType(search) // TODO: Add csearch?
	validAlias(alias.stringValue)
	if _, already := aliases[alias.stringValue]; already {
		panic(parseError(fmt.Sprintf("duplicate alias %q", alias.stringValue)))
	}
	aliases[alias.stringValue] = struct{}{}
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

func (p *Parser) parse(input string) (*Expr, error) {
	st := stack{}
	aliases := map[string]struct{}{}

	split(input, func(field string, isString bool) {

		if isString {
			st.pushString(field)
			return
		}

		switch field {
		case "true":
			st.pushTrue()
		case "false":
			st.pushFalse()
		case "not":
			st.not()
		case "and":
			st.binary(and)
		case "or":
			st.binary(or)
		case "float":
			st.float()
		case "int":
			st.integer()
		case "time":
			st.time()
		case "workflow":
			st.workflow()
		case "=":
			st.cmp(eq)
		case "!=":
			st.cmp(ne)
		case "<":
			st.cmp(lt)
		case "<=":
			st.cmp(le)
		case ">":
			st.cmp(gt)
		case ">=":
			st.cmp(ge)
		case "search":
			st.search(p)
		case "csearch":
			st.csearch(p)
		case "mentioned":
			st.mentioned(p)
		case "involved":
			st.involved()
		case "as":
			st.as(aliases)
		case "ilike":
			st.ilike()
		case "ilikepid":
			st.ilikePID()
		case "now":
			st.now()
		case "duration":
			st.duration()
		case "+":
			st.binary(add)
		case "-":
			st.binary(sub)
		case "/":
			st.binary(div)
		case "*":
			st.binary(mul)
		case "me":
			st.pushString(p.Me)
		default:
			if strings.HasPrefix(field, "$") {
				st.access(field[1:], p.Advisory)
			} else {
				st.pushString(field)
			}
		}
	})

	if len(st) != 1 {
		return nil, parseError(fmt.Sprintf(
			"invalid number of expression roots: expected 1 have %d", len(st)))
	}
	e := st[len(st)-1]
	if e.valueType != boolType {
		return nil, parseError("not a boolean expression")
	}
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
