package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntInSlice(t *testing.T) {
	intSlice := []int{1, 2, 3, 4, 5}

	assert.True(t, IsInSlice(1, intSlice))
	assert.False(t, IsInSlice(6, intSlice))
}

func TestStringInSlice(t *testing.T) {
	stringSlice := []string{"hello", "world", "here", "we", "go", "again"}

	assert.True(t, IsInSlice("hello", stringSlice))
	assert.False(t, IsInSlice("Hello", stringSlice))
	assert.False(t, IsInSlice("notExist", stringSlice))
}
