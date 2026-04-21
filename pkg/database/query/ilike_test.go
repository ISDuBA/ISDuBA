// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

package query

import (
	"reflect"
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
	expectedPairs := [][2]int{
		{3, 3},
		{7, 1},
		{8, 1},
	}
	if !reflect.DeepEqual(havePairs, expectedPairs) {
		t.Errorf("pairs: have: %v expected: %v", havePairs, expectedPairs)
	}
}
