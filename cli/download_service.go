package cli

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	afd "github.com/morozovcookie/afifiledownloader"
)

type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

var ErrInvalidDuration = errors.New("invalid duration")

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	if val, ok := v.(float64); ok {
		*d = Duration(time.Duration(val))

		return nil
	}

	if val, ok := v.(string); ok {
		t, err := time.ParseDuration(val)
		if err != nil {
			return err
		}

		*d = Duration(t)

		return nil
	}

	return ErrInvalidDuration
}

type Input struct {
	IgnoreSSLCertificates bool     `json:"ignore-ssl-certificates"`
	FollowRedirects       bool     `json:"follow-redirects"`
	URL                   string   `json:"url"`
	Output                string   `json:"output"`
	Timeout               Duration `json:"timeout"`
}

// Validate input HTTP URL
// Validate output TCP URL

func (i Input) Validate() (err error) {
	return nil
}

type StreamerCreator interface {
	CreateStreamer(address string) (s afd.Streamer, err error)
}

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
