package storage

import "io"

// Storage is the main storage interface
type Storage interface {
	ReadableStorage
	WriteableStorage
}

// ReadableStorage is the interface for read operations
type ReadableStorage interface {
	GetReader(id string) (io.ReadCloser, error)
	List(prefix string) ([]string, error)
	Has(id string) bool
}

// WriteableStorage is the interface for write operations
type WriteableStorage interface {
	GetWriter(id string) (io.WriteCloser, error)
	Delete(id string) error
}
