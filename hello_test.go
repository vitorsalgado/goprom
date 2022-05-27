package goprom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldReturnHello(t *testing.T) {
	assert.Equal(t, "hello", Hello())
}
