package main

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/vitorsalgado/goprom/internal/api"
	goprom "github.com/vitorsalgado/goprom/internal/std"
	"github.com/vitorsalgado/goprom/internal/std/config"
	"github.com/vitorsalgado/goprom/internal/std/storage"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf := config.Load()

	goprom.ConfigureEnv(conf)

	client := storage.NewRedisClient(conf)

	srv := api.NewSrv()
	server := srv.APIServer(conf.ServerAddr, client)

	ext := make(chan os.Signal, 1)
	signal.Notify(ext, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-ext

		c, fn := context.WithTimeout(ctx, 5*time.Second)
		defer fn()

		err := server.Shutdown(c)
		if err != nil {
			log.Fatal().Err(err).Msg("server shutdown failed")
			return
		}

		cancel()
	}()

	log.Info().Msgf("starting http server on address %s", conf.ServerAddr)
	log.Info().Msgf("debug: %t", conf.Debug)
	log.Info().Msgf("redis: %s", conf.RedisAddr)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Err(err).Msg("server shutdown failed")
	}
}
