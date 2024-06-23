package database

import (
	"context"
	"errors"
	"fmt"
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

// Allow - check if allow or block request by key
func (rl *LimiterRepositoryWithRedis) Allow(ctx context.Context, ratelimitKey string, limit int64, timeToBlock int64, windowTime int64) (bool, error) {
	// chave para bloquear
	keyToBlock := fmt.Sprintf("block:%s", ratelimitKey)
	// chave para armazenar o contador
	key := fmt.Sprintf("ratelimit[%d]:%s", time.Now().Unix(), ratelimitKey)

	// verifica se há um bloqueio com a cheve keyToBlock
	if rl.hasBlock(ctx, keyToBlock) {
		return false, nil
	}

	// usa pipeline para incrementar o contador e definir o tempo de verificaçao
	// ele seta contador e um tempo limite
	res, err := pipelineRedis(ctx, rl.client, key, windowTime)
	if handleErrByRedis(err) != nil {
		return false, errors.New("não foi possível incrementar o contador de bloqueio")
	}

	// pega o valor do contador atual
	count, err := res[0].(*redis.IntCmd).Result()
	if err != nil {
		return false, errors.New("não foi possível verificar o contador")
	}

	// verifica se o contador é maior ou igual ao limite
	if count >= int64(limit) {
		err := rl.doBlock(ctx, keyToBlock, timeToBlock)
		if handleErrByRedis(err) != nil {
			return false, errors.New("não foi possível bloquear a chave")
		}

		return false, nil
	}

	// retorna true se o contador for menor que o limite
	return true, nil
}

// handleErrByRedis - handle error by redis
func handleErrByRedis(err error) error {
	if err != nil && strings.TrimLeft(err.Error(), "redis: ") != "nil" {
		return err
	}

	return nil
}

// pipelineRedis - pipeline redis
func pipelineRedis(ctx context.Context, client *redis.Client, key string, windowTime int64) ([]redis.Cmder, error) {
	return client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HIncrBy(ctx, key, "count", 1)
		pipe.Expire(ctx, key, time.Duration(windowTime)*time.Second)
		return nil
	})
}

// hasBlock - check se há um bloqueio
func (rl *LimiterRepositoryWithRedis) hasBlock(ctx context.Context, key string) bool {
	value, err := rl.client.Get(ctx, key).Result()
	// caso não houver erros e a chave existir com o valor default 1
	if handleErrByRedis(err) == nil && value == "1" {
		return true
	}

	return false
}

// doBlock - bloqueia request
func (rl *LimiterRepositoryWithRedis) doBlock(ctx context.Context, key string, timeToBlock int64) error {
	// seta valor default 1 apenas para verificar se existe depois
	err := rl.client.SetNX(ctx, key, 1, time.Duration(timeToBlock)*time.Second).Err()
	if handleErrByRedis(err) != nil {
		return err
	}

	return nil
}
