package domain

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
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

	f := res.Val()
	if len(f) == 0 {
		return nil, nil
	}

	price, err := strconv.ParseFloat(f["price"], 64)
	if err != nil {
		return nil, err
	}

	return &Promotion{ID: f["id"], Price: price, ExpirationDate: f["expiration_date"]}, nil
}
