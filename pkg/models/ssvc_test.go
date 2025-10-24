// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package models

import "testing"

func TestSSVC(t *testing.T) {
	for _, input := range []struct {
		vector string
		pass   bool
	}{
		{"SSVCv2/E:N/A:N/T:P/P:M/B:M/M:L/D:T/2024-03-13T10:33:45Z/", true},
		{"SSVCv2/E:N/A:N/T:P/M:L/D:T/2024-03-13T10:34:39Z/", true},
		{"SSVCv2/E:N/A:N/T:P/M:L/D:T/2024-03-13T10:34:39Z", false},
		{"XXX/E:N/A:N/T:P/M:L/D:T/2024-03-13T10:34:39Z/", false},
		{"SSVCv2/E:N/A:N/T:P/M:L/D:T/XXX/", false},
		{"SSVCv2/2024-03-13T10:34:39Z/", false},
		{"SSVCv2/N/A:N/T:P/M:L/D:T/2024-03-13T10:34:39Z/", false},
		{"SSVCv2/💩:N/A:N/T:P/M:L/D:T/2024-03-13T10:34:39Z/", false},
		{"SSVCv2/E:N/A:N/T:P/M:💩/D:T/2024-03-13T10:34:39Z/", false},
	} {
		err := ValidateSSVCv2Vector(input.vector)
		if err != nil && input.pass {
			t.Errorf("%q failed to validate: %v", input.vector, err)
		}
		if err == nil && !input.pass {
			t.Errorf("expected %q to fail", input.vector)
		}
	}
}
