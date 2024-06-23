package entity

import (
	"context"
)

type LimiterRepository interface {
	Allow(ctx context.Context, ratelimitKey string, limit int64, timeToBlock int64, windowTime int64) (bool, error)
}
