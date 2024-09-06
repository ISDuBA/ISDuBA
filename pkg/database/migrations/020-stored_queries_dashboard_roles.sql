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
    ADD COLUMN dashboard bool NOT NULL DEFAULT FALSE,
    ADD COLUMN role stored_queries_roles;

CREATE TABLE default_query_exclusion (
    "user" text   NOT NULL,
    id  int  NOT NULL REFERENCES stored_queries(id) ON DELETE CASCADE,
    UNIQUE ("user", id)
);

GRANT INSERT, DELETE, SELECT, UPDATE ON default_query_exclusion TO {{ .User | sanitize }};
