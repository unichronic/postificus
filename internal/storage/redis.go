package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() error {
	redisAddr := os.Getenv("REDIS_URL")
	if redisAddr == "" {
		redisAddr = "localhost:6379" // Default fallback
	}

	var opts *redis.Options
	if strings.HasPrefix(redisAddr, "redis://") || strings.HasPrefix(redisAddr, "rediss://") {
		parsed, err := redis.ParseURL(redisAddr)
		if err != nil {
			return fmt.Errorf("invalid REDIS_URL: %w", err)
		}
		opts = parsed
	} else {
		opts = &redis.Options{Addr: redisAddr}
	}

	RedisClient = redis.NewClient(opts)

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
