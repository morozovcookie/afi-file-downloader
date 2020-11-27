package cli

import (
	"github.com/stretchr/testify/mock"

	afd "github.com/morozovcookie/afifiledownloader"
)

type StreamerCreator interface {
	CreateStreamer(address string) (s afd.Streamer, err error)
}

type MockStreamerCreator struct {
	mock.Mock
}

func (sc *MockStreamerCreator) CreateStreamer(address string) (s afd.Streamer, err error) {
	args := sc.Called(address)

	return args.Get(0).(afd.Streamer), args.Error(1)
}
