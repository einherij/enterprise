package utils

import (
	"io"
)

func CloserFunc(r io.Reader, close func() error) io.ReadCloser {
	return closer{Reader: r, closeFunc: close}
}

func (c closer) Close() error {
	return c.closeFunc()
}

type closer struct {
	io.Reader
	closeFunc func() error
}
