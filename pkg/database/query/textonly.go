// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

package query

import "strings"

// CreateTextSearchWhereClause creates an SQL WHERE clause to
// search matching strings for a given expression tree.
func CreateTextSearchWhereClause(e *Expr) (string, []any) {
	var (
		b            strings.Builder
		replacements []any
	)

	var recurse func(*Expr)

	op := func(curr *Expr, operator string) {
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
	recurse = func(curr *Expr) {
		switch curr.exprType {
		case and:
			op(curr, "AND")
		case or:
			op(curr, "OR")
		case not:
			b.WriteString("(NOT(")
			if len(curr.children) > 0 {
				recurse(curr.children[0])
			} else {
				b.WriteString("FALSE")
			}
			b.WriteString("))")
		case ilike:
			// TODO: Implement me!
		default:
			// All other conditions are assumed to be fulfilled.
			b.WriteString("TRUE")
		}
	}
	recurse(e)
	return b.String(), replacements
}
