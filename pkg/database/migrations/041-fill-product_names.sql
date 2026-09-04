-- This file is Free Software under the Apache-2.0 License
-- without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
--
-- SPDX-License-Identifier: Apache-2.0
--
-- SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
-- Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

CREATE TABLE products_name_texts (
    documents_id int NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    num          int NOT NULL,
    txt_id       int NOT NULL REFERENCES unique_texts(id) ON DELETE CASCADE,
    UNIQUE(documents_id, num)
);

CREATE INDEX ON products_name_texts(documents_id);
CREATE INDEX ON products_name_texts(txt_id);



GRANT INSERT, DELETE, SELECT, UPDATE ON products_name_texts TO {{ .User | sanitize }};



