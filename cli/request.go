package cli

import (
	"time"
)

const (
	DefaultMaxRedirects = 5
	DefaultTimeout      = Duration(time.Second)
)

type Request struct {
	IsIgnoreSSLCertificates bool     `json:"ignore-ssl-certificates"`
	IsFollowRedirects       bool     `json:"follow-redirects"`
	MaxRedirects            int64    `json:"max-redirects"`
	Method                  string   `json:"method"`
	URL                     string   `json:"url"`
	Timeout                 Duration `json:"timeout"`
}

type GetRequest struct {
	*Request

	Output string `json:"output"`
}

type HeadRequest struct {
	*Request

	Headers []string `json:"headers"`
}
