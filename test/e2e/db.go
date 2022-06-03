package e2e

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"github.com/vitorsalgado/goprom/internal/std/config"
	"github.com/vitorsalgado/goprom/internal/std/storage"
	"time"
)

func ConnectToRedis(conf *config.Config, d time.Duration) *redis.Client {
	client := storage.NewRedisClient(conf)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	ticker := time.NewTicker(1 * time.Second)
	timeout := time.After(d)

	defer func() {
		if r := recover(); r != nil {
			log.Error().Msgf("something bad happened ... %v", r)
		}
	}()

	defer cancel()
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("trying to connect ...")

			if client == nil {
				client = storage.NewRedisClient(conf)
			} else {
				if err := client.Ping(ctx); err == nil {
					fmt.Println("successfully connected to Redis")
					return client
				}
			}

		case <-timeout:
			log.Error().Msg("unable to establish connection with Redis instance. exiting ...")
			return nil
		}

		time.Sleep(1 * time.Second)
	}
}
