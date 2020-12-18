package http

import (
	"net/http"

	afifiledownloader "github.com/morozovcookie/afi-file-downloader"
)

//
type Client interface {
	//
	SetRedirects(redirects int64)

	//
	Do(req *http.Request, callback afifiledownloader.CallbackFunc) error
}
