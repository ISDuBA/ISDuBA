// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package config

import (
	"log/slog"
	"os"
	"strconv"
)

// envStore maps an env to a store function.
type envStore struct {
	name  string
	store func(string) error
}

func storeLevel(s string) (slog.Level, error) {
	var level slog.Level
	return level, level.UnmarshalText([]byte(s))
}

func parseFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func storeHumanSize(s string) (HumanSize, error) {
	var hs HumanSize
	return hs, hs.UnmarshalText([]byte(s))
}

// noparse returns an unparsed string.
func noparse(s string) (string, error) {
	return s, nil
}

// store returns a function to parse a string to return a function to store a value.
func store[T any](parse func(string) (T, error)) func(*T) func(string) error {
	return func(dst *T) func(string) error {
		return func(s string) error {
			x, err := parse(s)
			if err != nil {
				return err
			}
			*dst = x
			return nil
		}
	}
}

// fill iterates over the mapping and calls the store function
// of every env var that is found.
func storeFromEnv(stores ...envStore) error {
	for _, es := range stores {
		if v, ok := os.LookupEnv(es.name); ok {
			if err := es.store(v); err != nil {
				return err
			}
		}
	}
	return nil
}
