<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: MIT

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

# Prerequisites

 - A keycloak setup as described [here](./keycloak.md)

 - A recent version of go, see [here on how to download and install go](https://go.dev/doc/install)

 - A set of CSAF-Advisories, e.g. aquired via the [csaf_downloader tool](https://github.com/csaf-poc/csaf_distribution)
 
# Setup ISDuBA
Clone the repository:
```
git clone https://github.com/ISDuBA/ISDuBA.git
```
Switch into the directory
```
cd ISDuBA
```
## build the tools
Switch into the bulkimport directory and build it:
```
cd cmd/bulkimport
go build
```
Switch into the isdubad directory and build it:
```
cd ../isdubad
go build
```
Return to the main directory:
```
cd ../..
```
# Create isduba configuration
Create a configuration file for the tools used in this repository.
A detailed description of this configuration file can be found [here](./isdubad-config.md).
Create a configuration file:
```
vim isdubad.toml
```

An example configuration can be found [here](./isdubad.toml).

# Start `isdubad` to allow db creation
From the repositories main directory, start the isdubad program,
which creates the db and users according to the ./cmd/isdubad/isdubad -c isdubad.toml:
```
ISDUBA_DB_MIGRATE=true ./cmd/isdubad/isdubad -c isdubad.toml 
```

After the initial migration you can un-configure the `admin_` parts In
the configuration file adn start `isdubad` without the `ISDUBA_DB_MIGRATE`
env var set.

# Import advisories
Import some advisories into the database via the bulk importer:
- host: host from where you download your advisories from
- /path/to/advisories/to/import: location to download your advisories from
(An example would be the results of the csaf_downloader, located in localhost)
From the repositories main directory:
```
./cmd/bulkimport/bulkimport -database isdubad -user isdubad -password isdubad -host localhost /path/to/advisories/to/import
```

# Example use of `isdubad`
The following will define a TOKEN variable which holds the information 
about a user with name USERNAME and password USERPASSWORD as configured in keycloak.
(You can check whether the TOKEN is correct via e.g. jwt.io)
```
TOKEN=`curl -d 'client_id=auth'  -d 'username=USERNAME' -d 'password=USERPASSWORD' -d 'grant_type=password' 'http://127.0.0.1:8080/realms/isduba/protocol/openid-connect/token' | jq -r .access_token`
```
The contents of the Token can be checked via:
```
echo $TOKEN
```
