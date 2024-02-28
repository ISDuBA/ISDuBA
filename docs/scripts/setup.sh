#!/usr/bin/env bash

# This file is Free Software under the MIT License
# without warranty, see README.md and LICENSES/MIT.txt for details.
#
# SPDX-License-Identifier: MIT
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

set -e # to exit if a command in the script fails

./installgojava.sh # installs go and java

./installpostgres.sh # installs postgreSQL

./configurepostgres.sh # creates necessary postgres users and databases

./installkeycloak.sh # installs keycloak

./configurekeycloak.sh # configures keycloak

./keycloakonsystemstart.sh # adjust systemd to allow keycloak to start on systemstartup

./installplaywright.sh # Prepare frontend

./installisduba.sh # build the isdubad and bulkimporter tools
