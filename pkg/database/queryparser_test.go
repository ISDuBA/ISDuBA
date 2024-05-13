// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package database

import (
	"reflect"
	"testing"
)

func TestSplit(t *testing.T) {
	for _, x := range []struct {
		have     string
		expectVs []string
		expectBs []bool
	}{
		{` hello  `, []string{`hello`}, []bool{false}},
		{`hello`, []string{`hello`}, []bool{false}},
		{`hello world`, []string{`hello`, `world`}, []bool{false, false}},
		{`"hello" world`, []string{`hello`, `world`}, []bool{true, false}},
		{`"hello world"`, []string{`hello world`}, []bool{true}},
		{`hello\ world`, []string{`hello world`}, []bool{false}},
		{`"hello"\ world`, []string{`hello`, ` world`}, []bool{true, false}},
		{`"hello\" world"`, []string{`hello" world`}, []bool{true}},
		{`"hello world`, []string{`hello world`}, []bool{true}},
	} {
		var vs []string
		var bs []bool
		split(x.have, func(v string, b bool) {
			vs = append(vs, v)
			bs = append(bs, b)
		})
		if !reflect.DeepEqual(vs, x.expectVs) {
			t.Errorf("have %+q: expected %+q got %+q", x.have, x.expectVs, vs)
		}
		if !reflect.DeepEqual(bs, x.expectBs) {
			t.Errorf("have %+v: expected %+v got %+v", x.have, x.expectBs, bs)
		}
	}
}
