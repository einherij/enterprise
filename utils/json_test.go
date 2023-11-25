package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompareJSON(t *testing.T) {
	a := assert.New(t)
	diffs, err := CompareJSON(`{"wrong":"message", "wrong":"key"}`, `{"wrong":"msg", "wng":"key"}`)
	a.NoError(err)
	a.Len(diffs, 3)
}
