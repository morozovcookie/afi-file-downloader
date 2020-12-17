package http

import (
	"net/http"
)

type Client interface {
	SetRedirects(redirects int64)
	Do(req *http.Request, callback func(resp *Response) error) error
}
