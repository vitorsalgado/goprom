package handlers

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/vitorsalgado/goprom/internal/domain"
	"net/http"
	"strings"
)

type PromotionHandler struct {
	repo domain.PromotionRepository
}

func NewPromotionHandler(repo domain.PromotionRepository) *PromotionHandler {
	return &PromotionHandler{repo: repo}
}

func (h *PromotionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	id := parts[len(parts)-1]

	log.Debug().Msgf("getting promotion %s", id)

	promo, err := h.repo.GetByID(id)
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

	_, err = w.Write(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
