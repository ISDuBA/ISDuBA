#!/usr/bin/env bash

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

# This will work if the standard or custom has been set in env. If neither, this will fail.
if [[ -z "${KEYCLOAK_ADMIN}" ]]; then
  export KEYCLOAK_ADMIN="keycloak"
  echo "No Keycloak admin set. Assuming admin with name \"keycloak\""
else
  export KEYCLOAK_ADMIN="${KEYCLOAK_ADMIN}"
fi

if [[ -z "${KEYCLOAK_ADMIN_PASSWORD}" ]]; then
  export KEYCLOAK_ADMIN_PASSWORD="keycloak"
  echo "No Keycloak admin password set. Assuming admin with password \"keycloak\""
else
  export KEYCLOAK_ADMIN_PASSWORD="${KEYCLOAK_ADMIN_PASSWORD}"
fi

sudo /opt/keycloak/bin/kcadm.sh config credentials --server http://localhost:8080 --realm master --user "$KEYCLOAK_ADMIN" --password "$KEYCLOAK_ADMIN_PASSWORD"


if sudo /opt/keycloak/bin/kcadm.sh get 'http://localhost:8080/admin/realms/isduba/roles' | grep -F -q "\"name\" : \"$1\"", ; then
  echo "Role $1 already exists."
else
  # create role
  sudo /opt/keycloak/bin/kcadm.sh create roles --target-realm=isduba --set name=$1 \
      --set "description=$2"
fi
