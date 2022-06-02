package domain

import "time"

type (
	Promotion struct {
		ID             string
		Price          float64
		ExpirationDate time.Time
	}

	PromotionRepository interface {
		GetByID(id string) (*Promotion, error)
	}
)
