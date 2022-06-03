package loader

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vitorsalgado/goprom/internal/std/config"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

func TestFeedingPromotions(t *testing.T) {
	t.Run("should read all promotions from source file and feed then to the provided destination", func(t *testing.T) {
		wd, _ := os.Getwd()
		cfg := config.Load()
		cfg.PromotionsBulkLoadWorkers = 5
		cfg.PromotionsCsv = path.Join(wd, "_testdata", "promos.csv")
		cfg.PromotionsBulkCmdFilename = path.Join(wd, "_testdata", "promos_commands-%d.tmp")

		l := FakeLifecycle{}
		l.On("OnFinish", cfg.PromotionsCsv).Return(nil)

		s := NewStreamer(cfg)
		fake := FakeStreamer{real: s, m: mock.Mock{}}
		fake.m.On("Push", mock.Anything).Return(nil)

		loader := NewLoader(
			cfg, context.TODO(), &fake, NewSource(), &l)
		i, err := loader.Load()

		lines := make([]string, 0)
		count := 0

		for i := 0; i < cfg.PromotionsBulkLoadWorkers; i++ {
			cmds, _ := ioutil.ReadFile(fmt.Sprintf(cfg.PromotionsBulkCmdFilename, i))
			lines = append(lines, strings.Split(string(cmds), "\n")...)
		}

		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				count++
			}
		}

		assert.Nil(t, err)
		assert.Equal(t, int64(5), i)
		assert.Equal(t, 10, count) // 10 because there are 5 promotions and for each, we add a EXPIRE command
		assert.True(t, fake.m.AssertNumberOfCalls(t, "Push", cfg.PromotionsBulkLoadWorkers))
		assert.True(t, l.AssertNumberOfCalls(t, "OnFinish", 1))

		for i := 0; i < cfg.PromotionsBulkLoadWorkers; i++ {
			_ = os.Remove(fmt.Sprintf(cfg.PromotionsBulkCmdFilename, i))
		}
	})
}

type FakeStreamer struct {
	real Streamer
	m    mock.Mock
}

func (s *FakeStreamer) Stream(w io.StringWriter, chunk []string) error {
	return s.real.Stream(w, chunk)
}

func (s *FakeStreamer) Push(filename string) error {
	return s.m.Called(filename).Error(0)
}

type FakeLifecycle struct {
	mock.Mock
}

func (m *FakeLifecycle) OnFinish(filename string) error {
	return m.Called(filename).Error(0)
}
