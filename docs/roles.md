<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2025 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2025 Intevation GmbH <https://intevation.de>
-->

The actions any user can take depend on their role. Any user can have any combination of roles.
This also determines how workflowstates of advisories can be
set, [see this diagram for reference](./images/workflow.svg)
There are six roles:
 * [admin](#admin)
 * [auditor](#auditor)
 * [editor](#editor)
 * [importer](#importer)
 * [reviewer](#reviewer)
 * [source-manager](#source-manager)

The following is an overview of all roles within ISDuBA.

### admin
The `admin` role manages stored queries and is the role that can delete advisories that are set to delete.

### auditor
The `auditor` role represents users that may want to
audit how certain security information or advisories have been
handled by the organisation using an ISDuBA instance.

This role allows viewing of documents, comments, events and protocol data.

To make auditing easier, documents shall be set to state `archived` when they have
been worked upon.

Documents can also be marked for deletion.
Once an admin confirms that a document can be `deleted`,
the documents itself and the associated comments and other information
will be removed permanently from the current database.
(This is required by the workflow and additionally partly protects
personal data of users that interacted with the system.)

Adjust your operational instructions to keep documents in state `archived`
that you expect reasonably to be audited in the near future.

There may be rare cases where an audit is requested for an older state of the application.
In this situation we recommend to restore a full backup. How to create a backup is outlined
within [the security considerations' backup section](./security_considerations#docker/container_setup)

### editor
The `editor` role represents the advisory-evaluating user. They can set an SSVC.

### importer
The `importer` role is responsible for uploading advisories. It is expected to be used by tools that require
authentification.

### reviewer
The `reviewer` role reviews the outcome of the [editor's](#editor) evaluation and can then set the advisory
to archived (or delete).

### source-manager

The `source-manager` role manages sources, meaning which advisories are downloaded from where.
