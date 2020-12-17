package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	stdhttp "net/http"
	"os"
	"strings"
	"time"

	"github.com/morozovcookie/afi-file-downloader/http"
)

var (
	ErrInvalidDuration   = errors.New("invalid duration")
	ErrUnsupportedMethod = errors.New("unsupported method")
)

const (
	DefaultMaxRedirects = 5
	DefaultTimeout      = Duration(time.Second)

	MethodUnknown = "UNKNOWN"
	MethodGet     = "GET"
	MethodHead    = "HEAD"
	DefaultMethod = MethodGet
)

type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	if val, ok := v.(float64); ok {
		*d = Duration(time.Duration(val))

		return nil
	}

	if val, ok := v.(string); ok {
		t, err := time.ParseDuration(val)
		if err != nil {
			return err
		}

		*d = Duration(t)

		return nil
	}

	return ErrInvalidDuration
}

func (d Duration) Duration() time.Duration {
	return time.Duration(d)
}

func main() {
	var err error

	if err = serveRequest(os.Stdout, os.Stderr); err == nil {
		return
	}

	resp := &struct {
		Success      bool   `json:"success"`
		ErrorMessage string `json:"error-message"`
	}{
		Success:      false,
		ErrorMessage: err.Error(),
	}

	if errEncode := json.NewEncoder(os.Stdout).Encode(resp); errEncode != nil {
		_, _ = fmt.Fprintln(os.Stderr, errEncode)
	}
}

func serveRequest(w io.Writer, r io.Reader) error {
	req := &struct {
		Method string `json:"method"`
	}{
		Method: DefaultMethod,
	}

	if err := json.NewDecoder(r).Decode(req); err != nil {
		return err
	}

	handlers := map[string]func(w io.Writer, r io.Reader) error{
		MethodGet:  handleGetRequest,
		MethodHead: handleHeadRequest,
	}

	handler, ok := handlers[req.Method]
	if !ok {
		return ErrUnsupportedMethod
	}

	return handler(w, r)
}

type Request struct {
	IsIgnoreSSLCertificates bool     `json:"ignore-ssl-certificates"`
	IsFollowRedirects       bool     `json:"follow-redirects"`
	MaxRedirects            int64    `json:"max-redirects"`
	Method                  string   `json:"method"`
	URL                     string   `json:"url"`
	Timeout                 Duration `json:"timeout"`
}

type GetRequest struct {
	*Request

	Output string `json:"output"`
}

type HeadRequest struct {
	*Request

	Headers []string `json:"headers"`
}

func handleGetRequest(w io.Writer, r io.Reader) error {
	req := &GetRequest{
		Request: &Request{
			MaxRedirects: DefaultMaxRedirects,
			Method:       DefaultMethod,
			Timeout:      DefaultTimeout,
		},
	}

	if err := json.NewDecoder(r).Decode(req); err != nil {
		return err
	}

	callback := func(httpResp *http.Response) error {
		resp := &struct {
			Success       bool     `json:"success"`
			HTTPCode      int      `json:"http-code,omitempty"`
			ContentLength int64    `json:"content-length,omitempty"`
			ContentType   string   `json:"content-type,omitempty"`
			Redirects     []string `json:"redirects,omitempty"`
		}{
			Success:       true,
			HTTPCode:      httpResp.StatusCode,
			ContentLength: httpResp.ContentLength,
			ContentType:   httpResp.Header.Get("Content-Type"),
			Redirects:     httpResp.Redirects,
		}

		defer httpResp.Body.Close()

		if req.Output == "" {
			return json.NewEncoder(w).Encode(resp)
		}

		conn, err := net.Dial("tcp", req.Output)
		if err != nil {
			return err
		}

		defer conn.Close()

		if _, err = io.Copy(conn, httpResp.Body); err != nil {
			return err
		}

		return json.NewEncoder(w).Encode(resp)
	}

	return makeRequest(req.Request, callback)
}

func handleHeadRequest(w io.Writer, r io.Reader) error {
	req := &HeadRequest{
		Request: &Request{
			MaxRedirects: DefaultMaxRedirects,
			Method:       DefaultMethod,
			Timeout:      DefaultTimeout,
		},
	}

	if err := json.NewDecoder(r).Decode(req); err != nil {
		return err
	}

	callback := func(httpResp *http.Response) error {
		var (
			resp = &struct {
				Success       bool     `json:"success"`
				HTTPCode      int      `json:"http-code,omitempty"`
				ContentLength int64    `json:"content-length,omitempty"`
				Redirects     []string `json:"redirects,omitempty"`
				Headers       []string `json:"headers"`
			}{
				Success:       true,
				HTTPCode:      httpResp.StatusCode,
				ContentLength: httpResp.ContentLength,
				Redirects:     httpResp.Redirects,
				Headers:       make([]string, len(httpResp.Header)),
			}
			sb = strings.Builder{}
		)

		for name, value := range httpResp.Header {
			sb.WriteString(name)
			sb.WriteString(": ")
			sb.WriteString(strings.Join(value, ", "))

			resp.Headers = append(resp.Headers, sb.String())
			sb.Reset()
		}

		defer httpResp.Body.Close()

		return json.NewEncoder(w).Encode(resp)
	}

	return makeRequest(req.Request, callback)
}

func makeRequest(req *Request, callback func(httpResp *http.Response) error) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(req.Timeout.Duration()))
	defer cancel()

	httpReq, err := stdhttp.NewRequestWithContext(ctx, req.Method, req.URL, nil)
	if err != nil {
		return err
	}

	httpClient := http.NewClientFactory().Create(req.IsIgnoreSSLCertificates, req.IsFollowRedirects)
	httpClient.SetRedirects(req.MaxRedirects)

	return httpClient.Do(httpReq, callback)
}
