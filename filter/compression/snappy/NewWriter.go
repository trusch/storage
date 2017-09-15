package snappy

import (
	"io"

	"github.com/golang/snappy"
	"github.com/trusch/storage"
)

// NewWriter returns a new snappy writer
func NewWriter(base io.Writer) (io.WriteCloser, error) {
	snappyWriter := snappy.NewWriter(base)
	return storage.NewIOCoppler(snappyWriter, base), nil
}
