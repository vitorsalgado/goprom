package e2e

import (
	"github.com/go-redis/redis/v8"
	"github.com/vitorsalgado/goprom/internal/std/config"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	conf := config.Load()
	conf.RedisAddr = "localhost:6379"

	client := ConnectToRedis(conf, 20*time.Second)
	defer func(c *redis.Client) {
		_ = c.Close()
	}(client)

	c := m.Run()

	os.Exit(c)
}
