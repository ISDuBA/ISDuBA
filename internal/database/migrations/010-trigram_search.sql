-- This file is Free Software under the Apache-2.0 License
-- without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
--
-- SPDX-License-Identifier: Apache-2.0
--
-- SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
-- Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

DROP INDEX documents_texts_ts_idx;
DROP INDEX comments_ts_idx;

ALTER TABLE unique_texts DROP COLUMN ts;
ALTER TABLE comments     DROP COLUMN ts;

DROP FUNCTION to_tsvector_multilang;

CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX ON unique_texts USING gin(txt gin_trgm_ops);
CREATE INDEX ON comments     USING gin(message gin_trgm_ops);
