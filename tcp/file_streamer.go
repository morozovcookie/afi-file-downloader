package tcp

import (
	"io"
	"net"
)

type FileStreamer struct {
	conn io.WriteCloser
}

func NewFileStreamer(address string) (fs *FileStreamer, err error) {
	fs = &FileStreamer{}

	if fs.conn, err = net.Dial("tcp", address); err != nil {
		return nil, err
	}

	return fs, nil
}

func (fs *FileStreamer) ReadFrom(r io.Reader) (n int64, err error) {
	return io.Copy(fs.conn, r)
}

func (fs *FileStreamer) Close() (err error) {
	return fs.conn.Close()
}
