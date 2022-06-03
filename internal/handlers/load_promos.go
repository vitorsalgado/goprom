package handlers

import (
	"bufio"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"strings"
	"time"
)

type (
	// Streamer handles promotions file data chunks
	Streamer interface {
		Stream(w io.StringWriter, chunk []string) error
		Push() error
	}

	LoaderSource interface {
		File(filename string) (*os.File, error)
	}

	LoaderLifecycle interface {
		OnFinish(filename string) error
	}

	// LoadPromotionsHandler loads promotions from a source file into a data storage
	LoadPromotionsHandler struct {
		s                     Streamer
		lc                    LoaderLifecycle
		src                   LoaderSource
		promotionsCsv         string
		promotionsCmdFilename string
	}

	LoaderLocalFileSource struct {
	}

	LoaderDefaultLifecycle struct {
	}
)

const (
	columnID             = 0
	columnPrice          = 1
	columnExpirationDate = 2
)

// NewLoadPromotionsHandler initiates a new instance of LoadPromotionsHandler
func NewLoadPromotionsHandler(
	filename, cmds string, s Streamer, src LoaderSource, lc LoaderLifecycle,
) *LoadPromotionsHandler {
	return &LoadPromotionsHandler{
		promotionsCsv: filename, promotionsCmdFilename: cmds, s: s, src: src, lc: lc}
}

// Load loads promotions from a source into a data storage
func (p *LoadPromotionsHandler) Load() (int64, error) {
	log.Debug().Msgf("reading promotions file %s", p.promotionsCsv)
	pf, err := p.src.File(p.promotionsCsv)
	if err != nil {
		log.Error().Err(err).Msg("error opening promotions file")
		return -1, err
	}

	defer pf.Close()

	log.Debug().Msgf("creating commands file %s", p.promotionsCmdFilename)
	df, err := os.Create(p.promotionsCmdFilename)
	if err != nil {
		log.Error().Err(err).Msg("error creating redis commands file")
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
		log.Error().Err(err).Msg("error with promotions file scanner")
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

	err = p.lc.OnFinish(p.promotionsCsv)

	return c, err
}

func (src *LoaderLocalFileSource) File(filename string) (*os.File, error) {
	return os.Open(filename)
}

func (lc *LoaderDefaultLifecycle) OnFinish(filename string) error {
	parts := strings.Split(filename, ".csv")
	nm := parts[0]
	now := time.Now().UTC().Format("20060102150405")

	err := os.Rename(filename, fmt.Sprintf("%s--%s--imported.csv", nm, now))
	if err != nil {
		return err
	}

	return nil
}
