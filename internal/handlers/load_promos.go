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

	// LoadPromotionsHandler loads promotions from a source file into a data storage
	LoadPromotionsHandler struct {
		s                     Streamer
		promotionsCsv         string
		promotionsCmdFilename string
	}
)

const (
	columnID             = 0
	columnPrice          = 1
	columnExpirationDate = 2
)

// NewLoadPromotionsHandler initiates a new instance of LoadPromotionsHandler
func NewLoadPromotionsHandler(filename, cmds string, s Streamer) *LoadPromotionsHandler {
	return &LoadPromotionsHandler{promotionsCsv: filename, promotionsCmdFilename: cmds, s: s}
}

// Load loads promotions from a source into a data storage
func (p *LoadPromotionsHandler) Load() (int64, error) {
	log.Debug().Msgf("reading promotions file %s", p.promotionsCsv)
	pf, err := os.Open(p.promotionsCsv)
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

	parts := strings.Split(p.promotionsCsv, ".csv")
	nm := parts[0]
	now := time.Now().UTC().Format("20060102150405")

	err = os.Rename(p.promotionsCsv, fmt.Sprintf("%s--%s--imported.csv", nm, now))
	if err != nil {
		return -1, err
	}

	return c, nil
}
