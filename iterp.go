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

	"golang.org/x/exp/constraints"
)

// MapFunc transforms a single value
type MapFunc[T any, U any] = func(value T) U

// Map2Func transforms a value with an associated key/index
type Map2Func[K any, T any, U any] = func(index K, value T) U

// PredicateFunc returns true for values that should be selected.
type PredicateFunc[T any] = func(value T) bool

// PredicateFunc returns true for values that should be selected.
type Predicate2Func[K any, T any] = func(index K, value T) bool

// FoldLFunc merges a value into an accumulator from the left.
type FoldLFunc[T any, U any] = func(accumulator U, value T) U

// FoldRFunc merges a value into an accumulator from the right.
type FoldRFunc[T any, U any] = func(value T, accumulator U) U

// Summable is a constraint that matches types that can be added using the + operator.
type Summable interface {
	constraints.Integer | constraints.Float | constraints.Complex | ~string
}

// Numeric is a constraint that matches all numeric types.
type Numeric interface {
	constraints.Integer | constraints.Float | constraints.Complex
}

// Identity is a MapFunc that returns its argument unchanged.
func Identity[T any](v T) T { return v }

// Identity2 is a Map2Func that returns the second argument unchanged.
func Identity2[K any, T any](k K, v T) T { return v }

// True is a PredicateFunc that always returns true.
func True[T any](T) bool { return true }

// True2 is a Predicate2Func that always returns true.
func True2[K any, T any](K, T) bool { return true }

// False is a PredicateFunc that always returns false.
func False[T any](T) bool { return false }

// False2 is a Predicate2Func that always returns false.
func False2[K any, T any](K, T) bool { return false }

// Cons returns a sequence with head as the first element followed by the
// elements of tail.
func Cons[T any](head T, tail iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		if !yield(head) {
			return
		}
		for v := range tail {
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

// LenMax counts elements in the sequence by iterating it up to max elements.
// LenMax consumes non-replayable sequences. If max is not a positive int,
// LenMax iterates the entire sequence. The returned bool is true if the end of
// the sequence was reached.
func LenMax[T any](it iter.Seq[T], max int) (count int, ok bool) {
	if max <= 0 {
		return lenFast(it)
	}

	for range it {
		if count >= max {
			return count, false
		}
		count++
	}
	return count, true
}

// fast and dangerous
func lenFast[T any](it iter.Seq[T]) (int, bool) {
	count := 0
	for range it {
		count++
	}
	return count, true
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
func DropWhile[T any](it iter.Seq[T], p PredicateFunc[T]) iter.Seq[T] {
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
func TakeWhile[T any](it iter.Seq[T], p PredicateFunc[T]) iter.Seq[T] {
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
func Select[T any](it iter.Seq[T], p PredicateFunc[T]) iter.Seq[T] {
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

func Select2[K any, T any](it iter.Seq2[K, T], p Predicate2Func[K, T]) iter.Seq2[K, T] {
	return func(yield func(K, T) bool) {
		for k, v := range it {
			if p(k, v) {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// Find returns the first element in it for which p is true along with the bool
// true. If there is no such element, it returns a zero value and false.
func Find[T any](it iter.Seq[T], p PredicateFunc[T]) (T, bool) {
	for v := range it {
		if p(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}

// Find2 returns the first pair in it for which p is true along with the bool
// true. If there is no such pair, it returns zero values and false.
func Find2[K any, T any](it iter.Seq2[K, T], p Predicate2Func[K, T]) (K, T, bool) {
	for k, v := range it {
		if p(k, v) {
			return k, v, true
		}
	}
	var zerok K
	var zerot T
	return zerok, zerot, false
}

// Reject returns a subsequence of it without elements for which p is true.
func Reject[T any](it iter.Seq[T], p PredicateFunc[T]) iter.Seq[T] {
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

func RejectValue[T comparable](it iter.Seq[T], v T) iter.Seq[T] {
	return Reject(it, func(x T) bool { return x == v })
}

// Repeat returns a sequence that repeats the sequence it n times. If n is
// negative, the resulting sequence is infinite.  Only replayable sequences can
// be repeated.
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
func Map[T any, U any](it iter.Seq[T], f MapFunc[T, U]) iter.Seq[U] {
	return func(yield func(U) bool) {
		for v := range it {
			if !yield(f(v)) {
				return
			}
		}
	}
}

// FlatMap is a convenience function for Flatten(Map(it, f))
func FlatMap[T any, U any](it iter.Seq[T], f MapFunc[T, iter.Seq[U]]) iter.Seq[U] {
	return Flatten(Map(it, f))
}

// Map2 is like Map but operates on a sequence of pairs.  The resulting sequence
// has the same "left" values as the input and "right" values obtained by
// applying f to the input pairs.
func Map2[T any, U any, V any](it iter.Seq2[T, U], f Map2Func[T, U, V]) iter.Seq2[T, V] {
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
func FoldLeft[T any, U any](it iter.Seq[T], init U, f FoldLFunc[T, U]) U {
	acc := init
	for v := range it {
		acc = f(acc, v)
	}
	return acc
}

// Sum returns the sum of sequence elements.
func Sum[S Summable](it iter.Seq[S]) S {
	var sum S
	for v := range it {
		sum += v
	}
	return sum
}

// Unique returns a sequence that passes through unique elements of it the first
// time each is encountered. Unique uses memory proportional to the number of
// unique sequence elements.
//
// See Uniq for a more memory efficient way of removing duplicates from (sorted) sequences.
func Unique[T comparable](it iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		seen := make(map[T]struct{})
		for v := range it {
			if _, ok := seen[v]; !ok {
				seen[v] = struct{}{}
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Uniq returns a sequence of values from it with consecutive duplicates removed
// much like the uniq command line utility. The iteration uses constant memory.
//
// See Unique for building a sequence of truly unique values.
func Uniq[T comparable](it iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		var prev T
		first := true
		for v := range it {
			if first || v != prev {
				first = false
				prev = v
				if !yield(v) {
					return
				}
			}
		}
	}
}
