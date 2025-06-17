-- This file is Free Software under the Apache-2.0 License
-- without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
--
-- SPDX-License-Identifier: Apache-2.0
--
-- SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
-- Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

CREATE TYPE stored_queries_roles AS ENUM (
    'editor', 'reviewer', 'auditor', 'source-manager', 'importer', 'admin'
);

ALTER TABLE stored_queries
    ADD COLUMN dashboard    bool                    NOT NULL DEFAULT FALSE,
    ADD COLUMN role         stored_queries_roles;

CREATE TABLE default_query_exclusion (
    "user"  text    NOT NULL,
    id      int     NOT NULL REFERENCES stored_queries(id) ON DELETE CASCADE,
    UNIQUE ("user", id)
);

GRANT INSERT, DELETE, SELECT, UPDATE ON default_query_exclusion TO {{ .User | sanitize }};

DO $$
DECLARE
    default_definer constant varchar = 'system-default';
    default_advisory_columns text array default Array['cvss_v3_score', 'cvss_v2_score', 'comments', 'critical', 'id', 'recent', 'versions', 'title', 'publisher', 'ssvc', 'state', 'tracking_id'];
    default_event_columns text array default Array['cvss_v3_score', 'cvss_v2_score', 'comments', 'critical', 'id', 'title', 'publisher', 'ssvc', 'tracking_id', 'event', 'event_state', 'time', 'actor', 'comments_id'];
BEGIN
    INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role) VALUES(default_definer, true, 'Editor-advisories-global-default', 'New', '$state new workflow =', default_advisory_columns, '{"-recent"}', true, 'editor');
    INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role) VALUES(default_definer, true, 'Reviewer-advisories-global-default', 'Recently evaluated', '$state review workflow =', default_advisory_columns, '{"-critical"}', true, 'reviewer');
    INSERT INTO stored_queries (definer, global, name, description, query, columns, dashboard, role) VALUES(default_definer, true, 'Admin-advisories-global-default', 'To delete', '$state delete workflow =', default_advisory_columns, true, 'admin');
    INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role, kind) VALUES(default_definer, true, 'Importer-advisories-global-default', 'New', '$event import_document events = me mentioned me involved or and now 168h duration - $time <= $actor me !=', default_event_columns, '{"-time"}', true, 'importer', 'events');
    INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role) VALUES(default_definer, true, 'Source-manager-advisories-global-default', 'New', '$state new workflow =', default_advisory_columns, '{"-recent"}', true, 'source-manager');
    INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role, kind) VALUES(default_definer, true, 'Editor-recent-global-default', 'Recent changes', '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !=', default_event_columns, '{"-time"}', true, 'editor', 'events');
    INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role, kind) VALUES(default_definer, true, 'Reviewer-recent-global-default', 'Recent changes', '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !=', default_event_columns, '{"-time"}', true, 'reviewer', 'events');
    INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role, kind) VALUES(default_definer, true, 'Auditor-recent-global-default', 'Recent changes', '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !=', default_event_columns, '{"-time"}', true, 'auditor', 'events');
    INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role, kind) VALUES(default_definer, true, 'Importer-recent-global-default', 'Recent changes', '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !=', default_event_columns, '{"-time"}', true, 'importer', 'events');
    INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role, kind) VALUES(default_definer, true, 'Admin-recent-global-default', 'Recent changes', '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !=', default_event_columns, '{"-time"}', true, 'admin', 'events');
    INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role, kind) VALUES(default_definer, true, 'Source-manager-recent-global-default', 'Recent changes', '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !=', default_event_columns, '{"-time"}', true, 'source-manager', 'events');
END $$;
