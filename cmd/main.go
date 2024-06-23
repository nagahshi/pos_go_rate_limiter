package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"

	"github.com/nagahshi/pos_go_rate_limiter/configs"
	repository "github.com/nagahshi/pos_go_rate_limiter/internal/infra/database"
	server "github.com/nagahshi/pos_go_rate_limiter/internal/infra/web"
	"github.com/nagahshi/pos_go_rate_limiter/internal/usecase"
)

func main() {
	// LoadConfig - carrega as configurações do arquivo .env
	cfg, err := configs.LoadConfig("configs")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	limiter := usecase.NewLimiter(
		repository.NewLimiterRepositoryWithRedis(newRedisClient(cfg)),
		cfg.RateLimiterIP,
		cfg.RateLimiterToken,
		cfg.RateLimiterTimeout,
		cfg.RateLimiterWindowTime,
	)

	server := server.NewServer(cfg.PORT)
	server.AddHandler("/", func(w http.ResponseWriter, r *http.Request) {
		isAllowed, err := limiter.AllowToken(r.Context(), "abc123")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !isAllowed {
			http.Error(w, "Blocked", http.StatusTooManyRequests)
			return
		}

		fmt.Fprint(w, "Hello, World!")
	})

	server.Start()
}

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
