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
