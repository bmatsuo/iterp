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
	return func(yield func(Z) bool) {
		for i := start; ; i += step {
			if !yield(i) {
				return
			}
		}
	}
}
