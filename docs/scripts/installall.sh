#!/usr/bin/env bash

# This file is Free Software under the Apache-2.0 License
# without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
#
# SPDX-License-Identifier: Apache-2.0
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

set -e # to exit if a command in the script fails

cd "$HOME"

# TODO: Check whether sudo is necessary where used.

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
if [ ! -z "$1" ]; then # if a branch was given as a parameter
  BRANCH=$(git ls-remote --heads origin "refs/heads/$1" | wc -w) # 0 if branch does not exist
  if [ "$BRANCH" = "0"  ]; then
    echo "Could not find branch $1. Aborting..."
  else
    git checkout "$1"
    . setup.sh # Execute all the other setup scripts
  fi
else
  . setup.sh # Execute all the other setup scripts
fi

