// This file is Free Software under the MIT License
// without warranty, see README.md and LICENSES/MIT.txt for details.
//
// SPDX-License-Identifier: MIT
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

// Package version implements the burned-in version of the binaries.
package version

// SemVersion the version in semver.org format, MUST be overwritten during
// the linking stage of the build process
var SemVersion = "0.0.0"
