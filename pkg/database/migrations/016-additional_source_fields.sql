-- This file is Free Software under the Apache-2.0 License
-- without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
--
-- SPDX-License-Identifier: Apache-2.0
--
-- SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
-- Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

ALTER TABLE sources ADD strict_mode bool NULL;
ALTER TABLE sources ADD insecure bool NULL;
ALTER TABLE sources ADD signature_check bool NULL;
ALTER TABLE sources ADD age bigint NULL;
ALTER TABLE sources ADD COLUMN ignore_patterns text[];