package iterp

import (
	"iter"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChan(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		c := make(chan int)
		go func() {
			close(c)
		}()
		assert.Empty(t, slices.Collect(Chan(c)))
	})

	t.Run("finite", func(t *testing.T) {
		c := make(chan int)
		go func() {
			for i := range 5 {
				c <- i
			}
			close(c)
		}()
		expect := []int{0, 1, 2, 3, 4}
		vals := slices.Collect(Chan(c))
		assert.Equal(t, expect, vals)
	})

	t.Run("break", func(t *testing.T) {
		c := make(chan int)
		go func() {
			for i := 0; i < 5; i++ {
				c <- i
			}
			close(c)
		}()
		expect := []int{0, 1, 2, 3, 4}
		count := 0
		for v := range Chan(c) {
			assert.Equal(t, expect[count], v)
			count++
			break
		}
		assert.Equal(t, 1, count)
	})
}

func TestList(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Empty(t, slices.Collect(List[int]()))
	})

	t.Run("finite", func(t *testing.T) {
		expect := []int{0, 1, 2, 3, 4}
		vals := slices.Collect(List(expect...))
		assert.Equal(t, expect, vals)
	})

	t.Run("replayable", func(t *testing.T) {
		it := List(0, 1, 2, 3, 4)
		expect := []int{0, 1, 2, 3, 4}
		vals1 := slices.Collect(it)
		vals2 := slices.Collect(it)
		assert.Equal(t, expect, vals1)
		assert.Equal(t, expect, vals2)
	})

	t.Run("break", func(t *testing.T) {
		it := List(0, 1, 2, 3, 4)
		expect := []int{0, 1, 2, 3, 4}
		count := 0
		for v := range it {
			assert.Equal(t, expect[count], v)
			count++
			break
		}
		assert.Equal(t, 1, count)
	})
}

func TestEmpty(t *testing.T) {
	assert.Empty(t, slices.Collect(Empty[int]()))
}

func TestConcat(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Empty(t, slices.Collect(Concat[int]()))
	})

	t.Run("finite", func(t *testing.T) {
		expect := []int{0, 1, 2, 3, 4}
		vals := slices.Collect(Concat(List(expect[:2]...), List(expect[2:]...)))
		assert.Equal(t, expect, vals)
	})

	t.Run("break", func(t *testing.T) {
		expect := []int{0, 1, 2, 3, 4}
		vals := Concat(List(expect[:2]...), List(expect[2:]...))
		count := 0
		for v := range vals {
			assert.Equal(t, expect[count], v)
			count++
			if count == 1 {
				break
			}
		}
		assert.Equal(t, 1, count)
	})
}

func TestDropN(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Empty(t, slices.Collect(DropN(Empty[int](), 5)))
	})

	t.Run("negative n", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		expect := source
		vals := DropN(slices.Values(source), -1)
		assert.Equal(t, expect, slices.Collect(vals))
	})

	t.Run("finite n", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		expect := []int{2, 3, 4}
		vals := DropN(slices.Values(source), 2)
		assert.Equal(t, expect, slices.Collect(vals))
	})

	t.Run("n greater than length", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		vals := DropN(slices.Values(source), 10)
		assert.Empty(t, slices.Collect(vals))
	})

	t.Run("break", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		expect := []int{2, 3, 4}
		vals := DropN(slices.Values(source), 2)
		count := 0
		for v := range vals {
			assert.Equal(t, expect[count], v)
			count++
			break
		}
		assert.Equal(t, 1, count)
	})
}

func TestDropWhile(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Empty(t, slices.Collect(DropWhile(Empty[int](), func(_ int) bool { return false })))
	})

	t.Run("always true", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		vals := DropWhile(slices.Values(source), func(_ int) bool { return true })
		assert.Empty(t, slices.Collect(vals))
	})

	t.Run("always false", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		expect := source
		vals := DropWhile(slices.Values(source), func(_ int) bool { return false })
		assert.Equal(t, expect, slices.Collect(vals))
	})

	t.Run("suffix", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		expect := []int{2, 3, 4}
		vals := DropWhile(slices.Values(source), func(x int) bool { return x < 2 })
		assert.Equal(t, expect, slices.Collect(vals))
	})

	t.Run("break", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		expect := []int{2, 3, 4}
		vals := DropWhile(slices.Values(source), func(x int) bool { return x < 2 })
		count := 0
		for v := range vals {
			assert.Equal(t, expect[count], v)
			count++
			break
		}
		assert.Equal(t, 1, count)
	})
}

