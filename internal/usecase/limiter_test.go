package usecase

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockLimiterRepository struct {
	mock.Mock
}

func (m *MockLimiterRepository) HasBlockByKey(ctx context.Context, key string) bool {
	args := m.Called(ctx, key)
	return args.Bool(0)
}

func (m *MockLimiterRepository) CounterByKey(ctx context.Context, key string, windowTime int64) (int64, error) {
	args := m.Called(ctx, key, windowTime)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockLimiterRepository) DoBlockByKey(ctx context.Context, key string, timeToBlock int64) error {
	args := m.Called(ctx, key, timeToBlock)
	return args.Error(0)
}

func GetLimiterRepository() (*Limiter, *MockLimiterRepository) {
	mockRepository := &MockLimiterRepository{}
	maxIPRequests := int64(10)
	maxTokenRequests := int64(10)
	timeToBlock := int64(10)
	windowTime := int64(60)

	return NewLimiter(mockRepository, maxIPRequests, maxTokenRequests, timeToBlock, windowTime), mockRepository
}

func TestLimiter_AllowToken(t *testing.T) {
	ctx := context.Background()
	limiter, mockRepository := GetLimiterRepository()
	token := "validToken"

	mockRepository.On("HasBlockByKey", ctx, "block:"+token).Return(false)
	mockRepository.On("CounterByKey", ctx, "ratelimit["+strconv.FormatInt(time.Now().Unix(), 10)+"]:"+token, limiter.windowTime).Return(limiter.maxTokenRequests-1, nil)

	allowed, err := limiter.AllowToken(ctx, token)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !allowed {
		t.Error("Token should be Blocked")
	}
}

func TestLimiter_BlockToken(t *testing.T) {
	ctx := context.Background()
	limiter, mockRepository := GetLimiterRepository()
	token := "validToken"

	mockRepository.On("HasBlockByKey", ctx, "block:"+token).Return(true)
	allowed, err := limiter.AllowToken(ctx, token)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if allowed {
		t.Error("Token not should be blocked")
	}
}

func TestLimiter_TokenExceedsLimit(t *testing.T) {
	ctx := context.Background()
	limiter, mockRepository := GetLimiterRepository()
	token := "validToken"

	mockRepository.On("HasBlockByKey", ctx, "block:"+token).Return(true)
	mockRepository.On("CounterByKey", ctx, "ratelimit["+strconv.FormatInt(time.Now().Unix(), 10)+"]:"+token, limiter.windowTime).Return(limiter.maxTokenRequests+1, nil)

	allowed, err := limiter.AllowToken(ctx, token)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if allowed {
		t.Error("Token should be blocked")
	}
}

func TestLimiter_AllowIP(t *testing.T) {
	ctx := context.Background()
	limiter, mockRepository := GetLimiterRepository()
	IP := "validIP"

	mockRepository.On("HasBlockByKey", ctx, "block:"+IP).Return(false)
	mockRepository.On("CounterByKey", ctx, "ratelimit["+strconv.FormatInt(time.Now().Unix(), 10)+"]:"+IP, limiter.windowTime).Return(limiter.maxIPRequests-1, nil)

	allowed, err := limiter.AllowIP(ctx, IP)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !allowed {
		t.Error("IP should be Blocked")
	}
}

func TestLimiter_BlockIP(t *testing.T) {
	ctx := context.Background()
	limiter, mockRepository := GetLimiterRepository()
	IP := "validIP"

	mockRepository.On("HasBlockByKey", ctx, "block:"+IP).Return(true)
	allowed, err := limiter.AllowIP(ctx, IP)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if allowed {
		t.Error("IP not should be blocked")
	}
}

func TestLimiter_IPExceedsLimit(t *testing.T) {
	ctx := context.Background()
	limiter, mockRepository := GetLimiterRepository()
	IP := "validIP"

	mockRepository.On("HasBlockByKey", ctx, "block:"+IP).Return(true)
	mockRepository.On("CounterByKey", ctx, "ratelimit["+strconv.FormatInt(time.Now().Unix(), 10)+"]:"+IP, limiter.windowTime).Return(limiter.maxIPRequests+1, nil)

	allowed, err := limiter.AllowIP(ctx, IP)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if allowed {
		t.Error("IP should be blocked")
	}
}
