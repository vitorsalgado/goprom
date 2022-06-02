package handlers

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
)

type PingHandler struct {
}

func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

func (h *PingHandler) Ping(w http.ResponseWriter, _ *http.Request) {
	log.Debug().Msg("pong")
	_, _ = fmt.Fprint(w, "pong")
}
