// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package sources

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ISDuBA/ISDuBA/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// log writes a log message into the logs of a feed.
func (f *feed) log(m *Manager, level config.FeedLogLevel, format string, args ...any) {
	if f.invalid.Load() || level < config.FeedLogLevel(f.logLevel.Load()) {
		return
	}
	message := fmt.Sprintf(format, args...)
	const sql = `INSERT INTO feed_logs (feeds_id, lvl, msg) VALUES ($1, $2, $3)`
	if err := m.db.Run(
		context.Background(),
		func(ctx context.Context, con *pgxpool.Conn) error {
			_, err := con.Exec(ctx, sql, f.id, level.String(), message)
			return err
		}, 0,
	); err != nil {
		slog.Error("database error", "err", err)
	}
}
