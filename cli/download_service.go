package cli

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	afd "github.com/morozovcookie/afifiledownloader"
)

type DownloadService struct {
	df afd.DownloadFunc
	sc StreamerCreator
}

func NewDownloadService(df afd.DownloadFunc, sc StreamerCreator) *DownloadService {
	return &DownloadService{
		df: df,
		sc: sc,
	}
}

func (svc *DownloadService) Download(r io.Reader) (err error) {
	in := &Input{}

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

		s, err := svc.sc.CreateStreamer(in.Output)
		if err != nil {
			return err
		}

		defer s.Close()

		if _, err = io.Copy(s, res.Body); err != nil {
			return err
		}

		return nil
	}

	if err = svc.df(in.URL, time.Duration(in.Timeout), callback); err != nil {
		return err
	}

	return nil
}
