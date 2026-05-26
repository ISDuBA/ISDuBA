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
	expectedPairs := TextSections{
		{3, 3},
		{7, 1},
		{8, 1},
	}
	if !reflect.DeepEqual(havePairs, expectedPairs) {
		t.Errorf("pairs: have: %v expected: %v", havePairs, expectedPairs)
	}
}

func TestCarriageReturn(t *testing.T) {
	const txt = "PC contains a vulnerability that allows hpack table accounting errors could lead to " +
		"unwanted disconnects between clients and servers in exceptional cases. Three vectors were found " +
		"that allow the following DOS attacks: - Unbounded memory buffering in the HPACK parser - Unbounded " +
		"CPU consumption in the HPACK parser The unbounded CPU consumption is down to a copy that occurred " +
		"per-input-block in the parser, and because that could be unbounded due to the memory copy bug we " +
		"end up with a parsing loop, with n selected by the client. The unbounded memory buffering bugs: - " +
		"The header size limit check was behind the string reading code, so we needed to first buffer up to " +
		"a 4 gigabyte string before rejecting it as longer than 8 or 16kb. - HPACK varints have an encoding " +
		"quirk whereby an infinite number of 0’s can be added at the start of an integer. gRPC’s hpack " +
		"parser needed to read all of them before concluding a parse. - gRPC’s metadata overflow check " +
		"was performed per frame, so that the following sequence of frames could cause infinite buffering: " +
		"HEADERS: containing a: 1 CONTINUATION: containing a: 2 CONTINUATION: containing a: 3 etc…\r+ " +
		" - Unbounded memory buffering in the HPACK parser\r"
	expr := MustCompileILike(`%HEADERS%`)
	havePairs := expr.Search(txt)
	for _, pair := range havePairs {
		if have := txt[pair[0] : pair[0]+pair[1]]; have != "HEADERS" {
			t.Errorf("pair %+v results in %q not \"HEADERS\"", pair, have)
		}
	}
}
