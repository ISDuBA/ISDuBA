-- This file is Free Software under the Apache-2.0 License
-- without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
--
-- SPDX-License-Identifier: Apache-2.0
--
-- SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
-- Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

CREATE TYPE stored_queries_kind AS ENUM (
    'documents', 'advisories', 'events'
);

ALTER TABLE stored_queries ADD COLUMN kind stored_queries_kind NOT NULL DEFAULT 'advisories';
UPDATE stored_queries SET kind = 'documents' WHERE NOT advisories;
ALTER TABLE stored_queries DROP COLUMN advisories;
