<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

Calling installall.sh will, on a new Ubuntu 24.04 setup: 
 - install git and download the git repository of ISDuBA and then call setup.sh which will
   - install java and go via installgojava.sh
   - install PostgreSQL via installpostgres.sh
   - create keycloak user and adjust postgres user for postgres via configurepostgres.sh   
   - install keycloak 23.0.5 via installkeycloak.sh
   - adjust the keycloak configuration via configurekeycloak.sh
   - install node 20 and all frontend dependencies via installplaywright.sh
   - build the bulkimporter and isdubad tools via installisduba.sh
   - prepare the database via migrate.sh
   - start up the backend and frontend so they can be used via start_isduba.sh on localhost port 5173
 

To call installall (sudo privileges needed):
``` bash
    curl --fail -O https://raw.githubusercontent.com/ISDuBA/ISDuBA/main/docs/scripts/installall.sh
    installall.sh
```

Finally, the bulkimporter can be used to
import manually download documents, [see setup.md](./../setup.md#start-isdubad-to-allow-db-creation)

After that the setup is complete.
