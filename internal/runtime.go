package goprom

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vitorsalgado/goprom/internal/utils/config"
	"os"
)

func ConfigureRuntime(cfg *config.Config) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.Level(cfg.LogLevel))

	if cfg.Debug {
		log.Logger = log.Logger.With().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}
