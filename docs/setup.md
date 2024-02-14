The following document explains how to set up keycloak and postgresql to use
with as well as the basics of the tools provided through this repository on
an Ubuntu system, provided neither of the components have been previously
installed.

# Get Keycloak
Download Keycloak version 23.0.5, which has been used for development.
```
wget https://github.com/keycloak/keycloak/releases/download/23.0.5/keycloak-23.0.5.zip
```

# Unzip Keycloak
```
unzip keycloak-23.0.5.zip
```
(unzip or alternatively any other program that is capable of decompressing
 .zip archives may need to be installed first.)
```
mv keycloak-23.0.5 /opt/keycloak
```

# Get Postgresql 16
Download Postgresql version 16, which has been used for development.
```
apt install vim gnupg2 -y
curl -fsSL https://www.postgresql.org/media/keys/ACCC4CF8.asc| gpg --dearmor -o /etc/apt/trusted.g>
sh -c 'echo "deb https://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt>
apt update
apt install postgresql-16
```

# Create Postgresql keycloak user
Allow Keycloak to access the Postgresql databases.
The created user for keycloak will have the username and password 'keycloak'.
It is not recommended to use these for production,
as they are too easy to guess.
Similarly for the user.
```
su - postgres
```
Enter psql via:
```
psql
```
Create the Keycloak and your personal user, so Keycloak can access it later:
```
CREATE USER keycloak WITH PASSWORD 'keycloak';
CREATE USER username WITH PASSWORD 'userpassowrd';
```
Exit psql via:
```
\q 
```

# Create Postgres database
Create a Postgres and bsi database for Keycloak.

```
createdb -O keycloak -E 'UTF-8' keycloak
createdb -O username -E 'UTF-8' username
```

Exit the postgres user via:
```
exit
```

# Alter Keycloak config
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
The example-username and password are easy to guess and should not be used
for production.

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

# Initialize keycloak
Start Keycloak and allow it to configure itself.
```
bin/kc.sh start-dev
```

# Edit Postgres config
Edit the Postgresql configuration.

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
```
set listen_addresses = '*'
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

# Edit systemd keycloak.service
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

# Adjust systemd
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

# Adjust keycloak
Configure Keycloak.

Open Keycloaks Web-Interface, running on localhost:8080.
Via the admin console adjust the following if necessary:

- Create Isduba realm

- Create Users

### Create Clients: auth

Under Clients, create auth:

ID/Name: ```auth```

### Via Clients: auth:

- valid redirect url: ```/*```

- web origins url: ```/*```

- Tick the boxes Standard flow and Direct access grants

- Turn off ```consent required```

### Switch from "settings" to "client scopes" and click on auth-dedicated

#### Add mapper "User Attribute" with

- Name: TLP

- User Attribute: TLP

- Token Claim Name: TLP

- Claim JSON type: JSON

- For the switches, Multivalued should be turned off, the rest on

Create roles via Realm roles:
The following allows the role to handle
the WHITE and GREEN TLP levels of all publishers. Adjust as necessary:
Switch to the Attributes tab and set:

- Key: TLP

- Value: [{"publisher":"", "tlps":["WHITE, GREEN"]}]

### Assign Users their roles
Via ```Users``` via ```Role Mapping``` via ```Assign Role``` assign the users
their role.

# Setup Go
Download Go 1.22:
```
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
```
Extract it and place the new go version into the /usr/local directory:
```
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
```
## Edit your profile to always use this go:
```
vim /etc/profile
```
In there, add the line:
> export PATH=$PATH:/usr/local/go/bin

# Setup ISDuBA
Clone the repository:
```
git clone https://github.com/ISDuBA/ISDuBA.git
```
Switch into the directory
```
cd ISDuBA
```
## build the tools
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

# Create isduba configuration
Create a configuration file for the tools used in this repository.
Create a configuration file:
```
vim isduba-bsi.toml
```

Configure your setup, e.g. as follows:
```
[logging]
file="bsi.log"
level="debug"

[database]
migrate=true
user="username"
password="userpassword"
database="username"
host="localhost"
port=5432
admin_user="postgres"
admin_password="postgres"
admin_database="postgres"
```

# Import advisories
Import some advisories into the database via the bulk importer:
- host: host from where you download your advisories from
- advisories_to_import: location to download your advisories from
(An example would be the results of the csaf_downloader, located in localhost)
From the repositories main directory:
```
./cmd/bulkimport/bulkimport -database bsi -user bsi -password bsi -host localhost advisories_to_import
```

# Example use of isdubad
The following will define a TOKEN variable which holds the information 
about a user with name USERNAME and password USERPASSWORD.
Both the username and password are too easy to guess for production.
(You can check whether the TOKEN is correct via e.g. jwt.io)
```
TOKEN=`curl -d 'client_id=auth'  -d 'username=USERNAME' -d 'password=USERPASSWORD' -d 'grant_type=password' 'http://127.0.0.1:8080/realms/isduba/protocol/openid-connect/token' | jq -r .access_token`
```
The contents of the Token can be checked via:
```
echo $TOKEN
```
