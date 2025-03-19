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

## [setup.sh](https://github.com/ISDuBA/ISDuBA/blob/main/docs/scripts/setup.sh)
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




