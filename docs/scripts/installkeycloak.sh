#!/usr/bin/env bash
set -e

apt install -y unzip

wget https://github.com/keycloak/keycloak/releases/download/23.0.5/keycloak-23.0.5.zip

unzip keycloak-23.0.5.zip

mkdir -p /opt/

mv keycloak-23.0.5 /opt/keycloak

useradd keycloak
chown -R keycloak: /opt/keycloak
