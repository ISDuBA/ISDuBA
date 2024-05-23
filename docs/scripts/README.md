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
Usage: installall.sh [OPTIONS]"
where OPTIONS:"
  -h, --help                       show this help text (optional)
  -b, --branch=name                set up ISDuBA on branch 'name' instead of main (optional)
  -k, --keycloakRunning            signal the script that there is a keycloak running on port 8080 (optional)
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
   - the hereby created initial users are listed in [keycloak/users.txt](https://github.com/ISDuBA/ISDuBA/blob/main/docs/scripts/keycloak/users.txt)

### [keycloak/createRealm.sh](https://github.com/ISDuBA/ISDuBA/blob/main/docs/scripts/keycloak/createRealm.sh)
This script creates and adjusts the Keycloak ```isduba``` realm and requires Keycloak to be running to function.

### [keycloak/createGroup.sh](https://github.com/ISDuBA/ISDuBA/blob/main/docs/scripts/keycloak/createGroup.sh)
This script creates a group and requires Keycloak to be running to function.
```
Usage: createGroup.sh [OPTIONS] --name name"
where OPTIONS:"
  -h, --help                       show this help text and exit script (optional).
  -w, --white                      grant the group access to TLP:WHITE advisories (optional).
  -g, --green                      grant the group access to TLP:GREEN advisories (optional).
  -a, --amber                      grant the group access to TLP:AMBER advisories (optional).
  -r, --red                        grant the group access to TLP:RED advisories (optional).
  -p, --publisher=name             restrict access to advisories of the named publisher (optional).
      --noLogin                    do not attempt to log into keycloak. Requires active login to not cause errors (optional).
and
  -n, --name=name                  name of the group that is supposed to be created (mandatory).

```
### [keycloak/createRole.sh](https://github.com/ISDuBA/ISDuBA/blob/main/docs/scripts/keycloak/createRole.sh)
This script creates a role and requires Keycloak to be running to function.

```
 Usage: createRole name description [login]
 where:
  name: name of the role
  description: description of the role
  login=[true|false]: Whether to log into keycloak again (default:true).
```
### [keycloak/createUsers.sh](https://github.com/ISDuBA/ISDuBA/blob/main/docs/scripts/keycloak/createUsers.sh)
This script creates a user and assigns a role and group to them via the createUser and assignUserToRoleAndGroup scripts.
The user parameters are read from a file, see [keycloak/users.txt](https://github.com/ISDuBA/ISDuBA/blob/main/docs/scripts/keycloak/users.txt)
for an example.
```
 Usage: createUsers.sh --file filename [OPTIONS]
 where OPTIONS:
  -h, --help                       show this help text and exit script (optional).
  -f, --file=filename              name of the file that contains all user information (mandatory).
      --noLogin                    do not attempt to log into keycloak. Requires active login to not cause errors (optional).
```

### [keycloak/createUser.sh](https://github.com/ISDuBA/ISDuBA/blob/main/docs/scripts/keycloak/createUser.sh)
This script creates a user and requires Keycloak to be running to function.
```
 Usage: createUser.sh username first_name last_name email_address password [login]
 where:
  username: username of the user
  first_name: first name of the user
  last_name: surname of the user
  email_address: the users registrered email-address
  password: password of the user
  login=[true|false]: whether to login to keycloak again (default:true)
```
### [keycloak/assignUserToRoleAndGroup.sh](https://github.com/ISDuBA/ISDuBA/blob/main/docs/scripts/keycloak/assignUserToRoleAndGroup.sh)
This script assigns roles and groups to users and requires Keycloak to be running to function.
```
 Usage: assignUserToRoleAndGroup.sh username groupname rolename [login]
 where:
  username: username of the user
  groupname: group the user is going to be assigned to
  rolename: role the user is going to be assigned to
  login=[true|false]: whether to login to keycloak again (default:true)
```
### [installplaywright.sh](https://github.com/ISDuBA/ISDuBA/blob/main/docs/scripts/installplaywright.sh)
This script installs node, the npm dependencies and playwright for the client.

###  [installisduba.sh](https://github.com/ISDuBA/ISDuBA/blob/main/docs/scripts/installisduba.sh)
This script creates config files, installs make and executes the Makefile to build binaries, distributions and executes integrationtests.

### [migrate.sh](https://github.com/ISDuBA/ISDuBA/blob/main/docs/scripts/migrate.sh)
This script prepares a database for use with ISDuBA.

## Filling the database
The bulkimporter set up via [installisduba.sh](https://github.com/ISDuBA/ISDuBA/blob/groups_not_roles_scripts/docs/scripts/README.md#installisdubash) can be used to import documents.

## Starting the application
To get testing, start up keycloak, the frontend and the backend:

 * Keycloak is located in /opt/keycloak, start it up e.g. via ``` /opt/keycloak/bin/kc.sh start-dev```

 * The ISDuBA-backend can be started via the ```isdubad``` binary (in cmd/isdubad/) utilizing a config [toml](https://toml.io/en/) file given via the ```-c``` flag, e.g. the  ```isdubad.toml``` cloned from the docs directory into the main directory

 * The ISDuBA-frontend can be started in development mode via ```npm run dev``` in the client directory
