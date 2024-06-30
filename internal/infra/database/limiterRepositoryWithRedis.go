package database

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type LimiterRepositoryWithRedis struct {
	client *redis.Client
}

func NewLimiterRepositoryWithRedis(redisClient *redis.Client) *LimiterRepositoryWithRedis {
	return &LimiterRepositoryWithRedis{
		client: redisClient,
	}
}

func (rl *LimiterRepositoryWithRedis) CounterByKey(ctx context.Context, key string, windowTime int64) (int64, error) {
	// pipeline de incremento com expiração
	res, err := rl.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HIncrBy(ctx, key, "count", 1)
		pipe.Expire(ctx, key, time.Duration(windowTime)*time.Second)
		return nil
	})

	if handleErrInRedis(err) != nil {
		return 0, errors.New("não foi possível incrementar o contador de bloqueio")
	}

	if len(res) == 0 {
		return 0, errors.New("não foi possível buscar contador de bloqueio")
	}

	// pega o valor do contador atual
	return res[0].(*redis.IntCmd).Result()
}

// doBlock - bloqueia request
func (rl *LimiterRepositoryWithRedis) DoBlockByKey(ctx context.Context, key string, timeToBlock int64) error {
	// seta valor default 1 apenas para verificar se existe depois
	err := rl.client.SetNX(ctx, key, "blocked", time.Duration(timeToBlock)*time.Second).Err()
	if handleErrInRedis(err) != nil {
		return err
	}

	return nil
}

// hasBlock - check se há um bloqueio
func (rl *LimiterRepositoryWithRedis) HasBlockByKey(ctx context.Context, key string) bool {
	value, err := rl.client.Get(ctx, key).Result()
	// caso não houver erros e a chave existir com o valor default 1
	if handleErrInRedis(err) == nil && value == "blocked" {
		return true
	}

	return false
}

// handleErrByRedis - handle error by redis
func handleErrInRedis(err error) error {
	if err != nil && strings.TrimLeft(err.Error(), "redis: ") != "nil" {
		return err
	}

	return nil
}
