package http

import (
	"context"
	"io"
	"net/http"

	afifiledownloader "github.com/morozovcookie/afi-file-downloader"
)

type FileService struct{}

func NewFileService() *FileService {
	return &FileService{}
}

// FindFile download file by specified path.
func (svc *FileService) FindFile(ctx context.Context, name string, dst io.ReaderFrom) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, name, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if _, err = dst.ReadFrom(resp.Body); err != nil {
		return err
	}

	return nil
}

// GetFileStat return file statistics by specified path.
func (svc *FileService) GetFileStat(ctx context.Context, name string) (*afifiledownloader.FileStat, error) {
	return nil, nil
}
