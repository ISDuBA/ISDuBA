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
# $2: Attribute
# $3: publisher (optional)


# Simplified tlp chooser
tlp=""
publisher="*"
if ! [[ -z "$3" ]]; then
    publisher="$3"
fi
if [ "$2" = "1" ]; then
    tlp="\"TLP\" : [ \"{\\\"$publisher\\\": [\\\"WHITE\\\"]}\" ]"
fi
if [ "$2" = "2" ]; then
    tlp="\"TLP\" : [ \"{\\\"$publisher\\\": [\\\"GREEN\\\"]}\" ]"
fi
if [ "$2" = "3" ]; then
    tlp="\"TLP\" : [ \"{\\\"$publisher\\\": [\\\"WHITE\\\", \\\"GREEN\\\"]}\" ]"
fi

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

if sudo /opt/keycloak/bin/kcadm.sh get 'http://localhost:8080/admin/realms/isduba/groups' | grep -F -q "\"name\" : \"$1\"", ; then
  echo "Group $1 already exists."
else
  # create group
  sudo /opt/keycloak/bin/kcadm.sh create groups -r isduba -s name="$1"
fi

# Get the ID from name
GRP=$(sudo /opt/keycloak/bin/kcadm.sh get "http://localhost:8080/admin/realms/isduba/groups?search=$1")
IDS=(${GRP//,/ })
declare -i COUNTER=0
declare -i RESULT=-1
for i in "${IDS[@]}"
do
  if [ "$i" = "\"$1\"" ]; then
    if [[ "${IDS[$COUNTER-2]}" = "\"name\"" ]]; then
      RESULT=$COUNTER-3
    fi
  fi
COUNTER=$COUNTER+1
done

ID=${IDS[$RESULT]:1:-1}

if [ "$tlp" != "" ]; then
  if sudo /opt/keycloak/bin/kcadm.sh get 'http://localhost:8080/admin/realms/isduba/groups' | grep -F -q "\"name\" : \"$1\"", ; then
    sudo /opt/keycloak/bin/kcadm.sh update groups/$ID --target-realm isduba \
    --set "attributes={
      $tlp
    }"
  else
    echo "Failed to create group $1."
  fi
fi
