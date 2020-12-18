package tcp

import (
	"io"
	"net"

	afifiledownloader "github.com/morozovcookie/afi-file-downloader"
)

type Streamer struct {
	conn io.WriteCloser
}

func NewStreamer(address string) (afifiledownloader.FileStreamer, error) {
	var (
		s = &Streamer{}

		err error
	)

	if s.conn, err = net.Dial("tcp", address); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Streamer) Stream(r io.Reader) error {
	_, err := io.Copy(s.conn, r)

	return err
}

func (s *Streamer) Close() error {
	return s.conn.Close()
}
