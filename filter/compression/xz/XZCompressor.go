package xz

import (
	"io"

	"github.com/trusch/storage"
	"github.com/ulikunitz/xz"
)

// Compressor is a storage wrapper which applies gzip compression
type Compressor struct {
	base storage.Storage
}

// NewCompressor returns a new compressor instance
func NewCompressor(base storage.Storage) *Compressor {
	return &Compressor{base}
}

// GetReader returns a reader
func (compressor *Compressor) GetReader(id string) (io.ReadCloser, error) {
	baseReader, err := compressor.base.GetReader(id)
	if err != nil {
		return nil, err
	}
	xzReader, err := xz.NewReader(baseReader)
	if err != nil {
		return nil, err
	}
	return &storage.ReadCloser{Closer: baseReader, Reader: xzReader}, nil
}

// GetWriter returns a writer
func (compressor *Compressor) GetWriter(id string) (io.WriteCloser, error) {
	baseWriter, err := compressor.base.GetWriter(id)
	if err != nil {
		return nil, err
	}
	writer, err := xz.NewWriter(baseWriter)
	if err != nil {
		return nil, err
	}
	return &storage.WriteCloser{Closer: baseWriter, Writer: writer}, nil
}

// Has returns whether an entry exists
func (compressor *Compressor) Has(id string) bool {
	return compressor.base.Has(id)
}

// Delete deletes an entry
func (compressor *Compressor) Delete(id string) error {
	return compressor.base.Delete(id)
}

// List lists all stored objects, limited by a prefix
func (compressor *Compressor) List(prefix string) ([]string, error) {
	return compressor.base.List(prefix)
}
