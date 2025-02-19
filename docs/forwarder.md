<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->


# Forwarder

## Configuration
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
If set to false only that are manually forwarded are sent to the URL. This

`name` specifies the name of the forward target. This is required to let the
user distinguish between multiple targets.

## Filtering
Currently it is only possible to select which publisher should be forwarded. If
no publisher is configured all documents are uploaded.


## Forward request
The forwarder sends a POST request to the specified URL. The data is encoded
with `multipart/form-data`.
The form data contains the following fields:
- `advisory`: The JSON document that is forwarded.
- `validation_status`: The validation status of the document. This can be
`valid`, `invalid` or `not_validated`.

## Error handling
If the response to the forward request was `201` it will store the document as
successfully forwarded for this URL. If request failed it will retry to forward
the document at the next poll interval.
