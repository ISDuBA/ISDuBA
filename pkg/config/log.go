// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package config

import (
	"io"
	"log/slog"
	"os"
)

// Config applies the logging configuration to the default slog logger.
func (lg *Log) Config() error {
	var w io.Writer
	if lg.File == "" {
		w = os.Stderr
	} else {
		f, err := os.OpenFile(lg.File, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		w = f
	}
	opts := slog.HandlerOptions{
		AddSource: lg.Source,
		Level:     lg.Level,
	}
	var handler slog.Handler
	if lg.JSON {
		handler = slog.NewJSONHandler(w, &opts)
	} else {
		handler = slog.NewTextHandler(w, &opts)
	}
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return nil
}
