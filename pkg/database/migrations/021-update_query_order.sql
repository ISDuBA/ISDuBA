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
    default_advisory_columns character varying array default Array['cvss_v3_score', 'cvss_v2_score', 'comments', 'critical', 'id', 'recent', 'versions', 'title', 'publisher', 'ssvc', 'state', 'tracking_id'];
    default_event_columns character varying array default Array['cvss_v3_score', 'cvss_v2_score', 'comments', 'critical', 'id', 'title', 'publisher', 'ssvc', 'tracking_id', 'event', 'event_state', 'time', 'actor', 'comments_id'];
BEGIN
    IF EXISTS (SELECT 1
      FROM stored_queries
      WHERE definer = default_definer
      AND global = true
      AND name = 'Admin-advisories-global-default'
      AND description = 'To delete'
      AND query = '$state delete workflow ='
      AND columns = default_advisory_columns
      AND dashboard = true
      AND role = 'admin' ) THEN
          DELETE FROM stored_queries
          WHERE definer = default_definer
          AND global = true
          AND name = 'Admin-advisories-global-default'
          AND description = 'To delete'
          AND query = '$state delete workflow ='
          AND columns = default_advisory_columns
          AND dashboard = true
          AND role = 'admin';
          INSERT INTO stored_queries (definer, global, name, description, query, columns, dashboard, role) VALUES(default_definer, true, 'Admin-advisories-global-default', 'To delete', '$state delete workflow =', default_advisory_columns, true, 'admin');
    END IF;

    IF EXISTS (SELECT 1
      FROM stored_queries
      WHERE definer = default_definer
      AND global = true
      AND name = 'Admin-recent-global-default'
      AND description = 'Recent changes'
      AND query = '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !='
      AND columns = default_event_columns
      AND orders = '{"-time"}'
      AND dashboard = true
      AND role = 'admin' ) THEN
          DELETE FROM stored_queries
          WHERE definer = default_definer
          AND global = true
          AND name = 'Admin-recent-global-default'
          AND description = 'Recent changes'
          AND query = '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !='
          AND columns = default_event_columns
          AND orders = '{"-time"}'
          AND dashboard = true
          AND role = 'admin';
          INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role, kind) VALUES(default_definer, true, 'Admin-recent-global-default', 'Recent changes', '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !=', default_event_columns, '{"-time"}', true, 'admin', 'events');
    END IF;

    IF EXISTS (SELECT 1
      FROM stored_queries
      WHERE definer = default_definer
      AND global = true
      AND name = 'Reviewer-advisories-global-default'
      AND description = 'Recently evaluated'
      AND query = '$state review workflow ='
      AND columns = default_advisory_columns
      AND orders = '{"-critical"}'
      AND dashboard = true
      AND role = 'reviewer' ) THEN
          DELETE FROM stored_queries
          WHERE definer = default_definer
          AND global = true
          AND name = 'Reviewer-advisories-global-default'
          AND description = 'Recently evaluated'
          AND query = '$state review workflow ='
          AND columns = default_advisory_columns
          AND orders = '{"-critical"}'
          AND dashboard = true
          AND role = 'reviewer';
          INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role) VALUES(default_definer, true, 'Reviewer-advisories-global-default', 'Recently evaluated', '$state review workflow =', default_advisory_columns, '{"-critical"}', true, 'reviewer');
    END IF;

    IF EXISTS (SELECT 1
      FROM stored_queries
      WHERE definer = default_definer
      AND global = true
      AND name = 'Reviewer-recent-global-default'
      AND description = 'Recent changes'
      AND query = '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !='
      AND columns = default_event_columns
      AND orders = '{"-time"}'
      AND dashboard = true
      AND role = 'reviewer' ) THEN
          DELETE FROM stored_queries
          WHERE definer = default_definer
          AND global = true
          AND name = 'Reviewer-recent-global-default'
          AND description = 'Recent changes'
          AND query = '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !='
          AND columns = default_event_columns
          AND orders = '{"-time"}'
          AND dashboard = true
          AND role = 'reviewer';
          INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role, kind) VALUES(default_definer, true, 'Reviewer-recent-global-default', 'Recent changes', '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !=', default_event_columns, '{"-time"}', true, 'reviewer', 'events');
    END IF;

    IF EXISTS (SELECT 1
      FROM stored_queries
      WHERE definer = default_definer
      AND global = true
      AND name = 'Editor-advisories-global-default'
      AND description = 'New'
      AND query = '$state new workflow ='
      AND columns = default_advisory_columns
      AND orders = '{"-recent"}'
      AND dashboard = true
      AND role = 'editor' ) THEN
          DELETE FROM stored_queries
          WHERE definer = default_definer
          AND global = true
          AND name = 'Editor-advisories-global-default'
          AND description = 'New'
          AND query = '$state new workflow ='
          AND columns = default_advisory_columns
          AND orders = '{"-recent"}'
          AND dashboard = true
          AND role = 'editor';
          INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role) VALUES(default_definer, true, 'Editor-advisories-global-default', 'New', '$state new workflow =', default_advisory_columns, '{"-recent"}', true, 'editor');
    END IF;


    IF EXISTS (SELECT 1
      FROM stored_queries
      WHERE definer = default_definer
      AND global = true
      AND name = 'Editor-recent-global-default'
      AND description = 'Recent changes'
      AND query = '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !='
      AND columns = default_event_columns
      AND orders = '{"-time"}'
      AND dashboard = true
      AND role = 'editor' ) THEN
          DELETE FROM stored_queries
          WHERE definer = default_definer
          AND global = true
          AND name = 'Editor-recent-global-default'
          AND description = 'Recent changes'
          AND query = '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !='
          AND columns = default_event_columns
          AND orders = '{"-time"}'
          AND dashboard = true
          AND role = 'editor';
          INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role, kind) VALUES(default_definer, true, 'Editor-recent-global-default', 'Recent changes', '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !=', default_event_columns, '{"-time"}', true, 'editor', 'events');
    END IF;

    IF EXISTS (SELECT 1
      FROM stored_queries
      WHERE definer = default_definer
      AND global = true
      AND name = 'Source-manager-advisories-global-default'
      AND description = 'New'
      AND query = '$state new workflow ='
      AND columns = default_advisory_columns
      AND orders = '{"-recent"}'
      AND dashboard = true
      AND role = 'source-manager' ) THEN
          DELETE FROM stored_queries
          WHERE definer = default_definer
          AND global = true
          AND name = 'Source-manager-advisories-global-default'
          AND description = 'New'
          AND query = '$state new workflow ='
          AND columns = default_advisory_columns
          AND orders = '{"-recent"}'
          AND dashboard = true
          AND role = 'source-manager';
          INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role) VALUES(default_definer, true, 'Source-manager-advisories-global-default', 'New', '$state new workflow =', default_advisory_columns, '{"-recent"}', true, 'source-manager');
    END IF;

    IF EXISTS (SELECT 1
      FROM stored_queries
      WHERE definer = default_definer
      AND global = true
      AND name = 'Source-manager-recent-global-default'
      AND description = 'Recent changes'
      AND query = '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !='
      AND columns = default_event_columns
      AND orders = '{"-time"}'
      AND dashboard = true
      AND role = 'source-manager' ) THEN
          DELETE FROM stored_queries
          WHERE definer = default_definer
          AND global = true
          AND name = 'Source-manager-recent-global-default'
          AND description = 'Recent changes'
          AND query = '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !='
          AND columns = default_event_columns
          AND orders = '{"-time"}'
          AND dashboard = true
          AND role = 'source-manager';
          INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role, kind) VALUES(default_definer, true, 'Source-manager-recent-global-default', 'Recent changes', '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !=', default_event_columns, '{"-time"}', true, 'source-manager', 'events');
    END IF;

    IF EXISTS (SELECT 1
      FROM stored_queries
      WHERE definer = default_definer
      AND global = true
      AND name = 'Auditor-recent-global-default'
      AND description = 'Recent changes'
      AND query = '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !='
      AND columns = default_event_columns
      AND orders = '{"-time"}'
      AND dashboard = true
      AND role = 'auditor' ) THEN
          DELETE FROM stored_queries
          WHERE definer = default_definer
          AND global = true
          AND name = 'Auditor-recent-global-default'
          AND description = 'Recent changes'
          AND query = '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !='
          AND columns = default_event_columns
          AND orders = '{"-time"}'
          AND dashboard = true
          AND role = 'auditor';
          INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role, kind) VALUES(default_definer, true, 'Auditor-recent-global-default', 'Recent changes', '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !=', default_event_columns, '{"-time"}', true, 'auditor', 'events');
    END IF;

    IF EXISTS (SELECT 1
      FROM stored_queries
      WHERE definer = default_definer
      AND global = true
      AND name = 'Importer-advisories-global-default'
      AND description = 'New'
      AND query = '$event import_document events = me mentioned me involved or and now 168h duration - $time <= $actor me !='
      AND columns = default_event_columns
      AND orders = '{"-time"}'
      AND dashboard = true
      AND role = 'importer' ) THEN
          DELETE FROM stored_queries
          WHERE definer = default_definer
          AND global = true
          AND name = 'Importer-advisories-global-default'
          AND description = 'New'
          AND query = '$event import_document events = me mentioned me involved or and now 168h duration - $time <= $actor me !='
          AND columns = default_event_columns
          AND orders = '{"-time"}'
          AND dashboard = true
          AND role = 'importer';
          INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role, kind) VALUES(default_definer, true, 'Importer-advisories-global-default', 'New', '$event import_document events = me mentioned me involved or and now 168h duration - $time <= $actor me !=', default_event_columns, '{"-time"}', true, 'importer', 'events');
    END IF;

    IF EXISTS (SELECT 1
      FROM stored_queries
      WHERE definer = default_definer
      AND global = true
      AND name = 'Importer-recent-global-default'
      AND description = 'Recent changes'
      AND query = '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !='
      AND columns = default_event_columns
      AND orders = '{"-time"}'
      AND dashboard = true
      AND role = 'importer' ) THEN
          DELETE FROM stored_queries
          WHERE definer = default_definer
          AND global = true
          AND name = 'Importer-recent-global-default'
          AND description = 'Recent changes'
          AND query = '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !='
          AND columns = default_event_columns
          AND orders = '{"-time"}'
          AND dashboard = true
          AND role = 'importer';
          INSERT INTO stored_queries (definer, global, name, description, query, columns, orders, dashboard, role, kind) VALUES(default_definer, true, 'Importer-recent-global-default', 'Recent changes', '$event import_document events != me mentioned me involved or and now 168h duration - $time <= $actor me !=', default_event_columns, '{"-time"}', true, 'importer', 'events');
    END IF;
END $$;
