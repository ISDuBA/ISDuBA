// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

package query

import (
	"slices"
	"strconv"
	"strings"
)

// CreateTextSearchWhereClause creates an SQL WHERE clause to
// search matching strings for a given expression tree.
func CreateTextSearchWhereClause(column string, e *Expr) (string, []any) {
	var (
		b            strings.Builder
		replacements []any
	)

	textIndex := func(text string) (idx int) {
		replacement := any(LikeEscape(text))
		if idx = slices.Index(replacements, replacement); idx > -1 {
			return idx
		}
		idx = len(replacements)
		replacements = append(replacements, replacement)
		return
	}

	var recurse func(*Expr)
	opTerm := func(curr *Expr, operator string) {
		b.WriteByte('(')
		for i, child := range curr.children {
			if i > 0 {
				b.WriteString(operator)
			}
			b.WriteByte('(')
			recurse(child)
			b.WriteByte(')')
		}
		b.WriteByte(')')
	}
	notTerm := func(curr *Expr) {
		b.WriteString("(NOT(")
		if len(curr.children) > 0 {
			recurse(curr.children[0])
		} else {
			// XXX: This should not happend.
			b.WriteString("FALSE")
		}
		b.WriteString("))")
	}
	searchTerm := func(curr *Expr) {
		b.WriteByte('(')
		b.WriteString(column)
		b.WriteString(" ILIKE $")
		b.WriteString(strconv.Itoa(textIndex(curr.stringValue)))
		b.WriteByte(')')
	}
	// Flatten the tree.
	recurse = func(curr *Expr) {
		switch curr.exprType {
		case and:
			opTerm(curr, "AND")
		case or:
			opTerm(curr, "OR")
		case not:
			notTerm(curr)
		case search:
			searchTerm(curr)
		default:
			// All other conditions are assumed to be fulfilled.
			b.WriteString("TRUE")
		}
	}
	recurse(e)
	return b.String(), replacements
}
