package cli

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"time"
)

const (
	MethodGet  = "GET"
	MethodHead = "HEAD"
)

var (
	ErrUnsupportedMethod = errors.New("unsupported method")
)

type HandlerFunc func(w io.Writer, r io.Reader) (err error)

type FileService struct {
	handlers map[string]HandlerFunc
}

func NewFileService() *FileService {
	return &FileService{}
}

func (fs *FileService) Serve(w io.Writer, r io.Reader) (err error) {
	return encodeErr(w, fs.serve(w, r))
}

func (fs *FileService) serve(w io.Writer, r io.Reader) (err error) {
	req := &struct {
		Method string `json:"method"`
	}{
		Method: MethodGet,
	}

	if err = json.NewDecoder(r).Decode(req); err != nil {
		return err
	}

	handler, ok := fs.handlers[req.Method]
	if !ok {
		return ErrUnsupportedMethod
	}

	return handler(w, r)
}

func encodeErr(w io.Writer, err error) error {
	if err == nil {
		return nil
	}

	resp := &struct{
		Success      bool   `json:"success"`
		ErrorMessage string `json:"error-message"`
	}{
		ErrorMessage: err.Error(),
	}

	return json.NewEncoder(w).Encode(resp)
}

func (fs *FileService) findFileHandler(w io.Writer, r io.Reader) (err error) {
	req := &struct{
		IsIgnoreSSLCertificates bool     `json:"ignore-ssl-certificates"`  // <-- client parameter
		IsFollowRedirects       bool     `json:"follow-redirects"`         // <-- client parameter
		MaxRedirects            int64    `json:"max-redirects"`            // <-- client parameter
		URL                     string   `json:"url"`                      // <-- request parameter
		Output                  string   `json:"output"`                   // <-- streamer parameter
		Timeout                 Duration `json:"timeout"`                  // +
	}{}

	if err = json.NewDecoder(r).Decode(req); err != nil {
		return err
	}

	_, cancel := context.WithDeadline(context.Background(), time.Now().Add(req.Timeout.Duration()))
	defer cancel()

	// create client

	// create file_service

	// create streamer

	// call FindService()

	resp := &struct {
		Success       bool     `json:"success"`                   //
		HTTPCode      int      `json:"http-code,omitempty"`       // <-- low-level field
		ContentLength int64    `json:"content-length,omitempty"`  //
		ContentType   string   `json:"content-type,omitempty"`    // <-- low-level field
		Redirects     []string `json:"redirects,omitempty"`       // <-- low-level field
	}{
		Success: true,
	}

	return json.NewEncoder(w).Encode(resp)
}

func (fs *FileService) getFileStatHandler(w io.Writer, r io.Reader) (err error) {
	req := &struct {
		IsIgnoreSSLCertificates bool     `json:"ignore-ssl-certificates"`  // <-- client parameter
		IsFollowRedirects       bool     `json:"follow-redirects"`         // <-- client parameter
		MaxRedirects            int64    `json:"max-redirects"`            // <-- client parameter
		URL                     string   `json:"url"`                      // <-- request parameter
		Timeout                 Duration `json:"timeout"`                  // +
		Headers                 []string `json:"headers"`                  // <-- request parameter
	}{}

	if err = json.NewDecoder(r).Decode(req); err != nil {
		return err
	}

	_, cancel := context.WithDeadline(context.Background(), time.Now().Add(req.Timeout.Duration()))
	defer cancel()

	// create client

	// create file_service

	// call GetFileStat()

	resp := &struct{
		Success       bool     `json:"success"`                   //
		HTTPCode      int      `json:"http-code,omitempty"`       // <-- low-level field
		ContentLength int64    `json:"content-length"`            //
		Redirects     []string `json:"redirects,omitempty"`       // <-- low-level field
		Headers       []string `json:"headers,omitempty"`         //
	}{
		Success: true,
	}

	return json.NewEncoder(w).Encode(resp)
}
