<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

This guide describes how to set up ISDuBA for a development build on Ubuntu 24.04. These settings may not be suitable for production.

## Prerequisites

 - A set of CSAF-Advisories, e.g. aquired via the [csaf_downloader tool.](https://github.com/csaf-poc/csaf_distribution)
 
## Setup ISDuBA
This setup should be performed via the [installation scripts.](./scripts/README.md)

An example-configuration for isdubad can be found in [example_isdubad.toml](./example_isdubad.toml)

To manually start a database migration, use [the migration script.](./scripts/migrate.sh)

To create additional users, use the [createUsers script.](./scripts/keycloak/createUsers.sh)
A list of users created by the setup scripts can be found in [the users.txt.](./developer/users.txt)

To create additional groups, use the [createGroup script.](./scripts/keycloak/createGroup.sh)

The keycloak admin user created via the scripts will have the username and password ```keycloak```.

### Import advisories
Import some advisories into the database via the bulk importer:
- host: host from where you download your advisories from
- /path/to/advisories/to/import: location to download your advisories from
(An example would be the results of the csaf_downloader, located in localhost)
From the repositories main directory:
```
./cmd/bulkimport/bulkimport -database isduba -user isduba -password isduba -host localhost /path/to/advisories/to/import
```

### (Optional) Check if `isdubad` is correctly installed
The following will define a TOKEN variable which holds the information 
about a user with name USERNAME and password USERPASSWORD as configured in keycloak.
(You can check whether the TOKEN is correct via e.g. jwt.io, keycloak needs to be up and running.)
```
TOKEN=`curl -d 'client_id=auth'  -d 'username=USERNAME' -d 'password=USERPASSWORD' -d 'grant_type=password' 'http://127.0.0.1:8080/realms/isduba/protocol/openid-connect/token' | jq -r .access_token`
```
The contents of the Token can be checked via:
```
echo $TOKEN
```

### Run the application in a dev environment

To start the frontend, change into the client directory and use npm to start it:

```bash
npm run dev
```

This will start the client application and
print the URL the browser can be pointed to.

To start the backend, start it up using a config file, e.g. from the main directory:

```bash
  ./cmd/isdubad/isdubad -c isdubad.toml
```

Make sure to have keycloak running when trying to access the application.
(If set up via the script available under:)
``` bash
 sudo ./opt/keycloak/bin/kc.sh start-dev
```

(The isduba-keycloak-specific-config is configured in the client/.env.)
