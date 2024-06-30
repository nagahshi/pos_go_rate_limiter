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

func getConfig(t *testing.T) *configs.Conf {
	path, _ := filepath.Abs("../../../")
	cfg, err := configs.LoadConfig(path + "/")
	if err != nil {
		t.Errorf("Error loading config: %v", err)
	}

	return cfg
}

func TestRedisRepository_incremento_contador(t *testing.T) {
	ctx := context.Background()
	cfg := getConfig(t)
	key := "keyCounterTest1"

	limiter := NewLimiterRepositoryWithRedis(newRedisClient(cfg))

	// increment counter
	_, err := limiter.CounterByKey(ctx, key, cfg.RateLimiterWindowTime)
	if err != nil {
		t.Errorf("Error incrementing counter: %v", err)
	}

	// check if counter was incremented
	counter, err := limiter.CounterByKey(ctx, key, cfg.RateLimiterWindowTime)
	if err != nil {
		t.Errorf("Error incrementing counter: %v", err)
	}

	if counter != 2 {
		t.Error("Counter not working")
	}
}

func TestRedisRepository_bloqueio_manual(t *testing.T) {
	ctx := context.Background()
	cfg := getConfig(t)
	keyToBlock := "keyToBlockTest1"

	limiter := NewLimiterRepositoryWithRedis(newRedisClient(cfg))

	// then block
	err := limiter.DoBlockByKey(ctx, keyToBlock, cfg.RateLimiterTimeout)
	if err != nil {
		t.Errorf("Error blocking key: %v", err)
	}

	// check if blocked
	if !limiter.HasBlockByKey(ctx, keyToBlock) {
		t.Error("Block not working")
	}

	// wait block expire
	time.Sleep(time.Duration(cfg.RateLimiterTimeout+1) * time.Second)

	// check if unblocked
	if limiter.HasBlockByKey(ctx, keyToBlock) {
		t.Error("Após o timeout ainda continua bloqueado")
	}
}

func TestRedisRepository_bloqueio_via_contador(t *testing.T) {
	ctx := context.Background()
	keyToBlock := "keyToBlockTest2"
	cfg := getConfig(t)

	limiter := NewLimiterRepositoryWithRedis(newRedisClient(cfg))

	for i := 0; i < int(cfg.RateLimiterToken+1); i++ {
		counter, err := limiter.CounterByKey(ctx, "keyCounterTest2", cfg.RateLimiterWindowTime)
		if err != nil {
			t.Errorf("Error incrementing counter: %v", err)
		}

		if counter > cfg.RateLimiterToken {
			err = limiter.DoBlockByKey(ctx, keyToBlock, cfg.RateLimiterTimeout)
			if err != nil {
				t.Errorf("Error blocking key: %v", err)
			}

			t.Log(limiter.HasBlockByKey(ctx, keyToBlock))
			break
		}
	}

	// check if blocked
	if limiter.HasBlockByKey(ctx, keyToBlock) {
		t.Log("Block working")
	} else {
		t.Error("Block not working")
	}
}
