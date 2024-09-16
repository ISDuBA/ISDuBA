#!/usr/bin/env bash

# This file is Free Software under the Apache-2.0 License
# without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
#
# SPDX-License-Identifier: Apache-2.0
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

set -e # to exit if a command in the script fails

help() {
echo "Usage: createUsers.sh OPTIONS"
echo "where OPTIONS:"
echo "  -h, --help                       show this help text and exit script (optional)."
echo "  -f, --file=file                  name of the file that contains all user information (mandatory)."
echo "      --noLogin                    do not attempt to log into keycloak. Requires active login to not cause errors (optional)."
}


# If 1 or 0 arguments are given, then there
# cannot be the mandatory file option
if [[ $# -lt 2 ]]; then
  help
  exit 1
fi

login=true
while [[ $# -gt 0 ]]; do
  case "$1" in
    -h|--help)
      help
      exit 0
      ;;
    -f|--file)
      if [[ -n "$2" ]]; then
        file="$2"
        shift
      else
        echo "Error: No file given."
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

if $login; then
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

while read line; do
user=($line)
if [ "${#user[@]}" -eq 7 ]; then
  echo "creating user ${user[0]}..."
  ./createUser.sh  "${user[0]}" "${user[1]}" "${user[2]}" "${user[3]}" "${user[4]}"  false
  ./assignUserToRoleAndGroup.sh --name "${user[0]}" --group "${user[5]}" --role "${user[6]}" --noLogin
fi
done < $file

# test-user with all privileges etc.
./createUser.sh "test-user" "test-FirstName" "test-LastName" "test-Email@example.com" "test-user" false
./assignUserToRoleAndGroup.sh --name "test-user" --noLogin --role "editor" --group "all"
./assignUserToRoleAndGroup.sh --name "test-user" --noLogin --role "reviewer"
./assignUserToRoleAndGroup.sh --name "test-user" --noLogin --role "admin"
./assignUserToRoleAndGroup.sh --name "test-user" --noLogin --role "auditor"
./assignUserToRoleAndGroup.sh --name "test-user" --noLogin --role "source-manager"
./assignUserToRoleAndGroup.sh --name "test-user" --noLogin --role "importer"



