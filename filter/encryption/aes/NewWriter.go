package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
)

// NewWriter returns a new aes writer
func NewWriter(base io.Writer, key string) (io.WriteCloser, error) {
	k := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(k[:])
	if err != nil {
		return nil, err
	}
	iv := make([]byte, aes.BlockSize)
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	_, err = base.Write(iv)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewOFB(block, iv[:])
	writer := &cipher.StreamWriter{S: stream, W: base}
	return writer, nil
}
