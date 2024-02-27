#!/usr/bin/env bash

# This file is Free Software under the MIT License
# without warranty, see README.md and LICENSES/MIT.txt for details.
#
# SPDX-License-Identifier: MIT
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

set -e # to exit if a command in the script fails

# install PostgreSQL 16
apt install vim gnupg2 -y
curl -fsSL https://www.postgresql.org/media/keys/ACCC4CF8.asc| gpg --trust-model always --dearmor -o /etc/apt/trusted.gpg.d/postgresql.gpg
sh -c 'echo "deb https://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
apt update
apt install postgresql-16 -y
