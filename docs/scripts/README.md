<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

A collection of scripts which allow to set up ISDuBA on a Ubuntu 24.04 setup and simplify
some important setup-steps.

The following scripts exist:

 - installall.sh
   - Will install git if needed and download the ISDuBA repository in the current directory unless it already exists, in which case it will update it 
   - Will call the setup.sh script that calls all other scripts to set up a testing environment

 Usage: installall.sh [--help] [branch name]
 where:
  --help       show this help text
  branch name  set up ISDuBA on the 'branch name' branch instead of main

installall.sh can be downloaded via:
``` bash
    curl -O https://raw.githubusercontent.com/ISDuBA/ISDuBA/main/docs/scripts/installall.sh
```

 - setup.sh
   - Will call all other scripts (with the exception of installall)
    - Use this if you already cloned the repository, from the docs folder

 - installgojava.sh
  - installs the latest go and Java jdk 17

 - installkeycloak.sh
  - creates a keycloak user and installs keycloak

 - installpostgres.sh
  - installs a current PostgreSQL database via apt

 - configurepostgres.sh
  - configures a PostgreSQL for use with Keycloak

 - keycloak/configureKeycloak.sh
  - performs initial configuration steps to work with ISDuBA and starts keycloak on port 8080
  - calls upon the other scripts in the keycloak directory
   - the hereby created initial admin user will have both username and password ```keycloak```
   - the hereby created initial user will be the ``` editor ``` (role and group) ``` beate Bearbeiter ```
    - username and password for this user are ```beate```, the created editor group has access to all TLP WHITE advisories

 - keycloak/createGroup.sh

 # TODO: either needs to be expanded to cover all TLP levels and utilizing proper flags, or needs to be reworked
 Usage: createGroup name [tlp] [publisher]
 where tlp:
  1 - Group will only be able to see TLP WHITE advisories (default) 
  2 - Group will only be able to see TLP GREEN advisories
  3 - Group will be able to see TLP WHITE and TLP GREEN advisories

 and publisher:
  Name of the publishers whose advisories the group can see. (Can only be set if tlp has been set)

 - keycloak/createRole.sh

 Usage: createRole name description
 where name: name of the role
 description: description of the role

 - keycloak/createUser.sh
 Usage: createUser.sh username first_name last_name e-mail-adress password

 - keycloak/assignUserToRoleAndGroup.sh
 Usage: assignUserToRoleAndGroup.sh username groupname rolename

 - installplaywright.sh
  - installs node, the npm dependencies and playwright for the client

 - installisduba.sh
  - creates config files, installs make and executes the Makefile to build binaries, distributions and exectute integrationtests

 - migrate.sh
  - prepare a database for use with ISDuBA

 - start_isduba.sh
  - start backend and frontend in the background, avaible on localhost port 5173

The bulkimporter can be used to import manually downloaded documents.
