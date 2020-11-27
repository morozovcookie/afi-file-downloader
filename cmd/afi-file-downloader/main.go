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
// 1. Dockerfile (or werf.yaml)
// 2. Control redirects
// 3. SSL
// 4. Code docs
// 5. Project docs

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

	var svc *cli.DownloadService
	{
		downloader := createDownloadFunc(http.NewDownloader(), out)
		creator := func(address string) (s afd.Streamer, err error) {
			return tcp.NewStreamer(address)
		}

		svc = cli.NewDownloadService(downloader, creator)
	}

	if err = svc.Download(os.Stdin); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "download error: %v \n", err)

		return
	}

	if err = json.NewEncoder(os.Stdout).Encode(out); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "encode output error: %v \n", err)
	}
}

func createDownloadFunc(d *http.Downloader, out *Output) afd.DownloadFunc {
	return func(url string, timeout time.Duration, c afd.DownloadCallback) (err error) {
		if out.HTTPCode, out.ContentLength, out.ContentType, err = d.Download(url, timeout, c); err != nil {
			return err
		}

		return nil
	}
}
