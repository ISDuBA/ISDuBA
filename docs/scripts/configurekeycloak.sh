#!/usr/bin/env bash
set -e

DBURLB=#db-url=jdbc:postgresql://localhost/keycloak

DBURLA=db-url=jdbc:postgresql://localhost/keycloak

sed -i s,#db=postgres,db=postgres,g /opt/keycloak/conf/keycloak.conf

sed -i s,#db-username=keycloak,db-username=keycloak,g /opt/keycloak/conf/keycloak.conf

sed -i s,#db-password=password,db-password=keycloak,g /opt/keycloak/conf/keycloak.conf

sed -i s,$DBURLB,$DBURLA,g /opt/keycloak/conf/keycloak.conf

sed -i s,#hostname=myhostname,#hostname=isduba,g /opt/keycloak/conf/keycloak.conf

echo "Succesfully adjusted keycloaks configuration."
