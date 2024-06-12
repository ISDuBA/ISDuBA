-- This file is Free Software under the Apache-2.0 License
-- without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
--
-- SPDX-License-Identifier: Apache-2.0
--
-- SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
-- Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

ALTER TABLE documents ADD COLUMN critical float GENERATED ALWAYS AS (
    coalesce(max_cvss3_score(document), max_cvss2_score(document))) STORED;

CREATE INDEX documents_critical_idx ON documents(coalesce(critical, '0'::double precision) DESC);
