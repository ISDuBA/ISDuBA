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
# $1: username
# $2: group name
# $3: role name

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

# Get the user ID from name
USRID=$(sudo /opt/keycloak/bin/kcadm.sh get "http://localhost:8080/admin/realms/isduba/users?search=$1")
IDU=(${USRID//,/ })
declare -i COUNTERU=0
declare -i RESULTU=-1
for i in "${IDU[@]}"
do
  if [ "$i" = "\"$1\"" ]; then
    if [[ "${IDU[$COUNTERU-2]}" = "\"username\"" ]]; then
      RESULTU=$COUNTERU-3
    fi
  fi
COUNTERU=$COUNTERU+1
done

if [ $RESULTU = -1 ]; then
  echo "Couldn't find user $1"
  IDUSR=""
else
  IDUSR=${IDU[$RESULTU]:1:-1}
fi

# Get the group ID from name
GRPID=$(sudo /opt/keycloak/bin/kcadm.sh get "http://localhost:8080/admin/realms/isduba/groups?search=$2")
IDG=(${GRPID//,/ })
declare -i COUNTERG=0
declare -i RESULTG=-1
for j in "${IDG[@]}"
do
  if [ "$j" = \"$2\" ]; then
    if [[ "${IDG[$COUNTERG-2]}" = "\"name\"" ]]; then
      RESULTG=$COUNTERG-3
    fi
  fi
COUNTERG=$COUNTERG+1
done

if [ $RESULTG = -1 ]; then
  echo "Couldn't find group $2"
  IDGRP=""
else
  IDGRP=${IDG[$RESULTG]:1:-1}
fi

if [[ -z "$IDUSR" || -z "IDGRP" ]]; then
    echo "Failed to add user $1 to group $2"
else
  sudo /opt/keycloak/bin/kcadm.sh update users/$IDUSR/groups/$IDGRP -r isduba -s realm=isduba -s userId=$IDUSR -s groupId=$IDGRP -n
  echo "Added user $1 to group $2"
fi

sudo /opt/keycloak/bin/kcadm.sh add-roles -r isduba --uusername $1 --rolename $3
