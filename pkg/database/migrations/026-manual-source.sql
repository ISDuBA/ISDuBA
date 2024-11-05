-- This file is Free Software under the Apache-2.0 License
-- without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
--
-- SPDX-License-Identifier: Apache-2.0
--
-- SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
-- Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

INSERT INTO sources (id, name, url) VALUES (0, 'manual_imports', 'manual.invalid');
INSERT INTO feeds (label, sources_id, url) VALUES ('single', 0, 'https://manual.invalid');
INSERT INTO feeds (label, sources_id, url) VALUES ('bulk', 0, 'https://manual.invalid');
