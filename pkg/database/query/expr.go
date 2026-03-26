// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

package query

import (
	"fmt"
	"iter"
	"slices"
	"time"

	"github.com/ISDuBA/ISDuBA/pkg/itertools"
)

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
	ilikePName
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
	statusType
)

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

type parseError string

// Error implements [error].
func (pe parseError) Error() string {
	return string(pe)
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
	case statusType:
		return "status"
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

func (e *Expr) getStringValue() string { return e.stringValue }
func (e *Expr) getAlias() string       { return e.alias }

// Accesses returns a sequence over all database accessed columns in the expression tree.
func (e *Expr) Accesses() iter.Seq[string] {
	return itertools.Unique(itertools.Apply(itertools.Filter(
		e.all(),
		func(e *Expr) bool {
			return e.exprType == access && e.stringValue != ""
		}),
		(*Expr).getStringValue,
	))
}

// Aliases returns a sequence over all search aliases in the expression tree.
func (e *Expr) Aliases() iter.Seq[string] {
	return itertools.Unique(itertools.Apply(itertools.Filter(
		e.all(),
		func(e *Expr) bool {
			return e.exprType == search && e.alias != ""
		}),
		(*Expr).getAlias,
	))
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

// all returns a sequence over all expressions in the tree.
func (e *Expr) all() iter.Seq[*Expr] {
	return func(yield func(*Expr) bool) {
		var recursive func(*Expr) bool
		recursive = func(curr *Expr) bool {
			if curr == nil {
				return true
			}
			if !yield(curr) {
				return false
			}
			for _, child := range curr.children {
				if !recursive(child) {
					return false
				}
			}
			return true
		}
		recursive(e)
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

func (e *Expr) checkValueType(vt valueType) {
	if e.valueType != vt {
		panic(parseError(
			fmt.Sprintf("value type mismatch: %q %q", e.valueType, vt)))
	}
}

func (e *Expr) checkExprType(eTypes ...exprType) {
	if slices.Contains(eTypes, e.exprType) {
		return
	}
	panic(parseError(
		fmt.Sprintf("expression type mismatch: %q %v", e.exprType, eTypes)))
}
