package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	ErrDecodeInput = iota + 1
	ErrEncodeOutput
	ErrCreateRequest
	ErrDoRequest
	ErrBodyClose
	ErrDialConn
	ErrCloseConn
	ErrWriteOutput
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

	return errors.New("invalid duration")
}

type UtilityInput struct {
	IgnoreSSLCertificates bool     `json:"ignore-ssl-certificates"`
	FollowRedirects       bool     `json:"follow-redirects"`
	URL                   string   `json:"url"`
	Output                string   `json:"output"`
	Timeout               Duration `json:"timeout"`
}

type UtilityOutput struct {
	Success       bool     `json:"success"`
	HTTPCode      int      `json:"http-code"`
	ContentLength int64    `json:"content-length"`
	ContentType   string   `json:"content-type"`
	ErrorMessage  string   `json:"error-message"`
	Redirects     []string `json:"redirects"`
}

// 1. Receive input from stdin +
// 2. Unmarshal input to struct +
// 3. Create request +
// 4. Send request +
// 5. Receive response +
// 6. Forward response body to the output in background +
// 7. Create program result +
// 8. Marshal program result +
// 9. Print program result into stdout +
//
// Additional:
// 1. Control redirects
// 2. Tests
// 3. Validate input HTTP URL
// 4. Validate output TCP URL
func main() {
	var (
		in  = &UtilityInput{}
		out = &UtilityOutput{}
	)

	if err := json.NewDecoder(os.Stdin).Decode(in); err != nil {
		out.ErrorMessage = "decode input error: " + err.Error()
		writeOutput(os.Stdout, out)
		os.Exit(ErrDecodeInput)
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Duration(in.Timeout)))
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, in.URL, nil)
	if err != nil {
		out.ErrorMessage = "create request error: " + err.Error()
		writeOutput(os.Stdout, out)
		os.Exit(ErrCreateRequest)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		out.ErrorMessage = "do request error: " + err.Error()
		writeOutput(os.Stdout, out)
		os.Exit(ErrDoRequest)
	}

	defer func(cl io.Closer) {
		if closeErr := res.Body.Close(); closeErr != nil {
			_, _ = fmt.Fprintf(os.Stderr, "close body error: %v \n", err)
			os.Exit(ErrBodyClose)
		}
	}(res.Body)

	quit := make(chan struct{}, 1)

	go func(r io.Reader, url string, ch chan<- struct{}) {
		defer func(ch chan<- struct{}) {
			ch <- struct{}{}
			close(ch)
		}(ch)

		if url == "" {
			return
		}

		conn, err := net.Dial("tcp", url)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "dial conn error: %v \n", err)
			os.Exit(ErrDialConn)
		}

		defer func(cl io.Closer) {
			if closeErr := cl.Close(); closeErr != nil {
				_, _ = fmt.Fprintf(os.Stderr, "close tcp conn error: %v \n", err)
				os.Exit(ErrCloseConn)
			}
		}(conn)

		if _, err = io.Copy(conn, r); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "write output data error: %v \n", err)
			os.Exit(ErrWriteOutput)
		}
	}(res.Body, in.Output, quit)

	out.Success = true
	out.HTTPCode = res.StatusCode
	out.ContentLength = res.ContentLength
	out.ContentType = res.Header.Get("Content-Type")

	writeOutput(os.Stdout, out)

	<-quit
}

func writeOutput(w io.Writer, out *UtilityOutput) {
	if err := json.NewEncoder(w).Encode(out); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "encode output error: %v \n", err)
		os.Exit(ErrEncodeOutput)
	}
}
