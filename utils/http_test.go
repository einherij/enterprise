package utils

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CloneRequestBody(t *testing.T) {
	for _, tt := range []struct {
		name      string
		reqMethod string
		reqURL    string
		reqBody   io.Reader
		result    string
	}{
		{
			name:      "WithBody",
			reqMethod: http.MethodPost,
			reqURL:    "https://google.com",
			reqBody:   bytes.NewReader([]byte("request body content")),
			result:    "request body content",
		},
		{
			name:      "NoBody",
			reqMethod: http.MethodGet,
			reqURL:    "https://google.com",
			reqBody:   nil,
			result:    "",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			req, err := http.NewRequest(tt.reqMethod, tt.reqURL, tt.reqBody)
			a.NoError(err)
			bodyContent := CloneRequestBody(req)
			a.Equal(tt.result, bodyContent)

			// check that body is replaced
			if tt.result != "" {
				bodyBytes, err := ioutil.ReadAll(req.Body)
				a.NoError(err)
				a.Equal(tt.result, string(bodyBytes))
			}
		})
	}
}

func Test_CloneResponseBody(t *testing.T) {
	for _, tt := range []struct {
		name          string
		resStatus     string
		resStatusCode int
		resBody       io.ReadCloser
		result        string
	}{
		{
			name:          "WithBody",
			resStatus:     "200 OK",
			resStatusCode: 200,
			resBody:       ioutil.NopCloser(bytes.NewReader([]byte("response body content"))),
			result:        "response body content",
		},
		{
			name:          "NoBody",
			resStatus:     "200 OK",
			resStatusCode: 200,
			resBody:       nil,
			result:        "",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			res := &http.Response{Status: tt.resStatus, StatusCode: tt.resStatusCode, Body: tt.resBody}
			bodyContent := CloneResponseBody(res)
			a.Equal(tt.result, bodyContent)

			// check that body is replaced
			if tt.result != "" {
				bodyBytes, err := ioutil.ReadAll(res.Body)
				a.NoError(err)
				a.Equal(tt.result, string(bodyBytes))
			}
		})
	}
}

func Test_CloneResponseBody_Gzip(t *testing.T) {
	a := assert.New(t)
	var body bytes.Buffer
	gzWriter := gzip.NewWriter(&body)
	_, err := gzWriter.Write([]byte("multilne\nresponse\nbody\ncontent"))
	a.NoError(err)
	a.NoError(gzWriter.Close())
	res := &http.Response{Status: "200 OK", StatusCode: 200, Body: ioutil.NopCloser(&body)}
	res.Header = http.Header{
		contentEncodingHeader: []string{contentEncodingGZIP},
	}
	bodyContent := CloneResponseBody(res)
	a.Equal("multilne\nresponse\nbody\ncontent", bodyContent)
	anotherCopy := CloneResponseBody(res)
	a.Equal("multilne\nresponse\nbody\ncontent", anotherCopy)
}

func Test_GetRequestIP(t *testing.T) {
	for _, tt := range []struct {
		name string
		req  *http.Request
		ip   string
	}{
		{
			name: "EmptyIP",
			req:  &http.Request{},
			ip:   "",
		},
		{
			name: "ByHeader",
			req: &http.Request{
				Header: http.Header{
					xForwardedForHeader: []string{"192.168.1.1"},
				},
				RemoteAddr: "192.168.1.2",
			},
			ip: "192.168.1.1",
		},
		{
			name: "ByRemoteAddr",
			req: &http.Request{
				RemoteAddr: "192.168.1.1:443",
			},
			ip: "192.168.1.1",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.ip, GetRequestIP(tt.req))
		})
	}
}

func Test_SplitIPVersions(t *testing.T) {
	for _, tt := range []struct {
		name string
		ip   string
		ipv4 string
		ipv6 string
	}{
		{
			name: "EmptyIP",
			ip:   "",
			ipv4: "",
			ipv6: "",
		},
		{
			name: "BrokenIP",
			ip:   "afc",
			ipv4: "",
			ipv6: "",
		},
		{
			name: "IPv4",
			ip:   "192.168.1.1",
			ipv4: "192.168.1.1",
			ipv6: "",
		},
		{
			name: "IPv6",
			ip:   "1111:1111:1111:1111:1111:1111:1111:1111",
			ipv4: "",
			ipv6: "1111:1111:1111:1111:1111:1111:1111:1111",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ipv4, ipv6 := SplitIPVersions(tt.ip)
			assert.Equal(t, tt.ipv4, ipv4)
			assert.Equal(t, tt.ipv6, ipv6)
		})
	}
}

func Test_JoinIPVersions(t *testing.T) {
	for _, tt := range []struct {
		name string
		ip   string
		ipv4 string
		ipv6 string
	}{
		{
			name: "EmptyIP",
			ipv4: "",
			ipv6: "",
			ip:   "",
		},
		{
			name: "BrokenIP",
			ipv4: "afc",
			ipv6: "afc",
			ip:   "",
		},
		{
			name: "IPv4",
			ipv4: "192.168.1.1",
			ipv6: "",
			ip:   "192.168.1.1",
		},
		{
			name: "IPv6BrokenIPv4",
			ipv4: "afc",
			ipv6: "1111:1111:1111:1111:1111:1111:1111:1111",
			ip:   "1111:1111:1111:1111:1111:1111:1111:1111",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ip := JoinIPVersions(tt.ipv4, tt.ipv6)
			assert.Equal(t, tt.ip, ip)
		})
	}
}

func Test_FillEmptyIP(t *testing.T) {
	for _, tt := range []struct {
		name    string
		ipv4    string
		ipv6    string
		resIPv4 string
		resIPv6 string
	}{
		{
			name:    "EmptyIP",
			ipv4:    "",
			ipv6:    "",
			resIPv4: ZeroIPv4,
			resIPv6: ZeroIPv6,
		},
		{
			name:    "IPv4",
			ipv4:    "192.168.1.1",
			ipv6:    "",
			resIPv4: "192.168.1.1",
			resIPv6: ZeroIPv6,
		},
		{
			name:    "IPv6",
			ipv4:    "",
			ipv6:    "1111:1111:1111:1111:1111:1111:1111:1111",
			resIPv4: ZeroIPv4,
			resIPv6: "1111:1111:1111:1111:1111:1111:1111:1111",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ipv4, ipv6 := FillEmptyIP(tt.ipv4, tt.ipv6)
			assert.Equal(t, tt.resIPv4, ipv4)
			assert.Equal(t, tt.resIPv6, ipv6)
		})
	}
}
