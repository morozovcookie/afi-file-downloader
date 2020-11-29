package http

import (
	"context"
	"crypto/tls"
	"net/http"
)

type Requester struct {
	c *http.Client
}

// nolint: gosec
func NewRequester(isIgnoreSSLCertificates bool) (requester *Requester) {
	requester = &Requester{
		c: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}

	if !isIgnoreSSLCertificates {
		return requester
	}

	requester.c.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	return requester
}

func (r *Requester) MakeRequest(ctx context.Context, url string) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return r.c.Do(req)
}
