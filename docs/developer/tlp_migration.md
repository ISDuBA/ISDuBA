<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->

## Migrating to TLP 2.0 labels

ISDuBA uses the TLP 2.0 labels `CLEAR`, `GREEN`, `AMBER`, `AMBER+STRICT` and `RED`
for access rights now. CSAF 2.0 documents with the `WHITE` tlp label are still
accessible to users with `CLEAR` rights.


### TLP:WHITE alias for TLP:CLEAR

We removed the temporary mapping in `pkg/models/tlps.go`, which would expand
the `WHITE` rights by `CLEAR` and `AMBER` by `AMBER+STRICT`.

During advisory import the label `WHITE` is now aliased for the `CLEAR` access
right in the `PublishersTLPs.Allowed` method. This is because during import
we check the label of the actual document.

Stored advisories keep the original document label, raw bytes and signature,
however, we adjusted the `tlp` column of the `documents` table by altering the
expression to whenever the label is `WHITE` storing `CLEAR` in its place (see
`000-setup.sql` and `044-white_to_clear.sql` under `ISDuBA/pkg/database/migrations/`
for said changes).
This normalized internal representation of the document label is used during search
listing and access checks. This way stored CSAF 2.0 `WHITE` documents are accessible
for users with `CLEAR` right.


### Only TLP 2.0 rights are intended

Since `CLEAR` replaces `WHITE` in TLP 2.0, the group of tlp `WHITE` and all instances of
users having said groups rights need to be replaced with `CLEAR` in Keycloak,
otherwise users with `WHITE` in the `TLP` claim are unable to log in.
`WHITE` rights are also no longer accepted in `publishers_tlps` configurations.

With the removed mapping, having `AMBER` rights no longer automatically gives
access to `AMBER+STRICT`. It needs to be added as a separate group and
to the respective users (see [the realm file](../scripts/keycloak/isduba-realm.json) for
example configuration) and source manager needs it explicitly in `publishers_tlps` for import
(see [example configuration](../example_isdubad.toml)).


### Upgrading an existing installation

1. Run the database migration to change `documents.tlp` column for already existing
   CSAF 2.0 documents. Depending on the amount of entries, this can take some time.
2. Replace `WHITE` by `CLEAR` in Keycloak groups and in `publishers_tlps` config
   settings. This is important for the hard transition.
3. Add `AMBER+STRICT` to Keycloak groups and to `publishers_tlps`.
4. Check stored queries for `WHITE`, they no longer match any advisory:
   `SELECT id, name, query FROM stored_queries WHERE query LIKE '%WHITE%';`

