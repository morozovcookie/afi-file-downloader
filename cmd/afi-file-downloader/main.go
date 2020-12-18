package main

import (
	"encoding/json"
	"fmt"
	"os"

	afifiledownloader "github.com/morozovcookie/afi-file-downloader"
	"github.com/morozovcookie/afi-file-downloader/cli"
	"github.com/morozovcookie/afi-file-downloader/http"
	"github.com/morozovcookie/afi-file-downloader/tcp"
)

func main() {
	var (
		fileServiceCreator = func(insecure, redirects bool, redirectsCount int64) afifiledownloader.FileService {
			client := http.NewClientFactory().Create(insecure, redirects)
			client.SetRedirects(redirectsCount)

			return http.NewFileService(client)
		}

		fileSvc = cli.NewFileService(fileServiceCreator, tcp.NewStreamer)

		err error
	)

	if err = fileSvc.Serve(os.Stdout, os.Stdin); err == nil {
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
