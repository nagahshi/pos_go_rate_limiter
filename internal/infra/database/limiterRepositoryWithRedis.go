package database

import (
	"github.com/redis/go-redis/v9"
)

type LimiterRepositoryWithRedis struct {
	// RedisClient is a pointer to a RedisClient struct
	RedisClient *redis.Client
}
