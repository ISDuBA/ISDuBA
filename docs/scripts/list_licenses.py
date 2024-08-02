#!/usr/bin/env python3
"""Extract licensing info from the packages of an SPDX-2.3 SBOM JSON file.

An experimental example for the output of Github Action
anchore/sbom-action@v0 with: format: spdx-json output.

Give a filename to that output as parameter.


SPDX-License-Identifier: Apache-2.0

SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
"""
import json
import sys

if len(sys.argv) != 2:
    sys.exit(__doc__)

with open(sys.argv[1], 'rt', encoding="utf-8") as file:
    sbom = json.load(file)

    if sbom["SPDXID"] != "SPDXRef-DOCUMENT" or \
       sbom["spdxVersion"] != "SPDX-2.3":
        sys.exit("Not an SPDX-2.3 SBOM file")

    for p in sbom["packages"]:
        l = p["licenseConcluded"]
        if l == "NOASSERTION":
            l = p["licenseDeclared"]

        print("{},{}".format(p["name"], l))
