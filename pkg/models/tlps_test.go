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
			`{"*": [ "CLEAR", "GREEN" ]}`,
			`(((((tlp)=($1)))OR(((tlp)=($2)))))`,
			[]any{"CLEAR", "GREEN"},
		}, {
			`{}`,
			`(FALSE)`,
			[]any{},
		}, {
			`{"A": [ "CLEAR", "GREEN" ]}`,
			`(((((advisories.publisher)=($1)))AND(((((tlp)=($2)))OR(((tlp)=($3)))))))`,
			[]any{"A", "CLEAR", "GREEN"},
		}, {
			`{"A": [ "AMBER", "RED" ], "*": ["CLEAR"]}`,
			`(((((((advisories.publisher)=($1)))AND(((((tlp)=($2)))OR(((tlp)=($3)))))))OR(((((tlp)=($4)))AND((NOT (((advisories.publisher)=($1)))))))))`,
			[]any{"A", "AMBER", "RED", "CLEAR"},
		}, {
			`{"A": [ "AMBER", "RED" ], "*": ["CLEAR", "GREEN"]}`,
			`(((((((advisories.publisher)=($1)))AND(((((tlp)=($2)))OR(((tlp)=($3)))))))OR(((((((tlp)=($4)))OR(((tlp)=($5)))))AND((NOT (((advisories.publisher)=($1)))))))))`,
			[]any{"A", "AMBER", "RED", "CLEAR", "GREEN"},
		}, {
			`{"A": [ "AMBER" ], "B": ["RED"], "*": ["CLEAR"]}`,
			`(((((((((advisories.publisher)=($1)))AND(((tlp)=($2)))))OR(((((advisories.publisher)=($3)))AND(((tlp)=($4)))))))OR(((((tlp)=($5)))AND((NOT (((((advisories.publisher)=($1)))OR(((advisories.publisher)=($3)))))))))))`,
			[]any{"A", "AMBER", "B", "RED", "CLEAR"},
		}, {
			`{"*": ["CLEAR"], "A": [ "AMBER" ], "B": ["RED"]}`,
			`(((((((((advisories.publisher)=($1)))AND(((tlp)=($2)))))OR(((((advisories.publisher)=($3)))AND(((tlp)=($4)))))))OR(((((tlp)=($5)))AND((NOT (((((advisories.publisher)=($1)))OR(((advisories.publisher)=($3)))))))))))`,
			[]any{"A", "AMBER", "B", "RED", "CLEAR"},
		}, {
			`{"*": [ "CLEAR" ]}`,
			`(((tlp)=($1)))`,
			[]any{"CLEAR"},
		}, {
			`{"*": [ "AMBER", "AMBER+STRICT" ]}`,
			`(((((tlp)=($1)))OR(((tlp)=($2)))))`,
			[]any{"AMBER", "AMBER+STRICT"},
		}, {
			`{"*": [ "AMBER+STRICT" ]}`,
			`(((tlp)=($1)))`,
			[]any{"AMBER+STRICT"},
		},
	} {
		var ptlps PublishersTLPs
		if err := json.Unmarshal([]byte(x.input), &ptlps); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		expr := ptlps.AsExpr()

		builder, err := query.NewSQLBuilder(
			query.SQLBuilderExpr(expr),
		)

		if err != nil {
			t.Errorf("input: %s builder creation failed: %v", x.input, err)
		}

		have := builder.CreateWhereSQL()

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
		{"WHITE", "", true},
		{"CLEAR", TLPClear, false},
		{"GREEN", TLPGreen, false},
		{"AMBER", TLPAmber, false},
		{"AMBER+STRICT", TLPAmberStrict, false},
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
