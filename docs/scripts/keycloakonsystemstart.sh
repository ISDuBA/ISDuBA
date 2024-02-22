#!/usr/bin/env bash

# This file is Free Software under the MIT License
# without warranty, see README.md and LICENSES/MIT.txt for details.
#
# SPDX-License-Identifier: MIT
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

set -e # to exit if a command in the script fails

echo "[Unit]
Description=Keycloak
After=network.target

[Service]
Type=idle
User=keycloak
Group=keycloak
ExecStart=/opt/keycloak/bin/kc.sh start-dev
TimeoutStartSec=600
TimeoutStopSec=600

[Install]
WantedBy=multi-user.target
" > /etc/systemd/system/keycloak.service

systemctl enable keycloak
systemctl start keycloak
