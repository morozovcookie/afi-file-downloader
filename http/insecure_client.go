package http

import (
	"crypto/tls"
	"net/http"
)

type InsecureClient struct {
	requester *Requester
}

// nolint: gosec
func NewInsecureClient() Client {
	return &InsecureClient{
		requester: NewRequester(&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}),
	}
}

func (c *InsecureClient) SetRedirects(_ int64) {}

func (c *InsecureClient) Do(req *http.Request, callback func(resp *Response) error) error {
	resp, err := c.requester.MakeRequest(req)
	if err != nil {
		return err
	}

	return callback(resp)
}
