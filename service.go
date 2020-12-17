package afifiledownloader

import (
	"context"
	"io"
)

//----------------------------------------------------------------------------------------------------------------------

// FindFileFunc represent a function which implement file downloading by specified path.
type FindFileFunc func(ctx context.Context, name string, dst io.Writer) error

//----------------------------------------------------------------------------------------------------------------------

// FilePropertyKey is a file attribute name.
type FilePropertyKey string

// FilePropertyValue is a file attribute value.
type FilePropertyValue string

// FileStat is a file attributes container.
type FileStat map[FilePropertyKey]FilePropertyValue

// GetFileStatFunc represent a function which implement file statistics retrieving by specified path.
type GetFileStatFunc func(ctx context.Context, name string) (FileStat, error)

//----------------------------------------------------------------------------------------------------------------------

// StreamFileFunc represent a function which implement file streaming from source to destination.
type StreamFileFunc func(ctx context.Context, src io.Reader, dst io.Writer) error

//----------------------------------------------------------------------------------------------------------------------
