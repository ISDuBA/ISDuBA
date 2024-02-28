// This file is Free Software under the MIT License
// without warranty, see README.md and LICENSES/MIT.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package models

import (
	"encoding/json"
	"regexp"
	"strings"
	"testing"
)

func TestAsConditions(t *testing.T) {
	stripSpace := regexp.MustCompile(`\s+`)
	for _, x := range []struct {
		input    string
		expected string
	}{
		{`[{"publisher": "", "tlps": [ "WHITE", "GREEN" ]}]`, `$tlp WHITE = $tlp GREEN = or`},
		{`[{"publisher": "A", "tlps": [ "WHITE", "GREEN" ]}]`, `$tlp WHITE = $tlp GREEN = or $publisher "A" = and`},
		{`[{"publisher": "A", "tlps": [ "AMBER", "RED" ]}, {"publisher": "", "tlps": ["WHITE"]}]`,
			`$tlp AMBER = $tlp RED = or $publisher "A" = and $tlp WHITE = $publisher "A" != and or`},
		{`[{"publisher": "A", "tlps": [ "AMBER", "RED" ]}, {"publisher": "", "tlps": ["WHITE", "GREEN"]}]`,
			`$tlp AMBER = $tlp RED = or $publisher "A" = and $tlp WHITE = $tlp GREEN = or $publisher "A" != and or`},
		{`[{"publisher": "A", "tlps": [ "AMBER" ]}, {"publisher": "B", "tlps": ["RED"]}, {"publisher": "", "tlps": ["WHITE"]}]`,
			`$tlp AMBER = $publisher "A" = and $tlp RED = $publisher "B" = and or $tlp WHITE = $publisher "A" != $publisher "B" != and and or`},
		{`[{"publisher": "", "tlps": ["WHITE"]}, {"publisher": "A", "tlps": [ "AMBER" ]}, {"publisher": "B", "tlps": ["RED"]}]`,
			`$tlp AMBER = $publisher "A" = and $tlp RED = $publisher "B" = and or $tlp WHITE = $publisher "A" != $publisher "B" != and and or`},
	} {
		var ptlps PuplishersTLPs
		if err := json.Unmarshal([]byte(x.input), &ptlps); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		have := ptlps.AsConditions()
		have = strings.TrimSpace(have)
		have = stripSpace.ReplaceAllString(have, " ")
		if have != x.expected {
			t.Errorf("%s: have %q expect: %q", x.input, have, x.expected)
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
