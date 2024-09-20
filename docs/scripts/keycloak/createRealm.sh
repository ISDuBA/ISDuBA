#!/usr/bin/env bash

# This file is Free Software under the Apache-2.0 License
# without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
#
# SPDX-License-Identifier: Apache-2.0
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

set -e # to exit if a command in the script fails

sudo /opt/keycloak/bin/kcadm.sh create realms --set realm=isduba --set enabled=true --output
# same output as
sudo /opt/keycloak/bin/kcadm.sh get realms/isduba

id=$(sudo /opt/keycloak/bin/kcadm.sh create clients --target-realm=isduba --set clientId=auth --set enabled=true --id)

sudo /opt/keycloak/bin/kcadm.sh get realms/isduba/clients/$id

sudo /opt/keycloak/bin/kcadm.sh update clients/$id --target-realm=isduba \
  --set rootUrl=http://localhost:5173/ \
  --set 'redirectUris=["*"]' \
  --set 'attributes={
    "oidc.ciba.grant.enabled" : "false",
    "post.logout.redirect.uris" : "+",
    "oauth2.device.authorization.grant.enabled" : "false",
    "backchannel.logout.session.required" : "true",
    "backchannel.logout.revoke.offline.tokens" : "false" }' \
  --set 'webOrigins=["*"]' \
  --set 'adminUrl=http://localhost:5173/' \
  --set publicClient=true \
  --set standardFlowEnabled=true \
  --set directAccessGrantsEnabled=true \
  --set consentRequired=false

sudo /opt/keycloak/bin/kcadm.sh update clients/$id --target-realm=isduba \
  --set 'protocolMappers=[ {
    "name" : "TLP",
    "protocol" : "openid-connect",
    "protocolMapper" : "oidc-usermodel-attribute-mapper",
    "consentRequired" : false,
    "config" : {
      "aggregate.attrs" : "true",
      "introspection.token.claim" : "true",
      "userinfo.token.claim" : "true",
      "user.attribute" : "TLP",
      "id.token.claim" : "true",
      "lightweight.claim" : "true",
      "access.token.claim" : "true",
      "claim.name" : "TLP",
      "jsonType.label" : "JSON",
      "multivalued" : "true"
    } } ]'
