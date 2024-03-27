# This file is Free Software under the Apache-2.0 License
# without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
#
# SPDX-License-Identifier: Apache-2.0
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

#!/bin/bash

PATH=/opt/keycloak/bin:$PATH

adminuser=${KEYCLOAK_ADMIN}
adminpass=${KEYCLOAK_ADMIN_PASSWORD}
client_hostname_url=${CLIENT_HOSTNAME_URL}

# log into the master realm with admin rights, token saved in ~/.keycloak/kcadm.config
kcadm.sh config credentials --server http://isduba-keycloak:8080 --realm master --user "$adminuser" --password "$adminpass"

kcadm.sh create realms --set realm=isduba --set enabled=true --output
# same output as
kcadm.sh get realms/isduba

id=$(kcadm.sh create clients --target-realm=isduba --set clientId=auth --set enabled=true --id)
echo $id
kcadm.sh get realms/isduba/clients/$id

kcadm.sh update clients/$id --target-realm=isduba \
	--set rootUrl=$client_hostname_url \
	--set "redirectUris=[\"$client_hostname_url/*\"]" \
	--set 'attributes={
    "oidc.ciba.grant.enabled" : "false",
    "post.logout.redirect.uris" : "+",
    "oauth2.device.authorization.grant.enabled" : "false",
    "backchannel.logout.session.required" : "true",
    "backchannel.logout.revoke.offline.tokens" : "false" }' \
	--set 'webOrigins=["*"]' \
	--set "adminUrl=$client_hostname_url/" \
	--set publicClient=true \
	--set standardFlowEnabled=true \
	--set directAccessGrantsEnabled=true \
	--set consentRequired=false

kcadm.sh update clients/$id --target-realm=isduba \
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
      "jsonType.label" : "JSON"
    } } ]'

kcadm.sh create roles --target-realm=isduba --set name=bearbeiter \
	--set "description=Bearbeiter"
kcadm.sh update roles/bearbeiter --target-realm isduba \
	--set 'attributes={
    "TLP" : [ "[{\"publisher\":\"\", \"tlps\":[\"WHITE\", \"GREEN\"]}]" ]
  }'

#creating  a user
userid=$(kcadm.sh create users --target-realm isduba \
	--set username=alex --set enabled=true \
	--set firstName=Alex --set lastName=Klein \
	--set email=test@example.org \
	--set emailVerified=true)

password=GENPASSWORD
kcadm.sh set-password --target-realm isduba \
	--username alex --new-password "$password"

kcadm.sh add-roles -r isduba --uusername alex --rolename bearbeiter
