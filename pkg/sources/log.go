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
func (af *activeFeed) log(m *Manager, level config.FeedLogLevel, format string, args ...any) {
	if level < af.logLevel {
		return
	}
	message := fmt.Sprintf(format, args...)
	const sql = `INSERT INTO feed_logs (feeds_id, lvl, msg) ` +
		`SELECT $1, $2, $3 FROM feeds ` +
		`WHERE EXISTS(SELECT 1 FROM feeds WHERE id = $1)`
	if err := m.db.Run(
		context.Background(),
		func(ctx context.Context, con *pgxpool.Conn) error {
			_, err := con.Exec(ctx, sql, af.id, level.String(), message)
			return err
		}, 0,
	); err != nil {
		slog.Error("database error", "err", err)
	}
}
