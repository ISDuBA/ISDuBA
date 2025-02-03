// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package models

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/ISDuBA/ISDuBA/pkg/database/query"
)

func TestAsConditions(t *testing.T) {
	for _, x := range []struct {
		input        string
		expected     string
		replacements []any
	}{
		{
			`{"*": [ "WHITE", "GREEN" ]}`,
			`(((((tlp)=($1)))OR(((tlp)=($2)))))`,
			[]any{"WHITE", "GREEN"},
		}, {
			`{}`,
			`(FALSE)`,
			[]any{},
		}, {
			`{"A": [ "WHITE", "GREEN" ]}`,
			`(((((advisories.publisher)=($1)))AND(((((tlp)=($2)))OR(((tlp)=($3)))))))`,
			[]any{"A", "WHITE", "GREEN"},
		}, {
			`{"A": [ "AMBER", "RED" ], "*": ["WHITE"]}`,
			`(((((((advisories.publisher)=($1)))AND(((((tlp)=($2)))OR(((tlp)=($3)))))))OR(((((tlp)=($4)))AND((NOT (((advisories.publisher)=($1)))))))))`,
			[]any{"A", "AMBER", "RED", "WHITE"},
		}, {
			`{"A": [ "AMBER", "RED" ], "*": ["WHITE", "GREEN"]}`,
			`(((((((advisories.publisher)=($1)))AND(((((tlp)=($2)))OR(((tlp)=($3)))))))OR(((((((tlp)=($4)))OR(((tlp)=($5)))))AND((NOT (((advisories.publisher)=($1)))))))))`,
			[]any{"A", "AMBER", "RED", "WHITE", "GREEN"},
		}, {
			`{"A": [ "AMBER" ], "B": ["RED"], "*": ["WHITE"]}`,
			`(((((((((advisories.publisher)=($1)))AND(((tlp)=($2)))))OR(((((advisories.publisher)=($3)))AND(((tlp)=($4)))))))OR(((((tlp)=($5)))AND((NOT (((((advisories.publisher)=($1)))OR(((advisories.publisher)=($3)))))))))))`,
			[]any{"A", "AMBER", "B", "RED", "WHITE"},
		}, {
			`{"*": ["WHITE"], "A": [ "AMBER" ], "B": ["RED"]}`,
			`(((((((((advisories.publisher)=($1)))AND(((tlp)=($2)))))OR(((((advisories.publisher)=($3)))AND(((tlp)=($4)))))))OR(((((tlp)=($5)))AND((NOT (((((advisories.publisher)=($1)))OR(((advisories.publisher)=($3)))))))))))`,
			[]any{"A", "AMBER", "B", "RED", "WHITE"},
		},
	} {
		var ptlps PublishersTLPs
		if err := json.Unmarshal([]byte(x.input), &ptlps); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		expr := ptlps.AsExpr()
		builder := query.SQLBuilder{}
		have := builder.CreateWhere(expr)
		if x.expected != have {
			t.Errorf("input: %s have: %s, expected: %s", x.input, have, x.expected)
		}
		if !slices.Equal(x.replacements, builder.Replacements) {
			t.Errorf("input: %s have: %q expected: %q", x.input, builder.Replacements, x.replacements)
		}
	}
}

func TestUnmarshalText(t *testing.T) {
	for _, x := range []struct {
		input    string
		expected TLP
		error    bool
	}{
		{"WHITE", TLPWhite, false},
		{"GREEN", TLPGreen, false},
		{"AMBER", TLPAmber, false},
		{"RED", TLPRed, false},
		{"HEARD", "", true},
	} {
		var have TLP
		if err := have.UnmarshalText([]byte(x.input)); err != nil {
			if !x.error {
				t.Errorf("%s: should err but doesnt.", x.input)
			}
			continue
		}
		if have != x.expected {
			t.Errorf("%s: have %q expected %q", x.input, have, x.expected)
		}
	}
}
