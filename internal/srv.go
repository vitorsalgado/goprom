package goprom

import (
	"github.com/go-redis/redis/v8"
	"github.com/vitorsalgado/goprom/internal/domain"
	"github.com/vitorsalgado/goprom/internal/handlers"
	"github.com/vitorsalgado/goprom/internal/utils/middleware"
	"net/http"
	"time"
)

type Srv struct {
	mux *http.ServeMux
}

func NewSrv() *Srv {
	return &Srv{mux: http.NewServeMux()}
}

func (s *Srv) Configure(conf ...func(mux *http.ServeMux)) {
	for _, c := range conf {
		c(s.mux)
	}
}

// APIServer builds an HTTP server with default dependencies and routes
func (s *Srv) APIServer(addr string, client *redis.Client) *http.Server {
	s.mux.Handle("/", Dispatcher(
		handlers.NewPingHandler(),
		handlers.NewPromotionHandler(domain.NewPromotionRepository(client))))

	return &http.Server{
		Addr:              addr,
		Handler:           middleware.Recovery(s.mux),
		IdleTimeout:       30 * time.Second,
		WriteTimeout:      2 * time.Second,
		ReadTimeout:       2 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
	}
}
