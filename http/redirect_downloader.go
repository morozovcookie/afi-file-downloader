package http

import (
	"context"
	"errors"
	"net/http"
	"time"

	afd "github.com/morozovcookie/afifiledownloader"
)

type RedirectDownloader struct {
	c *http.Client

	maxRedirects int64
}

func NewRedirectDownloader(maxRedirects int64) *RedirectDownloader {
	return &RedirectDownloader{
		c: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
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

		req *http.Request
		res *http.Response
	)
	defer cancel()

	redirects = make([]string, 0, rd.maxRedirects)

	for {
		if leftRedirects < 0 {
			return 0, 0, "", nil, errors.New("too many redirects")
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
			return 0, 0, "", nil, errors.New("cyclic requests")
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
