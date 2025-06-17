// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package config

import (
	"errors"
	"fmt"
	"strings"
)

// FeedLogLevel represents a log level in feeds.
type FeedLogLevel int32

const (
	// DebugFeedLogLevel represents the debug log level in feeds.
	DebugFeedLogLevel FeedLogLevel = iota
	// InfoFeedLogLevel represents the info log level in feeds.
	InfoFeedLogLevel
	// WarnFeedLogLevel represents the warn log level in feeds.
	WarnFeedLogLevel
	// ErrorFeedLogLevel represents the error log level in feeds.
	ErrorFeedLogLevel
)

// String implements [fmt.Stringer].
func (fll FeedLogLevel) String() string {
	switch fll {
	case DebugFeedLogLevel:
		return "debug"
	case InfoFeedLogLevel:
		return "info"
	case WarnFeedLogLevel:
		return "warn"
	case ErrorFeedLogLevel:
		return "error"
	default:
		return fmt.Sprintf("unknown feed log level %d", fll)
	}
}

// ParseFeedLogLevel parses feed log levels.
func ParseFeedLogLevel(s string) (FeedLogLevel, error) {
	switch strings.ToLower(s) {
	case "debug":
		return DebugFeedLogLevel, nil
	case "info":
		return InfoFeedLogLevel, nil
	case "warn":
		return WarnFeedLogLevel, nil
	case "error":
		return ErrorFeedLogLevel, nil
	default:
		return 0, fmt.Errorf("unknown feed log level %q", s)
	}
}

// UnmarshalText implements [encoding.TextUnmarshaler].
func (fll *FeedLogLevel) UnmarshalText(b []byte) error {
	x, err := ParseFeedLogLevel(string(b))
	if err != nil {
		return err
	}
	*fll = x
	return nil
}

// MarshalText implements [encoding.TextMarshaler].
func (fll FeedLogLevel) MarshalText() ([]byte, error) {
	return []byte(fll.String()), nil
}

// Scan implements [sql.Scanner].
func (fll *FeedLogLevel) Scan(src any) error {
	if s, ok := src.(string); ok {
		x, err := ParseFeedLogLevel(s)
		if err != nil {
			return err
		}
		*fll = x
		return nil
	}
	return errors.New("unsupported type")
}
