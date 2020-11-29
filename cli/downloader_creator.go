package cli

import (
	afd "github.com/morozovcookie/afifiledownloader"
)

type DownloaderCreator func(isFollowRedirects bool, maxRedirects int64) afd.DownloadFunc
