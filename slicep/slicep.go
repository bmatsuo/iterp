package slicep

import (
	"slices"

	"github.com/bmatsuo/iterp"
	"github.com/bmatsuo/iterp/funcs"
)

func Map[T any, U any](s []T, f funcs.Map[T, U]) []U {
	result := make([]U, len(s))
	for i, v := range s {
		result[i] = f(v)
	}
	return result
}

func Map2[T any, U any](s []T, f funcs.Map2[int, T, U]) []U {
	result := make([]U, len(s))
	for i, v := range s {
		result[i] = f(i, v)
	}
	return result
}

func Select[T any](s []T, p funcs.Select[T]) []T {
	vals := slices.Values(s)
	vals = iterp.Select(vals, p)
	return slices.Collect(vals)
}

func Reject[T any](s []T, p funcs.Select[T]) []T {
	vals := slices.Values(s)
	vals = iterp.Reject(vals, p)
	return slices.Collect(vals)
}

func FoldLeft[T any, U any](s []T, init U, f funcs.FoldL[T, U]) U {
	return iterp.FoldLeft(slices.Values(s), init, f)
}

func FoldRight[T any, U any](s []T, init U, f funcs.FoldR[T, U]) U {
	// FoldRight is inefficient for generic sequences, but slices can be right-folded.
	acc := init
	for i := len(s) - 1; i >= 0; i-- {
		acc = f(s[i], acc)
	}
	return acc
}
