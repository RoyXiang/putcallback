package putio

import (
	"sync"

	put "github.com/putdotio/go-putio"
)

type Put struct {
	Client                *put.Client
	MaxTransfers          int
	DefaultDownloadFolder string
	mu                    sync.Mutex
}

type SortedTransfers []put.Transfer

type FileInfo struct {
	Name     string
	IsDir    bool
	FullPath string
}
