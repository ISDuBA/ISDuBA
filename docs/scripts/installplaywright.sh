#!/usr/bin/env bash

# This file is Free Software under the MIT License
# without warranty, see README.md and LICENSES/MIT.txt for details.
#
# SPDX-License-Identifier: MIT
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

set -e # to exit if a command in the script fails

curl -sL https://deb.nodesource.com/setup_20.x | bash - # install node

apt-install npm # install npm

apt-get install nodejs -y

cd ../../client # change directory

# install dependencies
npm install 

npx playwright install
