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
	"slices"
)

// AdvancedParser helps parsing database queries,
type AdvancedParser struct {
	// Mode indicates that only advisories should be considered.
	Mode ParserMode
	// MinSearchLength enforces a minimal lengths of search phrases.
	MinSearchLength int
	// Me is a replacement text for the "me" keyword.
	Me string

	// UsedSources are the sources found during parsing.
	UsedSources columnSources

	// aliases defined by 'as'.
	aliases map[string]struct{}
}

type (
	columnSource  string
	columnSources []columnSource
)

type documentColumn struct {
	name           string
	valueType      valueType
	modes          []ParserMode
	projectionOnly bool
	sources        columnSources
}

func (css columnSources) contains(cs columnSource) bool {
	return slices.Contains(css, cs)
}

func (css *columnSources) add(others ...columnSource) {
	for _, other := range others {
		if !slices.Contains(*css, other) {
			*css = append(*css, other)
		}
	}
}

// documentColumns are the documentColumns which can be accessed.
var documentColumns = []documentColumn{
	{"id", intType, docAdvEvtModes, false, columnSources{"documents"}},
	{"latest", boolType, docAdvEvtModes, false, columnSources{"documents"}},
	{"tracking_id", stringType, docAdvEvtModes, false, columnSources{"advisories"}},
	{"version", stringType, docAdvEvtModes, false, columnSources{"documents"}},
	{"publisher", stringType, docAdvEvtModes, false, columnSources{"advisories"}},
	{"current_release_date", timeType, docAdvEvtModes, false, columnSources{"documents"}},
	{"initial_release_date", timeType, docAdvEvtModes, false, columnSources{"documents"}},
	{"rev_history_length", intType, docAdvEvtModes, false, columnSources{"documents"}},
	{"title", stringType, docAdvEvtModes, false, columnSources{"documents"}},
	{"tlp", stringType, docAdvEvtModes, false, columnSources{"documents"}},
	{"ssvc", stringType, docAdvEvtModes, false, columnSources{"ssvc_history"}},
	{"cvss_v2_score", floatType, docAdvEvtModes, false, columnSources{"documents"}},
	{"cvss_v3_score", floatType, docAdvEvtModes, false, columnSources{"documents"}},
	{"critical", floatType, docAdvEvtModes, false, columnSources{"documents"}},
	{"four_cves", stringType, docAdvEvtModes, true, columnSources{"documents"}},
	{"comments", intType, docAdvEvtModes, false, columnSources{"advisories"}},
	{"tracking_status", statusType, docAdvEvtModes, false, columnSources{"documents"}},
	// Advisories only
	{"state", workflowType, advModes, false, columnSources{"advisories"}},
	{"recent", timeType, advModes, false, columnSources{"advisories"}},
	// ToDo: Column "versions" does not exist, but table versions does?
	{"versions", intType, advModes, false, columnSources{"documents", "advisories"}},
	// Events only
	{"event", eventsType, evtsModes, false, columnSources{"events_log"}},
	{"event_state", workflowType, evtsModes, false, columnSources{"events_log"}},
	{"time", timeType, evtsModes, false, columnSources{"events_log"}},
	{"actor", stringType, evtsModes, false, columnSources{"events_log"}},
	{"comments_id", intType, evtsModes, false, columnSources{"comments"}},
	{"message", stringType, evtsModes, false, columnSources{"comments"}},
}

var (
	// advancedBaseAction are the action available in every parser.
	advancedBaseAction = map[string]func(*AdvancedParser, *stack){
		"true":       (*AdvancedParser).pushTrue,
		"false":      (*AdvancedParser).pushFalse,
		"not":        (*AdvancedParser).pushNot,
		"and":        curry3((*AdvancedParser).pushBinary, and),
		"or":         curry3((*AdvancedParser).pushBinary, or),
		"float":      (*AdvancedParser).pushFloat,
		"integer":    (*AdvancedParser).pushInteger,
		"timestamp":  (*AdvancedParser).pushTimestamp,
		"workflow":   advancedPushEnum(workflowType, parseWorkflow),
		"events":     advancedPushEnum(eventsType, parseEvents),
		"status":     advancedPushEnum(statusType, parseStatus),
		"=":          curry3((*AdvancedParser).pushCmp, eq),
		"!=":         curry3((*AdvancedParser).pushCmp, ne),
		"<":          curry3((*AdvancedParser).pushCmp, lt),
		"<=":         curry3((*AdvancedParser).pushCmp, le),
		">":          curry3((*AdvancedParser).pushCmp, gt),
		">=":         curry3((*AdvancedParser).pushCmp, ge),
		"ilike":      (*AdvancedParser).pushILike,
		"ilikepname": (*AdvancedParser).pushILikePName,
		"ilikepid":   (*AdvancedParser).pushILikePID,
		"now":        (*AdvancedParser).pushNow,
		"duration":   (*AdvancedParser).pushDuration,
		"+":          curry3((*AdvancedParser).pushBinary, add),
		"-":          curry3((*AdvancedParser).pushBinary, sub),
		"/":          curry3((*AdvancedParser).pushBinary, div),
		"*":          curry3((*AdvancedParser).pushBinary, mul),
		"me":         (*AdvancedParser).pushMe,
		"mentioned":  (*AdvancedParser).pushMentioned,
		"involved":   (*AdvancedParser).pushInvolved,
		"search":     (*AdvancedParser).pushSearch,
		"as":         (*AdvancedParser).pushAs,
	}
	// advancedAction is for fast looking up actions along the parser mode.
	advancedAction = map[ParserMode]map[string]func(*AdvancedParser, *stack){
		DocumentMode: advancedBuildActions(DocumentMode),
		AdvisoryMode: advancedBuildActions(AdvisoryMode),
		EventMode:    advancedBuildActions(EventMode),
	}
)

