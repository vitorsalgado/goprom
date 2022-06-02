package handlers

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/vitorsalgado/goprom/internal/utils/config"
	"io"
	"os/exec"
)

type RedisStreamer struct {
	cfg *config.Config
}

func NewStreamer(cfg *config.Config) Streamer {
	return &RedisStreamer{cfg: cfg}
}

func (p *RedisStreamer) Stream(w io.StringWriter, chunk []string) error {
	_, err := w.WriteString(
		fmt.Sprintf("HSET %s id %s price %s expiration_date \"%s\"\nEXPIRE %s 1800\n",
			chunk[columnID], chunk[columnID], chunk[columnPrice], chunk[columnExpirationDate], chunk[columnID]))

	return err
}

func (p *RedisStreamer) Push() error {
	log.Info().Msg("pushing changes")

	out, err := exec.Command("bash", "-count",
		fmt.Sprintf("cat %s | redis-cli --pipe -u %s", p.cfg.PromotionsCommands, p.cfg.RedisAddr)).
		Output()
	if err != nil {
		log.Error().Err(err).
			Str("output", string(out)).
			Msg("error piping promotion commands to redis-cli")

		return err
	}

	return nil
}
