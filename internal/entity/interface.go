package entity

import (
	"context"
	"time"
)

type LimiterRepository interface {
	StoreRequest(ctx context.Context, identifier string, expiration time.Duration) (int64, error)
	StoreBlock(ctx context.Context, identifier string, expiration time.Duration) error
	HasBlock(ctx context.Context, identifier string) (bool, error)
}
