package utils

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStr string

func (ms *testStr) UnmarshalQuery(tagVal string) (err error) {
	*ms = testStr(tagVal)
	return nil
}

func (ms *testStr) MarshalQuery() (tagVal string, err error) {
	return string(*ms), nil
}

func TestDecodeQuery(t *testing.T) {
	a := assert.New(t)
	var parsedQuery = new(struct {
		Name  string  `query:"name,nm"`
		Value int     `query:"val , v"`
		Str   testStr `query:"str"`
	})
	q := make(url.Values)
	q.Set("nm", "123")
	q.Set("val", "123")
	q.Set("v", "12")
	q.Set("str", "123")
	p := NewQueryDecoder(q)
	a.NoError(p.DecodeQuery(parsedQuery))
	a.Equal("123", parsedQuery.Name)
	a.Equal(123, parsedQuery.Value)
	a.Equal(testStr("123"), parsedQuery.Str)
}

func TestEncodeQuery(t *testing.T) {
	a := assert.New(t)
	var parsedQuery = struct {
		Name  string  `query:"name,nm"`
		Value int     `query:"val , v"`
		Str   testStr `query:"str"`
	}{
		Name:  "TestName",
		Value: 5,
		Str:   "TestStr",
	}
	q := make(url.Values)
	e := NewQueryEncoder(q)
	a.NoError(e.EncodeQuery(&parsedQuery))
	a.Equal("TestName", q.Get("name"))
	a.Equal("5", q.Get("val"))
	a.Equal("TestStr", q.Get("str"))
}
