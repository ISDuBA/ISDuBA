!/usr/bin/env bash

# This file is Free Software under the Apache-2.0 License
# without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
#
# SPDX-License-Identifier: Apache-2.0
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

set -e # to exit if a command in the script fails

# arguments:
# $1: name
# $2: description

if sudo /opt/keycloak/bin/kcadm.sh get 'http://localhost:8080/admin/realms/isduba/roles' | grep -F -q "\"name\" : \"$1\"", ; then
  echo "Role $1 already exists."
else
  # create role
  sudo /opt/keycloak/bin/kcadm.sh create roles --target-realm=isduba --set name=$1 \
      --set "description=$2"
fi
