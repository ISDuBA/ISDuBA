<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

This guide describes how to set up keycloak on Ubuntu for a development build.
These settings may not be suitable for production.

# Prerequisites

 * A recent PostgreSQL installation
   * See [postgresql.md](./postgresql.md) for instructions on how to perform the setup.
<!---
   * Alternatively, use the [setup script]() // TODO
--->   
 * A recent version of java, e.g. Java 17

 * (Optional) Superuser privileges to allow keycloak to start on system-startup
 
# Get Keycloak

This section can be automated via the [keycloak setup script](./scripts/installkeycloak.sh)
and the [keycloak configuration script](./scripts/configurekeycloak.sh)

The creation of Realms and Users via keycloak needs to be done manually still.

Download a recent version of Keycloak.
Version 24.0.1 has been used for development.

```
wget https://github.com/keycloak/keycloak/releases/download/24.0.1/keycloak-24.0.1.zip
```

### Unzip Keycloak

```
unzip keycloak-24.0.1.zip
```

```
mv keycloak-24.0.1 /opt/keycloak
```

### Alter Keycloak config
Create a Keycloak user with access rights to your Keycloak
directory.
```
adduser --disabled-password --system --group --gecos "" keycloak
chown -R keycloak:keycloak /opt/keycloak
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

# (Optional) Keycloak on System-Startup
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

###  Adjust systemd
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

# Start Keycloak

Start Keycloak and allow it to configure itself.
```
bin/kc.sh start-dev &
```

# Adjust keycloak
Configure Keycloak.

Open Keycloaks Web-Interface, running on localhost:8080.
Via the admin console adjust the following if necessary:

- Create ```isduba``` realm

### Create Clients: auth

Under Clients, create auth:

ID/Name: ```auth```

### Via Clients: auth:

- `Root URL`: ```http://localhost:5173/``` 

- `Valid redirect URIs`: ```http://localhost:5173/*```

- `Valid post logout redirect URIs`: `+` or `/*`. `+` means that the value from `Valid redirect URIs` is taken.

- `Web origins`: ```*```

- `Admin URL`: ```http://localhost:5173/```

- Tick the boxes Standard flow and Direct access grants

- Turn off ```consent required```

### Switch from "settings" to "client scopes" and click on auth-dedicated

#### Add mapper "User Attribute" with

- Name: ```TLP```

- User Attribute: ```TLP```

- Token Claim Name: ```TLP```

- Claim JSON type: ```JSON```

- For the switches, Multivalued should be turned off, the rest on

Create roles via Realm roles:

E.g. 

- Name: ```bearbeiter```
- Description: ```bearbeiter```

### Add attributes

The following attribute allows the role to handle
the WHITE and GREEN TLP levels of all publishers. Adjust as necessary:
Switch to the Attributes tab and set:

- Key: ```TLP```

- Value: ```[{"publisher":"", "tlps":["WHITE", "GREEN"]}]```

## Create Users

Via ```Users``` use ```Create User``` to create a user.
USERNAME and USERPASSWORD are example credentials.
 
 - Username: ```USERNAME```
 - E-Mail verified: ```yes```

Then, set the password via ```Credentials```. This example uses the password
```USERPASSWORD```
Turn ```temporary``` off.

### Assign Users their roles
Via ```Users``` via ```Role Mapping``` via ```Assign Role``` assign the users
their role.


### (Optional) Adjust necessary profile information

By default, any user needs a registered E-Mail Adress, First Name, Last Name and Username.
These settings can be adjusted via  ```Realm Settings``` under ``` User Profile```. 
E.g. to allow users to log in with just their Username and Password, you can uncheck ```Required Field```
under E-Mail Adress, First Name and Last Name.
