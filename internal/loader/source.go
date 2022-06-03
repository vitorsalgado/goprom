package loader

import "os"

// Source abstracts promotions file source loading
type Source interface {
	File(filename string) (*os.File, error)
}

// LocalFileSource reads promotions from local filesystem
type LocalFileSource struct {
}

func (src *LocalFileSource) File(filename string) (*os.File, error) {
	return os.Open(filename)
}
