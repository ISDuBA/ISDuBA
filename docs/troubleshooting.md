<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2025 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2025 Intevation GmbH <https://intevation.de>
-->



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

### Check whether `isdubad` is correctly installed
The following will define a `TOKEN` variable which holds the information
about a user with name `USERNAME` and password `USERPASSWORD`
as configured in keycloak.

(You can check the `TOKEN` via jwt.io. Keycloak must be running and available.)

```sh
TOKEN=`curl -d 'client_id=auth'  -d 'username=USERNAME' -d 'password=USERPASSWORD' -d 'grant_type=password' 'http://127.0.0.1:8080/realms/isduba/protocol/openid-connect/token' | jq -r .access_token`
echo $TOKEN
```
