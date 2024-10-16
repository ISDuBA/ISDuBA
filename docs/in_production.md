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
