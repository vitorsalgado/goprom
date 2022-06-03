package storage

import (
	"github.com/go-redis/redis/v8"
	"github.com/vitorsalgado/goprom/internal/utils/config"
)

// NewRedisClient creates a new Redis client
func NewRedisClient(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: "",
		DB:       0,
	})
}
