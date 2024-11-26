<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

# Operations

## Security

ISDuBA is built to download CSAF documents from the internet.
The places where these are searched for are configured by users
with role `source-manager` and external documents like the
`provider-metadata.json`.

As regular operation the `isdubad` daemon does the downloading
in the backgroud.

This combination may be misused as a scanning device in form of blind
[Server Side Request Forgery (SSRF)](https://owasp.org/www-community/attacks/Server_Side_Request_Forgery).
_Blind_ because users may see that those scanning requests for CSAF things
fail, but do not get the contents back.

To reduce the risk, `isdubad` comes with a predefined set of rules which
IP adresses are blocked. Disallowed are typical internal network addresses
and localhost.  [See](./isdubad-config.md#section_general) for details.
If you need a connection to an internal service, for example when
running a provider that ISDuBA shall access,
you need to whitelist the IP address.
