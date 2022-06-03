package loader

import (
	"bufio"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/vitorsalgado/goprom/internal/std/config"
	"golang.org/x/sync/errgroup"
	"os"
	"strings"
	"sync/atomic"
)

type (
	// Handler loads promotions from a source file into a data storage
	Handler struct {
		s   Streamer
		lc  Lifecycle
		src Source
		cfg *config.Config
		ctx context.Context
	}
)

const (
	columnID             = 0
	columnPrice          = 1
	columnExpirationDate = 2
)

// NewLoader initiates a new instance of Handler
func NewLoader(
	cfg *config.Config, ctx context.Context, s Streamer, src Source, lc Lifecycle,
) *Handler {
	return &Handler{
		cfg: cfg, ctx: ctx, s: s, src: src, lc: lc}
}

// Load loads promotions from a source into a data storage
func (p *Handler) Load() (int64, error) {
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
	var c int64

	go func() {
		for scanner.Scan() {
			ch <- scanner.Text()
			atomic.AddInt64(&c, 1)
		}

		close(ch)
	}()

	for i := 0; i < p.cfg.PromotionsBulkLoadWorkers; i++ {
		d := i
		group.Go(func() error {
			df, e := os.Create(fmt.Sprintf(p.cfg.PromotionsBulkCmdFilename, d))
			if e != nil {
				log.Error().Err(e).Msg("error creating redis commands file")
				return e
			}

			cmds := bufio.NewWriter(df)
			defer cmds.Flush()
			defer df.Close()

			for {
				select {
				case chunk, ok := <-ch:
					if !ok {
						_ = cmds.Flush()
						return nil
					}

					_ = p.s.Stream(cmds, strings.Split(chunk, ","))

				case <-ctx.Done():
					_ = cmds.Flush()
					return nil
				}
			}
		})
	}

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
			return p.s.Push(fmt.Sprintf(p.cfg.PromotionsBulkCmdFilename, d))
		})
	}

	if err = group.Wait(); err != nil {
		return -1, err
	}

	err = p.lc.OnFinish(p.cfg.PromotionsCsv)

	return c, err
}
