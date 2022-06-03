package api

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/vitorsalgado/goprom/internal/domain"
	"net/http"
	"strings"
)

// PromotionHandler handles specific promotions requests
type PromotionHandler struct {
	repo domain.PromotionRepository
}

// NewPromotionHandler returns a new PromotionHandler instance
func NewPromotionHandler(repo domain.PromotionRepository) *PromotionHandler {
	return &PromotionHandler{repo: repo}
}

// GetByID returns a promotion by its id
// GET /promotions/:id
// 200 (OK)
func (h *PromotionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msgf("GET %s", r.RequestURI)

	parts := strings.Split(r.URL.Path, "/")
	id := parts[len(parts)-1]

	promo, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		log.Error().Err(err).Msgf("error getting promotion with id %s", id)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if promo == nil {
		http.NotFound(w, r)
		return
	}

	b, err := json.Marshal(promo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")

	_, err = w.Write(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
