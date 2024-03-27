#!/usr/bin/env bash

# This file is Free Software under the Apache-2.0 License
# without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
#
# SPDX-License-Identifier: Apache-2.0
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

set -e # to exit if a command in the script fails

# migrate
ISDUBA_DB_MIGRATE=true ./../../bin/isdubad -c ../../isdubad.toml &

touch isduba.log # to ensure file exists for grep
until grep -q -F "Starting web server" isduba.log
do
  sleep 1
done

mv isduba.log isduba_migrate.log

if ps -p $!; then
  echo "Migration successful"
  kill $!
else
 echo "Couldn't start migration, is isdubad already running?"
fi
