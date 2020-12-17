package http

import (
	"crypto/tls"
	"net/http"
)

type TLSRedirectClient struct {
	client    *http.Client
	redirects int64
}

func NewTLSRedirectClient() Client {
	return &TLSRedirectClient{
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

func (c *TLSRedirectClient) SetRedirects(redirects int64) {
	c.redirects = redirects
}

func (c *TLSRedirectClient) Do(req *http.Request, callback func(resp *Response) error) error {
	var (
		pp   = make(map[string]struct{}, c.redirects+1)
		left = c.redirects
		url  = req.URL.String()

		httpResp = &Response{
			Redirects: make([]string, 0, c.redirects),
		}

		err error
	)

	for {
		if left < 0 {
			return ErrToManyRedirects
		}

		pp[url] = struct{}{}

		if httpResp.Response, err = c.client.Do(req); err != nil {
			return err
		}

		if httpResp.Response.StatusCode == http.StatusMovedPermanently ||
			httpResp.Response.StatusCode == http.StatusFound {
			break
		}

		loc, err := httpResp.Response.Location()
		if err != nil {
			return err
		}

		url = loc.String()

		if _, ok := pp[url]; ok {
			return ErrCyclicRequests
		}

		httpResp.Redirects = append(httpResp.Redirects, url)

		left--
	}

	return callback(httpResp)
}
