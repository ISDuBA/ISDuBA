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

Some files in the repository are machine generated:

| generation result | input |
|-------------------|-------|
| [docs.go](../pkg/web/docs/docs.go)           | `pkg/web/*.go` |
| [swagger.json](../pkg/web/docs/swagger.json) | "              |
| `client/src/lib/workflow.ts` | `pkg/models/workflow.go` |
| `docs/images/workflow.svg`   | "                        |

If you change any of the input files
please use `go generate ./...` in the root folder
and commit updates results to the repository.

Regeneration requires `swaggo`,
to be installed via `go install github.com/swaggo/swag/cmd/swag@latest`.
This component will update the OpenAPI 2.0 documentation.

If you plan to add further machine generated files ensure that they
are marked with comments like
```
// THIS FILE IS MACHINE GENERATED. EDIT WITH CARE!
```
or ```DO NOT EDIT```
.