func advancedBuildActions(mode ParserMode) map[string]func(*AdvancedParser, *stack) {
	actions := maps.Clone(advancedBaseAction)
	for i := range documentColumns {
		col := &documentColumns[i]
		if !col.projectionOnly && slices.Contains(col.modes, mode) {
			actions["$"+col.name] = func(p *AdvancedParser, st *stack) {
				p.pushAccess(st, col)
			}
		}
	}
	return actions
}

func (*AdvancedParser) pushTrue(st *stack)  { st.push(True()) }
func (*AdvancedParser) pushFalse(st *stack) { st.push(False()) }

func (*AdvancedParser) pushNot(st *stack) {
	e := st.pop()
	e.checkValueType(boolType)
	st.push(&Expr{
		exprType:  not,
		valueType: boolType,
		children:  []*Expr{e},
	})
}

func (*AdvancedParser) pushBinary(st *stack, et exprType) {
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

func (p *AdvancedParser) pushAccess(st *stack, column *documentColumn) {
	p.UsedSources.add(column.sources...)
	st.push(&Expr{
		exprType:    access,
		valueType:   column.valueType,
		stringValue: column.name,
	})
}

func (*AdvancedParser) pushFloat(st *stack) {
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

func (*AdvancedParser) pushInteger(st *stack) {
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

func (*AdvancedParser) pushTimestamp(st *stack) {
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

func (*AdvancedParser) pushCmp(st *stack, et exprType) {
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

func advancedPushEnum(vtype valueType, parse func(string) string) func(*AdvancedParser, *stack) {
	return func(_ *AdvancedParser, st *stack) {
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

func (p *AdvancedParser) checkSearchLength(term string) {
	if p.MinSearchLength > 0 && len(term) < p.MinSearchLength {
		panic(parseError(
			fmt.Sprintf("search term too short (must be at least %d chars long)",
				p.MinSearchLength)))
	}
}

func (p *AdvancedParser) pushSearch(st *stack) {
	term := st.pop()
	term.checkValueType(stringType)
	p.checkSearchLength(term.stringValue)
	p.UsedSources.add("text_tables")
	st.push(&Expr{
		exprType:    search,
		valueType:   boolType,
		stringValue: term.stringValue,
	})
}

func (p *AdvancedParser) pushMentioned(st *stack) {
	term := st.pop()
	term.checkValueType(stringType)
	p.checkSearchLength(term.stringValue)
	p.UsedSources.add("events_log")
	st.push(&Expr{
		exprType:    mentioned,
		valueType:   boolType,
		stringValue: term.stringValue,
	})
}

func (p *AdvancedParser) pushInvolved(st *stack) {
	term := st.pop()
	term.checkValueType(stringType)
	p.UsedSources.add("events_log")
	st.push(&Expr{
		exprType:    involved,
		valueType:   boolType,
		stringValue: term.stringValue,
	})
}

func (p *AdvancedParser) pushILike(st *stack) {
	needle := st.pop()
	haystack := st.pop()
	needle.checkValueType(stringType)
	haystack.checkValueType(stringType)
	p.UsedSources.add("text_tables")
	st.push(&Expr{
		exprType:  ilike,
		valueType: boolType,
		children:  []*Expr{haystack, needle},
	})
}

func (p *AdvancedParser) pushILikePName(st *stack) {
	needle := st.pop()
	needle.checkValueType(stringType)
	p.UsedSources.add("documents", "text_tables")
	st.push(&Expr{
		exprType:  ilikePName,
		valueType: boolType,
		children:  []*Expr{needle},
	})
}

func (p *AdvancedParser) pushILikePID(st *stack) {
	needle := st.pop()
	needle.checkValueType(stringType)
	p.UsedSources.add("documents", "text_tables")
	st.push(&Expr{
		exprType:  ilikePID,
		valueType: boolType,
		children:  []*Expr{needle},
	})
}

func (*AdvancedParser) pushNow(st *stack) {
	st.push(&Expr{
		exprType:  now,
		valueType: timeType,
	})
}

func (p *AdvancedParser) pushMe(st *stack) {
	st.pushString(p.Me)
}

func (*AdvancedParser) pushDuration(st *stack) {
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

func (p *AdvancedParser) pushAs(st *stack) {
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
	p.UsedSources.add("text_tables")
	p.aliases[alias.stringValue] = struct{}{}
	srch.alias = alias.stringValue
}

func (p *AdvancedParser) parse(input string) (*Expr, error) {
	p.aliases = nil

	st := stack{}
	acts := advancedAction[p.Mode]

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
func (p *AdvancedParser) Parse(input string) (expr *Expr, err error) {
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
