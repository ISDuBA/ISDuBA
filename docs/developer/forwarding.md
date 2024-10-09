<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

There are two use cases about forwarding documents
to a different application.


## Use Case FW1: other specialised system

There is a specialised system that is separate from ISDuBA.
It also manages CSAF documents.
Those documents are a subset of all documents that get into the ISDuBA system.

They shall be forwarded automatically and
can happen right after a wanted document is within ISDuBA.


## Use case FW2: manually send to asset matching

Once a document is considered a base for further action,
an ISDuBA user shall be able to manually forward it to a different
system for asset matching.


## technical considerations

Each receiving point (or _target_) is an HTTP based endpoint
secured by TLS and optionally by a client certificate or
an HTTP header with an API token for access.

* The number of recieving systems is small and changes rarely.
* It is okay for the _source-manager_ role to manage.

The targets can be setup by a configuration file. It has to
specify publishers, credentials, and an URL for each target.
In the client the user interface must offer the user the possibility
to select documents and to select targets to forward the documents.

An open question is if the documents to forward have to fulfill
the same criteria of quality as for the downloader.
