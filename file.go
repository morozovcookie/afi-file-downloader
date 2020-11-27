package afifiledownloader

import (
	"io"
	"net/http"
	"time"

	"github.com/stretchr/testify/mock"
)

type DownloadCallback func(r *http.Response) (err error)

type DownloadFunc func(url string, d time.Duration, c DownloadCallback) (err error)

type Streamer interface {
	io.WriteCloser
}

type MockStreamer struct {
	mock.Mock
}

func (ms *MockStreamer) Write(p []byte) (n int, err error) {
	args := ms.Called(p)

	return args.Get(0).(int), args.Error(1)
}

func (ms *MockStreamer) Close() (err error) {
	return ms.Called().Error(0)
}
