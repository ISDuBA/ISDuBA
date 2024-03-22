#!/usr/bin/env bash

# This file is Free Software under the Apache-2.0 License
# without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
#
# SPDX-License-Identifier: Apache-2.0
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

set -e # to exit if a command in the script fails


LAB="#listen_addresses = 'localhost'" # Listen Adress Before
LAA="listen_addresses = '*'"          # Listen Adress After

# Alter PostgreSQL as postgres user
sudo -u postgres psql -c "ALTER USER postgres WITH PASSWORD 'postgres';"


# Create keycloak user if it does not exist
if sudo -u postgres psql -c "\du" | grep -q -F "keycloak"
then
  echo "User keycloak already exists."
else
  sudo -u postgres psql -c "CREATE USER keycloak WITH PASSWORD 'keycloak';"
fi

# Create keycloak database if it doesn't exist
if sudo -u postgres psql -c "\l" | grep -q -F "keycloak"
then
  echo "Database keycloak already exists."
else
  sudo -u postgres createdb -O keycloak -E 'UTF-8' keycloak
fi

# Adjust keycloak configuration
sudo -u postgres sed -i "s/$LAB/$LAA/g" /etc/postgresql/16/main/postgresql.conf
if ! sudo grep -q -F "# ISDuBA configuration" /etc/postgresql/16/main/pg_hba.conf;
then
sudo -u postgres tee -a /etc/postgresql/16/main/pg_hba.conf <<block_to_insert > /dev/null
# ISDuBA configuration
host    all             all             192.168.56.1/32         scram-sha-256
host    all             all             127.0.0.1/32            scram-sha-256
block_to_insert
fi

echo "Adjusted keycloak configuration"
