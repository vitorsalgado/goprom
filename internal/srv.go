package goprom

import (
	"context"
	"net/http"
)

type Srv struct {
	ctx context.Context
	Mux *http.ServeMux
}

func NewSrv(ctx context.Context, conf ...func(mux *http.ServeMux)) *Srv {
	mux := http.NewServeMux()

	for _, c := range conf {
		c(mux)
	}

	return &Srv{ctx: ctx, Mux: mux}
}
