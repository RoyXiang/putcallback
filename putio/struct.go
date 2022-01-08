package putio

import (
	"sync"

	put "github.com/putdotio/go-putio"
)

type Put struct {
	Client       *put.Client
	MaxTransfers int
	mu           sync.Mutex
}

type SortedTransfers []put.Transfer
