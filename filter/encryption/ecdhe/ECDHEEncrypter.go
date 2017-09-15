package ecdhe

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"

	"github.com/trusch/storage"

	log "github.com/sirupsen/logrus"
)

// Encrypter is a Storage which encrypts with aes
type Encrypter struct {
	base storage.Storage
	cert *x509.Certificate
	key  *ecdsa.PrivateKey
}

// NewEncrypter returns a new encrypter instance
func NewEncrypter(base storage.Storage, pubkeyPEM, privkeyPEM []byte) (*Encrypter, error) {
	store := &Encrypter{base: base}
	pubkeyBlock, _ := pem.Decode(pubkeyPEM)
	if pubkeyBlock == nil {
		log.Warn("no valid pem data in pubkey")
	} else {
		cert, err := x509.ParseCertificate(pubkeyBlock.Bytes)
		if err != nil {
			return nil, err
		}
		store.cert = cert
	}

	privkeyBlock, _ := pem.Decode(privkeyPEM)
	if privkeyBlock == nil {
		log.Warn("no valid pem data in privkey")
	} else {
		key, err := x509.ParseECPrivateKey(privkeyBlock.Bytes)
		if err != nil {
			return nil, err
		}
		store.key = key
	}

	if store.cert == nil && store.key == nil {
		return nil, errors.New("neither pub nor priv key supplied")
	}
	return store, nil
}

// GetReader returns a reader
func (encrypter *Encrypter) GetReader(id string) (io.ReadCloser, error) {
	if encrypter.key == nil {
		return nil, errors.New("no private key supplied, cannot decrypt")
	}

	baseReader, err := encrypter.base.GetReader(id)
	if err != nil {
		return nil, err
	}

	return NewReader(baseReader, encrypter.key)
}

// GetWriter returns a writer
func (encrypter *Encrypter) GetWriter(id string) (io.WriteCloser, error) {
	if encrypter.cert == nil {
		return nil, errors.New("no public key supplied, cannot encrypt")
	}

	baseWriter, err := encrypter.base.GetWriter(id)
	if err != nil {
		return nil, err
	}

	return NewWriter(baseWriter, encrypter.cert)
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
