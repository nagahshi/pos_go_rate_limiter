package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

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

// allow - verifica se a key(chave) pode fazer a requisição
func (l *Limiter) allow(ctx context.Context, key string, limit int64) (bool, error) {
	// chave para bloquear
	keyToBlock := fmt.Sprintf("block:%s", key)
	// chave para armazenar o contador
	keyToCounter := fmt.Sprintf("ratelimit[%d]:%s", time.Now().Unix(), key)

	// verifica se há um bloqueio com a cheve keyToBlock
	if l.repository.HasBlockByKey(ctx, keyToBlock) {
		return false, nil
	}

	// incrementa o contador
	count, err := l.repository.CounterByKey(ctx, keyToCounter, l.windowTime)
	if err != nil {
		return false, errors.New("não foi possível verificar o contador")
	}

	// verifica se o contador é maior ou igual ao limite
	if count >= int64(limit) {
		err := l.repository.DoBlockByKey(ctx, keyToBlock, l.timeToBlock)
		if err != nil {
			return false, errors.New("não foi possível bloquear a chave")
		}

		return false, nil
	}

	// retorna true se o contador for menor que o limite
	return true, nil
}

// AllowIP - verifica se o IP pode fazer a requisição
func (l *Limiter) AllowIP(ctx context.Context, IP string) (bool, error) {
	return l.allow(ctx, IP, l.maxIPRequests)
}

// AllowToken - verifica se o token pode fazer a requisição
func (l *Limiter) AllowToken(ctx context.Context, token string) (bool, error) {
	return l.allow(ctx, token, l.maxTokenRequests)
}

// GetMaxIPRequests - retorna o limite de requisições por IP
func (l *Limiter) GetMaxIPRequests() int64 {
	return l.maxIPRequests
}

// GetMaxTokenRequests - retorna o limite de requisições por token
func (l *Limiter) GetMaxTokenRequests() int64 {
	return l.maxTokenRequests
}
