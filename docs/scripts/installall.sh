#!/usr/bin/env bash

# This file is Free Software under the MIT License
# without warranty, see README.md and LICENSES/MIT.txt for details.
#
# SPDX-License-Identifier: MIT
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

set -e # to exit if a command in the script fails

apt-get install git -y # Needed to work with github repositories

git clone https://github.com/ISDuBA/ISDuBA.git # Clone ISDuBA repository

cd ISDuBA/docs/scripts # Change working directory to scripts

./setup.sh # Execute all the other scripts
