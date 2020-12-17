package http

import (
	"crypto/tls"
	"net/http"
)

type InsecureRedirectClient struct {
	requester *RedirectRequester
	redirects int64
}

// nolint: gosec
func NewInsecureRedirectClient() Client {
	return &InsecureRedirectClient{
		requester: NewRedirectRequester(&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: false,
				},
			},
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}),
	}
}

func (c *InsecureRedirectClient) SetRedirects(redirects int64) {
	c.redirects = redirects
}

func (c *InsecureRedirectClient) Do(req *http.Request, callback func(resp *Response) error) error {
	resp, err := c.requester.MakeRequest(req, c.redirects)
	if err != nil {
		return err
	}

	return callback(resp)
}
