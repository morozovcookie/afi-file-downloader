package http

import (
	"net/http"
)

type Response struct {
	*http.Response

	Redirects []string
}
