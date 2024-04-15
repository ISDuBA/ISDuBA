-- This file is Free Software under the Apache-2.0 License
-- without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
--
-- SPDX-License-Identifier: Apache-2.0
--
-- SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
-- Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

DROP VIEW extended_documents;

-- generator functions

CREATE FUNCTION max_cvss2_score(jsonb) RETURNS float AS $$
    SELECT max(a::float) FROM
        jsonb_path_query(
            $1, '$.vulnerabilities[*].scores[*].cvss_v2.baseScore') a
$$ LANGUAGE SQL IMMUTABLE;

CREATE FUNCTION max_cvss3_score(jsonb) RETURNS float AS $$
    SELECT max(a::float) FROM
        jsonb_path_query(
            $1, '$.vulnerabilities[*].scores[*].cvss_v3.baseScore') a
$$ LANGUAGE SQL IMMUTABLE;

CREATE FUNCTION first_four_cves(jsonb) RETURNS jsonb AS $$
    SELECT jsonb_path_query_array(
        $1, '$.vulnerabilities[0 to 3]."cve"')
$$ LANGUAGE SQL IMMUTABLE;

-- add new columns

ALTER TABLE documents ADD COLUMN cvss_v2_score float
    GENERATED ALWAYS AS (max_cvss2_score(document)) STORED;

ALTER TABLE documents ADD COLUMN cvss_v3_score float
    GENERATED ALWAYS AS (max_cvss3_score(document)) STORED;

ALTER TABLE documents ADD COLUMN four_cves jsonb
    GENERATED ALWAYS AS (first_four_cves(document)) STORED;
