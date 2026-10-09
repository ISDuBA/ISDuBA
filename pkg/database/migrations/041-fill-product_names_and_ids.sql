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

CREATE TABLE products_id_texts (
    documents_id int NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    num          int NOT NULL,
    txt_id       int NOT NULL REFERENCES unique_texts(id) ON DELETE CASCADE,
    UNIQUE(documents_id, num)
);

CREATE INDEX ON products_id_texts(documents_id);
CREATE INDEX ON products_id_texts(txt_id);

WITH pn AS (
  SELECT
    DISTINCT jsonb_path_query(document, '$.product_tree.**.product.name')::int num,
    id
  FROM
    documents
),
data AS (
  SELECT
    documents_id,
    pn.num AS num,
    txt_id
  FROM pn JOIN documents_texts dts
     ON pn.id = dts.documents_id AND pn.num = dts.num
)
INSERT INTO products_name_texts
  SELECT * FROM data;

WITH pid AS (
  SELECT
    DISTINCT jsonb_path_query(document, '$.product_tree.**.product.product_id')::int num,
    id
  FROM
    documents
),
data AS (
  SELECT
    documents_id,
    pid.num AS num,
    txt_id
  FROM pid JOIN documents_texts dts
     ON pid.id = dts.documents_id AND pid.num = dts.num
)
INSERT INTO products_id_texts
  SELECT * FROM data;

GRANT INSERT, DELETE, SELECT, UPDATE ON products_name_texts TO {{ .User | sanitize }};
GRANT INSERT, DELETE, SELECT, UPDATE ON products_id_texts   TO {{ .User | sanitize }};
