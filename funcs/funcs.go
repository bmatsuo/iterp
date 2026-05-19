package funcs

type Map[T any, U any] = func(T) U

type Map2[K any, T any, U any] = func(i K, v T) U

type Select[T any] = func(T) bool

type FoldL[T any, U any] = func(U, T) U

type FoldR[T any, U any] = func(T, U) U
