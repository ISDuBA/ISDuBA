<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

# ISDuBA

**Work in progress** **pre-alpha**

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

## Free Software

ISDuBA is Free Software licensed under the terms of the [Apache License, Version 2.0](./LICENSES/Apache-2.0.txt).
