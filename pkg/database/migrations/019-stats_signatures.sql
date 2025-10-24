-- This file is Free Software under the Apache-2.0 License
-- without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
--
-- SPDX-License-Identifier: Apache-2.0
--
-- SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
-- Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

ALTER TABLE documents ADD COLUMN signature bytea COMPRESSION lz4;
ALTER TABLE documents ADD COLUMN filename  varchar;

CREATE TABLE downloads (
    documents_id     int         REFERENCES documents(id) ON DELETE SET NULL,
    feeds_id         int         REFERENCES feeds(id)     ON DELETE SET NULL,
    time             timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    download_failed  bool,
    filename_failed  bool,
    schema_failed    bool,
    remote_failed    bool,
    checksum_failed  bool,
    signature_failed bool
);

CREATE INDEX ON downloads (time);

GRANT INSERT, DELETE, SELECT, UPDATE ON downloads TO {{ .User | sanitize }};

