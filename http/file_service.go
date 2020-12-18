package http

import (
	"context"
	"net/http"

	afifiledownloader "github.com/morozovcookie/afi-file-downloader"
)

type FileService struct {
	client Client
}

func NewFileService(client Client) afifiledownloader.FileService {
	return &FileService{
		client: client,
	}
}

func (svc *FileService) DownloadFile(ctx context.Context, name string, callback afifiledownloader.CallbackFunc) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, name, nil)
	if err != nil {
		return err
	}

	return svc.client.Do(req, callback)
}

func (svc *FileService) GetFileStat(ctx context.Context, name string, callback afifiledownloader.CallbackFunc) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, name, nil)
	if err != nil {
		return err
	}

	return svc.client.Do(req, callback)
}
