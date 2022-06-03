package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vitorsalgado/goprom/internal/utils/config"
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
		promotionsFilename := path.Join(wd, "_testdata", "promos.csv")
		promotionsCommandFilename := path.Join(wd, "_testdata", "promos_commands.tmp")

		l := FakeLifecycle{}
		l.On("OnFinish", promotionsFilename).Return(nil)

		s := NewStreamer(cfg)
		fake := FakeStreamer{real: s, m: mock.Mock{}}
		fake.m.On("Push").Return(nil)

		feeder := NewLoadPromotionsHandler(
			promotionsFilename, promotionsCommandFilename, &fake, &LoaderLocalFileSource{}, &l)
		i, err := feeder.Load()

		cmds, _ := ioutil.ReadFile(promotionsCommandFilename)
		lines := strings.Split(strings.TrimSpace(string(cmds)), "\n")

		assert.Nil(t, err)
		assert.Equal(t, int64(5), i)
		assert.Equal(t, i*2, int64(len(lines)))
		assert.True(t, fake.m.AssertNumberOfCalls(t, "Push", 1))
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

func (s *FakeStreamer) Push() error {
	return s.m.Called().Error(0)
}

type FakeLifecycle struct {
	mock.Mock
}

func (m *FakeLifecycle) OnFinish(filename string) error {
	return m.Called(filename).Error(0)
}
