package http

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	afd "github.com/morozovcookie/afifiledownloader"
)

type Downloader struct {
	c *http.Client
}

// nolint: gosec
func NewDownloader(isIgnoreSSLCertificates bool) (downloader *Downloader) {
	downloader = &Downloader{
		c: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}

	if !isIgnoreSSLCertificates {
		return downloader
	}

	downloader.c.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	return downloader
}

// nolint: bodyclose
func (d *Downloader) Download(
	url string,
	timeout time.Duration,
	c afd.DownloadCallback,
) (
	status int,
	contentLength int64,
	contentType string,
	err error,
) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(timeout))
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, 0, "", err
	}

	res, err := d.c.Do(req)
	if err != nil {
		return 0, 0, "", err
	}

	if err = c(res); err != nil {
		return 0, 0, "", err
	}

	return res.StatusCode, res.ContentLength, res.Header.Get("Content-Type"), nil
}
