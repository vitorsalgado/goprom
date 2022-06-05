package loader

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vitorsalgado/goprom/internal/std/config"
	"io"
	"os"
	"path"
	"testing"
)

func TestFeedingPromotions(t *testing.T) {
	t.Run("should read all promotions from source file and feed then to the provided destination", func(t *testing.T) {
		wd, _ := os.Getwd()
		cfg := config.Load()
		cfg.PromotionsBulkLoadWorkers = 5
		cfg.PromotionsCsv = path.Join(wd, "_testdata", "promos.csv")

		l := FakeLifecycle{}
		l.On("OnFinish", cfg.PromotionsCsv).Return(nil)

		s := NewStreamer(cfg)
		fake := FakeStreamer{real: s, m: mock.Mock{}}
		fake.m.On("Push", mock.Anything).Return(nil)

		loader := NewLoader(
			cfg, context.TODO(), &fake, NewSource(), &l)
		i, err := loader.Load()

		assert.Nil(t, err)
		assert.Equal(t, int64(5), i)
		assert.True(t, l.AssertNumberOfCalls(t, "OnFinish", 1))
	})
}

type FakeStreamer struct {
	real Streamer
	m    mock.Mock
}

func (s *FakeStreamer) Stream(w io.WriteCloser, chunk []string) error {
	return s.real.Stream(w, chunk)
}

type FakeLifecycle struct {
	mock.Mock
}

func (m *FakeLifecycle) OnFinish(filename string) error {
	return m.Called(filename).Error(0)
}
