---- This file is Free Software under the Apache-2.0 License
-- without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
--
-- SPDX-License-Identifier: Apache-2.0
--
-- SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
-- Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

-- create_advisory checks if the new document is newer than the old one.
CREATE OR REPLACE FUNCTION create_advisory() RETURNS trigger AS $$
    DECLARE
        old_id           int;
        old_rev_length   int;
        old_release_date timestamptz;
    BEGIN
        -- Ensure having an advisories record.
        INSERT INTO advisories (tracking_id, publisher)
            VALUES (NEW.tracking_id, NEW.publisher)
            ON CONFLICT (tracking_id, publisher) DO UPDATE SET state = 'new';

        SELECT id, rev_history_length, current_release_date
            INTO old_id, old_rev_length, old_release_date
            FROM documents
            WHERE latest AND tracking_id = NEW.tracking_id AND publisher = NEW.publisher;

        IF NOT FOUND THEN -- No latest -> we are
            UPDATE documents SET latest = TRUE WHERE id = NEW.id;
        ELSE
            -- Check if the new record is in fact newer than the old one.
            IF NEW.current_release_date > old_release_date OR
               (NEW.current_release_date = old_release_date AND
                NEW.rev_history_length > old_rev_length)
            THEN
                -- Take over lead.
                UPDATE documents SET latest = FALSE WHERE id = old_id;
                UPDATE documents SET latest = TRUE  WHERE id = NEW.id;
            END IF;
        END IF;
        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;
