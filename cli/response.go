package cli

type Response struct {
	Success       bool     `json:"success"`
	HTTPCode      int      `json:"http-code,omitempty"`
	ContentLength int64    `json:"content-length,omitempty"`
	Redirects     []string `json:"redirects,omitempty"`
}

type GetResponse struct {
	*Response

	ContentType string `json:"content-type,omitempty"`
}

type HeadResponse struct {
	*Response

	Headers []string `json:"headers"`
}
