#!/usr/bin/env bash

# This file is Free Software under the MIT License
# without warranty, see README.md and LICENSES/MIT.txt for details.
#
# SPDX-License-Identifier: Apache-2.0
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

set -e # to exit if a command in the script fails

sudo apt-get update

# Install Java
sudo apt install -y openjdk-17-jre-headless

# Install Go

latest_go="$(curl https://go.dev/VERSION\?m=text| head -1).linux-amd64.tar.gz"
wget -O /tmp/$latest_go https://dl.google.com/go/$latest_go
sudo rm -rf /usr/local/go # be sure that we do not have an old installation
sudo tar -C /usr/local -xzf /tmp/$latest_go

sudo rm -f /tmp/$latest_go

sudo ln -snf /usr/local/go/bin/go /usr/local/bin/go
