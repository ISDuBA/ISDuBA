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
The places where these are searched are configured by users and the downloading
itself is done by the `isdubad` daemon in the backgroud.
This combination may be misused as a scanning device in form of blind
[Server Side Request Forgery (SSRF)](https://owasp.org/www-community/attacks/Server_Side_Request_Forgery).
To reduced the risk `isdubad` comes with a predefined set of rules which URLs are
allowed to be visited. [See](./isdubad-config.md#section_general) for details.

