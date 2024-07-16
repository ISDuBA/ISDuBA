<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

This guide describes how to set up ISDuBA
for a *development* build on Ubuntu 24.04.
These settings are **not suitable for production**.


## Prerequisites

 - A set of CSAF-Advisories, e.g. aquired via the [csaf_downloader tool.](https://github.com/csaf-poc/csaf_distribution)


## Setup ISDuBA
The setup should be performed via the [installation scripts.](./scripts/README.md)

An example-configuration for `isdubad` can be found in [example_isdubad.toml](./example_isdubad.toml). Please edit to your needs.

Initially there is a migration needed to configure the database
by starting isdubad with the `ISDUBA_DB_MIGRATE` environment variable
set to true or by adjusting the toml-configuration file, e.g.

```sh
ISDUBA_DB_MIGRATE=true ./cmd/isdubad/isdubad -c isduba.toml
```

Create additional users via [createUsers script.](./scripts/keycloak/createUsers.sh) A list of users created by the setup scripts can be found in [the users.txt.](./developer/users.txt)

Create groups via [createGroup script.](./scripts/keycloak/createGroup.sh)
The restrictions set with the script are explained in [keycloak_values.md](./keycloak_values.md)

The keycloak admin user created via the scripts will
have the username password `keycloak`,
unless otherwise specified via the environment variable `KEYCLOAK_ADMIN`.
The password can be specified via the environment variable 
`KEYCLOAK_ADMIN_PASSWORD`, a file (`-f` option)
or directly (using the `-p` option).

If neither is set, then the script will try to see if
`docs/scripts/password.txt` contains a password.
If this is not set either, then a random password will be generated
and stored in `docs/scripts/password.txt`.


### Import advisories
Import the previously downloaded advisories into the database via the bulk importer:

An example for a local PostgreSQL:
- `~/downloaded_advisories`: location to download your advisories from, replace with your actual location

```sh
./cmd/bulkimport/bulkimport ~/downloaded_advisories
```

### (Optional) Check whether `isdubad` is correctly installed
The following will define a `TOKEN` variable which holds the information
about a user with name `USERNAME` and password `USERPASSWORD`
as configured in keycloak.

(You can check the `TOKEN` via jwt.io. Keycloak should be up and running.)

```sh
TOKEN=`curl -d 'client_id=auth'  -d 'username=USERNAME' -d 'password=USERPASSWORD' -d 'grant_type=password' 'http://127.0.0.1:8080/realms/isduba/protocol/openid-connect/token' | jq -r .access_token`
echo $TOKEN
```


### Run the application in a dev environment

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

Make sure to have keycloak running when trying to access the application.

(If set up via the script available under:)
```bash
sudo -u keycloak /opt/keycloak/bin/kc.sh start-dev
```

(The isduba-keycloak-specific-config is configured in `client/.env`.)
