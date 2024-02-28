-- This file is Free Software under the MIT License
-- without warranty, see README.md and LICENSES/MIT.txt for details.
--
-- SPDX-License-Identifier: Apache-2.0
--
-- SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
-- Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

DROP VIEW extended_documents;

CREATE FUNCTION revision_history_length(jsonb) RETURNS int AS $$
    SELECT jsonb_array_length(jsonb_path_query($1, '$.document.tracking.revision_history'))
$$ LANGUAGE SQL IMMUTABLE;

ALTER TABLE documents ADD COLUMN tlp text
    GENERATED ALWAYS AS (document #>> '{document,distribution,tlp,label}') STORED;

ALTER TABLE documents ADD COLUMN title text
    GENERATED ALWAYS AS (document #>> '{document,title}') STORED;

ALTER TABLE documents ADD COLUMN rev_history_length int
    GENERATED ALWAYS AS (revision_history_length(document)) STORED;

CREATE VIEW extended_documents AS SELECT
    *,
    (SELECT max(a::float) FROM
        jsonb_path_query(
            document, '$.vulnerabilities[*].scores[*].cvss_v2.baseScore') a)
        AS cvss_v2_score,
    (SELECT max(a::float) FROM
        jsonb_path_query(
            document, '$.vulnerabilities[*].scores[*].cvss_v3.baseScore') a)
        AS cvss_v3_score,
    (jsonb_path_query_array(
        document, '$.vulnerabilities[0 to 3]."cve"')) AS four_cves
    FROM documents;

GRANT SELECT ON extended_documents TO {{ .User | sanitize }};
