package loader

import (
	"fmt"
	"github.com/vitorsalgado/goprom/internal/std/config"
	"io"
	"os"
	"os/exec"
)

type Writer struct {
	decorated io.WriteCloser
	cmd       *exec.Cmd
}

func (rw *Writer) Write(chunk []byte) (int, error) {
	return rw.decorated.Write(chunk)
}

func (rw *Writer) Close() error {
	err := rw.decorated.Close()
	if err != nil {
		return err
	}

	return rw.cmd.Wait()
}

func RedisWriterFn(cfg *config.Config) WriterProvider {
	return func() (io.WriteCloser, error) {
		cmd := exec.Command("bash", "-c", fmt.Sprintf("redis-cli --pipe -u redis://%s", cfg.RedisAddr))
		in, err := cmd.StdinPipe()
		if err != nil {
			return nil, err
		}

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Start()
		if err != nil {
			return nil, err
		}

		return &Writer{decorated: in, cmd: cmd}, nil
	}
}
