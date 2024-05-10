#!/usr/bin/env bash

# This file is Free Software under the Apache-2.0 License
# without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
#
# SPDX-License-Identifier: Apache-2.0
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

set -e # to exit if a command in the script fails

# Create keycloak admin user
if [[ -z "${KEYCLOAK_ADMIN}" ]]; then
  export KEYCLOAK_ADMIN="keycloak"
  echo "No Keycloak admin set. Creating admin with name \"keycloak\""
else
  export KEYCLOAK_ADMIN="${KEYCLOAK_ADMIN}"
fi

if [[ -z "${KEYCLOAK_ADMIN_PASSWORD}" ]]; then
  export KEYCLOAK_ADMIN_PASSWORD="keycloak"
  echo "No Keycloak admin password set. Creating admin with password \"keycloak\""
else
  export KEYCLOAK_ADMIN_PASSWORD="${KEYCLOAK_ADMIN_PASSWORD}"
fi

# Alter the keycloak configuration
sudo sed --in-place=.orig -e 's/^#db=postgres/db=postgres/' \
            -e 's/^#db-username=/db-username=/' \
            -e 's/^#db-password=password/db-password=keycloak/' \
            -e 's/^#db-url=/db-url=/' /opt/keycloak/conf/keycloak.conf

nkid=""

# Give feedback after successful completion
echo "Succesfully adjusted keycloaks configuration."

# TODO: what if keycloak is running, but does not have an admin user yet?

if curl --silent --head -fsS http://localhost:8080/health/ready; then
  echo "keycloak is already running..."
else
  sudo --preserve-env=KEYCLOAK_ADMIN,KEYCLOAK_ADMIN_PASSWORD /opt/keycloak/bin/kc.sh start-dev --health-enabled=true &

  # wait for keycloak to start
  echo "Waiting for keycloak to start..."
  until curl --silent --head -fsS http://localhost:8080/health/ready
  do
    sleep 1
  done
  nkid=$!
fi

# log into the master realm with admin rights, token saved in ~/.keycloak/kcadm.config
sudo /opt/keycloak/bin/kcadm.sh config credentials --server http://localhost:8080 --realm master --user "$KEYCLOAK_ADMIN" --password "$KEYCLOAK_ADMIN_PASSWORD"

if sudo /opt/keycloak/bin/kcadm.sh get 'http://localhost:8080/admin/realms' | grep -F -q '"realm" : "isduba",' ; then
 echo "Realm isduba already exists."
else
./keycloak/createRealm.sh
fi

./keycloak/createGroup.sh 'editor' '1'

./keycloak/createRole.sh 'editor' 'Bearbeiter'

./keycloak/createUser.sh 'beate' 'beate' 'bear' 'bea@ISDuBA.isduba' 'beate' 'editor'

./keycloak/assignUserToRoleAndGroup.sh 'beate' 'editor' 'editor'

if [ ! -z "$nkid" ]; then
  if ps -p $nkid; then
    kill $nkid
  fi
fi
