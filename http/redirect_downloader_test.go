package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	afd "github.com/morozovcookie/afifiledownloader"
)

func TestRedirectDownloader_Download(t *testing.T) {
	tt := []struct {
		name    string
		enabled bool

		srvHandlerPattern string
		srvHandler        func(w http.ResponseWriter, r *http.Request)

		redirects int64

		url      func(string) string
		timeout  time.Duration
		callback afd.DownloadCallback

		wantErr bool

		expectedStatus        int
		expectedContentLength int64
		expectedContentType   string
		expectedRedirects     func(srv string) []string
	}{
		{
			name:    "pass",
			enabled: true,

			srvHandlerPattern: "/",
			srvHandler: func(w http.ResponseWriter, r *http.Request) {
				handlers := []func(w http.ResponseWriter, r *http.Request){
					func(w http.ResponseWriter, r *http.Request) {
						var (
							u = r.URL
							q = u.Query()
						)

						q.Add("redirect", "1")
						u.RawQuery = q.Encode()

						w.Header().Add("Location", u.String())
						w.WriteHeader(http.StatusMovedPermanently)
					},
					func(w http.ResponseWriter, r *http.Request) {
						w.Header().Add("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						if _, err := w.Write([]byte(`{}`)); err != nil {
							t.Fatal(err)
						}
					},
				}

				redirectParam := r.FormValue("redirect")
				if redirectParam == "" {
					handlers[0](w, r)
					return
				}

				redirect, err := strconv.ParseInt(redirectParam, 10, 64)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				handlers[redirect](w, r)
			},

			redirects: 1,

			url: func(srv string) string {
				return srv + "/index.html"
			},
			timeout: time.Hour,
			callback: func(r *http.Response) (err error) {
				return nil
			},

			expectedStatus:        http.StatusOK,
			expectedContentLength: int64(len([]byte(`{}`))),
			expectedContentType:   "application/json",
			expectedRedirects: func(srv string) []string {
				return []string{
					srv + "/index.html?redirect=1",
				}
			},
		},
		{
			name:    "too many redirects",
			enabled: true,

			srvHandlerPattern: "/",
			srvHandler: func(w http.ResponseWriter, r *http.Request) {
				handlers := []func(w http.ResponseWriter, r *http.Request){
					func(w http.ResponseWriter, r *http.Request) {
						var (
							u = r.URL
							q = u.Query()
						)

						q.Add("redirect", "1")
						u.RawQuery = q.Encode()

						w.Header().Add("Location", u.String())
						w.WriteHeader(http.StatusMovedPermanently)
					},
					func(w http.ResponseWriter, r *http.Request) {
						var (
							u = r.URL
							q = u.Query()
						)

						q.Add("redirect", "2")
						u.RawQuery = q.Encode()

						w.Header().Add("Location", u.String())
						w.WriteHeader(http.StatusMovedPermanently)
					},
				}

				redirectParam := r.FormValue("redirect")
				if redirectParam == "" {
					handlers[0](w, r)
					return
				}

				redirect, err := strconv.ParseInt(redirectParam, 10, 64)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				handlers[redirect](w, r)
			},

			redirects: 1,

			url: func(srv string) string {
				return srv + "/index.html"
			},
			timeout: time.Hour,
			callback: func(_ *http.Response) (err error) {
				return nil
			},

			wantErr: true,

			expectedRedirects: func(_ string) []string {
				return nil
			},
		},
		{
			name:    "create request error",
			enabled: true,

			srvHandlerPattern: "/",
			srvHandler: func(w http.ResponseWriter, r *http.Request) {
				handlers := []func(w http.ResponseWriter, r *http.Request){
					func(w http.ResponseWriter, r *http.Request) {
						var (
							u = r.URL
							q = u.Query()
						)

						q.Add("redirect", "1")
						u.RawQuery = q.Encode()

						w.Header().Add("Location", u.String())
						w.WriteHeader(http.StatusMovedPermanently)
					},
					func(w http.ResponseWriter, r *http.Request) {
						w.Header().Add("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						if _, err := w.Write([]byte(`{}`)); err != nil {
							t.Fatal(err)
						}
					},
				}

				redirectParam := r.FormValue("redirect")
				if redirectParam == "" {
					handlers[0](w, r)
					return
				}

				redirect, err := strconv.ParseInt(redirectParam, 10, 64)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				handlers[redirect](w, r)
			},

			redirects: 1,

			url: func(srv string) string {
				return "ffs^*&^*(U://"
			},
			timeout: time.Hour,
			callback: func(r *http.Response) (err error) {
				return nil
			},

			wantErr: true,

			expectedRedirects: func(srv string) []string {
				return nil
			},
		},
		{
			name:    "send request error",
			enabled: true,

			srvHandlerPattern: "/",
			srvHandler: func(w http.ResponseWriter, r *http.Request) {
				handlers := []func(w http.ResponseWriter, r *http.Request){
					func(w http.ResponseWriter, r *http.Request) {
						var (
							u = r.URL
							q = u.Query()
						)

						q.Add("redirect", "1")
						u.RawQuery = q.Encode()

						w.Header().Add("Location", u.String())
						w.WriteHeader(http.StatusMovedPermanently)
					},
					func(w http.ResponseWriter, r *http.Request) {
						w.Header().Add("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						if _, err := w.Write([]byte(`{}`)); err != nil {
							t.Fatal(err)
						}
					},
				}

				redirectParam := r.FormValue("redirect")
				if redirectParam == "" {
					handlers[0](w, r)
					return
				}

				redirect, err := strconv.ParseInt(redirectParam, 10, 64)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				handlers[redirect](w, r)
			},

			redirects: 1,

			url: func(srv string) string {
				return ""
			},
			timeout: time.Hour,
			callback: func(r *http.Response) (err error) {
				return nil
			},

			wantErr: true,

			expectedRedirects: func(srv string) []string {
				return nil
			},
		},
		{
			name:    "cyclic requests",
			enabled: true,

			srvHandlerPattern: "/",
			srvHandler: func(w http.ResponseWriter, r *http.Request) {
				handlers := []func(w http.ResponseWriter, r *http.Request){
					func(w http.ResponseWriter, r *http.Request) {
						w.Header().Add("Location", r.URL.String())
						w.WriteHeader(http.StatusMovedPermanently)
					},
					func(w http.ResponseWriter, r *http.Request) {
						w.Header().Add("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						if _, err := w.Write([]byte(`{}`)); err != nil {
							t.Fatal(err)
						}
					},
				}

				redirectParam := r.FormValue("redirect")
				if redirectParam == "" {
					handlers[0](w, r)
					return
				}

				redirect, err := strconv.ParseInt(redirectParam, 10, 64)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				handlers[redirect](w, r)
			},

			redirects: 1,

			url: func(srv string) string {
				return srv + "/index.html"
			},
			timeout: time.Hour,
			callback: func(r *http.Response) (err error) {
				return nil
			},

			wantErr: true,

			expectedRedirects: func(_ string) []string {
				return nil
			},
		},
		{
			name:    "callback error",
			enabled: true,

			srvHandlerPattern: "/",
			srvHandler: func(w http.ResponseWriter, r *http.Request) {
				handlers := []func(w http.ResponseWriter, r *http.Request){
					func(w http.ResponseWriter, r *http.Request) {
						var (
							u = r.URL
							q = u.Query()
						)

						q.Add("redirect", "1")
						u.RawQuery = q.Encode()

						w.Header().Add("Location", u.String())
						w.WriteHeader(http.StatusMovedPermanently)
					},
					func(w http.ResponseWriter, r *http.Request) {
						w.Header().Add("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						if _, err := w.Write([]byte(`{}`)); err != nil {
							t.Fatal(err)
						}
					},
				}

				redirectParam := r.FormValue("redirect")
				if redirectParam == "" {
					handlers[0](w, r)
					return
				}

				redirect, err := strconv.ParseInt(redirectParam, 10, 64)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				handlers[redirect](w, r)
			},

			redirects: 1,

			url: func(srv string) string {
				return srv + "/index.html"
			},
			timeout: time.Hour,
			callback: func(r *http.Response) (err error) {
				return errors.New("callback error")
			},

			wantErr: true,

			expectedRedirects: func(srv string) []string {
				return nil
			},
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

			downloader := NewRedirectDownloader(test.redirects)
			actualStatus, actualContentLength, actualContentType, actualRedirects, err := downloader.Download(
				test.url(srv.URL), test.timeout, test.callback)
			if (err != nil) != test.wantErr {
				t.Error(err)
				t.FailNow()
			}

			assert.Equal(t, test.expectedStatus, actualStatus)
			assert.Equal(t, test.expectedContentLength, actualContentLength)
			assert.Equal(t, test.expectedContentType, actualContentType)
			assert.Equal(t, test.expectedRedirects(srv.URL), actualRedirects)
		})
	}
}
