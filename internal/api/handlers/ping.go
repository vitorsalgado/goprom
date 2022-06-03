package handlers

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
)

// PingHandler is just a health check handler
type PingHandler struct {
}

// NewPingHandler returns a new PingHandler instance
func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

// Ping pong
// GET /ping
// 200 (OK)
func (h *PingHandler) Ping(w http.ResponseWriter, _ *http.Request) {
	log.Debug().Msg("pong")
	_, _ = fmt.Fprint(w, "pong")
}
