<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

A collection of scripts which allows to set up ISDuBA on a Ubuntu 24.04 setup and simplifies
some important setup-steps.

The following scripts exist:

## [installall.sh](https://github.com/ISDuBA/ISDuBA/blob/main/docs/scripts/installall.sh)
This script will install git if needed and download the ISDuBA repository in the current directory unless it already exists, in which case it will update it.
Then it will call the setup.sh script that calls all other scripts to set up a testing environment.

installall.sh can be downloaded via:
``` bash
    curl -O https://raw.githubusercontent.com/ISDuBA/ISDuBA/main/docs/scripts/installall.sh
```
Then you can make it executable (e.g. via chmod) and use it to set up the testing environment with default values:
 * Keycloak admin credentials: username: ```keycloak```, password: ```keycloak```
 * ISDuBA user credentials: username: ```beate```, password: ```beate```
 * Keycloak runs on localhost port 8080
 * ISDuBA-Frontend runs on localhost port 5173

``` 
 Usage: installall.sh [--help] [branch name]
 where:
  --help       show this help text
  branch name  set up ISDuBA on the 'branch name' branch instead of main
```

## Description of the scripts called
The installall.sh script will call other scripts to do the setup. Some of these scripts
can be reused manually to update or adjust their respective parts of the application.
Explanations of the scripts follow.

### [setup.sh](https://github.com/ISDuBA/ISDuBA/blob/main/docs/scripts/setup.sh)
This script will call all other scripts (with the exception of installall)
 - This script can be used in the scripts folder if you already cloned the repository

### [installgojava.sh](https://github.com/ISDuBA/ISDuBA/blob/main/docs/scripts/installgojava.sh)
This script installs the latest go and Java jdk 17.

### [installkeycloak.sh](https://github.com/ISDuBA/ISDuBA/blob/main/docs/scripts/installkeycloak.sh)
This script creates a keycloak user and installs keycloak.

### [installpostgres.sh](https://github.com/ISDuBA/ISDuBA/blob/main/docs/scripts/installpostgres.sh)
This script installs a current PostgreSQL database via apt.

### [configurepostgres.sh](https://github.com/ISDuBA/ISDuBA/blob/main/docs/scripts/configurepostgres.sh)
This script configures a PostgreSQL for use with Keycloak.

### [keycloak/configureKeycloak.sh](https://github.com/ISDuBA/ISDuBA/blob/main/docs/scripts/keycloak/configurekeycloak.sh)
This script performs initial configuration steps to work with ISDuBA and starts keycloak on port 8080.
To do this, it calls upon the other scripts in the keycloak directory
   - the hereby created initial admin user will have both username and password ```keycloak```
   - the hereby created initial user will be the ``` editor ``` (role and group) ``` beate Bearbeiter ```
   - username and password for the initial user are ```beate```, the created editor group has access to all TLP WHITE advisories

### [keycloak/createRealm.sh](https://github.com/ISDuBA/ISDuBA/blob/main/docs/scripts/keycloak/createRealm.sh)
This script creates and adjusts the Keycloak ```isduba``` realm.

### [keycloak/createGroup.sh](https://github.com/ISDuBA/ISDuBA/blob/main/docs/scripts/keycloak/createGroup.sh)
<!---
 TODO: either needs to be expanded to cover all TLP levels and utilizing proper flags, or needs to be reworked
-->
``` 
 Usage: createGroup name [tlp] [publisher]
 where tlp:
  1 - Group will only be able to see TLP WHITE advisories (default) 
  2 - Group will only be able to see TLP GREEN advisories
  3 - Group will be able to see TLP WHITE and TLP GREEN advisories

 and publisher:
  Name of the publishers whose advisories the group can see. (Can only be set if tlp has been set)
```
### [keycloak/createRole.sh](https://github.com/ISDuBA/ISDuBA/blob/main/docs/scripts/keycloak/createRole.sh)

```
 Usage: createRole name description
 where name: name of the role
 description: description of the role
```
### [keycloak/createUser.sh](https://github.com/ISDuBA/ISDuBA/blob/main/docs/scripts/keycloak/createUser.sh)
```
Usage: createUser.sh username first_name last_name email_address password
```
### [keycloak/assignUserToRoleAndGroup.sh](https://github.com/ISDuBA/ISDuBA/blob/main/docs/scripts/keycloak/assignUserToRoleAndGroup.sh)
```
Usage: assignUserToRoleAndGroup.sh username groupname rolename
```
### [installplaywright.sh](https://github.com/ISDuBA/ISDuBA/blob/main/docs/scripts/installplaywright.sh)
This script installs node, the npm dependencies and playwright for the client.

###  [installisduba.sh](https://github.com/ISDuBA/ISDuBA/blob/main/docs/scripts/installisduba.sh)
This script creates config files, installs make and executes the Makefile to build binaries, distributions and executes integrationtests.

### [migrate.sh](https://github.com/ISDuBA/ISDuBA/blob/main/docs/scripts/migrate.sh)
This script prepares a database for use with ISDuBA.

### [start_isduba.sh](https://github.com/ISDuBA/ISDuBA/blob/main/docs/scripts/start_isduba.sh)
This script starts backend and frontend in the background, available on localhost port 5173.


## Filling the database
The bulkimporter set up via [installisduba.sh](https://github.com/ISDuBA/ISDuBA/blob/groups_not_roles_scripts/docs/scripts/README.md#installisdubash) can be used to import documents.
