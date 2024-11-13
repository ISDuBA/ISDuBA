// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package sources

import (
	"fmt"
	"regexp"
)

// AsStrings returns a slice of strings from a slice of regular expressions.
func AsStrings(s []*regexp.Regexp) []string {
	if s == nil {
		return nil
	}
	slice := make([]string, len(s))
	for i, x := range s {
		slice[i] = x.String()
	}
	return slice
}

// AsRegexps returns a slice of regular expressions from a slice of strings.
func AsRegexps(s []string) ([]*regexp.Regexp, error) {
	if s == nil {
		return nil, nil
	}
	slice := make([]*regexp.Regexp, 0, len(s))
	for _, x := range s {
		// Ignore empty strings.
		if x == "" {
			continue
		}
		re, err := regexp.Compile(x)
		if err != nil {
			return nil, InvalidArgumentError(
				fmt.Sprintf("ignore pattern %q is not a valid regexp: %v", x, err))
		}
		slice = append(slice, re)
	}
	return slice, nil
}
