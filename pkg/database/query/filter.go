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

func concat[T any](seqs ...iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, seq := range seqs {
			for t := range seq {
				if !yield(t) {
					return
				}
			}
		}
	}
}

func unique[T comparable](seq iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		seen := make(map[T]struct{})
		for t := range seq {
			if _, found := seen[t]; !found {
				if !yield(t) {
					return
				}
				seen[t] = struct{}{}
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

func filter[T any](seq iter.Seq[T], accept func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for t := range seq {
			if accept(t) && !yield(t) {
				return
			}
		}
	}
}
