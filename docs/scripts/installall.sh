#!/usr/bin/env bash

# This file is Free Software under the Apache-2.0 License
# without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
#
# SPDX-License-Identifier: Apache-2.0
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

set -e # to exit if a command in the script fails

sudo apt-get install git -y # Needed to work with github repositories

# If repository already exists, update it instead of cloning it
if [ ! -d "ISDuBA" ] ; then
  git clone https://github.com/ISDuBA/ISDuBA.git
  cd ISDuBA/docs/scripts
else
  cd ISDuBA/docs/scripts
  git pull
fi

./setup.sh # Execute all the other setup scripts
