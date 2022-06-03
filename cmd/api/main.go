package main

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	goprom "github.com/vitorsalgado/goprom/internal"
	"github.com/vitorsalgado/goprom/internal/domain"
	"github.com/vitorsalgado/goprom/internal/handlers"
	"github.com/vitorsalgado/goprom/internal/utils/config"
	"github.com/vitorsalgado/goprom/internal/utils/middleware"
	"github.com/vitorsalgado/goprom/internal/utils/storage"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	conf := config.Load()

	client := storage.NewRedisClient(conf)
	srv := goprom.NewSrv(ctx, func(mux *http.ServeMux) {
		mux.Handle("/", goprom.Dispatcher(
			handlers.NewPingHandler(), handlers.NewPromotionHandler(domain.NewPromotionRepository(ctx, client))))
	})

	server := &http.Server{
		Addr:              conf.ServerAddr,
		Handler:           middleware.Recovery(srv.Mux),
		IdleTimeout:       30 * time.Second,
		WriteTimeout:      2 * time.Second,
		ReadTimeout:       2 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
	}

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

	log.Info().Msg("starting http server")
	log.Info().Msgf("debug is %s", os.Getenv("DEBUG"))
	log.Info().Msgf("debug (cfg) is %t", conf.Debug)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Err(err).Msgf("server shutdown failed")
	}

	<-ctx.Done()
}
