// This file is Free Software under the MIT License
// without warranty, see README.md and LICENSES/MIT.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package models

import "testing"

func TestSSVC(t *testing.T) {
	for _, vector := range []string{
		"SSVCv2/E:N/A:N/T:P/P:M/B:M/M:L/D:T/2024-03-13T10:33:45Z/",
		"SSVCv2/E:N/A:N/T:P/M:L/D:T/2024-03-13T10:34:39Z/",
	} {
		if err := ValidateSSVCv2Vector(vector); err != nil {
			t.Errorf("%q failed to validate: %v", vector, err)
		}
	}
}
