package handlers

import (
	"bufio"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/vitorsalgado/goprom/internal/utils/config"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	"strings"
	"time"
)

type (
	// Streamer handles promotions file data chunks
	Streamer interface {
		Stream(w io.StringWriter, chunk []string) error
		Push(filename string) error
	}

	LoaderSource interface {
		File(filename string) (*os.File, error)
	}

	LoaderLifecycle interface {
		OnFinish(filename string) error
	}

	// LoadPromotionsHandler loads promotions from a source file into a data storage
	LoadPromotionsHandler struct {
		s   Streamer
		lc  LoaderLifecycle
		src LoaderSource
		cfg *config.Config
		ctx context.Context
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
	cfg *config.Config, ctx context.Context, s Streamer, src LoaderSource, lc LoaderLifecycle,
) *LoadPromotionsHandler {
	return &LoadPromotionsHandler{
		cfg: cfg, ctx: ctx, s: s, src: src, lc: lc}
}

// Load loads promotions from a source into a data storage
func (p *LoadPromotionsHandler) Load() (int64, error) {
	log.Info().Msgf("loading promotions from source %s", p.cfg.PromotionsCsv)

	pf, err := p.src.File(p.cfg.PromotionsCsv)
	if err != nil {
		log.Error().Err(err).Msgf("error loading promotions from source %s", p.cfg.PromotionsCsv)
		return -1, err
	}

	defer pf.Close()

	scanner := bufio.NewScanner(pf)
	ch := make(chan string)

	group, ctx := errgroup.WithContext(p.ctx)

	for i := 0; i < p.cfg.PromotionsBulkLoadWorkers; i++ {
		d := i
		group.Go(func() error {
			df, e := os.Create(fmt.Sprintf(p.cfg.PromotionsBulkCmdFilename, d))
			if e != nil {
				log.Error().Err(e).Msg("error creating redis commands file")
				return e
			}

			cmds := bufio.NewWriter(df)

			for {
				select {
				case chunk, ok := <-ch:
					if !ok {
						_ = cmds.Flush()
						_ = df.Close()
						return nil
					}

					_ = p.s.Stream(cmds, strings.Split(chunk, ","))

				case <-ctx.Done():
					_ = cmds.Flush()
					_ = df.Close()
					return nil
				}
			}
		})
	}

	var c int64

	for scanner.Scan() {
		ch <- scanner.Text()
		c++
	}

	close(ch)

	if err = group.Wait(); err != nil {
		log.Error().Stack().Err(err).Msgf("error processing promotions")
		return -1, err
	}

	err = scanner.Err()
	if err != nil {
		log.Error().Err(err).Msg("error with promotions file scanner")
		return -1, err
	}

	group, ctx = errgroup.WithContext(ctx)

	for i := 0; i < p.cfg.PromotionsBulkLoadWorkers; i++ {
		d := i
		group.Go(func() error {
			err = p.s.Push(fmt.Sprintf(p.cfg.PromotionsBulkCmdFilename, d))
			if err != nil {
				return err
			}

			return nil
		})
	}

	if err = group.Wait(); err != nil {
		return -1, err
	}

	err = p.lc.OnFinish(p.cfg.PromotionsCsv)

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
