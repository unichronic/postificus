package storage

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() error {
	redisAddr := os.Getenv("REDIS_URL")
	if redisAddr == "" {
		host := os.Getenv("REDIS_HOST")
		port := os.Getenv("REDIS_PORT")
		if host != "" {
			if port == "" {
				port = "6379"
			}
			redisAddr = fmt.Sprintf("%s:%s", host, port)
		} else {
			redisAddr = "localhost:6379" // Default fallback
		}
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	ctx := context.Background()
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("✅ Connected to Redis")
	return nil
}

func CloseRedis() {
	if RedisClient != nil {
		RedisClient.Close()
		log.Println("🛑 Redis connection closed")
	}
}
