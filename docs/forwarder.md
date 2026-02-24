<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->


# Forwarder

The forwarder needs at least one target to get started.
A target is an external endpoint where the documents are send to.
An example implementation of such a target can befound [here](https://github.com/gocsaf/forwardertarget).

The target are fed in intervals. This can be configured with `update_interval`.
This secifies how often the database is checked for new documents. Defaults to `"5m"`.

To forward documents that are stored in the database, a forwarder target needs
to be configured.
An example configuration can look like this:
```TOML
[[forwarder.target]]
automatic = true
name = "First example target"
url = "http://example.com/api/v1/import"
publisher = "publisher-name"
header = [ "x-api-key:secret" ]
private_cert = "private-cert-file"
public_cert = "public-cert-file"
timeout = "5s"
```

If the `automatic` is set to true all documents that match the filter are
forwarded to endpoint. The backend polls for documents that were not uploaded
for the URL and forwards them. Already imported documents are also forwarded.
If set to false only those that are manually forwarded are sent to the URL.

`url` has to be unique for all the forwarder targets.

## Filtering
The first level of filtering is the `publisher`. If specified
only the documents for the given publisher are forwarded to
the target endpoint.
no publisher is configured, all documents are forwarded.
The second level is the `strategy`. This determines which documents are forwarded.
They are currently two strategies `all` and `new_major`.
As `all` implies all documents are forwarded, `new_major` sends new advisories, all not draft versions.
If you have semantical versioned documents new are documents are forwarded if they
are a major change in comparison to the former once.

The default strategy is `all`. You can adjust the strategy for each
forwarder target. If not specified there the global strategy from
forwarder is used.

## Forward request
The forwarder sends a POST request to the specified URL. The data is encoded
with `multipart/form-data`.
The form data contains the following fields:
- `advisory`: The JSON document that is forwarded.
- `validation_status`: The validation status of the document. This can be
`valid`, `invalid` or `not_validated`.

## Error handling
If the response to the forward request is `201`, then the document will be
recorded as successfully forwarded for the URL.
If the request fails, ISDuBA retries at the next poll interval.
