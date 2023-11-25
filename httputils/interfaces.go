package httputils

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/url"
)

/*
	interfaces to make mocks
*/

type HTTPServer interface {
	Close() error
	Shutdown(ctx context.Context) error
	RegisterOnShutdown(f func())
	ListenAndServe() error
	Serve(l net.Listener) error
	ServeTLS(l net.Listener, certFile string, keyFile string) error
	SetKeepAlivesEnabled(v bool)
	ListenAndServeTLS(certFile string, keyFile string) error
}

type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
	Do(req *http.Request) (*http.Response, error)
	Post(url string, contentType string, body io.Reader) (resp *http.Response, err error)
	PostForm(url string, data url.Values) (resp *http.Response, err error)
	Head(url string) (resp *http.Response, err error)
	CloseIdleConnections()
}

type HTTPHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}
