package domain

import (
	"github.com/vitorsalgado/goprom/internal/utils/config"
)

type (
	PromotionRedisRepository struct {
		cfg *config.Config
	}
)

func NewPromotionRepository() PromotionRepository {
	return &PromotionRedisRepository{}
}

func (r *PromotionRedisRepository) GetByID(id string) (*Promotion, error) {
	return nil, nil
}
