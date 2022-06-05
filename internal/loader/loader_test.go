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

		s := NewStreamer()

		fake := FakeStreamer{real: s, m: mock.Mock{}}
		fake.m.On("Push", mock.Anything).Return(nil)
		writer := &FakeWriter{}
		writer.On("Write", mock.Anything).Return(0, nil)
		writer.On("Close").Return(nil)

		loader := NewLoader(
			cfg, context.TODO(), func() (io.WriteCloser, error) { return writer, nil }, NewSource(), &l)
		i, err := loader.Load()

		assert.Nil(t, err)
		assert.Equal(t, int64(5), i)
		assert.True(t, writer.AssertNumberOfCalls(t, "Write", 5))
		assert.True(t, l.AssertNumberOfCalls(t, "OnFinish", 1))
	})
}

type FakeWriter struct {
	mock.Mock
}

func (fk *FakeWriter) Write(p []byte) (n int, err error) {
	args := fk.Called(p)
	return args.Get(0).(int), args.Error(1)
}

func (fk *FakeWriter) Close() error {
	return fk.Called().Error(0)
}

type FakeStreamer struct {
	real Streamer
	m    mock.Mock
}

func (s *FakeStreamer) Stream(w io.Writer, chunk []string) error {
	return s.real.Stream(w, chunk)
}

type FakeLifecycle struct {
	mock.Mock
}

func (m *FakeLifecycle) OnFinish(filename string) error {
	return m.Called(filename).Error(0)
}
