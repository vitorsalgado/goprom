package main

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vitorsalgado/goprom/internal/handlers"
	"github.com/vitorsalgado/goprom/internal/utils/config"
	"os"
	"path"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()
	wd, _ := os.Getwd()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.With().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	promoFile := path.Join(wd, cfg.Promotions)
	promoCommandsFile := path.Join(wd, cfg.PromotionsCommands)

	log.Info().Msgf("promotions file %s", promoFile)
	log.Info().Msgf("promotions commands file %s", promoCommandsFile)

	promo := handlers.NewPromotionFeedHandler(promoFile, promoCommandsFile, handlers.NewStreamer(cfg))
	n, err := promo.Feed()

	if err != nil {
		log.Fatal().Err(err).
			Msgf("an error occurred while feeding new promotions")
	}

	log.Info().
		Int64("promotions", n).
		Msgf("finished feeding %d promotions", n)
}
