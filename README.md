<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

# ISDuBA

**Work in progress** **early-development**

Plan: Develop a web application
for downloading and evaluating security advisories in CSAF 2.0 format
for internal use of a team responsible for a number of IT related topics.

Will be published as Free Software.

Envisoned is to use the following components:
 * Backend to be written in Go using
 * PostgreSQL as database
 * keycloak as identify provider
 * docker-compose setup example
 * [svelte-flowbite](https://flowbite-svelte.com/)
     for the single page web application frontend
 * [csaf_distribution](https://github.com/csaf-poc/csaf_distribution)
     for downloading advisories
 * [csaf_webview](https://github.com/csaf-poc/csaf_webview)
     for viewing documents
 * https://github.com/CERTCC/SSVC/tree/main/docs/ssvc-calc
 * https://github.com/rtfpessoa/diff2html


## How to get started
 * [Setup](docs/setup.md)


## What does the name ISDuBA mean?

The abbreviation expands to a German label, which translates to
  **I**nternal **s**ystem for **d**ownloading and evaluating **a**dvisories.

## Free Software

ISDuBA is Free Software licensed under the terms of the [Apache License, Version 2.0](./LICENSES/Apache-License-2.0.txt).
