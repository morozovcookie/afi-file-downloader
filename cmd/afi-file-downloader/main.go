package main

import (
	"context"
	"log"

	"github.com/morozovcookie/afi-file-downloader/http"
	"github.com/morozovcookie/afi-file-downloader/tcp"
)

func main() {
	fileSvc := http.NewFileService()
	fileStreamer, err := tcp.NewFileStreamer("")
	if err != nil {
		log.Fatal(err)
	}
	defer fileStreamer.Close()

	err = fileSvc.FindFile(context.Background(), "", fileStreamer)
	if err != nil {
		log.Fatal(err)
	}
}
