-- This file is Free Software under the MIT License
-- without warranty, see README.md and LICENSES/MIT.txt for details.
--
-- SPDX-License-Identifier: Apache-2.0
--
-- SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
-- Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

DROP VIEW extended_documents;

CREATE VIEW extended_documents AS SELECT
    *,
    (document #>> '{document,title}')                  AS title,
    (document #>> '{document,distribution,tlp,label}') AS tlp,
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
