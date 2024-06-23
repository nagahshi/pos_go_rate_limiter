package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/nagahshi/pos_go_rate_limiter/configs"
	"github.com/nagahshi/pos_go_rate_limiter/internal/entity"
	repository "github.com/nagahshi/pos_go_rate_limiter/internal/infra/database"
	server "github.com/nagahshi/pos_go_rate_limiter/internal/infra/web"
	"github.com/nagahshi/pos_go_rate_limiter/internal/usecase"
)

var repo entity.LimiterRepository

func main() {
	// LoadConfig - carrega as configurações do arquivo .env
	cfg, err := configs.LoadConfig("configs")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	repo = repository.NewLimiterRepositoryWithRedis(newRedisClient(cfg))
	limiter := usecase.NewLimiter(repo, 10, 10, 10*time.Second)

	server := server.NewServer(cfg.PORT)
	server.AddHandler("/", func(w http.ResponseWriter, r *http.Request) {
		blocked, err := limiter.DoRequest(context.Background(), "abc123", 10)
		if err != nil {
			log.Printf("Error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if blocked {
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
