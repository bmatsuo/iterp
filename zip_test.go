package iterp

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestZip(t *testing.T) {
	t.Run("no sequences", func(t *testing.T) {
		assert.Empty(t, slices.Collect(Zip[int]()))
	})

	t.Run("empty sequences", func(t *testing.T) {
		it1 := List(1, 2, 3)
		it2 := Empty[int]()
		assert.Empty(t, slices.Collect(Zip(it1, it2)))
	})

	t.Run("shortest sequence", func(t *testing.T) {
		it1 := Ints(0)
		it2 := List(1, 2, 3)
		seq := Zip(it1, it2)
		for vs := range seq {
			require.Len(t, vs, 2)
		}
		assert.Equal(t, [][]int{{0, 1}, {1, 2}, {2, 3}}, slices.Collect(seq))
	})

	t.Run("break", func(t *testing.T) {
		it1 := Ints(0)
		it2 := List(1, 2, 3)
		seq := Zip(it1, it2)
		for vs := range seq {
			require.Len(t, vs, 2)
			if vs[0] == 1 {
				break
			}
		}
	})
}

func TestZip2(t *testing.T) {
	t.Run("both empty", func(t *testing.T) {
		seq2 := Zip2(Empty[int](), Empty[int]())
		for range seq2 {
			t.Fail()
		}
	})

	t.Run("empty sequence", func(t *testing.T) {
		it1 := List(1, 2, 3)
		it2 := Empty[int]()
		seq2 := Zip2(it1, it2)
		for range seq2 {
			t.Fail()
		}

		seq2 = Zip2(it2, it1)
		for range seq2 {
			t.Fail()
		}
	})

	t.Run("shortest sequence", func(t *testing.T) {
		it1 := Ints(0)
		it2 := List(1, 2, 3)
		seq2 := Zip2(it1, it2)
		for v1, v2 := range seq2 {
			assert.Equal(t, v1+1, v2)
		}
		assert.Equal(t, 3, seqLen2(seq2))
	})

	t.Run("break", func(t *testing.T) {
		it1 := Ints(0)
		it2 := List(1, 2, 3)
		seq2 := Zip2(it1, it2)
		for v1, v2 := range seq2 {
			assert.Equal(t, v1+1, v2)
			if v1 == 1 {
				break
			}
		}
	})
}
