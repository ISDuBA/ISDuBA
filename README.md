<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

# ISDuBA

**Under development - some things already work**

A web application
for downloading and evaluating security advisories in CSAF 2.0 format.
Mainly ISDuBa wants to support teams that are responsible for
the IT security of a group of products.

ISDuBA uses the following components:
 * Go as programming language for the backend.
 * PostgreSQL as database
 * [keycloak](https://www.keycloak.org/) as identify provider
 * docker-compose setup example (planned)
 * [svelte-flowbite](https://flowbite-svelte.com/)
     for the single page web application frontend
 * [csaf_distribution](https://github.com/csaf-poc/csaf_distribution)
     for downloading advisories
 * [csaf_webview](https://github.com/csaf-poc/csaf_webview)
     for viewing documents


## How to get started
 * [Setup](docs/setup.md)


## What does the name ISDuBA mean?

The abbreviation expands to a German label, which translates to
  **I**nternal **s**ystem for **d**ownloading and evaluating **a**dvisories.


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
(at the time when they have been added)
the dependencies are upwards compatible with the ISDuBA main license.

### Dependencies

The top level dependencies can be seen from
 * [go.mod](./go.mod) for the `isduad` backend and server tools.
 * [package.json](./client/package.json) for the web application frontend.
 * The build and setup descriptions (linked above).

Use one of several available Free Software tools to examine indirect
dependencies and get a more complete list of component names and licenses.

For example use the SPDX-2.3 SBOM json file coming with an ISDuBA release
or use https://github.com/anchore/syft to create one.
Then run [list_licenses.py](./docs/scripts/list_licenses.py)
or `python3 -m json.tool` on it to see more.


