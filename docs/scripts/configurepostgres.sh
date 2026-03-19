#!/usr/bin/env bash

# This file is Free Software under the Apache-2.0 License
# without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
#
# SPDX-License-Identifier: Apache-2.0
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

set -e # to exit if a command in the script fails


# Alter PostgreSQL as postgres user
psql -c "ALTER USER postgres WITH PASSWORD 'postgres';"


# Create keycloak user if it does not exist
if psql -c "\du" | grep -q -F "keycloak"
then
  echo "User keycloak already exists."
else
  psql -c "CREATE USER keycloak WITH PASSWORD 'keycloak';"
fi

# Create keycloak database if it doesn't exist
if psql -c "\l" | grep -q -F "keycloak"
then
  echo "Database keycloak already exists."
else
  createdb -O keycloak -E 'UTF-8' keycloak
fi

# Adjust postgresql configuration
PG_HBA_PATH=$(psql -t -P format=unaligned -c 'SHOW hba_file;')

if ! grep -q -F "# ISDuBA configuration" $PG_HBA_PATH;
then
tee -a $PG_HBA_PATH <<block_to_insert > /dev/null
# ISDuBA configuration
host    all             all             127.0.0.1/32            scram-sha-256
host    all             all             ::1/128                 scram-sha-256
block_to_insert
fi

echo "Adjusted postgresql configuration"
