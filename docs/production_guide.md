<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

(Notes about running ISDuBA in Production)

## Backup

As generally recommended with any IT-system,
keep a backup of your ISDuBA instance's data and configuration.
This includes a backup of the external identity management system.

As documents and comments can be permanently deleted from the system,
consider your auditing needs for the backup strategy.
For instance make sure that in addition to incremental backups,
you have the ability to restore a full backup often enough for your
auditing needs.

It is recommended to store the `aes_key` that is specified in the
ISDuBA backend configuration, in a safe location. The used PostgreSQL database
besides the external identity management system, contains all relevant
application data. This data can be backed up using standard PostgreSQL tools.
For more information and backup strategies see the official
PostgreSQL documentation for backup and restore:
<https://www.postgresql.org/docs/current/backup.html>

## Audit work which was done with ISDuBA

In some cases there is a need to examine how a piece of security
information has been handled by the organisation
(that uses an ISDuBA instance).

This is what the role `auditor` is for.
It allows to see documents, comments, events and protocol data.

Documents shall be set to state `archived` when they have
been worked upon. This will allow later examination by an `auditor`.

Once an admin confirms that a document can be `deleted`,
the documents itself and the associated comments and other information
will be removed permanently from the current database.
(This is required by the workflow and additionally partly protects
personal data of users that interacted with the system.)

Adjust your operational instructions to keep documents in state `archived`
that you expect reasonably to be audited in the near future.

Then it will be rare case to get an audit request for an elder situation.
In this situation we recommend to restore a full backup
before the delete into a stand-a-lone system. Also need will be
a Keycloak with the users that had access to ISDuBA at this time
and the additional auditor users. The access to the system should be
tightly restricted to the auditing work.

## Docker/Container setup

This repo contains an example docker compose setup. All components are configured for ease of
setup, it is recommended to use the production setup for production usage.

### Keycloak

This setup starts Keycloak in a development mode and does not persist
configuration across rebuilds. It is preferred to use a Keycloak instance that
is already managed. The Keycloak setup container executes the script located in
`docker/keycloak/init.sh`. This script is used for automated testing; it is a
template on how to set up Keycloak and should not be used to automate the
production setup.

### Database

It is recommended to run the database outside the container or mount
`/var/lib/postgresql/data` to a persistent store. If the database is accessed
through a network, the default passwords should be changed.

### Application

The application does not provide a secure endpoint. It should be put behind a
reverse proxy with a valid TLS certificate.
