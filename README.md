<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

# ISDuBA

**Work in progress** **early-development**

Plan: Develop a web application
for downloading and evaluating security advisories in CSAF 2.0 format
for internal use of a team responsible for a number of IT related topics.

Using the following components:
 * Backend to be written in Go
 * PostgreSQL as database
 * [keycloak](https://www.keycloak.org/) as identify provider
 * docker-compose setup example (planned)
 * [svelte-flowbite](https://flowbite-svelte.com/)
     for the single page web application frontend
 * [csaf_distribution](https://github.com/csaf-poc/csaf_distribution)
     for downloading advisories
 * [csaf_webview](https://github.com/csaf-poc/csaf_webview)
     for viewing documents
 * https://github.com/CERTCC/SSVC/tree/main/docs/ssvc-calc (or alternative)
 * https://github.com/rtfpessoa/diff2html (or alternative)


## How to get started
 * [Setup](docs/setup.md)


## What does the name ISDuBA mean?

The abbreviation expands to a German label, which translates to
  **I**nternal **s**ystem for **d**ownloading and evaluating **a**dvisories.

## License

- `ISDuBA` is licensed as Free Software under the [Apache License, Version 2.0](./LICENSES/Apache-License-2.0.txt)

- See the specific source files
  for details, the license itself can be found in the directory `LICENSES/`.

- Contains third party Free Software components under licenses that to our best knowledge are compatible at time of adding the dependency, [3rdpartylicenses.md](3rdpartylicenses.md) has the details.

- The copyright is held by the Bundesamt für Sicherheit in der Informationstechnik (BSI)
