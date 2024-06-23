package usecase

import (
	"context"

	"github.com/nagahshi/pos_go_rate_limiter/internal/entity"
)

type Limiter struct {
	repository       entity.LimiterRepository
	maxIPRequests    int64
	maxTokenRequests int64
	timeToBlock      int64
	windowTime       int64
}

func NewLimiter(repository entity.LimiterRepository, maxIPRequests, maxTokenRequests int64, timeToBlock int64, windowTime int64) *Limiter {
	return &Limiter{
		repository:       repository,
		maxIPRequests:    maxIPRequests,
		maxTokenRequests: maxTokenRequests,
		timeToBlock:      timeToBlock,
		windowTime:       windowTime,
	}
}

func (l *Limiter) AllowIP(ctx context.Context, IP string) (bool, error) {
	return l.repository.Allow(ctx, IP, l.maxIPRequests, l.timeToBlock, l.windowTime)
}

func (l *Limiter) AllowToken(ctx context.Context, token string) (bool, error) {
	return l.repository.Allow(ctx, token, l.maxTokenRequests, l.timeToBlock, l.windowTime)
}
