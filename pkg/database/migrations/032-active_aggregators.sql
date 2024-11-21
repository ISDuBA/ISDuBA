-- This file is Free Software under the Apache-2.0 License
-- without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
--
-- SPDX-License-Identifier: Apache-2.0
--
-- SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
-- Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

ALTER TABLE aggregators ADD COLUMN active bool NOT NULL DEFAULT FALSE;

INSERT INTO aggregators (active, name, url) VALUES (FALSE, 'BSI CSAF Lister', 'https://wid.cert-bund.de/.well-known/csaf-aggregator/aggregator.json');
