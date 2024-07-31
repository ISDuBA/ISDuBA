// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package sources

import (
	"net/url"
	"time"
)

type location struct {
	updated   time.Time
	doc       *url.URL
	hash      *url.URL
	signature *url.URL
}

type activeFeed struct {
	id        int64
	url       string
	rolie     bool
	nextCheck time.Time
	locations []location
}

type activeSource struct {
	id    int64
	rate  *float64
	slots *int
	feeds []*activeFeed
}