func TestTakeN(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Empty(t, slices.Collect(TakeN(Empty[int](), 5)))
	})

	t.Run("negative n", func(t *testing.T) {
		vals := TakeN(List(0, 1, 2, 3, 4), -1)
		assert.Empty(t, slices.Collect(vals))
	})

	t.Run("finite n", func(t *testing.T) {
		expect := []int{0, 1}
		vals := TakeN(List(0, 1, 2, 3, 4), 2)
		assert.Equal(t, expect, slices.Collect(vals))
	})

	t.Run("n greater than length", func(t *testing.T) {
		expect := []int{0, 1, 2, 3, 4}
		vals := TakeN(List(expect...), 10)
		assert.Equal(t, expect, slices.Collect(vals))
	})

	t.Run("break", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		expect := []int{0, 1}
		vals := TakeN(List(source...), 2)
		count := 0
		for v := range vals {
			assert.Equal(t, expect[count], v)
			count++
			break
		}
		assert.Equal(t, 1, count)
	})
}

func TestTakeWhile(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Empty(t, slices.Collect(TakeWhile(Empty[int](), func(_ int) bool { return true })))
	})

	t.Run("always true", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		expect := source
		vals := TakeWhile(slices.Values(source), func(_ int) bool { return true })
		assert.Equal(t, expect, slices.Collect(vals))
	})

	t.Run("always false", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		vals := TakeWhile(slices.Values(source), func(_ int) bool { return false })
		assert.Empty(t, slices.Collect(vals))
	})

	t.Run("prefix", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		expect := []int{0, 1}
		vals := TakeWhile(slices.Values(source), func(x int) bool { return x < 2 })
		assert.Equal(t, expect, slices.Collect(vals))
	})

	t.Run("break", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		expect := []int{0, 1}
		vals := TakeWhile(slices.Values(source), func(x int) bool { return x < 2 })
		count := 0
		for v := range vals {
			assert.Equal(t, expect[count], v)
			count++
			break
		}
		assert.Equal(t, 1, count)
	})
}

func TestSelect(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Empty(t, slices.Collect(Select(Empty[int](), func(_ int) bool { return true })))
	})

	t.Run("always true", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		expect := source
		vals := Select(slices.Values(source), func(_ int) bool { return true })
		assert.Equal(t, expect, slices.Collect(vals))
	})

	t.Run("always false", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		vals := Select(slices.Values(source), func(_ int) bool { return false })
		assert.Empty(t, slices.Collect(vals))
	})

	t.Run("even", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		expect := []int{0, 2, 4}
		vals := Select(slices.Values(source), func(x int) bool { return x%2 == 0 })
		assert.Equal(t, expect, slices.Collect(vals))
	})

	t.Run("break", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		vals := Select(slices.Values(source), func(x int) bool { return x%2 == 0 })
		count := 0
		for v := range vals {
			assert.Equal(t, []int{0, 2, 4}[count], v)
			count++
			break
		}
		assert.Equal(t, 1, count)
	})
}

func TestReject(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Empty(t, slices.Collect(Reject(Empty[int](), func(_ int) bool { return true })))
	})

	t.Run("always true", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		vals := Reject(slices.Values(source), func(_ int) bool { return true })
		assert.Empty(t, slices.Collect(vals))
	})

	t.Run("always false", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		expect := source
		vals := Reject(slices.Values(source), func(_ int) bool { return false })
		assert.Equal(t, expect, slices.Collect(vals))
	})

	t.Run("even", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		expect := []int{1, 3}
		vals := Reject(slices.Values(source), func(x int) bool { return x%2 == 0 })
		assert.Equal(t, expect, slices.Collect(vals))
	})

	t.Run("break", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		vals := Reject(slices.Values(source), func(x int) bool { return x%2 == 0 })
		count := 0
		for v := range vals {
			assert.Equal(t, []int{1, 3}[count], v)
			count++
			if count == 1 {
				break
			}
		}
		assert.Equal(t, 1, count)
	})
}

