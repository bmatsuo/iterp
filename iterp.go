/*
Package iterp provides basic utilities for processing generic sequences.  The
supplemental packages iterp/slicep and iterp/mapp provide analogous utilities
for slices and maps with consistent APIs.

# Seq and Seq2

The iter.Seq and iter.Seq2 types are the primary abstractions in this package.
Sequences may be infinite and they may not be replayable. This package should
generally provide two implementations of each operations: one that operates on
Seq and one that operates on Seq2. The "normal" version of a function operates
on Seq values and the "2" variant of a function operates on Seq2 values. For
example, Map takes a Seq[T] as input and produces a Seq[U] as output while Map2
takes a Seq2[T,U] as input and produces a Seq2[T,V] as output.

# Map2 vs MapSeq2

Map2 returns an iter.Seq2 with the original sequence's "left" values paired with
the mapped "right" values. This definition of Map2 makes the function work
consistently with the Map2 functions for slices and maps which preserve the
indices/keys of the original container.

MapSeq2 in this package provides a more general mapping for iter.Seq2 than is
provided by Map2. But, the ability to rewrite the "left" value is not well
defined for maps and slices. So while MapSeq2 can be used to create new maps and
slices from existing ones, implementing this functionality is left up to
individual applications which can apply domain specific knowledge to reason
about collisions and/or "gaps" in the output keys/indices and determine
appropriate semantics as necessary.

# FoldLeft vs FoldRight

The FoldLeft and FoldRight functions aggregate a sequence into an accumulator
value. FoldLeft is defined as f(f(f(init, v1), v2), v3)... where v1, v2, v3...
are the sequence elements in order. On the other hand,FoldRight is defined as
f(v1, f(v2, f(v3,...  f(vn, init)))), processing elements in reverse order.
Because there is no memory-efficient way to right-fold a generic sequence the
FoldRight function is only provided for slices.
*/
package iterp

import (
	"iter"
	"slices"

	"github.com/bmatsuo/iterp/funcs"
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

// List2 returns an ordered sequence of pairs containing the given values paired with their indices.
func List2[T any](values ...T) iter.Seq2[int, T] {
	return slices.All(values)
}

// Empty returns a sequence with no elements.
func Empty[T any]() iter.Seq[T] {
	return func(yield func(T) bool) {}
}

// Empty2 returns a empty sequence of pairs.
func Empty2[T any, U any]() iter.Seq2[T, U] {
	return func(yield func(T, U) bool) {}
}

// Concat returns a sequence that concatenates its arguments.
func Concat[T any](its ...iter.Seq[T]) iter.Seq[T] {
	return Flatten(slices.Values(its))
}

// Concat2 returns a sequence that concatenates its arguments.
func Concat2[T any, U any](its ...iter.Seq2[T, U]) iter.Seq2[T, U] {
	return Flatten2(slices.Values(its))
}

// Flatten concatenates a sequence of sequences to produce a single sequence.
func Flatten[T any](its iter.Seq[iter.Seq[T]]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for it := range its {
			for v := range it {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Flatten2 concatenates a sequence of sequences of pairs to produce a single sequence of pairs.
func Flatten2[T any, U any](its iter.Seq[iter.Seq2[T, U]]) iter.Seq2[T, U] {
	return func(yield func(T, U) bool) {
		for it := range its {
			for v, w := range it {
				if !yield(v, w) {
					return
				}
			}
		}
	}
}

// DropN returns the suffix of it without the first n elements. If n is
// negative, the resulting sequence is it.
func DropN[T any](it iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		rem := n
		for v := range it {
			if rem > 0 {
				rem--
				continue
			}
			if !yield(v) {
				return
			}
		}
	}
}

// DropWhile returns the longest suffix of it for which p is false for the first
// element.
func DropWhile[T any](it iter.Seq[T], p funcs.Select[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		dropping := true
		for v := range it {
			if dropping && !p(v) {
				dropping = false
			}
			if !dropping {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// TakeN returns the longest prefix of it with at most n elements. If n is
// negative, the resulting sequence is empty.
func TakeN[T any](it iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		rem := n
		if n <= 0 {
			return
		}
		for v := range it {
			if rem <= 0 {
				return
			}
			rem--
			if !yield(v) {
				return
			}
		}
	}
}

// TakeWhile returns the longest prefix of it for which p is true for every
// element.
func TakeWhile[T any](it iter.Seq[T], p funcs.Select[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range it {
			if !p(v) {
				return
			}
			if !yield(v) {
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
			for range n {
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

// FlatMap is a convenience function for Flatten(Map(it, f))
func FlatMap[T any, U any](it iter.Seq[T], f funcs.Map[T, iter.Seq[U]]) iter.Seq[U] {
	return Flatten(Map(it, f))
}

// Map2 is like Map but operates on a sequence of pairs.  The resulting sequence
// has the same "left" values as the input and "right" values obtained by
// applying f to the input pairs.
func Map2[T any, U any, V any](it iter.Seq2[T, U], f funcs.Map2[T, U, V]) iter.Seq2[T, V] {
	return func(yield func(T, V) bool) {
		for v, w := range it {
			if !yield(v, f(v, w)) {
				return
			}
		}
	}
}

// MapSeq2 is like Map2 but the mapping function produces of the left and right
// values of the output sequence.
func MapSeq2[T any, U any, V any, W any](it iter.Seq2[T, U], f func(T, U) (V, W)) iter.Seq2[V, W] {
	return func(yield func(V, W) bool) {
		for v, w := range it {
			v2, w2 := f(v, w)
			if !yield(v2, w2) {
				return
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
