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
	"unicode"
)

// ILikeExpr is an compiled ILIKE expression.
type ILikeExpr []token

// CompileILike compiles an ILIKE search pattern.
func CompileILike(needle string) ILikeExpr {
	var (
		tokens ILikeExpr
		buf    []rune
	)
	flushLiteral := func() {
		if len(buf) > 0 {
			tokens = append(tokens, token{
				kind: litToken,
				lit:  slices.Clone(buf),
			})
			buf = buf[:0]
		}
	}
	escape := false
	for _, r := range needle {
		if escape {
			escape = false
			buf = append(buf, r)
			continue
		}
		switch r {
		case '\\':
			escape = true
		case '%':
			flushLiteral()
			if len(tokens) == 0 ||
				tokens[len(tokens)-1].kind == anyManyToken {
				tokens = append(tokens, token{kind: anyManyToken})
			}
		case '_':
			flushLiteral()
			tokens = append(tokens, token{kind: anyOneToken})
		default:
			buf = append(buf, r)
		}
	}
	if escape {
		buf = append(buf, '\\')
	}
	flushLiteral()
	return tokens
}

// Search searches ilike patterns in a haystack.
// Returns a list of matching positions.
func (expr ILikeExpr) Search(haystack string) [][2]int {
	if len(haystack) == 0 || len(expr) == 0 {
		//fmt.Println("nothing to match")
		return nil
	}
	var (
		runes     = []rune(haystack)
		positions [][2]int
	)
	for start := range runes {
		//fmt.Println(start)
		if end := expr.matchMinEnd(runes, start); end > start {
			positions = append(positions, [2]int{
				start,
				end - start,
			})
		}
	}
	return positions
}

type tokenKind int

const (
	litToken     = iota
	anyOneToken  // _
	anyManyToken // %
)

type token struct {
	kind tokenKind
	lit  []rune
}

func (expr ILikeExpr) matchMinEnd(hr []rune, start int) int {

	type state struct{ i, j int } // i in haystack runes, j in tokens

	seen := make(map[state]bool, 128)
	bestEnd := -1

	stack := make([]state, 1, 64)
	stack[0] = state{start, 0}

nextStack:
	for len(stack) > 0 {
		// pop
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if seen[cur] {
			continue
		}
		seen[cur] = true

		i, j := cur.i, cur.j

		if j == len(expr) {
			if bestEnd == -1 || i < bestEnd {
				bestEnd = i
			}
			continue
		}

		if bestEnd != -1 && i >= bestEnd {
			continue
		}

		t := expr[j]
		switch t.kind {
		case litToken:
			if i+len(t.lit) > len(hr) {
				continue
			}
			for k := 0; k < len(t.lit); k++ {
				if !runeEqualFold(hr[i+k], t.lit[k]) {
					continue nextStack
				}
			}
			stack = append(stack, state{i + len(t.lit), j + 1})
		case anyOneToken:
			if i >= len(hr) {
				continue
			}
			stack = append(stack, state{i + 1, j + 1})
		case anyManyToken:
			// % matches empty
			stack = append(stack, state{i, j + 1})
			// % matches one rune
			if i < len(hr) {
				stack = append(stack, state{i + 1, j})
			}
		}
	}
	return bestEnd
}

func runeEqualFold(a, b rune) bool {
	if a == b {
		return true
	}
	for r := unicode.SimpleFold(a); r != a; r = unicode.SimpleFold(r) {
		if r == b {
			return true
		}
	}
	return false
}
