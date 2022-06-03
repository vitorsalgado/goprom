package handlers

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vitorsalgado/goprom/internal/utils/config"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

func TestFeedingPromotions(t *testing.T) {
	t.Run("should read all promotions from source file and feed then to the provided destination", func(t *testing.T) {
		cfg := config.Load()
		cfg.PromotionsBulkLoadWorkers = 2
		cfg.PromotionsCsv = "./_testdata/promos.csv"
		cfg.PromotionsBulkCmdFilename = "./_testdata/promos_commands-%d.tmp"

		l := FakeLifecycle{}
		l.On("OnFinish", cfg.PromotionsCsv).Return(nil)

		s := NewStreamer(cfg)
		fake := FakeStreamer{real: s, m: mock.Mock{}}
		fake.m.On("Push", mock.Anything).Return(nil)

		loader := NewLoadPromotionsHandler(
			cfg, context.TODO(), &fake, &LoaderLocalFileSource{}, &l)
		i, err := loader.Load()

		cmds0, _ := ioutil.ReadFile(fmt.Sprintf(cfg.PromotionsBulkCmdFilename, 0))
		cmds1, _ := ioutil.ReadFile(fmt.Sprintf(cfg.PromotionsBulkCmdFilename, 1))
		lines0 := strings.Split(strings.TrimSpace(string(cmds0)), "\n")
		lines1 := strings.Split(strings.TrimSpace(string(cmds1)), "\n")

		assert.Nil(t, err)
		assert.Equal(t, int64(5), i)
		assert.Equal(t, i*2, int64(len(lines0)+len(lines1)))
		assert.True(t, fake.m.AssertNumberOfCalls(t, "Push", 2))
		assert.True(t, l.AssertNumberOfCalls(t, "OnFinish", 1))
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
