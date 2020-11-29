package http

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	afd "github.com/morozovcookie/afifiledownloader"
)

func TestDownloader_Download(t *testing.T) {
	tt := []struct {
		name    string
		enabled bool

		srvHandlerPattern string
		srvHandler        func(w http.ResponseWriter, r *http.Request)

		url      func(string) string
		timeout  time.Duration
		callback afd.DownloadCallback

		wantErr bool

		expectedStatus        int
		expectedContentLength int64
		expectedContentType   string
	}{
		{
			name:    "pass",
			enabled: true,

			srvHandlerPattern: "/",
			srvHandler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Content-Type", "application/json")

				if _, err := w.Write([]byte(`{}`)); err != nil {
					t.Error(err)
				}
			},

			url: func(srv string) string {
				return srv + "/index.html"
			},
			timeout: time.Second,
			callback: func(r *http.Response) (err error) {
				return nil
			},

			expectedStatus:        http.StatusOK,
			expectedContentLength: int64(len([]byte(`{}`))),
			expectedContentType:   "application/json",
		},
		{
			name:    "create request error",
			enabled: true,

			srvHandlerPattern: "/",
			srvHandler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Content-Type", "application/json")

				if _, err := w.Write([]byte(`{}`)); err != nil {
					t.Error(err)
				}
			},

			url: func(_ string) string {
				return "ffs^*&^*(U://"
			},
			timeout: time.Second,
			callback: func(r *http.Response) (err error) {
				return nil
			},

			wantErr: true,
		},
		{
			name:    "execute request error",
			enabled: true,

			srvHandlerPattern: "/",
			srvHandler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Content-Type", "application/json")

				if _, err := w.Write([]byte(`{}`)); err != nil {
					t.Error(err)
				}
			},

			url: func(_ string) string {
				return ""
			},
			timeout: time.Second,
			callback: func(r *http.Response) (err error) {
				return nil
			},

			wantErr: true,
		},
		{
			name:    "execute callback error",
			enabled: true,

			srvHandlerPattern: "/",
			srvHandler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Content-Type", "application/json")

				if _, err := w.Write([]byte(`{}`)); err != nil {
					t.Error(err)
				}
			},

			url: func(srv string) string {
				return srv + "/index.html"
			},
			timeout: time.Second,
			callback: func(r *http.Response) (err error) {
				return errors.New("callback error")
			},

			wantErr: true,
		},
		{
			name:    "redirect",
			enabled: true,

			srvHandlerPattern: "/",
			srvHandler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Location", r.URL.String())
				w.WriteHeader(http.StatusMovedPermanently)
			},

			url: func(srv string) string {
				return srv + "/index.html"
			},
			timeout: time.Second,
			callback: func(r *http.Response) (err error) {
				return nil
			},

			expectedStatus: http.StatusMovedPermanently,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			if !test.enabled {
				t.SkipNow()
			}

			mux := http.NewServeMux()
			mux.HandleFunc(test.srvHandlerPattern, test.srvHandler)

			srv := httptest.NewServer(mux)
			defer srv.Close()

			downloader := NewDownloader(false)
			actualStatus, actualContentLength, actualContentType, err := downloader.Download(
				test.url(srv.URL), test.timeout, test.callback)
			if (err != nil) != test.wantErr {
				t.Error(err)
				t.FailNow()
			}

			assert.Equal(t, test.expectedStatus, actualStatus)
			assert.Equal(t, test.expectedContentLength, actualContentLength)
			assert.Equal(t, test.expectedContentType, actualContentType)
		})
	}
}
