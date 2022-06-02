package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	t.Run("should return pong", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
		rr := httptest.NewRecorder()
		h := http.HandlerFunc(NewPingHandler().Ping)

		h.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "pong", rr.Body.String())
	})
}
