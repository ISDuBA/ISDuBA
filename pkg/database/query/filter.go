// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

package query

import (
	"iter"
)

func unique[T comparable](seqs ...iter.Seq[T]) iter.Seq[T] {
	have := make(map[T]struct{})
	return func(yield func(T) bool) {
		for _, seq := range seqs {
			for t := range seq {
				if _, ok := have[t]; ok {
					continue
				}
				have[t] = struct{}{}
				if !yield(t) {
					return
				}
			}
		}
	}
}

func enumerate[T any](seq iter.Seq[T]) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		i := 0
		for t := range seq {
			if !yield(i, t) {
				return
			}
			i++
		}
	}
}

func apply[S, T any](seq iter.Seq[S], fn func(S) T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for s := range seq {
			if !yield(fn(s)) {
				return
			}
		}
	}
}
