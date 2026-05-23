/*
Package mapp standard library's maps package by providing functions
corresponding to sequence processing utilities in the iterp package.
*/
package mapp

import (
	"github.com/bmatsuo/iterp"
)

func Map[K comparable, T any, U any](m map[K]T, f iterp.MapFunc[T, U]) map[K]U {
	result := make(map[K]U, len(m))
	for k, v := range m {
		result[k] = f(v)
	}
	return result
}

func Map2[K comparable, T any, U any](m map[K]T, f iterp.Map2Func[K, T, U]) map[K]U {
	result := make(map[K]U, len(m))
	for k, v := range m {
		result[k] = f(k, v)
	}
	return result
}
