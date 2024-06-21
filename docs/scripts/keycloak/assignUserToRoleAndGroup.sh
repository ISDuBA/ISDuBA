#!/usr/bin/env bash

# This file is Free Software under the Apache-2.0 License
# without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
#
# SPDX-License-Identifier: Apache-2.0
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

set -e # to exit if a command in the script fails

login=true

help() {
echo "Usage: assignUserToRoleAndGroup.sh name OPTIONS"
echo "where name:"
echo "  -n, --name=name                  username of the user to be added to roles or groups."
echo "where OPTIONS:"
echo "  -h, --help                       show this help text and exit script (optional)."
echo "  -g, --group=name                 name of the group the user should be added to (optional)."
echo "  -r, --role=name                  name of the role the user should be added to (optional)."
echo "      --noLogin                    do not attempt to log into keycloak. Requires active login to not cause errors (optional)."
}


while [[ $# -gt 0 ]]; do
  case "$1" in
    -h|--help)
      help
      exit 0
      ;;
    -n|--name)
      if [[ -n "$2" ]]; then
        name="$2"
        shift
      else
        echo "Error: No name given."
        exit 1
      fi
      ;;
    -g|--group)
      if [[ -n "$2" ]]; then
        group="$2"
        shift
      else
        echo "Error: No group given."
        exit 1
      fi
      ;;
    -r|--role)
      if [[ -n "$2" ]]; then
        role="$2"
        shift
      else
        echo "Error: No role given."
        exit 1
      fi
      ;;
    --noLogin)
      login=false
      ;;
    *)
      echo "Unknown option: $1"
      help
      exit 1
      ;;
  esac
  shift
done

if "$login"; then
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
fi

# Get the user ID from name
USRID=$(sudo /opt/keycloak/bin/kcadm.sh get "http://localhost:8080/admin/realms/isduba/users?search=$name")
IDU=(${USRID//,/ })
declare -i COUNTERU=0
declare -i RESULTU=-1
for i in "${IDU[@]}"
do
  if [ "$i" = "\"$name\"" ]; then
    if [[ "${IDU[$COUNTERU-2]}" = "\"username\"" ]]; then
      RESULTU=$COUNTERU-3
    fi
  fi
COUNTERU=$COUNTERU+1
done

if [ $RESULTU = -1 ]; then
  echo "Couldn't find user $name"
  exit 1
else
  IDUSR=${IDU[$RESULTU]:1:-1}
fi


if [ ! -z "$group" ]; then
  # Get the group ID from name
  GRPID=$(sudo /opt/keycloak/bin/kcadm.sh get "http://localhost:8080/admin/realms/isduba/groups?search=$group")
  IDG=(${GRPID//,/ })
  declare -i COUNTERG=0
  declare -i RESULTG=-1
  for j in "${IDG[@]}"
  do
    if [ "$j" = \"$group\" ]; then
      if [[ "${IDG[$COUNTERG-2]}" = "\"name\"" ]]; then
        RESULTG=$COUNTERG-3
      fi
    fi
  COUNTERG=$COUNTERG+1
  done

  if [ $RESULTG = -1 ]; then
    echo "Couldn't find group $group"
    IDGRP=""
  else
    IDGRP=${IDG[$RESULTG]:1:-1}
  fi

  if [[ -z "$IDUSR" || -z "$IDGRP" ]]; then
      echo "Failed to add user $name to group $group"
  else
    sudo /opt/keycloak/bin/kcadm.sh update users/$IDUSR/groups/$IDGRP -r isduba -s realm=isduba -s userId=$IDUSR -s groupId=$IDGRP -n
  fi
fi

if [ ! -z "$role" ]; then
  sudo /opt/keycloak/bin/kcadm.sh add-roles -r isduba --uusername $name --rolename $role
fi
