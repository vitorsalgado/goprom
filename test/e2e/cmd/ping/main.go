package main

import (
	"github.com/vitorsalgado/goprom/internal/utils/config"
	"github.com/vitorsalgado/goprom/test/e2e"
	"time"
)

func main() {
	e2e.ConnectToRedis(config.Load(), 30*time.Second)
}
