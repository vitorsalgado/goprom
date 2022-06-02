package db

import (
	"github.com/go-redis/redis/v8"
	"github.com/vitorsalgado/goprom/internal/utils/config"
)

func NewRedis(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: "",
		DB:       0,
	})
}
