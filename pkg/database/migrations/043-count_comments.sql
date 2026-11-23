-- This file is Free Software under the Apache-2.0 License
-- without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
--
-- SPDX-License-Identifier: Apache-2.0
--
-- SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
-- Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

CREATE OR REPLACE FUNCTION incr_comments() RETURNS trigger AS $$
    BEGIN
        UPDATE advisories
            SET comments = comments + 1
        WHERE
            id = (SELECT advisories_id
                FROM documents
                WHERE id = NEW.documents_id);
        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

UPDATE advisories a
SET comments = (
    SELECT COUNT(*)
    FROM documents d
    JOIN comments c ON c.documents_id = d.id
    WHERE d.advisories_id = a.id
);
