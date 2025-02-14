-- This file is Free Software under the Apache-2.0 License
-- without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
--
-- SPDX-License-Identifier: Apache-2.0
--
-- SPDX-FileCopyrightText: 2025 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
-- Software-Engineering: 2025 Intevation GmbH <https://intevation.de>

CREATE TYPE tracking_status AS ENUM (
    'draft', 'final', 'interim');

CREATE FUNCTION text_to_tracking_status(text) RETURNS tracking_status AS $$
    SELECT $1::tracking_status
$$ LANGUAGE SQL IMMUTABLE;

ALTER TABLE documents ADD COLUMN status tracking_status
    GENERATED ALWAYS AS (
    text_to_tracking_status(document #>> '{document,tracking,status}')) STORED;

ALTER TABLE documents
    DROP CONSTRAINT documents_advisories_id_version_rev_history_length_key,
    ADD UNIQUE (advisories_id, version, rev_history_length, status);

