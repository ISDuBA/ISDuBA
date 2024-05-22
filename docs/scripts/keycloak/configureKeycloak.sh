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

# create roles  name        description      login
./createRole.sh 'editor'    'Bearbeiter'     false

./createRole.sh 'reviewer'  'Reviewer'       false

./createRole.sh 'auditor'   'Auditor'        false

./createRole.sh 'importer'  'Importierer'    false

./createRole.sh 'admin'     'Administrator'  false

# create users  username    first name       last name       email address       password     login
./createUser.sh 'beate'     'Beate'          'Bear'          'bea@ISDuBA.isduba' 'beate'      false

./createUser.sh 'revy'      'Revy'           'Reviewer'      'rev@ISDuBA.isduba' 'revy'       false

./createUser.sh 'imke'      'Imke'           'Importer'      'imk@ISDuBA.isduba' 'imke'       false

./createUser.sh 'augustus'  'Augustus'       'Auditor'       'aug@ISDuBA.isduba' 'augustus'   false

./createUser.sh 'admin'     'Admin'          'Administrator' 'bea@ISDuBA.isduba' 'beate'      false

./createUser.sh 'adrian'    'Adrian'         'Karma'         'Adr@ISDuBA.isduba' 'adrian'     false

./createUser.sh 'alex'      'Alex'           'Klein'         'ale@ISDuBA.isduba' 'alex'       false

./createUser.sh 'max'       'Max'            'Maier'         'max@ISDuBA.isduba' 'max'        false

./createUser.sh 'white'     'White'          'Weiss'         'whi@ISDuBA.isduba' 'white'      false

./createUser.sh 'green'     'Green'          'Gruen'         'gru@ISDuBA.isduba' 'green'      false

./createUser.sh 'amber'     'Amber'          'Bernstein'     'amb@ISDuBA.isduba' 'amber'      false

./createUser.sh 'red'       'Red'            'Rot'           'red@ISDuBA.isduba' 'red'        false

./createUser.sh 'low'       'Low'            'Niedrig'       'low@ISDuBA.isduba' 'low'        false

./createUser.sh 'medium'    'Medium'         'Geisterseher'  'med@ISDuBA.isduba' 'medium'     false

./createUser.sh 'high'      'High'           'Hoch'          'hig@ISDuBA.isduba' 'high'       false

./createUser.sh 'all'       'All'            'Alle'          'all@ISDuBA.isduba' 'all'        false

# assign users to groups      username       group         role          login
./assignUserToRoleAndGroup.sh 'beate'        'all'         'editor'      false

./assignUserToRoleAndGroup.sh 'revy'         'all'         'reviewer'    false

./assignUserToRoleAndGroup.sh 'imke'         'all'         'importer'    false

./assignUserToRoleAndGroup.sh 'augustus'     'all'         'auditor'     false

./assignUserToRoleAndGroup.sh 'admin'        'all'         'admin'       false

./assignUserToRoleAndGroup.sh 'adrian'       'all'         'editor'      false

./assignUserToRoleAndGroup.sh 'alex'         'all'         'editor'      false

./assignUserToRoleAndGroup.sh 'max'          'all'         'editor'      false

./assignUserToRoleAndGroup.sh 'white'        'white'       'editor'      false

./assignUserToRoleAndGroup.sh 'green'        'green'       'editor'      false

./assignUserToRoleAndGroup.sh 'amber'        'amber'       'editor'      false

./assignUserToRoleAndGroup.sh 'red'          'red'         'editor'      false

./assignUserToRoleAndGroup.sh 'low'          'whitegreen'  'editor'      false

./assignUserToRoleAndGroup.sh 'medium'       'greenamber'  'editor'      false

./assignUserToRoleAndGroup.sh 'high'         'amberred'    'editor'      false

./assignUserToRoleAndGroup.sh 'all'          'all'         'editor'      false

if [ ! -z "$nkid" ]; then
  if ps -p $nkid; then
    kill $nkid
  fi
fi
