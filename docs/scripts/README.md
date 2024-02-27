Calling installall.sh will
 - install git and download the git repository of ISDuBA and then call setup.sh which will
   - install java and go via installgojava.sh
   - build the bulkimporter and isdubad tools via installisduba.sh
   - install keycloak 23.0.5 via installkeycloak.sh
   - adjust the keycloak configuration via configurekeycloak.sh
   - install PostgreSQL via installpostgres.sh
   - create keycloak user and adjust postgres user for postgres via configurepostgres.sh
   - enable keycloak to start on systemstartup via keycloakonsysrtemstart.sh
   - install node 20 and all frontend dependencies via installplaywright.sh
   

To call installall (as root):
``` bash
    curl --fail -O https://raw.githubusercontent.com/ISDuBA/ISDuBA/main/docs/scripts/installall.sh
    installall.sh
```
Keycloak will still need to be started, have a realm created and
configured via it's userinterface.
Then Isdubad needs to be started to allow for the database migration.
Finally, the bulkimporter can be used to
import manually download documents.

After that, everything is set up.
