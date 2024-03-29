#!/usr/bin/env bash

# This file is Free Software under the Apache-2.0 License
# without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
#
# SPDX-License-Identifier: Apache-2.0
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

set -e # to exit if a command in the script fails

sudo apt-get update # Make sure to be up-to-date

./installgojava.sh # installs go and java

./installpostgres.sh # installs postgreSQL

./configurepostgres.sh # creates necessary postgres users and databases

./installkeycloak.sh # installs keycloak

./configurekeycloak.sh # configures keycloak

./installplaywright.sh # Prepare frontend

./installisduba.sh # build the isdubad and bulkimporter tools

echo "All set up!"
