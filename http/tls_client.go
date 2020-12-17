package http

import (
	"net/http"
)

type TLSClient struct {
	requester *Requester
}

func NewTLSClient() Client {
	return &TLSClient{
		requester: NewRequester(&http.Client{
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}),
	}
}

func (c *TLSClient) SetRedirects(_ int64) {}

func (c *TLSClient) Do(req *http.Request, callback func(resp *Response) error) error {
	resp, err := c.requester.MakeRequest(req)
	if err != nil {
		return err
	}

	return callback(resp)
}
