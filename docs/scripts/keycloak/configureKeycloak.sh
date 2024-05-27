#!/usr/bin/env bash

# This file is Free Software under the Apache-2.0 License
# without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
#
# SPDX-License-Identifier: Apache-2.0
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

set -e # to exit if a command in the script fails

keycloak_running=false

while [[ $# -gt 0 ]]; do
  case "$1" in
    -k|--keycloakRunning)
      echo "Assuming keycloak is running..."
      keycloak_running=true
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
  shift
done

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
if $keycloak_running; then
  echo "keycloak is assumed to be running"
else
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
fi

# log into the master realm with admin rights, token saved in ~/.keycloak/kcadm.config
sudo /opt/keycloak/bin/kcadm.sh config credentials --server http://localhost:8080 --realm master --user "$KEYCLOAK_ADMIN" --password "$KEYCLOAK_ADMIN_PASSWORD"

if sudo /opt/keycloak/bin/kcadm.sh get 'http://localhost:8080/admin/realms' | grep -F -q '"realm" : "isduba",' ; then
 echo "Realm isduba already exists."
else
./createRealm.sh
fi

# create groups  name                tlps        login
./createGroup.sh --name 'white'      -w          --noLogin

./createGroup.sh --name 'green'      -g          --noLogin

./createGroup.sh --name 'amber'      -a          --noLogin

./createGroup.sh --name 'red'        -r          --noLogin

./createGroup.sh --name 'whitegreen' -w -g       --noLogin

./createGroup.sh --name 'greenamber' -g -a       --noLogin

./createGroup.sh --name 'amberred'   -a -r       --noLogin

./createGroup.sh --name 'all'        -w -g -a -r --noLogin

./createGroup.sh --name 'none'                   --noLogin

# create roles  name        description               login
./createRole.sh 'editor'    'Bearbeiter'              false

./createRole.sh 'reviewer'  'Reviewer'                false

./createRole.sh 'auditor'   'Auditor'                 false

./createRole.sh 'importer'  'Importierer'             false

./createRole.sh 'admin'     'Administrator'           false

./createRole.sh 'none'      'Role outside the system' false

# create Users    file containing users            login
./createUsers.sh  -f ./../../developers/users.txt  --noLogin

# end keycloak now that setup is done.
if [ ! -z "$nkid" ]; then
  if ps -p $nkid; then
    kill $nkid
  fi
fi
