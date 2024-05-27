<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

A collection of scripts which allows to set up ISDuBA on a Ubuntu 24.04 setup and simplifies
some important setup-steps.

## [installall.sh](https://github.com/ISDuBA/ISDuBA/blob/main/docs/scripts/installall.sh)
This script will install git if needed and download the ISDuBA repository in the current directory unless it already exists, in which case it will update it.
Then it will call the setup.sh script that calls all other scripts to set up a testing environment.

installall.sh can be downloaded via:
``` bash
    curl -O https://raw.githubusercontent.com/ISDuBA/ISDuBA/main/docs/scripts/installall.sh
```
Then you can make it executable (e.g. via chmod) and use it to set up the testing environment with default values:
 * Keycloak admin credentials: username: ```keycloak```, password: ```keycloak```
 * Keycloak runs on localhost port 8080
 * ISDuBA-Frontend runs on localhost port 5173

```
Usage: installall.sh [OPTIONS]
where OPTIONS:
  -h, --help                       show this help text (optional)
  -b, --branch=name                set up ISDuBA on branch 'name' instead of main (optional)
  -k, --keycloakRunning            signal the script that there is a keycloak running on port 8080 (optional)
```

## [setup.sh](https://github.com/ISDuBA/ISDuBA/blob/main/docs/scripts/setup.sh)
This script will call all other scripts (with the exception of installall.) (you can use this if you already cloned the repository)
