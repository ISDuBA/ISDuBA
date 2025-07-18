<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

The ```bulkimport```-tool  allows the automated import of one or multiple advisories into an ISDuBA database.


Usage:
```bulkimport [OPTIONS] dest ```

Where dest is either:

 * A single advisory file to import, or

 * A directory containing advisories directly or within subdirectories.

with the following supported options:

```
  -continue
       continue bulkimport even if an advisory was not imported successfully
  -database string
       database name (default "isduba")
  -delete
       delete successfully imported advisories
  -dry
       dont store values
  -host string
       database host (default "localhost")
  -importer string
       importing person (default "root")
  -move string
       move unsuccessfully imported advisories to this folder (create folder if it does not exist)
  -password string
       password (default "isduba")
  -port int
       database host (default 5432)
  -user string
       database user (default "isduba")
  -version
       show version information
```
