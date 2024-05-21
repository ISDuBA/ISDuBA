#!/usr/bin/env bash

# This file is Free Software under the Apache-2.0 License
# without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
#
# SPDX-License-Identifier: Apache-2.0
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

set -e # to exit if a command in the script fails

keycloak_running=false
while [[ $# -gt 0 ]]; do
  case "$1" in
    -k|--keycloakRunning)
      echo "Assuming keycloak is running..."
      keycloak_running=true
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
  shift
done

#./installgojava.sh # installs go and java

#./installpostgres.sh # installs postgreSQL

#./configurepostgres.sh # creates necessary postgres users and databases

#./installkeycloak.sh # installs keycloak

if $keycloak_running; then
  ./keycloak/configureKeycloak.sh -k # configures keycloak
else
  ./keycloak/configureKeycloak.sh # configures keycloak
fi

#./installplaywright.sh # prepare frontend

#./installisduba.sh # build the isdubad and bulkimporter tools

#./migrate.sh # prepare isduba database

echo "All set up!"
