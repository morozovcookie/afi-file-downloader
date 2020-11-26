package tcp

import (
	afd "github.com/morozovcookie/afifiledownloader"
)

type StreamerCreator struct{}

func NewStreamerCreator() *StreamerCreator {
	return &StreamerCreator{}
}

func (sc StreamerCreator) CreateStreamer(address string) (s afd.Streamer, err error) {
	return newStreamer(address)
}
