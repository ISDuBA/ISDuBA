<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

# Docker Setup
Use [docker compose](https://docs.docker.com/compose/install/) to build and start the whole stack with

```bash
cd docker
# Setting BUILD_VERSION is optional 
docker compose build --build-arg BUILD_VERSION=$(git describe --tags --always)
docker compose up -d
```

The default configuration is inside `docker/.env`.

To set the hostname of keycloak and the client change the respective environment variables, for example:
```bash
KC_HOSTNAME=keycloak-host CLIENT_HOST=client-host docker compose up -d
```

## Keycloak

The keycloak admin interface can be reached under http://localhost:8080.
By default, an admin user is created during setup:

* Username: admin
* Password: secret

Inside `docker/keycloak/init.sh` is an automated configuration script that configures the realm and creates a test user.

The password and username can be obtained with:
```bash
docker logs isduba-keycloak-setup | grep "Created user"
```

## Client application

The application can be reached under http://localhost:5371.

## Production usage considerations
### Keycloak
This setup starts Keycloak in a development mode and does not persist configuration across rebuilds. It is preferred to use a Keycloak instance that is already managed. The Keycloak setup container executes the script located in `docker/keycloak/init.sh`. This script is used for automated testing; it is a template on how to set up Keycloak and should not be used to automate the production setup.
### Database
It is recommended to run the database outside the container or mount `/var/lib/postgresql/data` to a persistent store. If the database is accessed through a network, the default passwords should be changed.
### Application
The application does not provide a secure endpoint. It should be put behind a reverse proxy with a valid TLS certificate.
