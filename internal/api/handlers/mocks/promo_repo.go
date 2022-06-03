package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/vitorsalgado/goprom/internal/domain"
)

type FakePromotionRepository struct {
	mock.Mock
}

func (m *FakePromotionRepository) GetByID(ctx context.Context, id string) (*domain.Promotion, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Promotion), args.Error(1)
}
