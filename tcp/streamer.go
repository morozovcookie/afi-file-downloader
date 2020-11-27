package tcp

import (
	"io"
	"net"
)

type Streamer struct {
	conn io.WriteCloser
}

func NewStreamer(address string) (s *Streamer, err error) {
	s = &Streamer{}

	if s.conn, err = net.Dial("tcp", address); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Streamer) Write(p []byte) (n int, err error) {
	return s.conn.Write(p)
}

func (s *Streamer) Close() (err error) {
	return s.conn.Close()
}
