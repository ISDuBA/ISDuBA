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

echo "This script is not meant to be started manually. It exists to support the integration tests."
echo "This script may not function as expected if isdubad, keycloak or the client have already been started."

./cmd/isdubad/isdubad &

isdubad=$!

echo "isdubad has been started:"
echo $isdubad

cd client

npm run dev &

client=$!

echo "the client has been started:"
echo $client

sudo /opt/keycloak/bin/kc.sh start-dev --health-enabled=true &

# wait for keycloak to start
echo "Waiting for keycloak to start..."
until curl --silent --head -fsS http://localhost:9000/health/ready
do
  sleep 5
done
keycloak=$!

echo "keycloak has been started:"
echo $keycloak
