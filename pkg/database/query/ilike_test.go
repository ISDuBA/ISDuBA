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
	have := CompileILike(`%Hal_\o%%\`)
	expected := ILikeExpr{
		{kind: anyManyToken},
		{kind: litToken, lit: []rune("Hal")},
		{kind: anyOneToken},
		{kind: litToken, lit: []rune("o")},
		{kind: anyManyToken},
		{kind: litToken, lit: []rune{'\\'}},
	}
	if !reflect.DeepEqual(have, expected) {
		t.Errorf("have: %+v expected: %+v", have, expected)
	}
}
