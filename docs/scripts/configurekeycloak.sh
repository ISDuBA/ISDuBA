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
DBN=db=postgres

DBU=db-username=keycloak 

DBPB=#db-password=password
DBPA=db-password=keycloak

DBURL=db-url=jdbc:postgresql://localhost/keycloak

sed -i s,#$DBN,$DBN,g /opt/keycloak/conf/keycloak.conf     # remove leading '#'

sed -i s,#$DBU,$DBU,g /opt/keycloak/conf/keycloak.conf     # remove leading '#'

sed -i s,$DBPB,$DBPA,g /opt/keycloak/conf/keycloak.conf    # exchange password

sed -i s,#$DBURL,$DBURL,g /opt/keycloak/conf/keycloak.conf # remove leading '#'


# Give feedback after successful completion
echo "Succesfully adjusted keycloaks configuration."
