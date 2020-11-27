package cli

import (
	afd "github.com/morozovcookie/afifiledownloader"
)

type StreamerCreator func(address string) (s afd.Streamer, err error)
