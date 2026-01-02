#!/usr/bin/env bash

# This file is Free Software under the Apache-2.0 License
# without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
#
# SPDX-License-Identifier: Apache-2.0
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

set -e # to exit if a command in the script fails

args=()
while [[ $# -gt 0 ]]; do
  case "$1" in
    -k|--keycloakRunning)
      echo "Assuming keycloak is running..."
      args+=(-k)
      ;;
    -q|--quick)
      args+=(-q)
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
  shift
done

./installutilities.sh # installs utilities like go and java

./installpostgres.sh # installs postgreSQL

sudo systemctl start postgresql.service # starts just installed postgres

sudo -u postgres bash ./configurepostgres.sh # creates necessary postgres users and databases

./installkeycloak.sh # installs keycloak

cd keycloak

./configureKeycloak.sh "${args[@]}"

cd ..

./installplaywright.sh # prepare frontend

./installisduba.sh # build the isdubad and bulkimporter tools

cd ../..

# migrate the database so it is up-to-date
ISDUBA_DB_MIGRATE=true ./cmd/isdubad/isdubad -c ./isduba.toml

echo "All set up!"
