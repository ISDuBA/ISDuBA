// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

package query

import (
	"cmp"
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
	// Ensure that the sections are ascending.
	slices.SortFunc(sections, func(a, b [2]int) int {
		return cmp.Compare(a[0], b[0])
	})
	return sections
}

// Shorten returns a shorten version of the given string.
// buffer is a buffer in number of runes around around the sections
// to give reading context. fill is a string to be used as filler for gaps
// (think "..."). delims is a pair for delimeters to mark the sections.
func (ts TextSections) Shorten(s string, buffer int, fill string, delims [2]string) string {

	var (
		instrs []func()
		out    []rune
	)

	xfer := func(in []rune, start, end int) func() {
		start, end = max(0, start), min(len(in), end)
		return func() {
			out = append(out, in[start:end]...)
		}
	}
	delim0 := func() { out = append(out, []rune(delims[0])...) }
	delim1 := func() { out = append(out, []rune(delims[1])...) }
	filler := func() { out = append(out, []rune(fill)...) }

	add := func(inst func()) { instrs = append(instrs, inst) }

	sX := []rune(s)

	for i, section := range ts {
		if i > 0 {
			last := ts[i-1]
			endLast := last[0] + last[1]
			gap := section[0] - endLast
			if gap > 2*buffer {
				add(filler)
				// TODO: Implement me!

			}
		}
		add(delim0)
		add(xfer(sX, section[0], section[0]+section[1]))
		add(delim1)
	}

	for _, instr := range instrs {
		instr()
	}

	return string(out)

	/*

		// (over) estimate number of used runes.
		n := 0
		for _, s := range ts {
			n += 2*buffer + s[1]
		}
		used := make(map[int]struct{}, n)
		sX := []rune(s)
		for _, s := range ts {
			start, end := s[0]-buffer, s[0]+s[1]+buffer
			for idx := start; idx < end; idx++ {
				if idx >= 0 && idx < len(sX) {
					used[idx] = struct{}{}
				}
			}
		}
		indices := slices.Sorted(maps.Keys(used))
		if len(indices) == 0 {
			return ""
		}
		// TODO: Add delims and fills.
		_ = fill
		_ = delims

		out := make([]rune, 0, len(indices))
		for _, idx := range indices {
			out = append(out, sX[idx])
		}
		return string(out)
	*/
}
