package config

import (
	goenv "github.com/Netflix/go-env"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
)

type environ struct {
	Debug                     bool   `env:"DEBUG,default=true"`
	LogLevel                  int8   `env:"LOG_LEVEL,default=1"`
	ServerAddr                string `env:"SERVER_ADDR,default=:8080"`
	PromotionsBulkLoadWorkers int    `env:"PROMOTIONS_BULK_LOAD_WORKERS,default=5"`
	PromotionsCsv             string `env:"PROMOTIONS,default=/data/promotions.csv"`
	PromotionsExpiration      int    `env:"PROMOTIONS_EXPIRATION,default=1800"`
	RedisAddr                 string `env:"REDIS_ADDR,default=redis:6379"`
}

// Config represent application configurations
type Config struct {
	environ
	WorkDir string
}

// Load loads configuration using environment variables
func Load() *Config {
	// ignoring dotenv error for missing .env file. the .env file is for local environment only.
	_ = godotenv.Load()

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal().Err(err).Msg("error getting the work directory")
		return nil
	}

	var env environ
	_, err = goenv.UnmarshalFromEnviron(&env)
	if err != nil {
		log.Fatal().Err(err).Msg("error loading environment variables into struct")
		return nil
	}

	return &Config{environ: env, WorkDir: wd}
}
