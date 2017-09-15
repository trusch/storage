package aes

import (
	"io"

	"github.com/trusch/storage"
)

// Encrypter is a Storage which encrypts with aes
type Encrypter struct {
	base storage.Storage
	key  string
}

// NewEncrypter returns a new encrypter instance
func NewEncrypter(base storage.Storage, key string) *Encrypter {
	return &Encrypter{base, key}
}

// GetReader returns a reader
func (encrypter *Encrypter) GetReader(id string) (io.ReadCloser, error) {
	baseReader, err := encrypter.base.GetReader(id)
	if err != nil {
		return nil, err
	}
	return NewReader(baseReader, encrypter.key)
}

// GetWriter returns a writer
func (encrypter *Encrypter) GetWriter(id string) (io.WriteCloser, error) {
	baseWriter, err := encrypter.base.GetWriter(id)
	if err != nil {
		return nil, err
	}
	return NewWriter(baseWriter, encrypter.key)
}

// Has returns whether an entry exists
func (encrypter *Encrypter) Has(id string) bool {
	return encrypter.base.Has(id)
}

// Delete deletes an entry
func (encrypter *Encrypter) Delete(id string) error {
	return encrypter.base.Delete(id)
}

// List lists all stored objects, limited by a prefix
func (encrypter *Encrypter) List(prefix string) ([]string, error) {
	return encrypter.base.List(prefix)
}
