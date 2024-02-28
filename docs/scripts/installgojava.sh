#!/usr/bin/env bash

# This file is Free Software under the MIT License
# without warranty, see README.md and LICENSES/MIT.txt for details.
#
# SPDX-License-Identifier: MIT
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

set -e # to exit if a command in the script fails

apt-get update

# Install Java
apt install -y openjdk-17-jre-headless

# Install Go
wget -O /tmp/go1.22.0.linux-amd64.tar.gz https://go.dev/dl/go1.22.0.linux-amd64.tar.gz 

cd /usr/local/

rm -rf go  && tar -xzf /tmp/go1.22.0.linux-amd64.tar.gz

rm -f /tmp/go1.22.0.linux-amd64.tar.gz

ln -snf /usr/local/go/bin/go /usr/local/bin/go
