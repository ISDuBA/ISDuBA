#!/usr/bin/env bash

# This file is Free Software under the Apache-2.0 License
# without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
#
# SPDX-License-Identifier: Apache-2.0
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

set -e # to exit if a command in the script fails

sudo apt-get update

# Install base packages
sudo apt install -y make bash curl sed tar

# Install Java
sudo apt install -y openjdk-21-jre-headless

# Install or update Go

# look up current go version
go_version="$(curl https://go.dev/VERSION\?m=text | head -1)"

if [[ -z "${KEYCLOAK_ADMIN_PASSWORD}" ]]; then
  sudo apt install xkcdpass
fi

# if go exists and is the newest version
if [ -x "$(command -v go version)" ] && [[ $(go version) == *"$go_version"* ]]; then
  echo "Newest go version already installed."
# if not, download the newest go version
else
  latest_go=$go_version".linux-amd64.tar.gz"

  curl -o /tmp/$latest_go https://dl.google.com/go/$latest_go
  sudo rm -rf /usr/local/go # remove any old installations
  sudo tar -C /usr/local -xzf /tmp/$latest_go

  sudo rm -f /tmp/$latest_go

  sudo ln -snf /usr/local/go/bin/go /usr/local/bin/go
  echo "Successfully installed $go_version."
fi

# install graphviz
sudo apt install -y graphviz
