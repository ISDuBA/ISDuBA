<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: MIT

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

# Setup

## Setup backend

The following document explains how to set up keycloak and postgresql to use
with as well as the basics of the tools provided through this repository on
an Ubuntu system, provided neither of the components have been previously
installed. Note that all IDs and passwords used in this setup are
easy to guess and should not be used in production, only for development.

### Prerequisites
A sufficiently new version of Java as well as an unzip-tool like unzip need to be installed.
You can install Java 17 via
```
sudo apt install openjdk-17-jre-headless
```
and unzip via
```
sudo apt install unzip
```

### Get Keycloak
Download Keycloak version 23.0.5, which has been used for development.
```
wget https://github.com/keycloak/keycloak/releases/download/23.0.5/keycloak-23.0.5.zip
```

### Unzip Keycloak
```
unzip keycloak-23.0.5.zip
```
(unzip or alternatively any other program that is capable of decompressing
 .zip archives may need to be installed first.)
```
mv keycloak-23.0.5 /opt/keycloak
```

### Get PostgreSQL 16
Download PostgreSQL version 16, which has been used for development.
```
apt install vim gnupg2 -y
curl -fsSL https://www.postgresql.org/media/keys/ACCC4CF8.asc| gpg --dearmor -o /etc/apt/trusted.gpg.d/postgresql.gpg
sh -c 'echo "deb https://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list' 
apt update
apt install postgresql-16
```

### Create PostgreSQL keycloak user
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

### Create Postgres database
Create a Postgres database for Keycloak.

```
createdb -O keycloak -E 'UTF-8' keycloak
```

Exit the postgres user via:
```
exit
```

### Alter Keycloak config
Create a Keycloak user with access rights to your Keycloak
directory.
```
useradd keycloak
chown -R keycloak: /opt/keycloak
```
Open the Keycloak config with a text-editor (like vim):
```
vim /opt/keycloak/conf/keycloak.conf
```
Paste the following for a valid configuration. 

```
# Basic settings for running in production. Change accordingly before deploying the server.

# Database

# The database vendor.
db=postgres

# The username of the database user.
db-username=keycloak

# The password of the database user.
db-password=keycloak

# The full database JDBC URL. 
# If not provided, a default URL is set based on the selected database vendor.
db-url=jdbc:postgresql://localhost/keycloak

# Observability

# If the server should expose healthcheck endpoints.
#health-enabled=true

# If the server should expose metrics endpoints.
#metrics-enabled=true

# HTTP

# The file path to a server certificate or certificate chain in PEM format.
#https-certificate-file=${kc.home.dir}conf/server.crt.pem

# The file path to a private key in PEM format.
#https-certificate-key-file=${kc.home.dir}conf/server.key.pem

# The proxy address forwarding mode if the server is behind a reverse proxy.
#proxy=reencrypt

# Do not attach route to cookies and rely on the session affinity capabilities from reverse proxy
#spi-sticky-session-encoder-infinispan-should-attach-route=false

# Hostname for the Keycloak server.
#hostname=isduba
```

### Initialize keycloak
Allow Keycloak to start on system-startup.

Create a systemd Keycloak file via a text editor, e.g. vim:
```
vim /etc/systemd/system/keycloak.service
```
Use the following configuration:

```
[Unit]
Description=Keycloak
After=network.target

[Service]
Type=idle
User=keycloak
Group=keycloak
ExecStart=/opt/keycloak/bin/kc.sh start-dev
TimeoutStartSec=600
TimeoutStopSec=600

[Install]
WantedBy=multi-user.target
```

### Adjust systemd
As superuser, enable keycloak to start on system-startup.

Enter superuser status.
```
sudo su
```

Enable Keycloak to start on system-startup.
```
systemctl enable keycloak
systemctl start keycloak
```

Start Keycloak in the background or restart your system.
```
bin/kc.sh start-dev &
```

### Edit Postgres config
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
Add the following two lines:
```
host    all             all             192.168.56.1/32         scram-sha-256
host    all             all             127.0.0.1/32            scram-sha-256
```

Exit the postgres user:
```
exit
```



### Adjust keycloak
Configure Keycloak.

Open Keycloaks Web-Interface, running on localhost:8080.
Via the admin console adjust the following if necessary:

- Create ```isduba``` realm

#### Create Clients: auth

Under Clients, create auth:

ID/Name: ```auth```

#### Via Clients: auth:

- valid redirect url: ```/*```

- web origins url: ```/*```

- Tick the boxes Standard flow and Direct access grants

- Turn off ```consent required```

#### Switch from "settings" to "client scopes" and click on auth-dedicated

##### Add mapper "User Attribute" with

- Name: ```TLP```

- User Attribute: ```TLP```

- Token Claim Name: ```TLP```

- Claim JSON type: ```JSON```

- For the switches, Multivalued should be turned off, the rest on

Create roles via Realm roles:

E.g. 

- Name: ```bearbeiter```
- Description: ```bearbeiter```

#### Add attributes

The following attribute allows the role to handle
the WHITE and GREEN TLP levels of all publishers. Adjust as necessary:
Switch to the Attributes tab and set:

- Key: ```TLP```

- Value: ```[{"publisher":"", "tlps":["WHITE", "GREEN"]}]```

### Create Users

Via ```Users``` use ```Create User``` to create a user.
USERNAME and USERPASSWORD are example credentials.
 
 - Username: ```USERNAME```
 - E-Mail verified: ```yes```

Then, set the password via ```Credentials```. This example uses the password
```USERPASSWORD```
Turn ```temporary``` off.

#### Assign Users their roles
Via ```Users``` via ```Role Mapping``` via ```Assign Role``` assign the users
their role.

### Setup Go
Download Go 1.22:
```
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
```
Extract it and place the new go version into the /usr/local directory:
```
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
```
#### Make the profile always use this version of go:
Open the profile with a text manager.
```
vim /etc/profile
```
In there, add the line:
```
export PATH=$PATH:/usr/local/go/bin
```
The system will now use go1.22 when go is called upon.

### Setup ISDuBA
Clone the repository:
```
git clone https://github.com/ISDuBA/ISDuBA.git
```
Switch into the directory
```
cd ISDuBA
```
#### build the tools
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
### Create isduba configuration
Create a configuration file for the tools used in this repository.
A detailed description of this configuration file can be found [here](./isdubad-config.md).
Create a configuration file:
```
cp docs/example_isdubad.toml isdubad.toml
vim isdubad.toml
```

### Start `isdubad` to allow db creation
From the repositories main directory, start the isdubad program,
which creates the db and users according to the ./cmd/isdubad/isdubad -c isdubad.toml:
```
ISDUBA_DB_MIGRATE=true ./cmd/isdubad/isdubad -c isdubad.toml 
```

After the initial migration you can un-configure the `admin_` parts In
the configuration file adn start `isdubad` without the `ISDUBA_DB_MIGRATE`
env var set.

### Import advisories
Import some advisories into the database via the bulk importer:
- host: host from where you download your advisories from
- /path/to/advisories/to/import: location to download your advisories from
(An example would be the results of the csaf_downloader, located in localhost)
From the repositories main directory:
```
./cmd/bulkimport/bulkimport -database isdubad -user isdubad -password isdubad -host localhost /path/to/advisories/to/import
```

### Example use of `isdubad`
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

## Setup client

### Prerequisites

A current Version of nodeJS LTS (version `20.11.1`).

### Install necessary packages

Assuming you are in the checked out repository

```bash
cd client
npm install
npx playwright install
```

### Run the client application in a dev environment

```bash
npm run dev -- --open
```

This will start the client application and opens a window in your default browser.