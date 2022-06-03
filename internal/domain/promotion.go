package domain

import (
	"context"
	"time"
)

type (
	Promotion struct {
		ID             string     `redis:"id"`
		Price          float64    `redis:"price"`
		ExpirationDate *time.Time `redis:"expiration_date"`
	}

	PromotionRepository interface {
		GetByID(ctx context.Context, id string) (*Promotion, error)
	}
)
