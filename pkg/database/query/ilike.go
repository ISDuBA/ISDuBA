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
	"regexp"
	"slices"
	"strings"
)

// TextSections is a list of positions in string.
// The [2]ints are to be read as [0] being the start position (zero based)
// and [1] the length of the covered section.
type TextSections [][2]int

// ILikeExpr is an compiled ILIKE expression.
type ILikeExpr struct{ *regexp.Regexp }

// MustCompileILike calls CompileILike and panics if the
// compiling failed.
func MustCompileILike(needles ...string) ILikeExpr {
	expr, err := CompileILike(needles...)
	if err != nil {
		panic(fmt.Sprintf("compiling ilike %v failed: %v", needles, err))
	}
	return expr
}

// CompileILike compiles an ILIKE search pattern.
func CompileILike(needles ...string) (ILikeExpr, error) {

	var pattern, buf strings.Builder

	flushLiteral := func() {
		if buf.Len() > 0 {
			pattern.WriteByte('(')
			pattern.WriteString(regexp.QuoteMeta(buf.String()))
			pattern.WriteByte(')')
			buf.Reset()
		}
	}
	pattern.WriteString(`(?i:`)
	for i, needle := range needles {
		if i > 0 {
			pattern.WriteByte('|')
		}
		pattern.WriteString(`(?:`)
		escape, many := false, false
		for _, r := range needle {
			if escape {
				escape = false
				buf.WriteRune(r)
				continue
			}
			switch r {
			case '\\':
				escape, many = true, false
			case '%':
				if !many {
					many = true
					flushLiteral()
					pattern.WriteString(".*")
				}
			case '_':
				many = false
				flushLiteral()
				pattern.WriteString(".")
			default:
				many = false
				buf.WriteRune(r)
			}
		}
		if escape {
			buf.WriteByte('\\')
		}
		flushLiteral()
		pattern.WriteByte(')')
	}
	pattern.WriteByte(')')
	expr, err := regexp.Compile(pattern.String())
	return ILikeExpr{expr}, err
}

// Search searches ilike patterns in a haystack.
// Returns a list of matching positions.
func (expr ILikeExpr) Search(haystack string) TextSections {
	all := expr.FindAllStringSubmatchIndex(haystack, -1)
	if len(all) == 0 {
		return nil
	}
	var sections TextSections
	for _, indices := range all {
		if len(indices) == 0 {
			continue
		}
		for indices = indices[2:]; len(indices) > 0; indices = indices[2:] {
			pair := [2]int{indices[0], indices[1] - indices[0]}
			if !slices.Contains(sections, pair) {
				sections = append(sections, pair)
			}
		}
	}
	return sections
}

// Shorten returns a shorten version of the given string.
// buffer is a buffer in number of runes around around the sections
// to give reading context. fill is a string to be used as filler for gaps
// (think "..."). delims is a pair for delimeters to mark the sections.
func (ts TextSections) Shorten(s string, buffer int, fill string, delims [2]string) string {
	// TODO: Implement me!
	_ = buffer
	_ = fill
	_ = delims
	return s
}
