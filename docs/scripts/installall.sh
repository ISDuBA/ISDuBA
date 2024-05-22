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

# Help function if --help was called
help() {
echo "Usage: installall.sh [OPTIONS]"
echo "where OPTIONS:"
echo "  -h, --help                       show this help text and exit script (optional)"
echo "  -b, --branch=name                set up on branch 'name' instead of main (optional)"
echo "  -k, --keycloakRunning            signal the script that there is a keycloak running"
echo "                                   on port 8080 (optional)"

}

# update, install git and get the repository
prepare() {
sudo apt-get update

sudo apt-get install git -y # Needed to work with github repositories

# If repository already exists, update it instead of cloning it
if [ ! -d "ISDuBA" ] ; then
  git clone https://github.com/ISDuBA/ISDuBA.git
  cd ISDuBA/docs/scripts
else
  cd ISDuBA/docs/scripts
  git pull
fi
}

# check for branch and check it out if it exists
checkout() {
  BRANCH=$(git ls-remote --heads origin "refs/heads/$1" | wc -w) # 0 if branch does not exist
    if [ "$BRANCH" = "0"  ]; then
      echo "Could not find branch $1. Aborting..."
      exit 1
    else
      git checkout "$1"
    fi
}

# TODO: Check whether sudo is necessary where used.

while [[ $# -gt 0 ]]; do
  case "$1" in
    -h|--help)
      help
      exit 0
      ;;
    -k|--keycloakRunning)
      echo "Assuming keycloak is running..."
      keycloak_running=true
      ;;
    -b|--branch)
      if [[ -n "$2" ]]; then
        prepare
        checkout "$2"
        shift
      else
        echo "Error: Branch requires a value."
        help
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

if [ -z "$1" ]; then # if a an argument was given, prepare was already called or the script finished
  prepare
fi

if $keycloak_running; then
  ./setup.sh -k # Execute all other setup scripts, assuming keycloak is running
else
  ./setup.sh # Execute all the other setup scripts
fi
