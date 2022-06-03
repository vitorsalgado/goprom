package e2e

import (
	"github.com/vitorsalgado/goprom/internal/utils/config"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	conf := config.Load()
	conf.RedisAddr = "localhost:6379"

	redis := ConnectToRedis(conf, 20*time.Second)
	defer redis.Close()

	c := m.Run()

	os.Exit(c)
}
