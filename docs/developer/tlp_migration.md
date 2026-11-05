<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->

## TLP 2.0 labels: temporary mapping

**Status: temporary.** We currently temporarily map the TLP 2 `CLEAR`
and `AMBER+STRICT` rights to the already present TLP 1 `WHITE` and `AMBER` rights.
This mapping will have to be removed in the future, once a proper migration
to the TLP 2.0 labels throughout ISDuBA and keycloak rights system has been accomplished.

To achive above mapping we extended the `pkg/models/tlps.go` file by a `temporaryTLPMapping` map, which
maps TLP label `WHITE` to `CLEAR` and `AMBER` to `AMBER+STRICT`, to essentially mock TLP 2.1 based
access rights with the already present TLP 1.0 right system we have. In the `mappedTLPs` function
we take a list of given TLPs and extend it by `CLEAR` or `AMBER+STRICT`, if `WHITE` or `AMBER` 
is present. `mappedTLPs` is then used in the `PublisherTLPs.Allowed` and `PublisherTLPs.AsExprPublisher` methods, 
otherwise import, upload and finding of TLP 2 based advisories is not possible.

## TODOs

- removing `temporaryTLPMapping`, `mappedTLPs` and its use in `PublisherTLPs.Allowed` and `PublisherTLPs.AsExprPublisher`
- adjustment of rights in keycloak for TLP 2
- plan for whether or how to handle `WHITE` and `CLEAR` equivalency
- how to handle `AMBER+STRICT` case specifically
- frontend changes (color for `CLEAR` and `AMBER+STRICT`?)



