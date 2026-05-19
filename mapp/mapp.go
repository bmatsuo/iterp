package mapp

import "github.com/bmatsuo/iterp/funcs"

func Map[K comparable, T any, U any](m map[K]T, f funcs.Map[T, U]) map[K]U {
	result := make(map[K]U, len(m))
	for k, v := range m {
		result[k] = f(v)
	}
	return result
}

func Map2[K comparable, T any, U any](m map[K]T, f funcs.Map2[K, T, U]) map[K]U {
	result := make(map[K]U, len(m))
	for k, v := range m {
		result[k] = f(k, v)
	}
	return result
}
