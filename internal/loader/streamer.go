package loader

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/vitorsalgado/goprom/internal/std/config"
	"io"
	"os/exec"
	"strconv"
	"time"
)

type RedisStreamer struct {
	cfg *config.Config
}

func NewStreamer(cfg *config.Config) Streamer {
	return &RedisStreamer{cfg: cfg}
}

func (p *RedisStreamer) Stream(w io.StringWriter, chunk []string) error {
	dt, err := time.Parse("2006-01-02 15:04:05 -0700 MST", chunk[columnExpirationDate])
	if err != nil {
		return err
	}

	price, err := strconv.ParseFloat(chunk[columnPrice], 64)
	if err != nil {
		return err
	}

	_, err = w.WriteString(
		fmt.Sprintf("HSET %s id %s price %s expiration_date \"%s\"\nEXPIRE %s 1800\n",
			chunk[columnID], chunk[columnID], fmt.Sprintf("%.2f", price), dt.Format(time.RFC3339Nano), chunk[columnID]))

	return err
}

func (p *RedisStreamer) Push(filename string) error {
	log.Info().Msg("pushing changes to Redis")

	out, err := exec.Command("bash", "-c",
		fmt.Sprintf("cat %s | redis-cli --pipe -u redis://%s", filename, p.cfg.RedisAddr)).
		Output()
	if err != nil {
		log.Error().Err(err).
			Str("output", string(out)).
			Msg("error piping promotion commands to redis-cli")

		return err
	}

	return nil
}
