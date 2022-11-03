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
	ID          int64
	Name        string
	IsDir       bool
	Size        int64
	FullPath    string
	ContentType string
}
