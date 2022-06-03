package loader

import "os"

// Source abstracts promotions file source loading
type Source interface {
	File(filename string) (*os.File, error)
}

// localFileSource reads promotions from local filesystem
type localFileSource struct {
}

func NewSource() Source {
	return &localFileSource{}
}

func (src *localFileSource) File(filename string) (*os.File, error) {
	return os.Open(filename)
}
