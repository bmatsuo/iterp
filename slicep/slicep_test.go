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

func TestDropWhile(t *testing.T) {
	assert.Equal(t, []int(nil), DropWhile([]int(nil), func(x int) bool { return true }))
	assert.Equal(t, []int{}, DropWhile([]int{}, func(x int) bool { return true }))

	s := []int{1, 2, 3, 4, 5}
	z := DropWhile(s, func(x int) bool { return x < 3 })
	assert.Equal(t, []int{3, 4, 5}, z)
	assert.Equal(t, 1, s[0])
	assert.Equal(t, 4, s[3])
	assert.Equal(t, 3, cap(z))
}

func TestDeleteWhile(t *testing.T) {
	assert.Equal(t, []int(nil), DeleteWhile([]int(nil), func(x int) bool { return true }))
	assert.Equal(t, []int{}, DeleteWhile([]int{}, func(x int) bool { return true }))

	s := []int{1, 2, 3, 4, 5}
	z := DeleteWhile(s, func(x int) bool { return x < 3 })
	assert.Equal(t, []int{3, 4, 5}, z)
	assert.Equal(t, 3, s[0])
	assert.Equal(t, 0, s[3])
	assert.Equal(t, 5, cap(z))
}

func TestDropRightWhile(t *testing.T) {
	assert.Equal(t, []int(nil), DropRightWhile([]int(nil), func(x int) bool { return true }))
	assert.Equal(t, []int{}, DropRightWhile([]int{}, func(x int) bool { return true }))

	s := []int{1, 2, 3, 4, 5}
	z := DropRightWhile(s, func(x int) bool { return x > 3 })
	assert.Equal(t, []int{1, 2, 3}, z)
	assert.Equal(t, 4, s[3])
	assert.Equal(t, 5, cap(z))
}

func TestDeleteRightWhile(t *testing.T) {
	assert.Equal(t, []int(nil), DeleteRightWhile([]int(nil), func(x int) bool { return true }))
	assert.Equal(t, []int{}, DeleteRightWhile([]int{}, func(x int) bool { return true }))

	s := []int{1, 2, 3, 4, 5}
	z := DeleteRightWhile(s, func(x int) bool { return x > 3 })
	assert.Equal(t, []int{1, 2, 3}, z)
	assert.Equal(t, 0, s[3])
	assert.Equal(t, 5, cap(z))
}
