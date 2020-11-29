package http

import (
	"context"
	"time"

	afd "github.com/morozovcookie/afifiledownloader"
)

type Downloader struct {
	requester *Requester
}

func NewDownloader(isIgnoreSSLCertificates bool) *Downloader {
	return &Downloader{
		requester: NewRequester(isIgnoreSSLCertificates),
	}
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

	resp, err := d.requester.MakeRequest(ctx, url)
	if err != nil {
		return 0, 0, "", err
	}

	if err = c(resp); err != nil {
		return 0, 0, "", err
	}

	return resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), nil
}
