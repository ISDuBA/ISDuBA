#!/usr/bin/env bash

# This file is Free Software under the MIT License
# without warranty, see README.md and LICENSES/MIT.txt for details.
#
# SPDX-License-Identifier: MIT
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

set -e # to exit if a command in the script fails

# Alter the keycloak configuration
DBURLB=#db-url=jdbc:postgresql://localhost/keycloak

DBURLA=db-url=jdbc:postgresql://localhost/keycloak

sed -i s,#db=postgres,db=postgres,g /opt/keycloak/conf/keycloak.conf

sed -i s,#db-username=keycloak,db-username=keycloak,g /opt/keycloak/conf/keycloak.conf

sed -i s,#db-password=password,db-password=keycloak,g /opt/keycloak/conf/keycloak.conf

sed -i s,$DBURLB,$DBURLA,g /opt/keycloak/conf/keycloak.conf

sed -i s,#hostname=myhostname,#hostname=isduba,g /opt/keycloak/conf/keycloak.conf

# Give feedback after successful completion
echo "Succesfully adjusted keycloaks configuration."
