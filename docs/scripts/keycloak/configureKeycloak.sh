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
file="./../password.txt"
live="9000"

while [[ $# -gt 0 ]]; do
  case "$1" in
    -k|--keycloakRunning)
      echo "Assuming keycloak is running..."
      keycloak_running=true
      ;;
    -f|--file)
      if [[ -n "$2" ]]; then
        file="$2"
        shift
      else
        echo "Error: Options -f and --file require an argument."
        exit 1
      fi
      ;;
    -l|--live)
      if [[ -n "$2" ]]; then
        live="$2"
        shift
      else
        echo "Error: Options -l and --live require an argument."
        exit 1
      fi
      ;;
    -p|--password)
      if [[ -n "$2" ]]; then
        export KEYCLOAK_ADMIN_PASSWORD="$2"
        shift
      else
        echo "Error: Options -p and --pw require an argument."
        exit 1
      fi
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
  echo "No Keycloak admin set. Trying to create admin with name \"keycloak\"."
  echo "Note that if a keycloak user has been previously set to something other than \"keycloak\", you need to change the environment variable \"KEYCLOAK_ADMIN\" for this script to work."
else
  export KEYCLOAK_ADMIN="${KEYCLOAK_ADMIN}"
fi

if [[ -z "${KEYCLOAK_ADMIN_PASSWORD}" ]]; then
  if [ -f $file ]; then
    export KEYCLOAK_ADMIN_PASSWORD=$(<./../password.txt)
  else
    password=$(xkcdpass --min 8 -d .)
    echo $password > ./../password.txt
    export KEYCLOAK_ADMIN_PASSWORD="$password"
    echo "No Keycloak admin password set. Trying to create new password. (See password.txt)"
    echo "Note that if a keycloak password has been previously set, this script won't work."
  fi
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
  if curl --silent --head -fsS http://localhost:$live/health/ready; then
    echo "keycloak is already running..."
  else
    sudo --preserve-env=KEYCLOAK_ADMIN,KEYCLOAK_ADMIN_PASSWORD /opt/keycloak/bin/kc.sh start-dev --health-enabled=true &

    # wait for keycloak to start
    echo "Waiting for keycloak to start..."
    until curl --silent --head -fsS http://localhost:$live/health/ready
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
./createUsers.sh  -f ./../../developer/users.txt  --noLogin

# end keycloak now that setup is done.
if [ ! -z "$nkid" ]; then
  if ps -p $nkid; then
    kill $nkid
  fi
fi
