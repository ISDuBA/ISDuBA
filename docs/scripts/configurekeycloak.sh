#!/usr/bin/env bash

# This file is Free Software under the Apache-2.0 License
# without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
#
# SPDX-License-Identifier: Apache-2.0
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

set -e # to exit if a command in the script fails

# Alter the keycloak configuration
sudo sed --in-place=.orig -e 's/^#db=postgres/db=postgres/' \
            -e 's/^#db-username=/db-username=/' \
            -e 's/^#db-password=password/db-password=keycloak/' \
            -e 's/^#db-url=/db-url=/' /opt/keycloak/conf/keycloak.conf


# Give feedback after successful completion
echo "Succesfully adjusted keycloaks configuration."

# create keycloak admin-user
export KEYCLOAK_ADMIN="keycloak"

export KEYCLOAK_ADMIN_PASSWORD="keycloak"

# TODO: what if keycloak is running, but does not have an admin user yet?

if curl --head --silent http://localhost:8080/ | grep -F -q "Location: http://localhost:8080/admin/"; then
  echo "keycloak is already running..."
else
  sudo --preserve-env=KEYCLOAK_ADMIN,KEYCLOAK_ADMIN_PASSWORD /opt/keycloak/bin/kc.sh start-dev &

  # wait for keycloak to start
  echo "Waiting for keycloak to start..."
  until curl --head --silent http://localhost:8080/ | grep -F -q "Location: http://localhost:8080/admin/"
  do
    sleep 1
  done
fi

adminuser=keycloak
adminpass=keycloak

# log into the master realm with admin rights, token saved in ~/.keycloak/kcadm.config
sudo /opt/keycloak/bin/kcadm.sh config credentials --server http://localhost:8080 --realm master --user "$adminuser" --password "$adminpass"

if sudo /opt/keycloak/bin/kcadm.sh get 'http://localhost:8080/admin/realms' | grep -F -q '"realm" : "isduba",' ; then
 echo "Realm isduba already exists."
else
./keycloak/createRealm.sh
fi

./keycloak/createRole.sh 'bearbeiter' 'Bearbeiter' '' '\"WHITE\", \"GREEN\"'

./keycloak/createUser.sh 'beate' 'beate' 'bear' 'bea@ISDuBA.isduba' 'beate' 'bearbeiter'
