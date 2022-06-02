package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Run("it should load with default values when there's no .env file", func(t *testing.T) {
		config := Load()
		assert.Equal(t, "redis://localhost:6379", config.RedisAddr)
	})
}
