<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

# Operations

General hints towards operations.

## Sections

- [Security](#section_security)
- [Performance](#section_performance)

### <a name="section_security">Security</a>

As a precaution, place the backend machine that runs `isdubad`
in a network setup that it does not have access to internal services.

ISDuBA is built to download CSAF documents from the internet.
The places where these are searched for are configured by users
with role `source-manager` and external documents like the
`provider-metadata.json` files.

As regular operation the `isdubad` daemon does the downloading
in the background.

This combination may be misused as a scanning device in form of blind
[Server Side Request Forgery (SSRF)](https://owasp.org/www-community/attacks/Server_Side_Request_Forgery).
_Blind_ because users may see that those scanning requests for CSAF contents
on other ports fail, but do not get the contents back.

To reduce the risk, `isdubad` comes with a predefined set of rules which
IP adresses to block. Disallowed are typical internal network addresses
and localhost.  [See](./isdubad-config.md#section_general) for details.
If you need a connection to an internal service, for example when
running a provider that ISDuBA shall access,
you must whitelist the IP address in that configuration.

### <a name="section_performance">Performance</a>

The performance of `isdubad` is mainly influenced by the
configuration of the PostgreSQL database. Most distributions of GNU/linux
pre-configure it quiet conservatively, using a minimum of RAM and assuming
it to be run on classical hard disk drives (HDDs).
To improve performance you should run `isdubad` on a system with
a reasonable (at least 8GiB) amount of memory and solid-state drives (SSDs).
You should consider using tools like [PGTune](https://pgtune.leopard.in.ua/)
to fine-tune your PostgreSQL installation.

One observation that was made during the development is that 
improvements can be made by setting `random_page_cost = 1.0`.
In its current state (PG15) the
query planner of PostgreSQL does not know how many data is already in
memory and estimates the costs of loading them incorrectly.
This affects the use of indices in particular. Without the mentioned flag
it often chooses to use linear scans, resulting in significant
slowdown in e.g. searching.
