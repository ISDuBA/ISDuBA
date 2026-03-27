// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

package itertools

import (
	"slices"
	"strconv"
	"testing"
)

func TestConcat(t *testing.T) {
	a := slices.Values([]int{1, 2, 3})
	b := slices.Values([]int{1, 2, 3})
	c := slices.Values([]int{1, 2, 3})
	have := slices.Collect(Concat(a, b, c))
	want := []int{1, 2, 3, 1, 2, 3, 1, 2, 3}
	if !slices.Equal(have, want) {
		t.Errorf("have: %v want: %v", have, want)
	}
}

func TestConcatEnumerate(t *testing.T) {
	a := slices.Values([]int{1, 2, 3})
	b := slices.Values([]int{1, 2, 3})
	c := slices.Values([]int{1, 2, 3})
	var have []int
	for i, v := range Enumerate(Concat(a, b, c)) {
		if i == 3 {
			break
		}
		have = append(have, v)
	}
	want := []int{1, 2, 3}
	if !slices.Equal(have, want) {
		t.Errorf("have: %v want: %v", have, want)
	}
}

func TestUnique(t *testing.T) {
	a := slices.Values([]int{1, 2, 3, 1, 4, 1, 2, 3})
	have := slices.Collect(Unique(a))
	want := []int{1, 2, 3, 4}
	if !slices.Equal(have, want) {
		t.Errorf("have: %v want: %v", have, want)
	}
}

func TestUniqueBreak(t *testing.T) {
	a := slices.Values([]int{1, 2, 3, 1, 4, 1, 2, 3})
	var have []int
	for i, v := range Enumerate(Unique(a)) {
		if i == 3 {
			break
		}
		have = append(have, v)
	}
	want := []int{1, 2, 3}
	if !slices.Equal(have, want) {
		t.Errorf("have: %v want: %v", have, want)
	}
}

func TestApply(t *testing.T) {
	a := slices.Values([]string{"1", "2", "3"})
	have := slices.Collect(Apply(a, func(s string) int {
		x, _ := strconv.Atoi(s)
		return x
	}))
	want := []int{1, 2, 3}
	if !slices.Equal(have, want) {
		t.Errorf("have: %v want: %v", have, want)
	}
}

func TestApplyBreak(t *testing.T) {
	a := slices.Values([]string{"1", "2", "3"})
	var have []int
	for i, v := range Enumerate(Apply(a, func(s string) int {
		x, _ := strconv.Atoi(s)
		return x
	})) {
		if i == 2 {
			break
		}
		have = append(have, v)
	}
	want := []int{1, 2}
	if !slices.Equal(have, want) {
		t.Errorf("have: %v want: %v", have, want)
	}
}

func TestNot(t *testing.T) {
	for _, want := range []bool{false, true} {
		accept := func(x bool) bool { return x == want }
		not := Not(accept)
		if not(want) == accept(want) {
			t.Errorf("have: %t want: %t", not(want), want)
		}
	}
}

func TestFilter(t *testing.T) {
	a := slices.Values([]int{1, 2, 3})
	have := slices.Collect(Filter(a, func(x int) bool { return x != 2 }))
	want := []int{1, 3}
	if !slices.Equal(have, want) {
		t.Errorf("have: %v want: %v", have, want)
	}
}

func TestFilterBreak(t *testing.T) {
	a := slices.Values([]int{0, 1, 2, 3})
	var have []int
	for i, v := range Enumerate(Filter(a, func(x int) bool {
		return x != 2
	})) {
		if i == 1 {
			break
		}
		have = append(have, v)
	}
	want := []int{0}
	if !slices.Equal(have, want) {
		t.Errorf("have: %v want: %v", have, want)
	}
}
