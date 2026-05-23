package iterp

import (
	"iter"

	"golang.org/x/exp/constraints"
)

// Ints returns an infinite sequence of integers from start.
func Ints[Z constraints.Integer](start Z) iter.Seq[Z] {
	return IntsStep(start, 1)
}

// IntsStep returns an infinite sequence of integers from start and increasing by step.
func IntsStep[Z constraints.Integer](start Z, step Z) iter.Seq[Z] {
	return AffineStep(start, 1, step)
}

// AffineStep returns a sequence of repeated applications of the affine function
// f(z) = a*z + b to the initial value start.
func AffineStep[Z Numeric](start Z, a, b Z) iter.Seq[Z] {
	return RepeatedFunc(start, func(z Z) Z { return a*z + b })
}

// RepeatedFunc returns a sequence of repeated applications of the function:
// start, f(start), f(f(start)), ...
func RepeatedFunc[T any](start T, f func(T) T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := start; ; i = f(i) {
			if !yield(i) {
				return
			}
		}
	}
}
