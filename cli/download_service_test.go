package cli

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	afd "github.com/morozovcookie/afifiledownloader"
)

func TestDownloadService_Download(t *testing.T) {
	defaultCallback := func(url string, d time.Duration, c afd.DownloadCallback) (err error) {
		res := &http.Response{
			Status:        http.StatusText(http.StatusOK),
			StatusCode:    http.StatusOK,
			Body:          ioutil.NopCloser(bytes.NewBufferString(`{}`)),
			ContentLength: int64(len(`{}`)),
		}

		if err := c(res); err != nil {
			return err
		}

		return nil
	}

	tt := []struct {
		name   string
		enable bool

		df afd.DownloadFunc

		sc                *MockStreamerCreator
		creatorInput      []interface{}
		creatorOutputFunc func() []interface{}

		in io.Reader

		wantErr bool
	}{
		{
			name:   "pass",
			enable: true,

			df: defaultCallback,

			sc: new(MockStreamerCreator),
			creatorInput: []interface{}{
				"127.0.0.1:5000",
			},
			creatorOutputFunc: func() []interface{} {
				s := new(afd.MockStreamer)
				s.
					On("Write", []interface{}{[]byte(`{}`)}...).
					Return([]interface{}{len(`{}`), (error)(nil)}...)
				s.
					On("Close").
					Return([]interface{}{(error)(nil)}...)

				return []interface{}{
					s,
					(error)(nil),
				}
			},

			in: bytes.NewBufferString(`{"url":"http://127.0.0.1:8080/index.html","timeout":"1s","output":"127.0.0.1:5000"}`),
		},
		{
			name:   "empty output",
			enable: true,

			df: defaultCallback,

			sc:           new(MockStreamerCreator),
			creatorInput: []interface{}{},
			creatorOutputFunc: func() []interface{} {
				return []interface{}{}
			},

			in: bytes.NewBuffer([]byte(`{"url":"http://127.0.0.1:8080/index.html","timeout":"1s"}`)),
		},
		{
			name:   "decode error",
			enable: true,

			df: func(url string, d time.Duration, c afd.DownloadCallback) (err error) {
				return nil
			},

			sc:           new(MockStreamerCreator),
			creatorInput: []interface{}{},
			creatorOutputFunc: func() []interface{} {
				return []interface{}{}
			},

			in: bytes.NewBuffer([]byte(`{"timeout":null}`)),

			wantErr: true,
		},
		{
			name:   "validate error",
			enable: true,

			df: defaultCallback,

			sc:           new(MockStreamerCreator),
			creatorInput: []interface{}{},
			creatorOutputFunc: func() []interface{} {
				return []interface{}{}
			},

			in: bytes.NewBufferString(`{"timeout":"1s","output":"127.0.0.1:5000"}`),

			wantErr: true,
		},
		{
			name:   "download error",
			enable: true,

			df: func(url string, d time.Duration, c afd.DownloadCallback) (err error) {
				return errors.New("download error")
			},

			sc:           new(MockStreamerCreator),
			creatorInput: []interface{}{},
			creatorOutputFunc: func() []interface{} {
				return []interface{}{}
			},

			in: bytes.NewBufferString(`{"url":"http://127.0.0.1:8080/index.html","timeout":"1s"}`),

			wantErr: true,
		},
		{
			name:   "create streamer error",
			enable: true,

			df: defaultCallback,

			sc: new(MockStreamerCreator),
			creatorInput: []interface{}{
				"127.0.0.1:5000",
			},
			creatorOutputFunc: func() []interface{} {
				return []interface{}{
					(*afd.MockStreamer)(nil),
					errors.New("dial network error"),
				}
			},

			in: bytes.NewBufferString(`{"url":"http://127.0.0.1:8080/index.html","timeout":"1s","output":"127.0.0.1:5000"}`),

			wantErr: true,
		},
		{
			name:   "copy body error",
			enable: true,

			df: defaultCallback,

			sc: new(MockStreamerCreator),
			creatorInput: []interface{}{
				"127.0.0.1:5000",
			},
			creatorOutputFunc: func() []interface{} {
				s := new(afd.MockStreamer)
				s.
					On("Write", []interface{}{[]byte(`{}`)}...).
					Return([]interface{}{0, errors.New("write error")}...)
				s.
					On("Close").
					Return([]interface{}{(error)(nil)}...)

				return []interface{}{
					s,
					(error)(nil),
				}
			},

			in: bytes.NewBufferString(`{"url":"http://127.0.0.1:8080/index.html","timeout":"1s","output":"127.0.0.1:5000"}`),

			wantErr: true,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			if !test.enable {
				t.SkipNow()
			}

			test.sc.
				On("CreateStreamer", test.creatorInput...).
				Return(test.creatorOutputFunc()...)

			err := NewDownloadService(test.df, test.sc).Download(test.in)
			if (err != nil) != test.wantErr {
				t.Error(err)
				t.FailNow()
			}
		})
	}
}
