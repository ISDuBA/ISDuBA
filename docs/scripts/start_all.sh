#!/usr/bin/env bash

# This file is Free Software under the Apache-2.0 License
# without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
#
# SPDX-License-Identifier: Apache-2.0
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

set -e # to exit if a command in the script fails

echo "This script is designed to be used when the client, the backend and the keycloak are not running."

cd ../..

./cmd/isdubad/isdubad &

isdubad=$!

echo "isdubad is running:"
echo $isdubad

cd client

npm run dev &

client=$!

echo "the client is running:"
echo $client

sudo /opt/keycloak/bin/kc.sh start-dev --health-enabled=true &

# wait for keycloak to start
echo "Waiting for keycloak to start..."
until curl --silent --head -fsS http://localhost:9000/health/ready
do
  sleep 5
done
keycloak=$!

echo "keycloak is running:"
echo $keycloak
