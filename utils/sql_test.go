package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSQLReplaceArgs(t *testing.T) {
	a := assert.New(t)
	params := []string{"more", "info"}
	query := SQLReplaceArgs("any(?)random?query", params)
	a.Equal("any(more,info)random?query", query)

	query = SQLReplaceArgs("any(?)random?query", ToSQLStringValue(params...))
	a.Equal("any('more','info')random?query", query)
}
