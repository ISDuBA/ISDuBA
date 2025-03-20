<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

A collection of scripts which allows to set up ISDuBA on a Ubuntu 24.04 setup and simplifies
some important setup-steps. 

These scripts install the following:
 * [Postgresql 15](https://www.postgresql.org/docs/release/15.0/)
 * [Keycloak 26](https://www.keycloak.org/docs/latest/release_notes/index.html#keycloak-26-0-0)
 * make
 * bash
 * curl
 * sed
 * tar
 * Java 21 runtime environment
 * the latest go version
 * graphviz
 * swag
 * xkcdpass
 * playwright
 * node
 * ca-certificates
 
Since these scripts need to install all dependencies, sudo privileges are required. 

### important notes

 - A list of users created by the setup scripts can be found in [the users.txt.](./developer/users.txt) including their usernames and passwords. Editing this file before using the createUsers or the setup script will change which users are created.
 
 - The Keycloak admin user created via the scripts will have the username `keycloak`, unless otherwise specified via the environment variable `KEYCLOAK_ADMIN`. The password can be specified via the environment variable  `KEYCLOAK_ADMIN_PASSWORD`, a file (`-f` option) or directly (using the `-p` option).
  - If neither option is used, then the script will try to see if `docs/scripts/password.txt` exists and contains a password. If this is not set either, then a random password will be generated and stored in `docs/scripts/password.txt`.
  
## [installall.sh](https://github.com/ISDuBA/ISDuBA/blob/main/docs/scripts/installall.sh)
This script will download the ISDuBA repository in the current directory unless it already exists, in which case it will update it.
Then it will call the [setup.sh script](./setup.sh) that calls all other scripts to set up a testing environment.

installall.sh can be used via:
``` bash
    curl -O https://raw.githubusercontent.com/ISDuBA/ISDuBA/main/docs/scripts/installall.sh
    chmod +x
    ./installall.sh
```
This sets up the testing environment with default values:
 * Keycloak admin credentials: username: ```keycloak```, password will be randomly generated. A file containing the password will be set into the scripts folder.
 * Alternatively, set the environment variables KEYCLOAK_ADMIN and KEYCLOAK_ADMIN_PASSWORD to manually set username and password.
 * Keycloak is configured to listen on port 8080
 * ISDuBA-Frontend is configured to listen on localhost port 5173

```
Usage: installall.sh [OPTIONS]
where OPTIONS:
  -h, --help                       show this help text and exit script (optional)
  -b, --branch=name                set up on branch 'name' instead of main (optional)
  -k, --keycloakRunning            signal the script that there is a keycloak running on port 8080 (optional)
```

### [setup.sh](https://github.com/ISDuBA/ISDuBA/blob/main/docs/scripts/setup.sh)
This script will call other scripts to finish the setup. Use this if you already cloned the repository.

#### Additional options

Should use of a custom keycloak instance be desired, it may be necessary to signal this to ISDuBA.
The setup scripts utilize Keycloak's health checks to determine whether Keycloak is running. The port to use may change depending on your Keycloak version or admin's configuration.
The default for the current version of 26 is port 9000.
This means it may be necessary to call docs/scripts/keycloak/configurekeycloak.sh with the -l/--live flag to manually set a port, e.g. for Keycloak 24:

```bash
  ./configureKeycloak.sh --live 8080
```

Not setting the correct port without the -k/--keycloakRunning option will cause the script to try and call the wrong port over and over until stopped.

Some further notes for development can be found [in our development documentation](./development.md)

## Called Scripts

The following will briefly explain what every other script does. It's not necessary to call any script again after installall.sh or setup.sh have been completed.

#### [configurepostgres.sh](./configurepostgres.sh)
 Prepares the postgres database for use with ISDuBA.

#### [installisduba.sh](./installisduba.sh)
 Creates a config file for ISDuBA and executes the Makefile.
 
#### [installkeycloak.sh](./installkeycloak.sh)
 Installs keycloak version 26, creates a keycloak user on your system and gives them the ownership over it.
 
#### [installplaywright.sh](./installplaywright.sh)
 Installs ca-certificates, curl, gnupg and playwright with dependencies.
 
#### [installpostgres.sh](./installpostgres.sh)
  Installs Postgresql 15.
  
#### [installutilities.sh](./installutilities.sh)
  Installs 
   * [make](https://www.gnu.org/software/make/)
   * [bash](https://www.gnu.org/software/bash/)
   * [curl](https://curl.se/)
   * [sed](https://www.gnu.org/software/sed/)
   * [tar](https://www.gnu.org/software/tar/)
   * [Java (openjdk-21-jre-headless)](https://openjdk.org/projects/jdk/21/)
   * [Go](https://go.dev/)
   * [xkcdpass](https://pypi.org/project/xkcdpass/)
   * [graphviz](https://graphviz.org/)
   * [swag](https://github.com/swaggo/swag).

#### [list_licenses.py](./list_licenses.py)
 Extract licensing info from the packages of an SPDX-2.3 SBOM JSON file.
 
#### [start_all.sh](./start_all.sh)
 Requires sudo privileges. Starts isdubad, the isduba client and keycloak. Is not called by the setup scripts and can be called independently to quickly start all components necessary.

### Keycloak-scripts

All these following scripts will adjust Keycloak to create a development setup with example users and groups as outlined in [the users example file list users.txt.](../../developer/users.txt).
Aside from [configureKeycloak.sh](#configurekeycloak.sh), these scripts require Keycloak to be running and you to be able to log into keycloak.
 
#### [assignUserToRoleAndGroup.sh](./keycloak/assignUserToRoleAndGroup.sh)
 Can be used to give ISDuBA users a role or a group. 

```
 Usage: assignUserToRoleAndGroup.sh name OPTIONS"
 where name:
   -n, --name=name                  username of the user to be added to roles or groups.
 where OPTIONS:
   -h, --help                       show this help text and exit script (optional).
   -g, --group=name                 name of the group the user should be added to (optional).
   -r, --role=name                  name of the role the user should be added to (optional).
       --noLogin                    do not attempt to log into keycloak. Requires active 
login to not cause errors (optional).
```

#### [configureKeycloak.sh](./keycloak/configureKeycloak.sh)
 If not signalled to be already running, starts Keycloak, creates an admin user if it doesn't exist yet and calls all other Keycloak scripts with default values.
 
```
 Usage: configureKeycloak.sh [OPTIONS]
 where OPTIONS:
 -h, --help                       Show this help text and exit script (optional).
 -k, --keycloakRunning            Skip checks on whether keycloak is running.
 -f, --file                       Specify file storing the keycloak admin password. (optional, default: ./../password.txt)
 -l, --live                       Specify the port which accepts keycloak health checks. (Optional, default: 9000)
 -p, --password                   Specify the keycloak admin password directly (optional).
```

#### [createGroup.sh](./keycloak/createGroup.sh)
 [Create a group](./../keycloak.md#groups) for ISDuBAs Keycloak. 
```
Usage: createGroup.sh [OPTIONS] --name name
 where OPTIONS:
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

#### [createRealm.sh](./keycloak/createRealm.sh)
 Creates an `isduba` Keycloak realm that can be used to manage all ISDuBA-users.
  
#### [createRole.sh](./keycloak/createRole.sh)
 [Create a role] for ISDuBAs Keycloak. ISDuBA uses a set set of roles that are created during the initial setup, so there should be no reason to call this script manually.
``` 
 Usage: createRole name description [login]
 where:
  name: name of the role
  description: description of the role
  login=[true|false]: Whether to log into keycloak again (default:true).
```
#### [createUser.sh](./keycloak/createUser.sh)
 Creates a singular user for ISDuBA, along with credentials.
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
#### [createUsers.sh](./keycloak/createUsers.sh)
 Reads users to be created from a file and calls the [createUser skript](#createuser.sh) to create them.
```
Usage: createUsers.sh OPTIONS"
where OPTIONS:"
  -h, --help                       show this help text and exit script (optional).
  -f, --file=file                  name of the file that contains all user information (mandatory).
      --noLogin                    do not attempt to log into keycloak. Requires active login to not cause errors (optional).
```
