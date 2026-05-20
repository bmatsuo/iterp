/*
Package slicep supplements the standard library's slices package by providing
functions corresponding to sequence processing utilities in the iterp package.
Functions in this package should match corresponding iterp function signatures
and semantics.

In addition the FoldRight function can be reasonably implemented for slices
while it is not practical (and not available) on general sequences.
*/
package slicep

import (
	"slices"

	"github.com/bmatsuo/iterp"
	"github.com/bmatsuo/iterp/funcs"
)

// Map applies f to each element of s and returns a slice of the results.
func Map[S ~[]T, T any, U any](s S, f funcs.Map[T, U]) []U {
	result := make([]U, len(s))
	for i, v := range s {
		result[i] = f(v)
	}
	return result
}

// Map2 is like Map but operates on a slice of pairs.
func Map2[S ~[]T, T any, U any](s S, f funcs.Map2[int, T, U]) []U {
	result := make([]U, len(s))
	for i, v := range s {
		result[i] = f(i, v)
	}
	return result
}

// Select returns a slice containing all elements of s that satisfy the predicate p.
func Select[S ~[]T, T any](s S, p funcs.Select[T]) S {
	// TODO: heuristic for when pre-allocation is bad
	return slices.AppendSeq(
		make(S, 0, len(s)),
		iterp.Select(slices.Values(s), p),
	)
}

// Reject returns a slice containing all elements of s that do not satisfy the predicate p.
func Reject[S ~[]T, T any](s S, p funcs.Select[T]) S {
	// TODO: heuristic for when pre-allocation is bad
	return slices.AppendSeq(
		make(S, 0, len(s)),
		iterp.Reject(slices.Values(s), p),
	)
}

// FoldLeft aggregates s from left to right by merging elements into an accumulator initialized to init.
func FoldLeft[S ~[]T, T any, U any](s S, init U, f funcs.FoldL[T, U]) U {
	acc := init
	for _, v := range s {
		acc = f(acc, v)
	}
	return acc
}

// FoldRight aggregates s from right to left by merging elements into an accumulator initialized to init.
func FoldRight[S ~[]T, T any, U any](s S, init U, f funcs.FoldR[T, U]) U {
	// FoldRight is inefficient for generic sequences, but slices can be right-folded.
	acc := init
	for i := len(s) - 1; i >= 0; i-- {
		acc = f(s[i], acc)
	}
	return acc
}
