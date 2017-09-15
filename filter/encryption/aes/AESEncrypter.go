package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"

	"github.com/trusch/storage"
)

// Encrypter is a Storage which encrypts with aes
type Encrypter struct {
	base storage.Storage
	key  [sha256.Size]byte
}

// NewEncrypter returns a new encrypter instance
func NewEncrypter(base storage.Storage, key string) *Encrypter {
	keyBytes := sha256.Sum256([]byte(key))
	return &Encrypter{base, keyBytes}
}

// GetReader returns a reader
func (encrypter *Encrypter) GetReader(id string) (io.ReadCloser, error) {
	baseReader, err := encrypter.base.GetReader(id)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(encrypter.key[:])
	if err != nil {
		return nil, err
	}
	iv := make([]byte, aes.BlockSize)
	bs, err := baseReader.Read(iv[:])
	if bs != aes.BlockSize {
		return nil, errors.New("ciphertext to short")
	}
	if err != nil {
		return nil, err
	}
	stream := cipher.NewOFB(block, iv[:])
	reader := &cipher.StreamReader{S: stream, R: baseReader}
	return &storage.ReadCloser{Closer: baseReader, Reader: reader}, nil
}

// GetWriter returns a writer
func (encrypter *Encrypter) GetWriter(id string) (io.WriteCloser, error) {
	baseWriter, err := encrypter.base.GetWriter(id)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(encrypter.key[:])
	if err != nil {
		return nil, err
	}
	iv := make([]byte, aes.BlockSize)
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	_, err = baseWriter.Write(iv)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewOFB(block, iv[:])
	writer := &cipher.StreamWriter{S: stream, W: baseWriter}
	return writer, nil
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
