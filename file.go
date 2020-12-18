package afifiledownloader

import (
	"context"
	"io"
)

//
type FileDownloadResponse interface {
	io.ReadCloser

	//
	StatusCode() int

	//
	ContentLength() int64

	//
	Redirects() []string

	//
	ContentType() string

	//
	Headers() map[string][]string
}

//
type CallbackFunc func(resp FileDownloadResponse) error

//
type FileService interface {
	//
	DownloadFile(ctx context.Context, name string, callback CallbackFunc) error

	//
	GetFileStat(ctx context.Context, name string, callback CallbackFunc) error
}

//
type FileStreamer interface {
	//
	Stream(r io.Reader) error

	//
	Close() error
}
