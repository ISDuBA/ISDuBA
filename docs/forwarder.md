<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->


# Forwarder configuration

- [`[forwarder]`](#global) Global
- [`[[forwarder.target]]`](#target) Target

## <a name="global"></a> `[forwarder]` Global

- `update_interval`: Secifies how often the database is checked for new documents. Defaults to `"5m"`.
- `strategy`: Filtering strategy. See [Filtering](#filtering) for details. Defaults to `"all"`.

If `external_url` in [`[web]`](./example_isdubad.toml#section_web) is configured
this URL is postfixed with `/api/documents/{id}` (with `id` being the internal ISDuBA id of the document) and is send to the targets to enable them to download the document over the API of the ISDuBA server. Set this to the URL where your ISDuBA server is reachable. The documents will be sent either way. Defaults to not set.

## <a name="target"></a> `[[forwarder.target]]` Target
The forwarder needs at least one target to get started.
A target is an external endpoint where the documents are send to.
An example implementation of such a target can befound [here](https://github.com/gocsaf/forwardertarget).
The target are fed in intervals. This can be configured with `update_interval`.
This To forward documents that are stored in the database, a forwarder target needs
to be configured.
Only documents that are successfully imported into the database are forwarded.

- `automatic`: Specifies if the target automatically receives new documents. If disabled the target only receives documents on manual forwarding.
- `url`: The URL of the forward target and has to be unique for all the forwarder targets.
- `name`: The name of target. This value will be displayed on manual forwarding the document.
- `publisher`: Specifies the publisher of the documents that need to be forwarded.
- `header`: List all headers that are sent to the target. The format is `key:value`.
- `private_cert`: The location of the private client certificate.
- `public_cert`: The location of the public client certificate.
- `timeout`: Sets the http client timeout. Set this value if the network is unstable.

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

The backend polls for documents that were not uploaded
for the configure URL and forwards them. Already imported documents are also forwarded.
If set to false only those that are manually forwarded are sent to the URL.

## <a name="filtering"></a> Filtering
The first level of filtering is the `publisher`. If specified
only the documents for the given publisher are forwarded to
the target endpoint.
If no publisher is configured all documents are forwarded.
The second level is the `strategy`. This determines which documents are forwarded.
They are currently two strategies `all` and `new_major`.
`all` implies all documents are forwarded, `new_major` sends new advisories, all not draft versions.
If you have semantical versioned documents new are documents are forwarded if they
are a major change in comparison to the former once.

You can adjust the strategy for each forwarder target. If no target specific strategy
is given the global strategy will be used.

## Forward request
The forwarder sends a POST request to the specified URL. The data is encoded
with `multipart/form-data`.
The form data contains the following fields:
- `advisory`: The JSON document that is forwarded.
- `validation_status`: The validation status of the document. This can be
`valid`, `invalid` or `not_validated`.
- `document_url`: The API endpoint URL where to download the document from. Only send if
`[web]`/`external_url` is configured (see above).

## Error handling
If the response to the forward request is `201`, then the document will be
recorded as successfully forwarded for the URL.
If the request fails, ISDuBA retries at the next poll interval.

## Architecture

The forwarding subsystem consists of three parts.
The central component is the **Manager**. It reacts to direct
upload request from the API and triggers the none automatic forwarders.
It also gets signaled by the **Poller** if new documents are
integrated into the database. The poller checks for new
documents at the rate configured by `update_interval`.
The manager filters the list of new documents given by the
poller by publisher and strategy configured for the automatic
targets and writes upload requests in a database queue for each
**Forwarder**. These forwarders poll from their respective
upload requests and try to forward the documents to the
configured target URLs. The results of these upload attempts are
stored back in the queue. If they were successfull the
document is never forwarded again by this forwarder.
Same holds in the case of explicit rejection by the endpoint.
If there is e.g. a network error the not uploaded documents
stay in a pending state and are tried to upload later again.

![Architecture text](./images/forwarder.svg)
