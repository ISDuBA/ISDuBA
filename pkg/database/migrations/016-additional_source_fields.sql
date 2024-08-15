-- This file is Free Software under the Apache-2.0 License
-- without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
--
-- SPDX-License-Identifier: Apache-2.0
--
-- SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
-- Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

ALTER TABLE sources ADD COLUMN strict_mode     bool;
ALTER TABLE sources ADD COLUMN insecure        bool;
ALTER TABLE sources ADD COLUMN signature_check bool;
ALTER TABLE sources ADD COLUMN age             interval;
ALTER TABLE sources ADD COLUMN ignore_patterns text[];
