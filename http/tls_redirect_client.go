package http

import (
	"net/http"

	afifiledownloader "github.com/morozovcookie/afi-file-downloader"
)

type TLSRedirectClient struct {
	requester *RedirectRequester
	redirects int64
}

func NewTLSRedirectClient() Client {
	return &TLSRedirectClient{
		requester: NewRedirectRequester(&http.Client{
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}),
	}
}

func (c *TLSRedirectClient) SetRedirects(redirects int64) {
	c.redirects = redirects
}

func (c *TLSRedirectClient) Do(req *http.Request, callback afifiledownloader.CallbackFunc) error {
	resp, err := c.requester.MakeRequest(req, c.redirects)
	if err != nil {
		return err
	}

	return callback(resp)
}
