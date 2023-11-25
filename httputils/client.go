package httputils

import (
	"net/http"
	"time"
)

func NewHighLoadClient(requestTimeout time.Duration) http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxConnsPerHost = 100
	transport.MaxIdleConns = 100000
	transport.MaxIdleConnsPerHost = 50000
	return http.Client{
		Transport: transport,
		Timeout:   requestTimeout,
	}
}
