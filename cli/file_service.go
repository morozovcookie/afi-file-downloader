package cli

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"time"

	afifiledownloader "github.com/morozovcookie/afi-file-downloader"
)

var ErrUnsupportedMethod = errors.New("unsupported method")

const (
	MethodGet     = "GET"
	MethodHead    = "HEAD"
	DefaultMethod = MethodGet
)

type FileServiceCreatorFunc func(insecure, redirects bool, redirectsCount int64) afifiledownloader.FileService

type FileStreamerCreatorFunc func(address string) (afifiledownloader.FileStreamer, error)

type HandlerFunc func(w io.Writer, r io.Reader) error

type Service struct {
	fileServiceCreator FileServiceCreatorFunc
	streamerCreator    FileStreamerCreatorFunc

	handlers map[string]HandlerFunc
}

func NewFileService(fileServiceCreator FileServiceCreatorFunc, streamerCreator FileStreamerCreatorFunc) *Service {
	svc := &Service{
		fileServiceCreator: fileServiceCreator,
		streamerCreator:    streamerCreator,
	}

	svc.handlers = map[string]HandlerFunc{
		MethodGet:  svc.DownloadFileHandler,
		MethodHead: svc.GetFileStatHandler,
	}

	return svc
}

func (svc *Service) Serve(w io.Writer, r io.Reader) error {
	req := &struct {
		Method string `json:"method"`
	}{
		Method: DefaultMethod,
	}

	if err := json.NewDecoder(r).Decode(req); err != nil {
		return err
	}

	handler, ok := svc.handlers[req.Method]
	if !ok {
		return ErrUnsupportedMethod
	}

	return handler(w, r)
}

func (svc *Service) DownloadFileHandler(w io.Writer, r io.Reader) error {
	req := &GetRequest{
		Request: &Request{
			MaxRedirects: DefaultMaxRedirects,
			Method:       DefaultMethod,
			Timeout:      DefaultTimeout,
		},
	}

	if err := json.NewDecoder(r).Decode(req); err != nil {
		return err
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(req.Timeout.Duration()))
	defer cancel()

	callback := func(httpResp afifiledownloader.FileDownloadResponse) error {
		resp := &GetResponse{
			Response: &Response{
				Success:       true,
				HTTPCode:      httpResp.StatusCode(),
				ContentLength: httpResp.ContentLength(),
				Redirects:     httpResp.Redirects(),
			},
			ContentType: httpResp.ContentType(),
		}

		defer httpResp.Close()

		if req.Output == "" {
			return json.NewEncoder(w).Encode(resp)
		}

		streamer, err := svc.streamerCreator(req.Output)
		if err != nil {
			return err
		}

		defer streamer.Close()

		if err = streamer.Stream(httpResp); err != nil {
			return err
		}

		return json.NewEncoder(w).Encode(resp)
	}

	return svc.
		fileServiceCreator(req.IsIgnoreSSLCertificates, req.IsFollowRedirects, req.MaxRedirects).
		DownloadFile(ctx, req.URL, callback)
}

func (svc *Service) GetFileStatHandler(w io.Writer, r io.Reader) error {
	req := &HeadRequest{
		Request: &Request{
			MaxRedirects: DefaultMaxRedirects,
			Method:       DefaultMethod,
			Timeout:      DefaultTimeout,
		},
	}

	if err := json.NewDecoder(r).Decode(req); err != nil {
		return err
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(req.Timeout.Duration()))
	defer cancel()

	callback := func(httpResp afifiledownloader.FileDownloadResponse) error {
		var (
			resp = &HeadResponse{
				Response: &Response{
					Success:       true,
					HTTPCode:      httpResp.StatusCode(),
					ContentLength: httpResp.ContentLength(),
					Redirects:     httpResp.Redirects(),
				},
				Headers: make([]string, len(httpResp.Headers())),
			}
			sb = strings.Builder{}
		)

		for name, value := range httpResp.Headers() {
			sb.WriteString(name)
			sb.WriteString(": ")
			sb.WriteString(strings.Join(value, ", "))

			resp.Headers = append(resp.Headers, sb.String())
			sb.Reset()
		}

		defer httpResp.Close()

		return json.NewEncoder(w).Encode(resp)
	}

	return svc.
		fileServiceCreator(req.IsIgnoreSSLCertificates, req.IsFollowRedirects, req.MaxRedirects).
		GetFileStat(ctx, req.URL, callback)
}
