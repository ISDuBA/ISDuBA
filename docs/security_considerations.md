<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2025 Intevation GmbH <https://intevation.de>
-->

Within the following, some security and maintenance considerations when running ISDuBA are outlined.

# Application

The application does not provide a secure endpoint. It should be put behind a
reverse proxy with a valid TLS certificate.

As a precaution, place the backend machine that runs `isdubad`
in a network setup that it does not have access to internal services.

ISDuBA is built to download CSAF documents from the internet.
The places where these are searched for are configured by users
with role `source-manager` and external documents like the
`provider-metadata.json` files.

As regular operation the `isdubad` daemon does the downloading
in the background.

This combination may be misused as a scanning device in form of blind
[Server Side Request Forgery (SSRF)](https://owasp.org/www-community/attacks/Server_Side_Request_Forgery).
_Blind_ because users may see that those scanning requests for CSAF contents
on other ports fail, but do not get the contents back.

To reduce the risk, `isdubad` comes with a predefined set of rules which
IP adresses to block. Disallowed are typical internal network addresses
and localhost.  [See](./isdubad-config.md#section_general) for details.
If you need a connection to an internal service, for example when
running a provider that ISDuBA shall access,
you must whitelist the IP address in that configuration.

# Backup

As generally recommended with any IT-system,
keep a backup of your ISDuBA instance's data and configuration.
This includes a backup of the external identity management system.

As documents and comments can be permanently deleted from the system,
carefully consider the backup strategy.
For instance make sure that in addition to incremental backups,
you have the ability to restore a full backup, to be able to restore an earlier state.

It is recommended to store the `aes_key`, which is specified in the
ISDuBA backend configuration, in a safe location. The used PostgreSQL database
contains all relevant application data, aside from those managed via the identity management system. 
The PostgreSQL data can be backed up using standard PostgreSQL tools.
For more information and backup strategies see the official
PostgreSQL documentation for backups and restoration:
<https://www.postgresql.org/docs/current/backup.html>

# Docker/Container setup

This repo contains guides for docker compose and development setups. 
Within these, all components are configured for ease of setup, which leads to security risks.
It is strongly recommended to use the production setup for production usage.

## Docker maintenance

### Keycloak

Keycloak starts in development mode and the configuration does not persist
across rebuilds. It is preferred to use a Keycloak instance that
is already managed. The Keycloak setup container executes the script located in
`docker/keycloak/init.sh`. This script is used for automated testing; it is a
template on how to set up Keycloak and should not be used to automate the
production setup.

### Database

It is recommended to run the database outside the container or mount
`/var/lib/postgresql/data` to a persistent store. If the database is accessed
through a network, the default passwords should be changed.

