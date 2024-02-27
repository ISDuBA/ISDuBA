#!/usr/bin/env bash

# This file is Free Software under the MIT License
# without warranty, see README.md and LICENSES/MIT.txt for details.
#
# SPDX-License-Identifier: MIT
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

set -e # to exit if a command in the script fails

# build the isdubad tool
cd ../../cmd/isdubad/

go build

# build the bulkimporter
cd ../bulkimport/

go build

# go back into the scripts directory
cd ../../docs/scripts

# create the isdubad configuration
cp ../example_isdubad.toml ../../isdubad.toml
