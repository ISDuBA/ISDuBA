// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// HumanSize de-serializes sizes from integer strings
// with suffix "k" (1000), "K" (1024), "m", "M", "g", "G".
// With no suffix given bytes are assumed.
type HumanSize int64

// FeedLogLevel represents a log level in feeds.
type FeedLogLevel int32

// ForwarderStrategy is the filter strategy used by a forwarder.
type ForwarderStrategy int

const (
	// ForwarderStrategyAll forwards all documents to a target.
	ForwarderStrategyAll ForwarderStrategy = iota
	// ForwarderStrategyImportant only forwards important documents to a target.
	ForwarderStrategyImportant
)

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

// UnmarshalText implements [encoding.TextUnmarshaler].
func (hs *HumanSize) UnmarshalText(b []byte) error {
	scale := int64(1)
	if l := len(b); l > 0 {
		switch b[l-1] {
		case 'k':
			scale = 1000
		case 'K':
			scale = 1024
		case 'm':
			scale = 1000 * 1000
		case 'M':
			scale = 1024 * 1024
		case 'g':
			scale = 1000 * 1000 * 1000
		case 'G':
			scale = 1024 * 1024 * 1024
		default:
			goto noUnits
		}
		b = b[:l-1]
	}
noUnits:
	x, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	*hs = HumanSize(scale * x)
	return nil
}

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

// String implements [fmt.Stringer].
func (fs ForwarderStrategy) String() string {
	switch fs {
	case ForwarderStrategyAll:
		return "all"
	case ForwarderStrategyImportant:
		return "important"
	default:
		return fmt.Sprintf("unknown forward strategy %d", fs)
	}
}

// MarshalText implements [encoding.TextMarshaler].
func (fs ForwarderStrategy) MarshalText() ([]byte, error) {
	return []byte(fs.String()), nil
}

// ParseForwarderStrategy parses the forward stratey.
func ParseForwarderStrategy(s string) (ForwarderStrategy, error) {
	switch strings.ToLower(s) {
	case "all":
		return ForwarderStrategyAll, nil
	case "important":
		return ForwarderStrategyImportant, nil
	default:
		return 0, fmt.Errorf("unknown forward strategy %q", s)
	}
}

// UnmarshalText implements [encoding.TextUnmarshaler].
func (fs *ForwarderStrategy) UnmarshalText(b []byte) error {
	x, err := ParseForwarderStrategy(string(b))
	if err != nil {
		return err
	}
	*fs = x
	return nil
}
