package putio

import (
	"context"
	"log"
	"sort"
	"strings"
)

func (s SortedTransfers) Len() int {
	return len(s)
}

func (s SortedTransfers) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortedTransfers) Less(i, j int) bool {
	isFirstSeeding := s[i].Status == StatusSeeding
	if s[i].Status == s[j].Status {
		if !isFirstSeeding {
			return s[i].ID < s[j].ID
		}
		if s[i].IsPrivate != s[j].IsPrivate {
			return s[i].IsPrivate
		}
		return s[i].SecondsSeeding < s[j].SecondsSeeding
	}
	firstWeight, secondWeight := 3, 3
	if strings.HasSuffix(s[i].Status, Ongoing) {
		firstWeight--
		if !isFirstSeeding {
			firstWeight--
		}
	}
	if strings.HasSuffix(s[j].Status, Ongoing) {
		secondWeight--
		if isFirstSeeding {
			secondWeight--
		}
	}
	return firstWeight < secondWeight
}

func (put *Put) CleanupTransfers() {
	put.mu.Lock()
	defer put.mu.Unlock()

	ctx := context.Background()
	transfers, err := put.Client.Transfers.List(ctx)
	if err != nil {
		return
	}
	sort.Sort(SortedTransfers(transfers))

	count := put.MaxTransfers - 1
	var idsToBeCanceled []int64
	var numToBeCanceled, numToBeCleaned int
	for _, transfer := range transfers {
		if strings.HasSuffix(transfer.Status, Ongoing) {
			if transfer.Status != StatusSeeding || count > 0 {
				count--
			} else {
				idsToBeCanceled = append(idsToBeCanceled, transfer.ID)
			}
		} else {
			numToBeCleaned++
		}
	}
	numToBeCanceled = len(idsToBeCanceled)

	if numToBeCanceled > 0 {
		err = put.Client.Transfers.Cancel(ctx, idsToBeCanceled...)
		if err == nil {
			numToBeCleaned += numToBeCanceled
		}
	}

	if numToBeCleaned > 0 {
		_ = put.Client.Transfers.Clean(ctx)
		log.Printf("Transfers cleaned, %d canceled.", numToBeCanceled)
	}
}
