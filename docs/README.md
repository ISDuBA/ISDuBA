<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2025 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2025 Intevation GmbH <https://intevation.de>
-->

(See [the parent README](../README.md) for an overview over what ISDuBA is.)

# Get Started

Decide how to test or setup your own ISDuBA instance.

 1. Want to try ISDuBA for yourself? [Try our Docker setup](#docker-setup)
 2. Want to support the ISDuBA project with your own code? [Here's how to setup a development instance of ISDuBA](#development-setup)
 3. Want to use ISDuBA for yourself or your organization? [Here's how to setup ISDuBA for production](#production-setup)
 4. Having set up an instance of ISDuBA, you can read about what to do now within the [first steps guide](./first_steps.md)

When starting the application, you will be prompted to safe your aes_key. This can be ignored for test or development instances and is further explained in [the aes-keys section of the security considerations documentation](./security_considerations.md#aes-keys).

## Docker-setup

**The docker setup is not meant to be used in production.**

Use [docker compose](https://docs.docker.com/compose/install/) to build and start the whole stack with

```bash
cd docker
docker compose build
docker compose up -d
```

The default configuration is inside `docker/.env` and can be used as provided.

A user `user` with password `user` with all roles will also be created. This user has
the authorization to handle TLP WHITE and TLP GREEN advisories.

The application can then be reached under <http://localhost:5371>.

#### Docker-Keycloak (optional)
To try out different roles or users not included in the default Docker image, changes must be made through Keycloak.
The docker image uses a Keycloak which can be reached under <http://localhost:8080>.
The default admin-user set for the docker-Keycloak-instance is:
 * Username: admin
 * Password: secret

To find out how to create and manage users, read the [Keycloak documentation](./keycloak.md). Note that the
scripts are not designed for use with docker.

## Development-setup

The setup should be performed via the [installation scripts.](./scripts/README.md) on a Ubuntu 24.04 OS.


#### Run the application in a dev environment

To start the frontend via a `vite` dev-server:

```bash
cd client
npm run dev
```

This will start the client application and
print the URL a browser could be pointed to.

ISDuBA's backend is called `isdubad` and is located under `/cmd/isdubad/isdubad`.
An example-configuration for `isdubad` can be found in [example_isdubad.toml](./example_isdubad.toml). This example can be used as is. What each value represents is further 
explained in [the config documentation](isdubad-config.md). Per default, `isdubad` will
expect the config file to be named isduba.toml and to be within your working directory. The setup-scripts will create a usable example-`isduba.toml` within the main directory.

Otherwise, you can point isdubad towards the configuration file via the -c option. An example:

```bash
./cmd/isdubad/isdubad -c isduba.toml
```

The keycloak server set up via the installation scripts, needed to be able to login and authorize yourself within ISDuBA, can be started with: 
```bash
sudo -u keycloak /opt/keycloak/bin/kc.sh start-dev
```
Note that keycloak might take a while to start up.


After having made changes, the new application can be build via the Makefile:

```sh
make all
```

#### Upgrading
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

## Production-setup

You can download the latest stable release from [github.](https://github.com/ISDuBA/ISDuBA/releases/)

#### Building

Alternatively, to build the application the latest Golang version, NodeJS 20 and standard build
tools, like GNU Make are required. At the root of the repo run `make dist` to
build the frontend and backend; this will result in a Tar-file inside `dist/`
that contains the application. It can also be useful to look inside the
Dockerfile of the application to see how individual components of the
application can be built. If there are no special requirements it can be
enough to use the already built tar-file from the release page.

#### Setup

The Tar-file can be copied and extracted on a production server. This file
contains the `isdubad` backend, which can be run on any modern amd64 linux
system and the frontend which is contained in the `web/` folder. No further
dependencies are required to start the application. By default the backend will
serve the contents of `web/`. However, a [PostgreSQL database](#configuring-postgres) and [Keycloak instance](./keycloak.md) are still necessary to properly access ISDuBA.

See the [Keycloak documentation](./keycloak.md) on how to set up keycloak for your ISDuBA instance.

##### Configuring Postgres
 In your Postgres database, create a `keycloak` user with password `keycloak` as well as a database `keycloak` which will be owned by the user `keycloak`:
```
psql -c "CREATE USER keycloak WITH PASSWORD 'keycloak';"
createdb -O keycloak -E 'UTF-8' keycloak
```


 Next up, Postgres' [client authentification](https://www.postgresql.org/docs/current/auth-pg-hba-conf.html) configuration file has to be adjusted, by adding ISDuBA-directed configuration. Simply add

```
host    all             all             127.0.0.1/32            scram-sha-256
host    all             all             ::1/128                 scram-sha-256
```

to the end of the file. Read Postgres' [pg_hba.conf file documentation](https://www.postgresql.org/docs/current/auth-pg-hba-conf.html) for more information.

#### Running

For a quick start copy `example_isdubad.toml` to the
folder where the application is contained. Configure the postgres and Keycloak
settings and rename the file to `isduba.toml`. The application can now be
started. For exposing the application to the network it is recommended to use
a TLS-terminating reverse proxy.

##### Further documentation for production

See [security_considerations](./security_considerations.md) for security and maintenance considerations.

If you need help to know how to configure keycloak as an identity management for ISDuBA, read [our keycloak documentation.](./keycloak.md)

Where and how to configure the ISDuBA application is outlined [in isdubad-config.md.](./isdubad-config.md)

If other problems still persist, see if they are outlined [in the troubleshooting guide.](./troubleshooting.md)

