package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/vitorsalgado/goprom/internal/domain"
)

type FakePromotionRepository struct {
	mock.Mock
}

func (m *FakePromotionRepository) GetByID(id string) (*domain.Promotion, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Promotion), args.Error(1)
}
