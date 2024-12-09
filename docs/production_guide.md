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
This includes a backup of the external identify management system.

As documents and comments can be permanently deleted from the system,
consider your auditing needs for the backup strategy.
For instance make sure that in addition to incremental backups,
you have the ability to restore a full backup often enough for your
auditing needs.

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

## Baremetal setup

### Building

To build the application the latest Golang version, NodeJS 20 and standard build
tools, like GNU Make are required. At the root of the repo run `make dist` to
build the frontend and backend; this will result in a Tar-file inside `dist/`
that contains the application. It can also be useful to look inside the
Dockerfile of the application to see how individual components of the
application can be built. If if there are no special requirements it can be
enough to use already built Tar-file from the release page.

### Running

The Tar-file can be copied and extracted on a production server. This file
contains the `isdubad` backend, which can be run on any modern amd64 linux
system and the frontend which is contained in the `web/` folder. No further
dependencies are required to start the application. By default the backend will
serve the contents of `web/`.

For a quick start copy `example_isdubad.toml` to the
folder where the application is conatined. Configure the postgres and keycloak
settings and rename the file to `isduba.toml`. The application can now be
started. For exposing the application to the network it is recommended to use
a TLS-terminating reverse proxy.

## Docker/Container setup

This repo conatins an example docker compose setup. The isduba container is
built for production usage. All other components are configured for ease of
setup, it is recommended to use another configuration or setup own
containers/server for production usage.

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
