package entity

import (
	"context"
)

type LimiterRepository interface {
	HasBlockByKey(ctx context.Context, key string) bool
	DoBlockByKey(ctx context.Context, key string, timeToBlock int64) error
	CounterByKey(ctx context.Context, key string, windowTime int64) (int64, error)
}
