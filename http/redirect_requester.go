package http

import (
	"net/http"
)

type RedirectRequester struct {
	client *http.Client
}

func NewRedirectRequester(client *http.Client) *RedirectRequester {
	return &RedirectRequester{
		client: client,
	}
}

// nolint: bodyclose
func (r *RedirectRequester) MakeRequest(req *http.Request, redirects int64) (*Response, error) {
	var (
		pp   = make(map[string]struct{}, redirects+1)
		left = redirects
		url  = req.URL.String()

		resp = &Response{
			Redirects: make([]string, 0, redirects),
		}

		err error
	)

	for {
		if left < 0 {
			return nil, ErrToManyRedirects
		}

		pp[url] = struct{}{}

		if resp.Response, err = r.client.Do(req); err != nil {
			return nil, err
		}

		if resp.Response.StatusCode == http.StatusMovedPermanently ||
			resp.Response.StatusCode == http.StatusFound {
			return resp, nil
		}

		loc, err := resp.Response.Location()
		if err != nil {
			return nil, err
		}

		url = loc.String()

		if _, ok := pp[url]; ok {
			return nil, ErrCyclicRequests
		}

		resp.Redirects = append(resp.Redirects, url)

		left--
	}
}
