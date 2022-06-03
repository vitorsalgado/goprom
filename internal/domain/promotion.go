package domain

import (
	"context"
)

type (
	Promotion struct {
		ID             string  `json:"id"`
		Price          float64 `json:"price"`
		ExpirationDate string  `json:"expiration_date"`
	}

	PromotionRepository interface {
		GetByID(ctx context.Context, id string) (*Promotion, error)
	}
)

const (
	PromotionDatetimeFormat = "2006-01-02 15:04:05"
)
