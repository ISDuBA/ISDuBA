#!/usr/bin/env bash

# This file is Free Software under the Apache-2.0 License
# without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
#
# SPDX-License-Identifier: Apache-2.0
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

set -e # to exit if a command in the script fails

cd ../..

mkdir -p bin

go build -o ./bin/isdubad ./cmd/isdubad

go build -o ./bin/bulkimport ./cmd/bulkimport
echo "Successfully created isdubad and bulkimport binaries."

# create the isdubad configuration
cp docs/example_isdubad.toml isdubad.toml
echo "Successfully created example-configuration isdubad.toml."