func TestRepeat(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Empty(t, slices.Collect(Repeat(Empty[int](), 5)))
	})

	t.Run("zero repeats", func(t *testing.T) {
		assert.Empty(t, slices.Collect(Repeat(List(0, 1, 2, 3, 4), 0)))
	})

	t.Run("finite repeats", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		expect := append(source, source...)
		vals := slices.Collect(Repeat(slices.Values(source), 2))
		assert.Equal(t, expect, vals)
	})

	t.Run("infinite repeats", func(t *testing.T) {
		source := []int{0, 1}
		expect := []int{0, 1, 0, 1, 0, 1, 0}
		inf := Repeat(slices.Values(source), -1)
		vals := slices.Collect(TakeN(inf, 7))
		assert.Equal(t, expect, vals)
	})

	t.Run("break infinite", func(t *testing.T) {
		source := []int{0, 1}
		inf := Repeat(slices.Values(source), -1)
		count := 0
		for v := range inf {
			assert.Equal(t, source[count%2], v)
			count++
			if count == 7 {
				break
			}
		}
		assert.Equal(t, 7, count)
	})

	t.Run("break finite", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		fin := Repeat(slices.Values(source), 2)
		count := 0
		for v := range fin {
			assert.Equal(t, source[count%5], v)
			count++
			if count == 7 {
				break
			}
		}
		assert.Equal(t, 7, count)
	})
}

func TestLeft(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		vals := Left(slices.All[[]int](nil))
		assert.Empty(t, slices.Collect(vals))
	})

	t.Run("finite", func(t *testing.T) {
		source := []string{"a", "b", "c"}
		expect := []int{0, 1, 2}
		vals := Left(slices.All(source))
		assert.Equal(t, expect, slices.Collect(vals))
	})

	t.Run("break", func(t *testing.T) {
		source := []string{"a", "b", "c"}
		vals := Left(slices.All(source))
		assert.Equal(t, 3, seqLen(vals))
		for v := range vals {
			assert.Equal(t, 0, v)
			break
		}
	})
}

func TestRight(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		vals := Right(slices.All[[]int](nil))
		assert.Empty(t, slices.Collect(vals))
	})

	t.Run("finite", func(t *testing.T) {
		source := []string{"a", "b", "c"}
		vals := Right(slices.All(source))
		assert.Equal(t, source, slices.Collect(vals))
	})

	t.Run("break", func(t *testing.T) {
		source := []string{"a", "b", "c"}
		vals := Right(slices.All(source))
		assert.Equal(t, 3, seqLen(vals))
		for v := range vals {
			assert.Equal(t, "a", v)
			break
		}
	})
}

func TestMap(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Empty(t, slices.Collect(Map(Empty[int](), func(x int) int { return x * 2 })))
	})

	t.Run("finite", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		expect := []int{0, 2, 4, 6, 8}
		vals := slices.Collect(Map(slices.Values(source), func(x int) int { return x * 2 }))
		assert.Equal(t, expect, vals)
	})

	t.Run("break", func(t *testing.T) {
		source := []int{1, 2, 3, 4}
		sum := 0
		expected := 2
		for doubled := range Map(slices.Values(source), func(x int) int {
			return x * 2
		}) {
			sum += doubled
			break
		}
		assert.Equal(t, expected, sum)
	})
}

func TestMap2(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Empty(t, slices.Collect(Right(Map2(Empty2[int, int](), func(i int, x int) int { return i + x }))))
	})

	t.Run("finite", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		expect := []int{0, 2, 4, 6, 8}
		vals := slices.Collect(Right(Map2(slices.All(source), func(i int, x int) int { return i + x })))
		assert.Equal(t, expect, vals)
	})

	t.Run("break", func(t *testing.T) {
		source := []int{1, 2, 3, 4}
		sum := 0
		expected := 2
		for _, doubled := range Map2(slices.All(source), func(_ int, x int) int {
			return x * 2
		}) {
			sum += doubled
			break
		}
		assert.Equal(t, expected, sum)
	})
}

func TestFoldLeft(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		result := FoldLeft(Empty[int](), 0, func(acc int, x int) int { return acc + x })
		assert.Equal(t, 0, result)
	})

	t.Run("finite", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		expect := 10
		result := FoldLeft(slices.Values(source), 0, func(acc int, x int) int { return acc + x })
		assert.Equal(t, expect, result)
	})

	t.Run("order", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		result := FoldLeft(slices.Values(source), []int(nil), func(acc []int, x int) []int { return append(acc, x) })
		assert.Equal(t, source, result)
	})
}

func seqLen[T any](it iter.Seq[T]) int {
	count := 0
	for range it {
		count++
	}
	return count
}

func seqLen2[T any, U any](it iter.Seq2[T, U]) int {
	return seqLen(Left(it))
}
