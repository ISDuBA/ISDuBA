<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

# SSVC calculator

## Logic: Decision tree

The SSVC calculator of this application uses the decision tree "CISA Coordinator v2.0.3". To change the
logic of the calculator it is necessary to put the JSON with a different decision tree into the folder
[client/src/lib/Advisories/SSVC/](/client/src/lib/Advisories/SSVC/). Then replace `CISA-Coordinator` in the import
of [SSVCCalculator.ts](/client/src/lib/Advisories/SSVC/SSVCCalculator.ts) with the name of the JSON.

## Svelte module

The main component which allows the calculation of the SSVC vector is
[SSVCCalculator.svelte](/client/src/lib/Advisories/SSVC/SSVCCalculator.svelte). To replace it with another module
the new component has to implement the exported attributes used in the calculator, e.g. `disabled` because the
parent component decides if a user is allowed to calculate the SSVC.

Furthermore, the component has to throw the event `updateSSVC` to tell the parent when the calculation is finished.
However, the result of the calculation is saved by the module itself by using the endpoint `/api/ssvc/` of the
backend API.