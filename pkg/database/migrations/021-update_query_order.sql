-- This file is Free Software under the Apache-2.0 License
-- without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
--
-- SPDX-License-Identifier: Apache-2.0
--
-- SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
-- Software-Engineering: 2024 Intevation GmbH <https://intevation.de>


DO $$
DECLARE
    default_definer constant varchar = 'system-default';
    default_advisory_columns text array default Array['cvss_v3_score', 'cvss_v2_score', 'comments', 'critical', 'id', 'recent', 'versions', 'title', 'publisher', 'ssvc', 'state', 'tracking_id'];
    default_event_columns text array default Array['cvss_v3_score', 'cvss_v2_score', 'comments', 'critical', 'id', 'title', 'publisher', 'ssvc', 'tracking_id', 'event', 'event_state', 'time', 'actor', 'comments_id'];
BEGIN
    DELETE FROM stored_queries WHERE definer = default_definer;
    INSERT INTO stored_queries (definer, global, name, description, query, columns, dashboard, role) VALUES(default_definer, true, 'Admin-advisories-global-default', 'To delete', '$state delete workflow =', default_advisory_columns, true, 'admin');
    INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role, kind) VALUES(default_definer, true, 'Admin-recent-global-default', 'Recent changes', '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !=', default_event_columns, '{"-time"}', true, 'admin', 'events');
    INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role) VALUES(default_definer, true, 'Reviewer-advisories-global-default', 'Recently evaluated', '$state review workflow =', default_advisory_columns, '{"-critical"}', true, 'reviewer');
    INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role, kind) VALUES(default_definer, true, 'Reviewer-recent-global-default', 'Recent changes', '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !=', default_event_columns, '{"-time"}', true, 'reviewer', 'events');
    INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role, kind) VALUES(default_definer, true, 'Editor-recent-global-default', 'Recent changes', '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !=', default_event_columns, '{"-time"}', true, 'editor', 'events');
    INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role) VALUES(default_definer, true, 'Editor-advisories-global-default', 'New', '$state new workflow =', default_advisory_columns, '{"-recent"}', true, 'editor');
    INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role) VALUES(default_definer, true, 'Source-manager-advisories-global-default', 'New', '$state new workflow =', default_advisory_columns, '{"-recent"}', true, 'source-manager');
    INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role, kind) VALUES(default_definer, true, 'Source-manager-recent-global-default', 'Recent changes', '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !=', default_event_columns, '{"-time"}', true, 'source-manager', 'events');
    INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role, kind) VALUES(default_definer, true, 'Auditor-recent-global-default', 'Recent changes', '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !=', default_event_columns, '{"-time"}', true, 'auditor', 'events');
    INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role, kind) VALUES(default_definer, true, 'Importer-advisories-global-default', 'New', '$event import_document events = me mentioned me involved or and now 168h duration - $time <= $actor me !=', default_event_columns, '{"-time"}', true, 'importer', 'events');
    INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role, kind) VALUES(default_definer, true, 'Importer-recent-global-default', 'Recent changes', '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !=', default_event_columns, '{"-time"}', true, 'importer', 'events');
END $$;
