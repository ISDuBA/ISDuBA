-- This file is Free Software under the Apache-2.0 License
-- without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
--
-- SPDX-License-Identifier: Apache-2.0
--
-- SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
-- Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

-- Used for searching
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TABLE versions (
    version     int PRIMARY KEY,
    description text NOT NULL,
    time        timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE workflow AS ENUM (
    'new', 'read', 'assessing',
    'review', 'archived', 'delete');

CREATE TABLE advisories (
    tracking_id  text NOT NULL,
    publisher    text NOT NULL,
    state        workflow NOT NULL DEFAULT 'new',
    -- comments and recent are cached here for performance.
    comments     int NOT NULL DEFAULT 0,
    recent       timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(tracking_id, publisher),
    CHECK(comments >= 0)
);

CREATE INDEX advisories_recent_idx ON advisories(recent);

CREATE FUNCTION utc_timestamp(text) RETURNS timestamp with time zone AS $$
    SELECT $1::timestamp with time zone AT time zone 'utc'
$$ LANGUAGE SQL IMMUTABLE;

CREATE FUNCTION revision_history_length(jsonb) RETURNS int AS $$
    SELECT jsonb_array_length(jsonb_path_query($1, '$.document.tracking.revision_history'))
$$ LANGUAGE SQL IMMUTABLE;

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

CREATE TABLE documents (
    id          int PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    latest      boolean,
    -- The 'kinda' primary key if we could version seriously.
    tracking_id text NOT NULL
                GENERATED ALWAYS AS (document #>> '{document,tracking,id}') STORED,
    publisher   text NOT NULL
                GENERATED ALWAYS AS (document #>> '{document,publisher,name}') STORED,
    version     text NOT NULL
                GENERATED ALWAYS AS (document #>> '{document,tracking,version}') STORED,
    -- Tracking dates
    current_release_date timestamptz
                GENERATED ALWAYS AS (
                utc_timestamp(document #>> '{document,tracking,current_release_date}')) STORED,
    initial_release_date timestamptz
                GENERATED ALWAYS AS (
                utc_timestamp(document #>> '{document,tracking,initial_release_date}')) STORED,
    -- Often used
    tlp         text
                GENERATED ALWAYS AS (document #>> '{document,distribution,tlp,label}') STORED,
    title       text
                GENERATED ALWAYS AS (document #>> '{document,title}') STORED,
    rev_history_length int
                GENERATED ALWAYS AS (revision_history_length(document)) STORED,
    cvss_v2_score float
                GENERATED ALWAYS AS (max_cvss2_score(document)) STORED,
    cvss_v3_score float
                GENERATED ALWAYS AS (max_cvss3_score(document)) STORED,
    critical    float
                GENERATED ALWAYS AS (
                    coalesce(max_cvss3_score(document), max_cvss2_score(document))) STORED,
    four_cves   jsonb
                GENERATED ALWAYS AS (first_four_cves(document)) STORED,
    ssvc        text,
    -- The data
    document    jsonb COMPRESSION lz4 NOT NULL,
    original    bytea COMPRESSION lz4 NOT NULL,
    signature   bytea COMPRESSION lz4,
    filename    varchar,

    FOREIGN KEY (tracking_id, publisher)
        REFERENCES advisories(tracking_id, publisher)
        ON DELETE CASCADE
        DEFERRABLE INITIALLY DEFERRED,
    UNIQUE (tracking_id, publisher, version, rev_history_length)
);

CREATE UNIQUE INDEX only_one_latest_constraint
    ON documents (tracking_id, publisher)
    WHERE latest;

-- create_advisory checks if the new document is newer than the old one.
CREATE FUNCTION create_advisory() RETURNS trigger AS $$
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

-- delete_advisory tries to re-establish an advisory after the last head was deleted.
CREATE FUNCTION delete_advisory() RETURNS trigger AS $$
    DECLARE
        lead_id int;
    BEGIN
        -- Update is only needed if deleted one was latest.
        IF OLD.latest THEN
            SELECT id
                INTO lead_id
                FROM documents
                WHERE tracking_id = OLD.tracking_id AND publisher = OLD.publisher
                ORDER BY current_release_date DESC, rev_history_length DESC;
            IF FOUND THEN
                UPDATE documents SET latest = TRUE WHERE id = lead_id;
            ELSE -- No documents for advisory -> Delete advisory.
                DELETE FROM advisories WHERE (tracking_id, publisher) = (OLD.tracking_id, OLD.publisher);
            END IF;
        END IF;
        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insert_document
    AFTER INSERT
    ON documents
    FOR EACH ROW EXECUTE FUNCTION create_advisory();

CREATE TRIGGER delete_document AFTER DELETE ON documents
    FOR EACH ROW EXECUTE FUNCTION delete_advisory();

CREATE INDEX current_release_date_idx ON documents (current_release_date);
CREATE INDEX initial_release_date_idx ON documents (initial_release_date);

CREATE INDEX documents_cvss2_idx ON documents(coalesce(cvss_v2_score, '0'::double precision) DESC);
CREATE INDEX documents_cvss3_idx ON documents(coalesce(cvss_v3_score, '0'::double precision) DESC);
CREATE INDEX documents_critical_idx ON documents(coalesce(critical, '0'::double precision) DESC);

CREATE TABLE unique_texts (
    id  int PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    txt text COMPRESSION lz4 NOT NULL,
    EXCLUDE USING HASH (txt WITH =)
);

CREATE INDEX ON unique_texts USING gin(txt gin_trgm_ops);

CREATE TABLE documents_texts (
    documents_id int NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    num          int NOT NULL,
    txt_id       int NOT NULL REFERENCES unique_texts(id) ON DELETE CASCADE,
    UNIQUE(documents_id, num)
);

CREATE INDEX ON documents_texts(documents_id);
CREATE INDEX ON documents_texts(txt_id);

CREATE TABLE comments (
    id           int PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    documents_id int NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    time         timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentator  varchar NOT NULL,
    message      varchar(10000)
);

CREATE INDEX ON comments(documents_id);
CREATE INDEX ON comments USING gin(message gin_trgm_ops);

-- Trigger functions to update cached comment count per advisory.
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

CREATE TYPE events AS ENUM (
    'import_document', 'delete_document',
    'state_change',
    'add_sscv', 'change_sscv', 'delete_sscv',
    'add_comment', 'change_comment', 'delete_comment'
);

CREATE TABLE events_log (
    event        events NOT NULL,
    state        workflow,
    time         timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    actor        varchar,
    documents_id int REFERENCES documents(id) ON DELETE SET NULL,
    comments_id  int REFERENCES comments(id) ON DELETE SET NULL
);

CREATE INDEX events_log_time_idx ON events_log(time);
CREATE INDEX ON events_log(documents_id);

-- Trigger to update cached recent value of advisory.
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
                SET recent = greatest(recent, NEW.time)
                WHERE publisher = p AND tracking_id = t;
        END IF;
        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_recent
    AFTER INSERT OR UPDATE
    ON events_log
    FOR EACH ROW EXECUTE FUNCTION upd_recent();

--
-- user defined stored queries
--
CREATE TYPE stored_queries_kind AS ENUM (
    'documents', 'advisories', 'events'
);

CREATE TYPE stored_queries_roles AS ENUM (
    'editor', 'reviewer', 'auditor', 'source-manager', 'importer', 'admin'
);

CREATE TABLE stored_queries (
    id          int                 PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    kind        stored_queries_kind NOT NULL DEFAULT 'advisories',
    definer     varchar             NOT NULL,
    global      boolean             NOT NULL DEFAULT FALSE,
    name        varchar             NOT NULL,
    description varchar,
    query       varchar             NOT NULL,
    num         int                 NOT NULL GENERATED BY DEFAULT AS IDENTITY,
    columns     varchar[]           NOT NULL,
    orders      varchar[],
    dashboard   bool                NOT NULL DEFAULT FALSE,
    role        stored_queries_roles,
    CHECK(name <> ''),
    UNIQUE (definer, name),
    UNIQUE (definer, num) DEFERRABLE INITIALLY DEFERRED
);

CREATE TABLE default_query_exclusion (
    "user"  text    NOT NULL,
    id      int     NOT NULL REFERENCES stored_queries(id) ON DELETE CASCADE,
    UNIQUE ("user", id)
);

---
--- sources
---
CREATE TABLE sources (
    id                     int     PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    name                   varchar NOT NULL UNIQUE,
    url                    varchar NOT NULL,
    active                 bool    NOT NULL DEFAULT FALSE,
    rate                   float,
    slots                  int,
    headers                text[],
    strict_mode            bool,
    insecure               bool,
    signature_check        bool,
    age                    interval,
    ignore_patterns        text[],
    client_cert_public     bytea,
    client_cert_private    bytea,
    client_cert_passphrase bytea,
    CHECK(name <> ''),
    CHECK(url <> ''),
    CHECK(rate IS NULL OR rate > 0.0),
    CHECK(slots IS NULL OR slots >= 1)
);

CREATE TYPE feed_logs_level AS ENUM (
    'debug', 'info', 'warn', 'error');

CREATE TABLE feeds (
    id         int             PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    label      varchar         NOT NULL,
    sources_id int             NOT NULL REFERENCES sources(id) ON DELETE CASCADE,
    url        varchar         NOT NULL,
    rolie      bool            NOT NULL DEFAULT FALSE,
    log_lvl    feed_logs_level NOT NULL DEFAULT 'info',
    CHECK(label <> ''),
    CHECK(url <> ''),
    UNIQUE(label, sources_id)
);

CREATE TABLE changes (
    url      varchar     NOT NULL,
    time     timestamptz NOT NULL,
    feeds_id int         NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    PRIMARY KEY(url, feeds_id),
    CHECK(url <> '')
);

CREATE TABLE feed_logs (
    feeds_id int             NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    lvl      feed_logs_level NOT NULL DEFAULT 'info',
    time     timestamptz     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    msg      text            NOT NULL
);

CREATE INDEX ON feed_logs(time);

CREATE TABLE downloads (
    documents_id     int         REFERENCES documents(id) ON DELETE SET NULL,
    feeds_id         int         REFERENCES feeds(id)     ON DELETE SET NULL,
    time             timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    download_failed  bool,
    filename_failed  bool,
    schema_failed    bool,
    remote_failed    bool,
    checksum_failed  bool,
    signature_failed bool,
    duplicate_failed bool
);

CREATE INDEX ON downloads (time);

-- Track CVEs for documents.
CREATE TABLE unique_cves (
    id  int PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    cve text NOT NULL UNIQUE
);

CREATE TABLE documents_cves (
    documents_id int NOT NULL REFERENCES documents(id)   ON DELETE CASCADE,
    cve_id       int NOT NULL REFERENCES unique_cves(id) ON DELETE CASCADE,
    UNIQUE(documents_id, cve_id)
);

CREATE FUNCTION extract_cves() RETURNS TRIGGER AS $$
    BEGIN
        DELETE FROM documents_cves WHERE documents_id = NEW.id;
        WITH cves_from_document AS (
            SELECT
                id AS doc_id,
                jsonb_array_elements_text(jsonb_path_query_array(document, '$.vulnerabilities."cve"')) AS cve
            FROM documents
            WHERE id = NEW.id
        ),
        inserted AS (
            INSERT INTO unique_cves (cve)
            SELECT cve FROM cves_from_document
            ON     CONFLICT DO NOTHING
            RETURNING id, cve
        ),
        selected AS (
            SELECT id, unique_cves.cve AS cve
            FROM unique_cves JOIN cves_from_document ON unique_cves.cve = cves_from_document.cve
        ),
        resolved AS (
            SELECT * FROM selected
            UNION ALL
            SELECT * FROM inserted
        )
        INSERT INTO documents_cves (documents_id, cve_id)
        SELECT doc_id, resolved.id
        FROM cves_from_document JOIN resolved ON cves_from_document.cve = resolved.cve;
        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER extract_cves_trigger_insert AFTER INSERT
    ON documents
    FOR EACH ROW
    EXECUTE FUNCTION extract_cves();

CREATE TRIGGER extract_cves_trigger_update AFTER UPDATE
    ON documents
    FOR EACH ROW
    WHEN (NEW.document <> OLD.document)
    EXECUTE FUNCTION extract_cves();

--
-- permissions
--
GRANT SELECT ON versions                                        TO {{ .User | sanitize }};
GRANT INSERT, DELETE, SELECT, UPDATE ON advisories              TO {{ .User | sanitize }};
GRANT INSERT, DELETE, SELECT, UPDATE ON documents               TO {{ .User | sanitize }};
GRANT INSERT, DELETE, SELECT, UPDATE ON documents_texts         TO {{ .User | sanitize }};
GRANT INSERT, DELETE, SELECT, UPDATE ON unique_texts            TO {{ .User | sanitize }};
GRANT INSERT, DELETE, SELECT, UPDATE ON comments                TO {{ .User | sanitize }};
GRANT INSERT, DELETE, SELECT, UPDATE ON events_log              TO {{ .User | sanitize }};
GRANT INSERT, DELETE, SELECT, UPDATE ON stored_queries          TO {{ .User | sanitize }};
GRANT INSERT, DELETE, SELECT, UPDATE ON default_query_exclusion TO {{ .User | sanitize }};
GRANT INSERT, DELETE, SELECT, UPDATE ON sources                 TO {{ .User | sanitize }};
GRANT INSERT, DELETE, SELECT, UPDATE ON feeds                   TO {{ .User | sanitize }};
GRANT INSERT, DELETE, SELECT, UPDATE ON changes                 TO {{ .User | sanitize }};
GRANT INSERT, DELETE, SELECT, UPDATE ON feed_logs               TO {{ .User | sanitize }};
GRANT INSERT, DELETE, SELECT, UPDATE ON downloads               TO {{ .User | sanitize }};
GRANT INSERT, DELETE, SELECT, UPDATE ON unique_cves             TO {{ .User | sanitize }};
GRANT INSERT, DELETE, SELECT, UPDATE ON documents_cves          TO {{ .User | sanitize }};

--
-- default queries
--
DO $$
DECLARE
    -- See explanation for definer in docs/developer/queries.md
    default_definer constant varchar = 'system-default';
    default_advisory_columns text array default Array['cvss_v3_score', 'cvss_v2_score', 'comments', 'critical', 'id', 'recent', 'versions', 'title', 'publisher', 'ssvc', 'state', 'tracking_id'];
    default_event_columns text array default Array['cvss_v3_score', 'cvss_v2_score', 'comments', 'critical', 'id', 'title', 'publisher', 'ssvc', 'tracking_id', 'event', 'event_state', 'time', 'actor', 'comments_id'];
BEGIN
    INSERT INTO stored_queries (definer, global, name, description, query, columns, dashboard, role) VALUES(default_definer, true, 'Admin-advisories-global-default', 'To delete', '$state delete workflow =', default_advisory_columns, true, 'admin');
    INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role, kind) VALUES(default_definer, true, 'Admin-recent-global-default', 'Recent changes', '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !=', default_event_columns, '{"-time"}', true, 'admin', 'events');
    INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role) VALUES(default_definer, true, 'Reviewer-advisories-global-default', 'Recently evaluated', '$state review workflow =', default_advisory_columns, '{"-critical"}', true, 'reviewer');
    INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role, kind) VALUES(default_definer, true, 'Reviewer-recent-global-default', 'Recent changes', '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !=', default_event_columns, '{"-time"}', true, 'reviewer', 'events');
    INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role) VALUES(default_definer, true, 'Editor-advisories-global-default', 'New', '$state new workflow =', default_advisory_columns, '{"-recent"}', true, 'editor');
    INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role, kind) VALUES(default_definer, true, 'Editor-recent-global-default', 'Recent changes', '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !=', default_event_columns, '{"-time"}', true, 'editor', 'events');
    INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role) VALUES(default_definer, true, 'Source-manager-advisories-global-default', 'New', '$state new workflow =', default_advisory_columns, '{"-recent"}', true, 'source-manager');
    INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role, kind) VALUES(default_definer, true, 'Source-manager-recent-global-default', 'Recent changes', '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !=', default_event_columns, '{"-time"}', true, 'source-manager', 'events');
    INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role, kind) VALUES(default_definer, true, 'Auditor-recent-global-default', 'Recent changes', '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !=', default_event_columns, '{"-time"}', true, 'auditor', 'events');
    INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role, kind) VALUES(default_definer, true, 'Importer-advisories-global-default', 'New', '$event import_document events = me mentioned me involved or and now 168h duration - $time <= $actor me !=', default_event_columns, '{"-time"}', true, 'importer', 'events');
    INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role, kind) VALUES(default_definer, true, 'Importer-recent-global-default', 'Recent changes', '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !=', default_event_columns, '{"-time"}', true, 'importer', 'events');
END $$;
