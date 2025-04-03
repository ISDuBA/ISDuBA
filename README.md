<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

# ISDuBA

A web application
for downloading and evaluating security advisories in the CSAF 2.0 format.
ISDuBA is designed to support teams that are responsible
for the IT security of a group of products.

We appreciate your problem reports, please check the list of issues first.

To gain insight into the architecture of ISDuBA take a look at the
[overview diagram](docs/images/overview-diagram.svg).


Components used by ISDuBA:

- [PostgreSQL](https://www.postgresql.org/) as database
- [keycloak](https://www.keycloak.org/) as identity provider
- [svelte-flowbite](https://flowbite-svelte.com/)
  and TypeScript for the single page web application frontend
- Go as programming language for the backend.
- a downloading kernel that is close to
  [gocsaf](https://github.com/gocsaf/csaf)
- an extended version of
  [csaf_webview](https://github.com/csaf-poc/csaf_webview)

## How to get started

- [Setup](docs/README.md)
- [Important Security Information](docs/security_considerations.md)
- [Troubleshooting](docs/troubleshooting.md)

## What does the name ISDuBA mean?

The abbreviation expands to **I**nternes **S**ystem zum **D**ownload **u**nd
  der **B**ewertung von **A**dvisories, a german label which translates to
`Internal system for downloading and evaluation of advisories`.

## License

ISDuBA is Free Software.

Source code written for ISDuBA was placed under the
[Apache License, Version 2.0](./LICENSES/Apache-2.0.txt).

```
 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
```

ISDuBA depends on third party Free Software components which have their
own right holders and licenses. To our best knowledge
(at the time when they were added)
the dependencies are upwards compatible with the ISDuBA main license.

### Dependencies

The top level dependencies can be seen from

- [go.mod](./go.mod) for the `isdubad` backend and server tools.
- [package.json](./client/package.json) for the web application frontend.
- The build and setup descriptions (linked above).

Use one of several available Free Software tools to examine indirect
dependencies and get a more complete list of component names and licenses.

For example use the SPDX-2.3 SBOM json file coming with an ISDuBA release
or use <https://github.com/anchore/syft> to create one.
Then run [list_licenses.py](./docs/scripts/list_licenses.py) on it
or `python3 -m json.tool`, to see more.
