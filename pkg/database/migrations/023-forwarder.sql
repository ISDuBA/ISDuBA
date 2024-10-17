-- This file is Free Software under the Apache-2.0 License
-- without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
--
-- SPDX-License-Identifier: Apache-2.0
--
-- SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
-- Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

CREATE TABLE forwarded_documents (
    url          varchar NOT NULL,
    documents_id int NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    UNIQUE(url, documents_id)
);

GRANT INSERT, DELETE, SELECT, UPDATE ON forwarded_documents TO {{ .User | sanitize }};
