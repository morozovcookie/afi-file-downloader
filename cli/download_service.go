package cli

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type DownloadService struct {
	dc DownloaderCreator
	sc StreamerCreator
}

func NewDownloadService(dc DownloaderCreator, sc StreamerCreator) *DownloadService {
	return &DownloadService{
		dc: dc,
		sc: sc,
	}
}

func (svc *DownloadService) Download(r io.Reader) (err error) {
	in := &Input{
		MaxRedirects: DefaultMaxRedirects,
		Timeout:      DefaultTimeout,
	}

	if err = json.NewDecoder(r).Decode(in); err != nil {
		return err
	}

	if err = in.Validate(); err != nil {
		return err
	}

	callback := func(res *http.Response) (err error) {
		defer res.Body.Close()

		if in.Output == "" {
			return nil
		}

		s, err := svc.sc(in.Output)
		if err != nil {
			return err
		}

		defer s.Close()

		if _, err = io.Copy(s, res.Body); err != nil {
			return err
		}

		return nil
	}

	err = svc.dc(in.IsFollowRedirects, in.MaxRedirects, in.IsIgnoreSSLCertificates)(
		in.URL, time.Duration(in.Timeout), callback)
	if err != nil {
		return err
	}

	return nil
}
