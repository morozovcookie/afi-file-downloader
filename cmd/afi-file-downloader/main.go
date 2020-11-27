package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	afd "github.com/morozovcookie/afifiledownloader"
	"github.com/morozovcookie/afifiledownloader/cli"
	"github.com/morozovcookie/afifiledownloader/http"
	"github.com/morozovcookie/afifiledownloader/tcp"
)

type Output struct {
	Success       bool     `json:"success"`
	HTTPCode      int      `json:"http-code,omitempty"`
	ContentLength int64    `json:"content-length,omitempty"`
	ContentType   string   `json:"content-type,omitempty"`
	ErrorMessage  string   `json:"error-message,omitempty"`
	Redirects     []string `json:"redirects,omitempty"`
}

// Additional:
// 1. Control redirects
// 2. SSL
// 3. Code docs
// 4. Project docs

func main() {
	var (
		out = &Output{Success: true}

		err error
	)

	defer func(err *error) {
		if *err == nil {
			return
		}

		if encodeErr := json.NewEncoder(os.Stdout).Encode(&Output{ErrorMessage: (*err).Error()}); encodeErr != nil {
			_, _ = fmt.Fprintf(os.Stderr, "encode output error: %v \n", encodeErr)
		}
	}(&err)

	svc := cli.NewDownloadService(downloaderCreator(out), tcp.NewStreamer)

	if err = svc.Download(os.Stdin); err != nil {
		return
	}

	if err = json.NewEncoder(os.Stdout).Encode(out); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "encode output error: %v \n", err)
	}
}

func downloaderCreator(out *Output) cli.DownloaderCreator {
	return func(isFollowRedirects bool, maxRedirects int) afd.DownloadFunc {
		if isFollowRedirects {
			return func(url string, timeout time.Duration, c afd.DownloadCallback) (err error) {
				out.HTTPCode, out.ContentLength, out.ContentType, out.Redirects, err = http.
					NewRedirectDownloader(maxRedirects).
					Download(url, timeout, c)
				if err != nil {
					return err
				}

				return nil
			}
		}

		return func(url string, timeout time.Duration, c afd.DownloadCallback) (err error) {
			out.HTTPCode, out.ContentLength, out.ContentType, err = http.NewDownloader().Download(url, timeout, c)
			if err != nil {
				return err
			}

			return nil
		}
	}
}
