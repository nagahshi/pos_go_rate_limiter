package database

import (
	"context"
	"log"
	"path/filepath"
	"testing"
	"time"

	"github.com/nagahshi/pos_go_rate_limiter/configs"
	"github.com/redis/go-redis/v9"
)

func newRedisClient(cfg *configs.Conf) *redis.Client {
	// Create a Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort, // Redis server address
		Password: cfg.RedisPassword,                   // Redis password
		DB:       cfg.RedisDatabaseIndex,              // Redis database index
	})

	ctx := context.Background()
	// Test the Redis connection
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
		return nil
	}

	return redisClient
}

func TestRedisRepository_bloqueio_manual(t *testing.T) {
	ctx := context.Background()
	path, _ := filepath.Abs("../../../")

	cfg, err := configs.LoadConfig(path + "/")
	if err != nil {
		t.Errorf("Error loading config: %v", err)
	}

	limiter := NewLimiterRepositoryWithRedis(newRedisClient(cfg))

	// then block
	err = limiter.DoBlockByKey(ctx, "keyToBlock", cfg.RateLimiterTimeout)
	if err != nil {
		t.Errorf("Error blocking key: %v", err)
	}

	t.Log(limiter.HasBlockByKey(ctx, "keyToBlock"))

	// check if blocked
	if !limiter.HasBlockByKey(ctx, "keyToBlock") {
		t.Error("Block not working")
	}

	// wait block expire
	time.Sleep(time.Duration(cfg.RateLimiterTimeout) * time.Second)

	// check if unblocked
	if limiter.HasBlockByKey(ctx, "keyToBlock") {
		t.Error("Após o timeout ainda continua bloqueado")
	}
}

func TestRedisRepository_bloqueio_via_contador(t *testing.T) {
	ctx := context.Background()
	path, _ := filepath.Abs("../../../")
	keyToBlock := "keyToBlock"
	requests := 10
	windowTime := 1
	cfg, err := configs.LoadConfig(path + "/")
	if err != nil {
		t.Errorf("Error loading config: %v", err)
	}

	limiter := NewLimiterRepositoryWithRedis(newRedisClient(cfg))

	for i := 0; i < requests+1; i++ {
		counter, err := limiter.CounterByKey(ctx, keyToBlock, int64(windowTime))
		if err != nil {
			t.Errorf("Error incrementing counter: %v", err)
		}

		if counter > int64(requests) {
			t.Log("then block")
			// then block
			err = limiter.DoBlockByKey(ctx, keyToBlock, cfg.RateLimiterTimeout)
			if err != nil {
				t.Errorf("Error blocking key: %v", err)
			}

			t.Log("Blocked")
			t.Log(limiter.HasBlockByKey(ctx, keyToBlock))
			break
		}
	}

	t.Log(limiter.HasBlockByKey(ctx, keyToBlock))
	// check if blocked
	if limiter.HasBlockByKey(ctx, keyToBlock) {
		t.Log("Block working")
	} else {
		t.Error("Block not working")
	}
}
