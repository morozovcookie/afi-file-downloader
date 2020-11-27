package cli

import (
	"encoding/json"
	"errors"
	"regexp"
	"time"
)

var ErrInvalidDuration = errors.New("invalid duration")

type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

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

var (
	ErrInvalidURL    = errors.New("invalid url address")
	ErrInvalidOutput = errors.New("invalid output address")
)

func (i Input) Validate() (err error) {
	if err = validateURL(i.URL); err != nil {
		return err
	}

	if err = validateOutput(i.Output); err != nil {
		return err
	}

	return nil
}

const (
	URLRegex = `(?m)^((([^:/?#]+):)?(//([^/?#]*))?([^?#]*)(\?([^#]*))?(#(.*))?)$`

	HostPortRegex = `(?m)^((((25[0-5])|(2[0-4]\d{1})|([0-1]?\d{1,2}))\.){3}((25[0-5])|(2[0-4]\d{1})|` +
		`([0-1]?\d{1,2})){1}(:((6553[0-5])|(655[0-2]\d{1})|(65[0-4]\d{2})|(6[0-4]\d{3})|([1-5]\d{4})|` +
		`([1-9]\d{3})|([1-9]\d{2})|([1-9]\d{1})|([1-9])))?)$`
	HostnamePortRegex = `(?m)^(((([\d\w]|[\d\w][\d\w\-]*[\d\w])\.)*([\d\w]|[\d\w][\d\w\-]*[\d\w]))` +
		`(:((6553[0-5])|(655[0-2]\d{1})|(65[0-4]\d{2})|(6[0-4]\d{3})|([1-5]\d{4})|([1-9]\d{3})|` +
		`([1-9]\d{2})|([1-9]\d{1})|([1-9])))?)$`
)

func validateURL(s string) (err error) {
	if s == "" {
		return ErrInvalidURL
	}

	if ok := regexp.MustCompile(URLRegex).MatchString(s); !ok {
		return ErrInvalidURL
	}

	return nil
}

func validateOutput(s string) (err error) {
	if s == "" {
		return nil
	}

	if ok := regexp.MustCompile(HostPortRegex).MatchString(s); ok {
		return nil
	}

	if ok := regexp.MustCompile(HostnamePortRegex).MatchString(s); ok {
		return nil
	}

	return ErrInvalidOutput
}
