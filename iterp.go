/*
Package iterp provides basic utilities for processing generic sequences.

The supplemental packages iterp/slicep and iterp/mapp provide analogous
utilities for slices and maps with consistent APIs.
*/
package iterp

import (
	"iter"
	"slices"

	"github.com/bmatsuo/iterp/funcs"
	"golang.org/x/exp/constraints"
)

// Chan wraps c as a sequence so it can be passed to sequence processing functions.
// Because the channel is stateful the resulting sequence is not replayable.
func Chan[T any](c <-chan T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range c {
			if !yield(v) {
				return
			}
		}
	}
}

// List returns an ordered sequence containing the given values.
func List[T any](values ...T) iter.Seq[T] {
	return slices.Values(values)
}

// Ints returns an infinite sequence of integers from start.
func Ints[Z constraints.Integer](start Z) iter.Seq[Z] {
	return func(yield func(Z) bool) {
		for i := start; ; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

// IntsStep returns an infinite sequence of integers from start and increasing by step.
func IntsStep[Z constraints.Integer](start Z, step Z) iter.Seq[Z] {
	return func(yield func(Z) bool) {
		for i := start; ; i += step {
			if !yield(i) {
				return
			}
		}
	}
}

// Empty returns a sequence with no elements.
func Empty[T any]() iter.Seq[T] {
	return func(yield func(T) bool) {}
}

// Concat returns a sequence that concatenates its arguments.
func Concat[T any](its ...iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, it := range its {
			for v := range it {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Repeat returns a sequence that repeats it n times. If n is negative, the resulting sequence is infinite.
func Repeat[T any](it iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		if n < 0 {
			for {
				for v := range it {
					if !yield(v) {
						return
					}
				}
			}
		} else {
			for i := 0; i < n; i++ {
				for v := range it {
					if !yield(v) {
						return
					}
				}
			}
		}
	}
}

// Left creates a sequence of "left" elements from a sequence of pairs.
func Left[T any, U any](it iter.Seq2[T, U]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range it {
			if !yield(v) {
				return
			}
		}
	}
}

// Right creates a sequence of "right" elements from a sequence of pairs.
func Right[T any, U any](it iter.Seq2[T, U]) iter.Seq[U] {
	return func(yield func(U) bool) {
		for _, w := range it {
			if !yield(w) {
				return
			}
		}
	}
}

// Map returns a sequence resulting from applying f to elements of it
func Map[T any, U any](it iter.Seq[T], f funcs.Map[T, U]) iter.Seq[U] {
	return func(yield func(U) bool) {
		for v := range it {
			if !yield(f(v)) {
				return
			}
		}
	}
}

// Map2 is like Map but operates on a sequence of pairs.
func Map2[T any, U any, V any](it iter.Seq2[T, U], f funcs.Map2[T, U, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v, w := range it {
			if !yield(f(v, w)) {
				return
			}
		}
	}
}

// Select returns a subsequence of it with all elements for which p is true.
func Select[T any](it iter.Seq[T], p funcs.Select[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range it {
			if p(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Reject returns a subsequence of it without elements for which p is true.
func Reject[T any](it iter.Seq[T], p funcs.Select[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range it {
			if !p(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// FoldLeft aggregates it into an accumulator initialized to init.
//
// Note that a right fold over a generic sequence is very inefficient and it is
// not provided here. Slices can be right-folded.
func FoldLeft[T any, U any](it iter.Seq[T], init U, f funcs.FoldL[T, U]) U {
	acc := init
	for v := range it {
		acc = f(acc, v)
	}
	return acc
}
