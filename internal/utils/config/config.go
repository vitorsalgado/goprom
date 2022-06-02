package config

import (
	goenv "github.com/Netflix/go-env"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
)

// Env holds environment variables used by the application
type Env struct {
	Promotions         string `env:"PROMOTIONS"`
	PromotionsCommands string `env:"PROMOTIONS_CMDS"`
	RedisAddr          string `env:"REDIS_ADDR,default=redis://localhost:6379"`
}

// Config represent application configurations
type Config struct {
	Env
	WorkDir string
}

// Load loads configuration using environment variables
func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		// ignoring dotenv error for missing .env file. the .env file is for local environment only.
		log.Trace().Err(err).Msg("no .env file. moving forward")
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal().Err(err).Msg("error getting the work directory")
		return nil
	}

	var env Env
	_, err = goenv.UnmarshalFromEnviron(&env)
	if err != nil {
		log.Fatal().Err(err).Msg("error loading environment variables into struct")
		return nil
	}

	return &Config{Env: env, WorkDir: wd}
}
