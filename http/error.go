package http

import (
	"errors"
)

var (
	ErrToManyRedirects = errors.New("download error: too many redirects")
	ErrCyclicRequests  = errors.New("download error: cyclic requests")
)
