// This file is Free Software under the MIT License
// without warranty, see README.md and LICENSES/MIT.txt for details.
//
// SPDX-License-Identifier: MIT
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package config

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
)

// Configure sets up the global web server attributes.
func (w *Web) Configure() {
	// If there is a fighting env var, warn the user.
	if ginMode, ok := os.LookupEnv("GIN_MODE"); ok && ginMode != w.GinMode {
		slog.Warn(
			"GIN_MODE ev var conflicts configuration. The configuration always wins.",
			"env", ginMode,
			"cfg", w.GinMode)
	}
	gin.SetMode(w.GinMode)
}
