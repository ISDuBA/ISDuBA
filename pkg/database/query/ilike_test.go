// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

package query

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestCompileILike(t *testing.T) {
	expr := MustCompileILike(`%Hal_\o%%\`)

	have := expr.String()

	const expected = `(?i:(?:.*(Hal).(o).*(\\)))`

	if have != expected {
		t.Errorf("have: %v expected: %v", have, expected)
	}

	havePairs := expr.Search(`xxxhallo\`)
	expectedPairs := TextSections{
		{3, 3},
		{7, 1},
		{8, 1},
	}
	if !reflect.DeepEqual(havePairs, expectedPairs) {
		t.Errorf("pairs: have: %v expected: %v", havePairs, expectedPairs)
	}
}

func indexUTF8(s string) func(int) int {
	idx := 0
	pos := make(map[int]int, len(s))
	for i := range s {
		pos[i] = idx
		idx++
	}
	return func(idx int) int { return pos[idx] }
}

func TestCarriageReturn(t *testing.T) {

	txtRaw, err := os.ReadFile("search_response.txt")
	if err != nil {
		t.Errorf("cannot load file: %v\n", err)
	}

	txt := string(txtRaw)
	runes := []rune(txt)
	index := indexUTF8(txt)

	/*
		fmt.Println(len(runes), len(txtRaw))
		for _, r := range runes {
			if r > 127 {
				fmt.Printf("\trune: %c\n", r)
			}
		}
	*/
	expr := MustCompileILike(`%header%`)
	havePairs := expr.Search(txt)
	for _, pair := range havePairs {
		have := txt[pair[0] : pair[0]+pair[1]]
		if strings.ToUpper(have) != "HEADER" {
			t.Errorf("pair %+v results in %q not \"HEADER\"", pair, have)
		}
		start := index(pair[0])
		end := index(pair[0] + pair[1])
		runeSliced := string(runes[start:end])
		if strings.ToUpper(runeSliced) != "HEADER" {
			t.Errorf("rune slicinf %+v results in %q not \"HEADER\"", pair, runeSliced)
		}
	}
}
