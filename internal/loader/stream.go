package loader

import (
	"fmt"
	"github.com/vitorsalgado/goprom/internal/domain"
	"io"
	"strconv"
	"time"
)

// Streamer handles promotions file data chunks
type Streamer interface {
	Stream(w io.Writer, chunk []string) error
}

type WriterProvider func() (io.WriteCloser, error)

type redisStreamer struct {
}

func NewStreamer() Streamer {
	return &redisStreamer{}
}

func (p *redisStreamer) Stream(w io.Writer, chunk []string) error {
	dt, err := time.Parse("2006-01-02 15:04:05 -0700 MST", chunk[columnExpirationDate])
	if err != nil {
		return err
	}

	price, err := strconv.ParseFloat(chunk[columnPrice], 64)
	if err != nil {
		return err
	}

	_, err = w.Write(
		[]byte(fmt.Sprintf("HSET %s id %s price %s expiration_date \"%s\"\nEXPIRE %s 1800\n",
			chunk[columnID], chunk[columnID], fmt.Sprintf("%.2f", price), dt.Format(domain.PromotionDatetimeFormat), chunk[columnID])))

	return err
}
