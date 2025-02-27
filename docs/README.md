<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2025 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2025 Intevation GmbH <https://intevation.de>
-->

# Get Started

Learn how to test or setup your own ISDuBA instance.

 1. Want to try ISDuBA for yourself? [Try our Docker setup](#docker-setup)
 2. Want to support the ISDuBA project with your own code? [Here's how to setup a development instance of ISDuBA](#development-setup)
 3. Want to use ISDuBA for yourself or your organization? [Here's how to setup ISDuBA for production](#production-setup)

See [the production guide](https://github.com/ISDuBA/ISDuBA/blob/main/docs/production_guide.md) for security and maintenance considerations.

### Docker-setup

**The docker setup is not meant to be used in production.**

Use [docker compose](https://docs.docker.com/compose/install/) to build and start the whole stack with

```bash
cd docker
#### Setting BUILD_VERSION is optional
docker compose build --build-arg BUILD_VERSION=$(git describe --tags --always)
docker compose up -d
```

The default configuration is inside `docker/.env`.

To set the hostname of Keycloak and the client change the respective environment variables, for example:

```bash
KC_HOSTNAME=keycloak-host CLIENT_HOST=client-host docker compose up -d
```

##### Keycloak

The Keycloak admin interface can be reached under <http://localhost:8080>.
By default, an admin user is created during setup:

- Username: admin
- Password: secret

Inside `docker/keycloak/init.sh` is an automated configuration script that configures the realm and creates a test user.

The password and username can be obtained with:

```bash
docker logs isduba-keycloak-setup | grep "Created user"
```

##### Client application

The application can be reached under <http://localhost:5371>.


### Development-setup

The setup should be performed via the [installation scripts.](./scripts/README.md) on a Ubuntu 24.04 OS.

An example-configuration for `isdubad` can be found in [example_isdubad.toml](./example_isdubad.toml). Please edit to your needs.


##### Upgrading
When upgrading from an older version, a migration is needed to 
configure the database by starting isdubad with the 
`ISDUBA_DB_MIGRATE` environment variable set to true or
by adjusting the toml-configuration file, e.g.

<!-- MARKDOWN-AUTO-DOCS:START (CODE:src=../docs/scripts/setup.sh&lines=53-53) -->
<!-- The below code snippet is automatically added from ../docs/scripts/setup.sh -->
```sh
ISDUBA_DB_MIGRATE=true ./cmd/isdubad/isdubad -c ./isduba.toml
```
<!-- MARKDOWN-AUTO-DOCS:END -->

##### Additional tasks

Groups and users can be managed directly via Keycloak or the scripts:

 - Create additional users via [the createUsers script.](./scripts/keycloak/createUsers.sh)
  - A list of users created by the setup scripts can be found in [the users.txt.](./developer/users.txt) Editing this file before using the createUsers or the setup script will change which users are created.

 - Create groups via [the createGroup script.](./scripts/keycloak/createGroup.sh)
  - The restrictions set with the script are explained in [keycloak_values.md](./keycloak_values.md)

The Keycloak admin user created via the scripts will
have the username password `keycloak`,
unless otherwise specified via the environment variable `KEYCLOAK_ADMIN`.
The password can be specified via the environment variable 
`KEYCLOAK_ADMIN_PASSWORD`, a file (`-f` option)
or directly (using the `-p` option).

If neither is set, then the script will try to see if
`docs/scripts/password.txt` exists and contains a password.
If this is not set either, then a random password will be generated
and stored in `docs/scripts/password.txt`.

##### Run the application in a dev environment

To start the frontend via a `vite` dev-server:

```bash
cd client
npm run dev
```

This will start the client application and
print the URL a browser could be pointed to.

With a previously created configuration file (named e.g. `isduba.toml`) you could start the backend from the main directory:

```bash
./cmd/isdubad/isdubad -c isduba.toml
```

Make sure to have Keycloak running when trying to access the application.

(If set up via the script available under:)
```bash
sudo -u keycloak /opt/keycloak/bin/kc.sh start-dev
```

##### Notice when using versions of Keycloak other than a default installation of Keycloak 25

The setup scripts utilize Keycloak's health checks to determine whether Keycloak is running. The port to use may change depending on your Keycloak version or admin's configuration.
The default for the current version of 25 is port 9000.
This means it may be necessary to call docs/scripts/keycloak/configurekeycloak.sh with the -l/--live flag to manually set a port, e.g. for Keycloak 24:

```bash
  ./configureKeycloak.sh --live 8080
```

Not setting the correct port without the -k/--keycloakRunning option will cause the script to try and call the wrong port over and over until stopped.


### Production-setup

You can download the latest stable release from [github.](https://github.com/ISDuBA/ISDuBA/releases/)

##### Building

Alternatively, to build the application the latest Golang version, NodeJS 20 and standard build
tools, like GNU Make are required. At the root of the repo run `make dist` to
build the frontend and backend; this will result in a Tar-file inside `dist/`
that contains the application. It can also be useful to look inside the
Dockerfile of the application to see how individual components of the
application can be built. If there are no special requirements it can be
enough to use the already built tar-file from the release page.

##### Running

The Tar-file can be copied and extracted on a production server. This file
contains the `isdubad` backend, which can be run on any modern amd64 linux
system and the frontend which is contained in the `web/` folder. No further
dependencies are required to start the application. By default the backend will
serve the contents of `web/`.

For a quick start copy `example_isdubad.toml` to the
folder where the application is contained. Configure the postgres and Keycloak
settings and rename the file to `isduba.toml`. The application can now be
started. For exposing the application to the network it is recommended to use
a TLS-terminating reverse proxy.
