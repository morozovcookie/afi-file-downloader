package afifiledownloader

import (
	"io"
	"net/http"
	"time"
)

type DownloadCallback func(r *http.Response) (err error)

type DownloadFunc func(url string, d time.Duration, c DownloadCallback) (err error)

type Streamer interface {
	io.WriteCloser
}
