package handlers

import (
	"bufio"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"strings"
)

type (
	Streamer interface {
		Stream(w io.StringWriter, chunk []string) error
		Push() error
	}

	PromotionFeedHandler struct {
		s        Streamer
		filename string
		cmds     string
	}
)

const (
	columnID             = 0
	columnPrice          = 1
	columnExpirationDate = 2
)

func NewPromotionFeedHandler(filename, cmds string, s Streamer) *PromotionFeedHandler {
	return &PromotionFeedHandler{filename: filename, cmds: cmds, s: s}
}

func (p *PromotionFeedHandler) Feed() (int64, error) {
	log.Info().Msgf("reading promotions file %s", p.filename)
	pf, err := os.Open(p.filename)
	if err != nil {
		log.Error().Err(err).Msg("error opening promotions file")
		return -1, err
	}

	defer pf.Close()

	log.Info().Msgf("reading promotions file %s", p.cmds)
	df, err := os.Create(p.cmds)
	if err != nil {
		log.Error().Err(err).Msg("error opening redis commands file")
		return -1, err
	}

	scanner := bufio.NewScanner(pf)
	cmds := bufio.NewWriter(df)

	var c int64

	for scanner.Scan() {
		_ = p.s.Stream(cmds, strings.Split(scanner.Text(), ","))
		c++
	}

	err = scanner.Err()
	if err != nil {
		log.Error().Err(err).Msg("error with scanner")
		return -1, err
	}

	err = cmds.Flush()
	if err != nil {
		return -1, err
	}

	err = p.s.Push()
	if err != nil {
		return -1, err
	}

	return c, nil
}
