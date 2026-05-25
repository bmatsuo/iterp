package iterp

import (
	"iter"
	"maps"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCons(t *testing.T) {
	t.Run("empty tail", func(t *testing.T) {
		expect := []int{0}
		vals := slices.Collect(Cons(0, Empty[int]()))
		assert.Equal(t, expect, vals)
	})

	t.Run("finite tail", func(t *testing.T) {
		expect := []int{0, 1}
		vals := slices.Collect(Cons(0, Cons(1, Empty[int]())))
		assert.Equal(t, expect, vals)
	})

	t.Run("break", func(t *testing.T) {
		vals := Cons(0, Cons(1, Empty[int]()))
		count := 0
		// there are two parts of the iteration which can be broken. Skipping
		// the first allows us to head both.
		for range vals {
			if count > 0 {
				break
			}
			count++
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

func TestLenMax(t *testing.T) {
	n, ok := LenMax(Empty[int](), 3)
	assert.Equal(t, 0, n)
	assert.True(t, ok)

	n, ok = LenMax(List(0, 1, 2), 3)
	assert.Equal(t, 3, n)
	assert.True(t, ok)

	n, ok = LenMax(List(0, 1, 2, 3), 3)
	assert.Equal(t, 3, n)
	assert.False(t, ok)

	n, ok = LenMax(List(0, 1, 2, 3), 0)
	assert.Equal(t, 4, n)
	assert.True(t, ok)
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

func TestConcat2(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Empty(t, maps.Collect(Concat2[int, int]()))
	})

	t.Run("finite", func(t *testing.T) {
		m1 := map[string]string{
			"a": "a1",
			"b": "b1",
		}
		m2 := map[string]string{
			"b": "b2",
			"c": "c2",
		}
		expect := map[string]string{
			"a": "a1",
			"b": "b2",
			"c": "c2",
		}
		assert.Equal(t, expect, maps.Collect(Concat2(maps.All(m1), maps.All(m2))))
	})

	t.Run("break", func(t *testing.T) {
		entered := false
		for range Concat2(List2(1, 2, 3)) {
			entered = true
			break
		}
		assert.True(t, entered)
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
		vals := Select(slices.Values(source), True)
		assert.Equal(t, expect, slices.Collect(vals))
	})

	t.Run("always false", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		vals := Select(slices.Values(source), False)
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

func TestSelect2(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Empty(t, maps.Collect(Select2(Empty2[int, string](), True2)))
	})

	t.Run("always true", func(t *testing.T) {
		source := []string{"a", "b", "c"}
		expect := map[int]string{
			0: "a",
			1: "b",
			2: "c",
		}
		vals := Select2(slices.All(source), True2)
		assert.Equal(t, expect, maps.Collect(vals))
	})

	t.Run("always false", func(t *testing.T) {
		source := []string{"a", "b", "c"}
		vals := Select2(slices.All(source), False2)
		assert.Empty(t, maps.Collect(vals))
	})

	t.Run("index even", func(t *testing.T) {
		source := []string{"a", "b", "c", "d"}
		expect := map[int]string{
			0: "a",
			2: "c",
		}
		vals := Select2(slices.All(source), func(i int, _ string) bool { return i%2 == 0 })
		assert.Equal(t, expect, maps.Collect(vals))
	})

	t.Run("break", func(t *testing.T) {
		source := []string{"a", "b", "c", "d"}
		vals := Select2(slices.All(source), func(i int, _ string) bool { return i%2 == 1 })
		ran := false
		for i, v := range vals {
			assert.Equal(t, 1, i)
			assert.Equal(t, "b", v)
			ran = true
			break
		}
		assert.True(t, ran)
	})
}

func TestFind(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		v, ok := Find(Empty[int](), func(_ int) bool { return true })
		assert.False(t, ok)
		assert.Equal(t, 0, v)
	})

	t.Run("found", func(t *testing.T) {
		source := []string{"a", "b", "cde"}
		v, ok := Find(slices.Values(source), func(x string) bool { return len(x) > 1 })
		assert.True(t, ok)
		assert.Equal(t, "cde", v)
	})

	t.Run("not found", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		v, ok := Find(slices.Values(source), func(x int) bool { return x > 10 })
		assert.False(t, ok)
		assert.Equal(t, 0, v)
	})
}

func TestFind2(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		k, v, ok := Find2(Empty2[int, string](), func(_ int, _ string) bool { return true })
		assert.False(t, ok)
		assert.Equal(t, 0, k)
		assert.Equal(t, "", v)
	})

	t.Run("found", func(t *testing.T) {
		source := []string{"a", "b", "cde"}
		k, v, ok := Find2(slices.All(source), func(i int, x string) bool { return len(x) > 1 })
		assert.True(t, ok)
		assert.Equal(t, 2, k)
		assert.Equal(t, "cde", v)
	})

	t.Run("not found", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		k, v, ok := Find2(slices.All(source), False2)
		assert.False(t, ok)
		assert.Equal(t, 0, k)
		assert.Equal(t, 0, v)
	})
}

