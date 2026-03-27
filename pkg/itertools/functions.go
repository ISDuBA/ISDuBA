// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

// Package itertools is a set of functions
// used commonly on sequences.
package itertools

import (
	"iter"
)

// Concat concatenates a list of sequences into one one.
func Concat[T any](seqs ...iter.Seq[T]) iter.Seq[T] {
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

// Unique removes duplicates from a given sequence and only
// delivers unique ones on their first appearance.
func Unique[T comparable](seq iter.Seq[T]) iter.Seq[T] {
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

// Enumerate returns a [iter.Seq2] injecting i = 0...n
// as the first value.
func Enumerate[T any](seq iter.Seq[T]) iter.Seq2[int, T] {
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

// Apply applies a given function to the given sequence
// and returns a sequence of the return values of the function.
func Apply[S, T any](seq iter.Seq[S], fn func(S) T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for s := range seq {
			if !yield(fn(s)) {
				return
			}
		}
	}
}

// Not is a convenience function to negate a [Filter] accept function.
func Not[T any](accept func(T) bool) func(T) bool {
	return func(t T) bool { return !accept(t) }
}

// Filter returns a sequence that only contains values that
// are accepted by the given function.
func Filter[T any](seq iter.Seq[T], accept func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for t := range seq {
			if accept(t) && !yield(t) {
				return
			}
		}
	}
}
