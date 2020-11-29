package http

import (
	"context"
	"errors"
	"net/http"
	"time"

	afd "github.com/morozovcookie/afifiledownloader"
)

var (
	ErrToManyRedirects = errors.New("download error: too many redirects")
	ErrCyclicRequests  = errors.New("download error: cyclic requests")
)

type RedirectDownloader struct {
	requester *Requester

	maxRedirects int64
}

func NewRedirectDownloader(maxRedirects int64, isIgnoreSSLCertificates bool) *RedirectDownloader {
	return &RedirectDownloader{
		requester: NewRequester(isIgnoreSSLCertificates),

		maxRedirects: maxRedirects,
	}
}

// nolint: bodyclose
func (rd *RedirectDownloader) Download(
	url string,
	timeout time.Duration,
	c afd.DownloadCallback,
) (
	status int,
	contentLength int64,
	contentType string,
	redirects []string,
	err error,
) {
	var (
		path          = make(map[string]struct{}, rd.maxRedirects+1)
		leftRedirects = rd.maxRedirects
		reqURL        = url

		ctx, cancel = context.WithDeadline(context.Background(), time.Now().Add(timeout))

		resp *http.Response
	)

	defer cancel()

	redirects = make([]string, 0, rd.maxRedirects)

	for {
		if leftRedirects < 0 {
			return 0, 0, "", nil, ErrToManyRedirects
		}

		path[reqURL] = struct{}{}

		if resp, err = rd.requester.MakeRequest(ctx, reqURL); err != nil {
			return 0, 0, "", nil, err
		}

		if isRedirectChainEnd(resp.StatusCode) {
			break
		}

		u, _ := resp.Location()
		reqURL = u.String()

		if _, ok := path[reqURL]; ok {
			return 0, 0, "", nil, ErrCyclicRequests
		}

		redirects = append(redirects, reqURL)
		leftRedirects--
	}

	if err = c(resp); err != nil {
		return 0, 0, "", nil, err
	}

	return resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), redirects, nil
}

func isRedirectChainEnd(status int) bool {
	return status != http.StatusMovedPermanently && status != http.StatusFound
}
