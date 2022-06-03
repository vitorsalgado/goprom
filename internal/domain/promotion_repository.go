package domain

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type (
	PromotionRedisRepository struct {
		r *redis.Client
	}
)

func NewPromotionRepository(client *redis.Client) PromotionRepository {
	return &PromotionRedisRepository{r: client}
}

func (repo *PromotionRedisRepository) GetByID(ctx context.Context, id string) (*Promotion, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	res := repo.r.HGetAll(ctx, id)
	if res.Err() != nil {
		return nil, res.Err()
	}

	promo := &Promotion{}
	err := res.Scan(promo)
	if err != nil {
		return nil, err
	}

	return promo, nil
}
