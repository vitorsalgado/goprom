package loader

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Lifecycle exposes hooks to load process
type Lifecycle interface {
	OnFinish(filename string) error
}

type DefaultLifecycle struct {
}

func (lc *DefaultLifecycle) OnFinish(filename string) error {
	parts := strings.Split(filename, ".csv")
	nm := parts[0]
	now := time.Now().UTC().Format("20060102150405")

	err := os.Rename(filename, fmt.Sprintf("%s--%s--imported.csv", nm, now))
	if err != nil {
		return err
	}

	return nil
}
