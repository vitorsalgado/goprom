package main

import (
	"context"
	goprom "github.com/vitorsalgado/goprom/internal"
	"github.com/vitorsalgado/goprom/internal/domain"
	"github.com/vitorsalgado/goprom/internal/handlers"
	"net/http"
)

func main() {
	ctx, _ := context.WithCancel(context.Background())
	srv := goprom.NewSrv(ctx, func(mux *http.ServeMux) {
		mux.Handle("/", goprom.Dispatcher(
			handlers.NewPingHandler(), handlers.NewPromotionHandler(domain.NewPromotionRepository())))
	})
	server := http.Server{Addr: ":8080", Handler: srv.Mux}

	_ = server.ListenAndServe()
}
