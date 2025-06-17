-- This file is Free Software under the Apache-2.0 License
-- without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
--
-- SPDX-License-Identifier: Apache-2.0
--
-- SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
-- Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

ALTER TABLE sources
    ADD COLUMN checksum         bytea,
    ADD COLUMN checksum_ack     timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP - '1 second'::interval,
    ADD COLUMN checksum_updated timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP;

