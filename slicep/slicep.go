/*
Package slicep supplements the standard library's slices package by providing
functions corresponding to sequence processing utilities in the iterp package.
Functions in this package should match corresponding iterp function signatures
and semantics.

In addition the FoldRight function can be reasonably implemented for slices
while it is not practical (and not available) on general sequences.

# Destructive functions

Some operations may be possible to perform in-place on their input, avoiding
allocations but destroying the input state as a side effect. These functions are
named with the InPlace prefix to clearly indicate their destructive nature.
*/
package slicep

import (
	"iter"
	"slices"

	"github.com/bmatsuo/iterp"
)

// Map applies f to each element of s and returns a slice of the results.
func Map[S ~[]T, T any, U any](s S, f iterp.MapFunc[T, U]) []U {
	result := make([]U, len(s))
	for i, v := range s {
		result[i] = f(v)
	}
	return result
}

// Map2 is like Map but operates on a slice of pairs.
func Map2[S ~[]T, T any, U any](s S, f iterp.Map2Func[int, T, U]) []U {
	result := make([]U, len(s))
	for i, v := range s {
		result[i] = f(i, v)
	}
	return result
}

// Select returns a slice containing all elements of s that satisfy the predicate p.
func Select[S ~[]T, T any](s S, p iterp.PredicateFunc[T]) S {
	// TODO: heuristic for when pre-allocation is bad
	return makeSubseq(s, iterp.Select(slices.Values(s), p))
}

// InPlaceSelect selects elements s matching the predicate p and discards others.
func InPlaceSelect[S ~[]T, T any](s S, p iterp.PredicateFunc[T]) S {
	return keepSubseq(s, iterp.Select(slices.Values(s), p))
}

// Reject returns a slice containing all elements of s that do not satisfy the predicate p.
func Reject[S ~[]T, T any](s S, p iterp.PredicateFunc[T]) S {
	// TODO: heuristic for when pre-allocation is bad
	return makeSubseq(s, iterp.Reject(slices.Values(s), p))
}

// RejectValue returns a slice containing all elements of s that are not equal to v.
func RejectValue[S ~[]T, T comparable](s S, v T) S {
	return makeSubseq(s, iterp.RejectValue(slices.Values(s), v))
}

// InPlaceRejectValue returns a slice containing all elements of s that are not equal
// to v. InPlaceRejectValue reuses the memory of s and compacts all results.
func InPlaceRejectValue[S ~[]T, T comparable](s S, v T) S {
	return keepSubseq(s, iterp.RejectValue(slices.Values(s), v))
}

// makeSubseq preallocates and collects it's values into a newly pre-allocated
// slice.
func makeSubseq[S ~[]T, T any](s S, it iter.Seq[T]) S {
	return slices.AppendSeq(make(S, 0, len(s)), it)
}

// keepSubseq returns a slice in the memory range of s containing the elements
// of it. The caller must ensure that it produces a subsequence of s and that
// the resulting slice is not used after the original slice is modified.
func keepSubseq[S ~[]T, T any](s S, it iter.Seq[T]) S {
	dst := s[:0]
	for v := range it {
		if len(dst) >= len(s) {
			panic("keepSubseq: subsequence is too long to fit in original slice")
		}

		dst = append(dst, v)
	}

	return slices.Delete(s, len(dst), len(s))
}

// FoldLeft aggregates s from left to right by merging elements into an accumulator initialized to init.
func FoldLeft[S ~[]T, T any, U any](s S, init U, f iterp.FoldLFunc[T, U]) U {
	acc := init
	for _, v := range s {
		acc = f(acc, v)
	}
	return acc
}

// FoldRight aggregates s from right to left by merging elements into an accumulator initialized to init.
func FoldRight[S ~[]T, T any, U any](s S, init U, f iterp.FoldRFunc[T, U]) U {
	// FoldRight is inefficient for generic sequences, but slices can be right-folded.
	acc := init
	for i := len(s) - 1; i >= 0; i-- {
		acc = f(s[i], acc)
	}

	return acc
}

// DropWhile returns the longest of suffix of s such that the first element does
// not satisfy the predicate p. DropWhile returns a value consistent with the
// nilness of s.
func DropWhile[S ~[]T, T any](s S, p iterp.PredicateFunc[T]) S {
	i, _, ok := iterp.Find2(slices.All(s), func(i int, v T) bool { return !p(v) })
	if !ok {
		return s[:0]
	}

	return s[i:]
}

// InPlaceDropWhile deletes the longest prefix of s which all satisfy the predicate
// p, shifting values over to reuse space. InPlaceDropWhile returns a value
// consistent with the nilness of s.
func InPlaceDropWhile[S ~[]T, T any](s S, p iterp.PredicateFunc[T]) S {
	i, _, ok := iterp.Find2(slices.All(s), func(i int, v T) bool { return !p(v) })
	if !ok {
		return slices.Delete(s, 0, len(s))
	}

	return slices.Delete(s, 0, i)
}

// DropRightWhile returns a the longest prefix of s such that the final element
// does not satisfy the predicate p. DropRightWhile returns a value consistent
// with the nilness of s.
func DropRightWhile[S ~[]T, T any](s S, p iterp.PredicateFunc[T]) S {
	i, _, ok := iterp.Find2(slices.Backward(s), func(i int, v T) bool { return !p(v) })
	if !ok {
		return s[:0]
	}

	return s[:i+1]
}

// InPlaceDropRightWhile deletes the longest suffix of s which all satisfy the
// predicate p, zeroing out the deleted elements to avoid memory leaks.
// InPlaceDropRightWhile returns a value consistent with the nilness of s.
func InPlaceDropRightWhile[S ~[]T, T any](s S, p iterp.PredicateFunc[T]) S {
	i, _, ok := iterp.Find2(slices.Backward(s), func(i int, v T) bool { return !p(v) })
	if !ok {
		return slices.Delete(s, 0, len(s))
	}
	return slices.Delete(s, i+1, len(s))
}
