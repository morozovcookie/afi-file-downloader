package http

import (
	"context"
	"net/http"
	"time"

	afd "github.com/morozovcookie/afifiledownloader"
)

type Downloader struct{}

func NewDownloader() *Downloader {
	return &Downloader{}
}

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

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, 0, "", err
	}

	if err = c(res); err != nil {
		return 0, 0, "", err
	}

	return res.StatusCode, res.ContentLength, res.Header.Get("Content-Type"), nil
}
