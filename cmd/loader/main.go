package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	goprom "github.com/vitorsalgado/goprom/internal"
	"github.com/vitorsalgado/goprom/internal/handlers"
	"github.com/vitorsalgado/goprom/internal/utils/config"
	"github.com/vitorsalgado/goprom/internal/utils/storage"
	"golang.org/x/sys/unix"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	f, _ := os.Create("promo.lock")
	if err := unix.Flock(int(f.Fd()), unix.LOCK_EX); err != nil {
		log.Fatal().Err(err).Msgf("cannot acquire exclusive lock. maybe there is another job running")
		return
	}

	defer func() {
		if err := unix.Flock(int(f.Fd()), unix.LOCK_UN); err != nil {
			log.Fatal().Err(err).Msgf("error releasing lock")
		}
	}()

	start := time.Now()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

		<-exit

		cancel()
	}()

	_ = godotenv.Load()
	cfg := config.Load()

	goprom.ConfigureEnv(cfg)

	// testing if redis is reachable
	redisCtx, redisCancelFn := context.WithTimeout(ctx, 10*time.Second)
	defer redisCancelFn()

	client := storage.NewRedisClient(cfg)
	st := client.Ping(redisCtx)
	if st.Err() != nil {
		log.Fatal().Err(st.Err()).Msgf("unable to connect to redis")
		return
	}

	if _, err := os.Stat(cfg.PromotionsCsv); err != nil {
		log.Info().Msgf("promotions file %s does not exists", cfg.PromotionsCsv)
		os.Exit(0)
		return
	}

	promo := handlers.NewLoadPromotionsHandler(
		cfg, ctx, handlers.NewStreamer(cfg), &handlers.LoaderLocalFileSource{}, &handlers.LoaderDefaultLifecycle{})
	n, err := promo.Load()

	if err != nil {
		log.Fatal().Err(err).
			Msgf("an error occurred while feeding new promotions")
	}

	elapsed := time.Since(start)

	log.Info().Msg(
		fmt.Sprintf("finished feeding %d promotions. took %f seconds", n, elapsed.Seconds()))
}