func TestReject(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Empty(t, slices.Collect(Reject(Empty[int](), func(_ int) bool { return true })))
	})

	t.Run("always true", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		vals := Reject(slices.Values(source), True)
		assert.Empty(t, slices.Collect(vals))
	})

	t.Run("always false", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		expect := source
		vals := Reject(slices.Values(source), False)
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

func TestRejectValue(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Empty(t, slices.Collect(RejectValue(Empty[int](), 0)))
	})

	t.Run("not found", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		expect := source
		vals := RejectValue(slices.Values(source), 10)
		assert.Equal(t, expect, slices.Collect(vals))
	})

	t.Run("found", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		expect := []int{0, 1, 3, 4}
		vals := RejectValue(slices.Values(source), 2)
		assert.Equal(t, expect, slices.Collect(vals))
	})

	t.Run("break", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		expect := []int{0, 1, 3, 4}
		vals := RejectValue(slices.Values(source), 2)
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
		vals := Left(Empty2[int, int]())
		assert.Empty(t, slices.Collect(vals))
	})

	t.Run("finite", func(t *testing.T) {
		expect := []int{0, 1, 2}
		vals := Left(List2("a", "b", "c"))
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
		vals := Right(Empty2[int, string]())
		assert.Empty(t, slices.Collect(vals))
	})

	t.Run("finite", func(t *testing.T) {
		source := []string{"a", "b", "c"}
		vals := Right(slices.All(source))
		assert.Equal(t, source, slices.Collect(vals))
	})

	t.Run("break", func(t *testing.T) {
		vals := Right(List2("a", "b", "c"))
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

	t.Run("identity", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		assert.Equal(t, source, slices.Collect(Map(slices.Values(source), Identity)))
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

func TestFlatMap(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Empty(t, slices.Collect(FlatMap(Empty[int](), func(x int) iter.Seq[int] { return List(x) })))
	})

	t.Run("finite", func(t *testing.T) {
		source := []int{0, 1, 2}
		expect := []int{1, 2, 2}
		vals := slices.Collect(FlatMap(slices.Values(source), func(x int) iter.Seq[int] { return Repeat(List(x), x) }))
		assert.Equal(t, expect, vals)
	})
}

func TestMap2(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Empty(t, slices.Collect(Right(Map2(Empty2[int, int](), func(i int, x int) int { return i + x }))))
	})

	t.Run("identity", func(t *testing.T) {
		source := []int{0, 1, 2, 3, 4}
		assert.Equal(t, source, slices.Collect(Right(Map2(slices.All(source), Identity2))))
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

func TestMapSeq2(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Empty(t, slices.Collect(Right(MapSeq2(Empty2[int, int](), func(i int, x int) (int, int) { return i, x }))))
	})

	t.Run("finite", func(t *testing.T) {
		source := map[string]int{
			"a": 0,
			"b": 1,
			"c": 2,
		}
		got := maps.Collect(MapSeq2(maps.All(source), func(k string, x int) (string, int) {
			return strings.ToUpper(k), x * 2
		}))
		expect := map[string]int{
			"A": 0,
			"B": 2,
			"C": 4,
		}
		assert.Equal(t, expect, got)
	})

	t.Run("break", func(t *testing.T) {
		source := []int{1, 2, 3, 4}
		sum := 0
		expected := 2
		for i, v := range MapSeq2(slices.All(source), func(i int, x int) (int, int) {
			return i, 2 * x
		}) {
			sum += i + v
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

func TestSum(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		result := Sum(Empty[int]())
		assert.Equal(t, 0, result)
	})

	t.Run("finite", func(t *testing.T) {
		source := []string{"a", "b", "c", "d"}
		expect := "abcd"
		result := Sum(slices.Values(source))
		assert.Equal(t, expect, result)
	})
}

func TestUnique(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Empty(t, slices.Collect(Unique(Empty[int]())))
	})

	t.Run("finite", func(t *testing.T) {
		source := []int{0, 1, 2, 2, 3, 4, 4}
		expect := []int{0, 1, 2, 3, 4}
		vals := slices.Collect(Unique(slices.Values(source)))
		assert.Equal(t, expect, vals)
	})

	t.Run("break", func(t *testing.T) {
		source := []int{0, 1, 2, 2, 3, 4, 4}
		expect := []int{0, 1, 2, 3}
		vals := Unique(slices.Values(source))
		got := []int(nil)
		for v := range vals {
			got = append(got, v)
			if len(got) >= len(expect) {
				break
			}
		}
		assert.Equal(t, expect, got)
	})
}

func TestUniq(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Empty(t, slices.Collect(Uniq(Empty[int]())))
	})

	t.Run("finite", func(t *testing.T) {
		source := []int{0, 1, 2, 2, 1, 3}
		expect := []int{0, 1, 2, 1, 3}
		vals := slices.Collect(Uniq(slices.Values(source)))
		assert.Equal(t, expect, vals)
	})

	t.Run("break", func(t *testing.T) {
		source := []int{0, 1, 2, 2, 1, 3}
		expect := []int{0, 1, 2}
		vals := Uniq(slices.Values(source))
		got := []int(nil)
		for v := range vals {
			got = append(got, v)
			if len(got) >= len(expect) {
				break
			}
		}
		assert.Equal(t, expect, got)
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
