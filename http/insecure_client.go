package http

import (
	"crypto/tls"
	"net/http"
)

type InsecureClient struct {
	client *http.Client
}

func NewInsecureClient() Client {
	return &InsecureClient{
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

func (c *InsecureClient) SetRedirects(_ int64) {}

func (c *InsecureClient) Do(req *http.Request, callback func(resp *Response) error) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	return callback(&Response{
		Response: resp,
	})
}
