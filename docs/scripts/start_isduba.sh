#!/usr/bin/env bash

# This file is Free Software under the Apache-2.0 License
# without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
#
# SPDX-License-Identifier: Apache-2.0
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

set -e # to exit if a command in the script fails

./../../bin/isdubad -c ../../isdubad.toml &

# TODO: race conditions
if curl --head --silent http://localhost:5173/ | grep -F -q "HTTP/1.1 200 OK"; then
  echo "Port 5173 is already being used. Is isduba already running?"
else
 cd ./../../client/

 npm run dev &

fi
