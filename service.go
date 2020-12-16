package afifiledownloader

import (
	"context"
	"io"
)

type FileStat map[string]string

// FileService represent a 2nd edition of service for file operating.
type FileService interface {
	// FindFile download file by specified path.
	FindFile(ctx context.Context, name string, dst io.ReaderFrom) error		// <-- http implementation

	// GetFileStat return file statistics by specified path.
	GetFileStat(ctx context.Context, name string) (*FileStat, error)		// <-- http implementation
}
