#!/usr/bin/env bash

# This file is Free Software under the Apache-2.0 License
# without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
#
# SPDX-License-Identifier: Apache-2.0
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

set -e # to exit if a command in the script fails

# Alter the keycloak configuration
sudo sed --in-place=.orig -e 's/^#db=postgres/db=postgres/' \
            -e 's/^#db-username=/db-username=/' \
            -e 's/^#db-password=password/db-password=keycloak/' \
            -e 's/^#db-url=/db-url=/' /opt/keycloak/conf/keycloak.conf


# Give feedback after successful completion
echo "Succesfully adjusted keycloaks configuration."

role='bearbeiter'

publisher=''

tlp='\"WHITE\", \"GREEN\"'

echo $A

./keycloak/login.sh

./keycloak/createRealm

./keycloak/createRole.sh $role 'Bearbeiter' $publisher $tlp | echo

./keycloak/createUser.sh 'beate' 'beate' 'bear' 'bea@ISDuBA.isduba' 'beate' $role | echo
