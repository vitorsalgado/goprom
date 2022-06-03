package domain

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type (
	PromotionRedisRepository struct {
		ctx context.Context
		r   *redis.Client
	}
)

func NewPromotionRepository(ctx context.Context, client *redis.Client) PromotionRepository {
	return &PromotionRedisRepository{ctx: ctx, r: client}
}

func (repo *PromotionRedisRepository) GetByID(id string) (*Promotion, error) {
	res := repo.r.HGetAll(repo.ctx, id)
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
