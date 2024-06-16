-- This file is Free Software under the Apache-2.0 License
-- without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
--
-- SPDX-License-Identifier: Apache-2.0
--
-- SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
-- Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

ALTER TABLE advisories ADD COLUMN comments int NOT NULL DEFAULT 0;
ALTER TABLE advisories ADD COLUMN recent timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP;

ALTER TABLE advisories ADD CONSTRAINT advisories_comments_check CHECK (comments >= 0);

CREATE INDEX events_log_time_idx ON events_log(time);
CREATE INDEX advisories_recent_idx ON advisories(recent);

UPDATE advisories a SET comments = (
    SELECT count(*) FROM comments c JOIN documents d ON c.documents_id = d.id
    WHERE d.publisher = a.publisher AND d.tracking_id = a.tracking_id);

UPDATE advisories a SET recent = (
    SELECT max(time) from events_log el JOIN documents d ON el.documents_id = d.id
    WHERE d.publisher = a.publisher AND d.tracking_id = a.tracking_id);

CREATE FUNCTION incr_comments() RETURNS trigger AS $$
    DECLARE
        p text;
        t text;
    BEGIN
        SELECT publisher, tracking_id
            INTO p, t
            FROM documents
            WHERE id = NEW.documents_id;
        IF FOUND THEN
            UPDATE advisories
                SET comments = comments + 1
                WHERE publisher = p AND tracking_id = t;
        END IF;
        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE FUNCTION decr_comments() RETURNS trigger AS $$
    DECLARE
        p text;
        t text;
    BEGIN
        SELECT publisher, tracking_id
            INTO p, t
            FROM documents
            WHERE id = OLD.documents_id;
        IF FOUND THEN
            UPDATE advisories
                SET comments = greatest(0, comments - 1)
                WHERE publisher = p AND tracking_id = t;
        END IF;
        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER increment_comments
    AFTER INSERT
    ON comments
    FOR EACH ROW EXECUTE FUNCTION incr_comments();

CREATE TRIGGER decrement_comments
    AFTER DELETE
    ON comments
    FOR EACH ROW EXECUTE FUNCTION decr_comments();

CREATE FUNCTION upd_recent() RETURNS trigger AS $$
    DECLARE
        p text;
        t text;
    BEGIN
        SELECT publisher, tracking_id
            INTO p, t
            FROM documents
            WHERE id = NEW.documents_id;
        IF FOUND THEN
            UPDATE advisories
                SET recent = greatest(recent, NEW.recent)
                WHERE publisher = p AND tracking_id = t;
        END IF;
        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_recent
    AFTER INSERT OR UPDATE
    ON comments
    FOR EACH ROW EXECUTE FUNCTION upd_recent();
