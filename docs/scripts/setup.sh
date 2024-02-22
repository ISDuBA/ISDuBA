#!/usr/bin/env bash
set -e

./installgojava.sh

./installkeycloak.sh

./configurekeycloak.sh

./installpostgres.sh

./configurepostgres.sh

./keycloakonsystemstart.sh
