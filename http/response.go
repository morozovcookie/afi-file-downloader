package http

import (
	"net/http"
)

type Response struct {
	*http.Response

	redirects []string
}

func (r *Response) StatusCode() int {
	return r.Response.StatusCode
}

func (r *Response) ContentLength() int64 {
	return r.Response.ContentLength
}

func (r *Response) Redirects() []string {
	return r.redirects
}

func (r *Response) ContentType() string {
	return r.Header.Get("Content-Type")
}

func (r *Response) Close() error {
	return r.Body.Close()
}

func (r *Response) Read(p []byte) (n int, err error) {
	return r.Body.Read(p)
}

func (r *Response) Headers() map[string][]string {
	return r.Response.Header
}
