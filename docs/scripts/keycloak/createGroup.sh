#!/usr/bin/env bash

# This file is Free Software under the Apache-2.0 License
# without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
#
# SPDX-License-Identifier: Apache-2.0
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

set -e # to exit if a command in the script fails

# Simplified tlp chooser

help() {
echo "Usage: createGroup.sh OPTIONS"
echo "where OPTIONS:"
echo "  -h, --help                       show this help text and exit script (optional)."
echo "  -n, --name=name                  name of the group that is supposed to be created (mandatory)."
echo "  -w, --white                      grant the group access to TLP:WHITE advisories (optional)."
echo "  -g, --green                      grant the group access to TLP:GREEN advisories (optional)."
echo "  -a, --amber                      grant the group access to TLP:AMBER advisories (optional)."
echo "  -r, --red                        grant the group access to TLP:RED advisories (optional)."
echo "  -p, --publisher=name             restrict access to advisories of the named publisher (optional)."
echo "      --noLogin                    do not attempt to log into keycloak. Requires active login to not cause errors (optional)."
}

publisher="*"
tlpw=false
tlpg=false
tlpa=false
tlpr=false
comma=""
login=true

# Check which tlp levels should be added
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
    -w|--white)
      tlpw=true
      ;;
    -g|--green)
      tlpg=true
      ;;
    -a|--amber)
      tlpa=true
      ;;
    -r|--red)
      tlpr=true
      ;;
    --noLogin)
      login=false
      ;;
    -p|--publisher)
      if [[ -n "$2" ]]; then
        publisher="$2"
        shift
      else
        echo "Error: Empty publisher name given."
        exit 1
      fi
      ;;
    *)
      echo "Unknown option: $1"
      help
      exit 1
      ;;
  esac
  shift
done


# Add the TLP levels
tlp="]}\" ]"
if $tlpr; then
  tlp="\\\"RED\\\"$tlp"
  comma=", "
fi
if $tlpa; then
  tlp="\\\"AMBER\\\"$comma$tlp"
  comma=", "
fi
if $tlpg; then
  tlp="\\\"GREEN\\\"$comma$tlp"
  comma=", "
fi
if $tlpw; then
  tlp="\\\"WHITE\\\"$comma$tlp"
fi

# Add beginning part
tlp="\"TLP\" : [ \"{\\\"$publisher\\\": [$tlp"

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

if sudo /opt/keycloak/bin/kcadm.sh get 'http://localhost:8080/admin/realms/isduba/groups' | grep -F -q "\"name\" : \"$name\"", ; then
  echo "Group $name already exists."
else
  # create group
  sudo /opt/keycloak/bin/kcadm.sh create groups -r isduba -s name="$name"
fi

# Get the ID from name
GRP=$(sudo /opt/keycloak/bin/kcadm.sh get "http://localhost:8080/admin/realms/isduba/groups?search=$name")
IDS=(${GRP//,/ })
declare -i COUNTER=0
declare -i RESULT=-1
for i in "${IDS[@]}"
do
  if [ "$i" = "\"$name\"" ]; then
    if [[ "${IDS[$COUNTER-2]}" = "\"name\"" ]]; then
      RESULT=$COUNTER-3
    fi
  fi
COUNTER=$COUNTER+1
done

ID=${IDS[$RESULT]:1:-1}

if [ "$tlp" != "" ]; then
  if sudo /opt/keycloak/bin/kcadm.sh get 'http://localhost:8080/admin/realms/isduba/groups' | grep -F -q "\"name\" : \"$name\"", ; then
    sudo /opt/keycloak/bin/kcadm.sh update groups/$ID --target-realm isduba \
    --set "attributes={
      $tlp
    }"
  else
    echo "Failed to create group $name."
  fi
fi
