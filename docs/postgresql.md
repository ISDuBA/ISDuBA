<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

Everything described here can also be done via the [postgres install](./scripts/installpostgres.sh)
and [postgres configuration](./scripts/configurepostgres.sh) scripts.

This guide describes how to set up a postgres database for a development build on Ubuntu 24.04. These settings may not be suitable for production.

# Get PostgreSQL
Download PostgreSQL version 15 or newer.
PostgreSQL 16 has been used for development.
```
apt install postgresql
```

# Create PostgreSQL keycloak user
Allow Keycloak to access the PostgreSQL databases.
The created user for keycloak will have the username and password 'keycloak'.
```
su - postgres
```
Enter psql via:
```
psql
```
Create the Keycloak user so Keycloak can access it later:
```
CREATE USER keycloak WITH PASSWORD 'keycloak';
```
Give your postgres an explicit password so it can be accessed later:
```
ALTER USER postgres WITH PASSWORD 'postgres';
```
Exit psql via:
```
\q 
```

# Create Postgres database
Create a Postgres database for Keycloak.

```
createdb -O keycloak -E 'UTF-8' keycloak
```

Exit the postgres user via:
```
exit
```

# Edit Postgres config
Edit the PostgreSQL configuration.

Change to the postgres user and change into the postgres directory:
```
su - postgres
cd /etc/postgresql/16/main/
```

Open the postgresql.conf:
```
vim postgresql.conf
```
Change the following line:
> #listen_addresses = 'localhost'

to

```
listen_addresses = '*'
```
Open the pg_hba.conf:
```
vim pg_hba.conf
```
Add lines that allow logins via the network,
replace `192.168.56.1/32` in the following example
with the subnet you want to allow access from:
```
host    all             all             192.168.56.1/32         scram-sha-256
host    all             all             127.0.0.1/32            scram-sha-256
```

Use a method to reload that config file, e.g. `pg_ctl reload`.

Exit the postgres user:
```
exit
```
