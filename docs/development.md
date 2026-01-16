<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2025 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2025 Intevation GmbH <https://intevation.de>
-->

# Development

Find some design considerations and history in `developer/`.
Not everything was implemented as it was proposed in there.


## Generated files

Some source code files are machine generated. At the moment is only:

* [docs.go](../pkg/web/docs/docs.go)

If you change the source files please regenerate the generated files
with `go generate ./...` in the root folder and add the updated files
to the version control. This will update the swagger documentation.

If you plan to add further machine generated files ensure that they
are marked with comments like
```
// THIS FILE IS MACHINE GENERATED. EDIT WITH CARE!
```
or ```DO NOT EDIT```
.
