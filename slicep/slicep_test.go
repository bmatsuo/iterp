package slicep

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	s := []int{1, 2, 3}
	assert.Equal(t, []string{"1", "2", "3"}, Map(s, func(x int) string { return strconv.Itoa(x) }))
}

func TestMap2(t *testing.T) {
	s := []int{1, 2, 3}
	strs := Map2(s, func(i int, x int) string { return strconv.Itoa(i) + ":" + strconv.Itoa(x) })
	assert.Equal(t, []string{"0:1", "1:2", "2:3"}, strs)
}

func TestSelect(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	assert.Equal(t, []int{2, 4}, Select(s, func(x int) bool { return x%2 == 0 }))
}

func TestReject(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	assert.Equal(t, []int{1, 3, 5}, Reject(s, func(x int) bool { return x%2 == 0 }))
}

func TestFoldLeft(t *testing.T) {
	s := []int{1, 2, 3}
	assert.Equal(t, "123", FoldLeft(s, "", func(acc string, x int) string { return acc + strconv.Itoa(x) }))
}

func TestFoldRight(t *testing.T) {
	s := []int{1, 2, 3}
	assert.Equal(t, "123", FoldRight(s, "", func(x int, acc string) string { return strconv.Itoa(x) + acc }))
}
