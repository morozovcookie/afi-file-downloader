package http

import (
	"crypto/tls"
	"net/http"
)

type TLSClient struct {
	client *http.Client
}

func NewTLSClient() Client {
	return &TLSClient{
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: false,
				},
			},
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

func (c *TLSClient) SetRedirects(_ int64) {}

func (c *TLSClient) Do(req *http.Request, callback func(resp *Response) error) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	return callback(&Response{
		Response: resp,
	})
}
