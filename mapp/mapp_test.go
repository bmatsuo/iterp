package mapp

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	assert.Equal(t, map[string]string{"a": "1", "b": "2", "c": "3"}, Map(m, func(x int) string { return strconv.Itoa(x) }))
}

func TestMap2(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	strs := Map2(m, func(k string, v int) string { return k + ":" + strconv.Itoa(v) })
	assert.Equal(t, map[string]string{"a": "a:1", "b": "b:2", "c": "c:3"}, strs)
}
