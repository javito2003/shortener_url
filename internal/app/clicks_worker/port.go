package clicks_worker

import (
	"context"
)

type ClickData struct {
	ShortCode string
	Increment int64
}

type ClickCacheReader interface {
	FetchAndClear(ctx context.Context) (map[string]int64, error)
}

type LinkBulkUpdater interface {
	IncrementClickCounts(ctx context.Context, counts map[string]int64) error
}
