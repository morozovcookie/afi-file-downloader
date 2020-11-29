package http

import (
	"context"
	"crypto/tls"
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
	c *http.Client

	maxRedirects int64
}

// nolint: gosec
func NewRedirectDownloader(maxRedirects int64, isIgnoreSSLCertificates bool) (downloader *RedirectDownloader) {
	downloader = &RedirectDownloader{
		c: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		maxRedirects: maxRedirects,
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

		req *http.Request
		res *http.Response
	)

	defer cancel()

	redirects = make([]string, 0, rd.maxRedirects)

	for {
		if leftRedirects < 0 {
			return 0, 0, "", nil, ErrToManyRedirects
		}

		path[reqURL] = struct{}{}

		if req, err = http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil); err != nil {
			return 0, 0, "", nil, err
		}

		if res, err = rd.c.Do(req); err != nil {
			return 0, 0, "", nil, err
		}

		if isRedirectChainEnd(res.StatusCode) {
			break
		}

		u, _ := res.Location()
		reqURL = u.String()

		if _, ok := path[reqURL]; ok {
			return 0, 0, "", nil, ErrCyclicRequests
		}

		redirects = append(redirects, reqURL)
		leftRedirects--
	}

	if err = c(res); err != nil {
		return 0, 0, "", nil, err
	}

	return res.StatusCode, res.ContentLength, res.Header.Get("Content-Type"), redirects, nil
}

func isRedirectChainEnd(status int) bool {
	return status != http.StatusMovedPermanently && status != http.StatusFound
}
