package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/vitorsalgado/goprom/internal/domain"
	"github.com/vitorsalgado/goprom/internal/handlers/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPromotionHandler(t *testing.T) {
	t.Run("should return the promotion with provided when it exists in the database", func(t *testing.T) {
		id := "test"
		expected := &domain.Promotion{}
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/promotions/%s", id), nil)
		rr := httptest.NewRecorder()
		repo := mocks.FakePromotionRepository{}
		repo.On("GetByID", id).Return(expected, nil)

		h := http.HandlerFunc(NewPromotionHandler(&repo).GetByID)
		h.ServeHTTP(rr, req)

		promo := &domain.Promotion{}
		err := json.Unmarshal(rr.Body.Bytes(), &promo)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, expected, promo)
	})
}
