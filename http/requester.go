package http

import (
	"net/http"
)

type Requester struct {
	client *http.Client
}

func NewRequester(client *http.Client) *Requester {
	return &Requester{
		client: client,
	}
}

// nolint: bodyclose
func (r *Requester) MakeRequest(req *http.Request) (*Response, error) {
	var (
		resp = &Response{}

		err error
	)

	if resp.Response, err = r.client.Do(req); err != nil {
		return nil, err
	}

	return resp, nil
}
