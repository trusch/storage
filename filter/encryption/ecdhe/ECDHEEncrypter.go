package ecdhe

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
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

	// extract ephemeral public key
	var (
		ephLen [1]byte
		ephPub []byte
	)
	_, err = baseReader.Read(ephLen[:1])
	if err != nil {
		return nil, err
	}
	ephPub = make([]byte, int(ephLen[0]))
	_, err = baseReader.Read(ephPub[:])
	if err != nil {
		return nil, err
	}
	x, y := elliptic.Unmarshal(elliptic.P256(), ephPub)
	ok := elliptic.P256().IsOnCurve(x, y) // Rejects the identity point too.
	if x == nil || !ok {
		return nil, errors.New("Invalid public key")
	}

	// compute shared secret by multiplying ephemeral pubkey with own private key
	priv := encrypter.key
	x, _ = priv.Curve.ScalarMult(x, y, priv.D.Bytes())
	if x == nil {
		return nil, errors.New("Failed to generate encryption key")
	}
	shared := sha256.Sum256(x.Bytes())

	// start AES mode
	block, err := aes.NewCipher(shared[:])
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
	if encrypter.cert == nil {
		return nil, errors.New("no public key supplied, cannot encrypt")
	}

	baseWriter, err := encrypter.base.GetWriter(id)
	if err != nil {
		return nil, err
	}

	// create ephemeral ec key
	ephemeral, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	// compute shared secret by multiplying the public key with ephemeral private key
	pub := encrypter.cert.PublicKey.(*ecdsa.PublicKey)
	x, _ := pub.Curve.ScalarMult(pub.X, pub.Y, ephemeral.D.Bytes())
	if x == nil {
		return nil, errors.New("Failed to generate encryption key")
	}
	shared := sha256.Sum256(x.Bytes())

	// write ephemeral public key
	ephPub := elliptic.Marshal(pub.Curve, ephemeral.PublicKey.X, ephemeral.PublicKey.Y)
	_, err = baseWriter.Write([]byte{byte(len(ephPub))})
	if err != nil {
		return nil, err
	}
	_, err = baseWriter.Write(ephPub)
	if err != nil {
		return nil, err
	}

	// start AES mode
	block, err := aes.NewCipher(shared[:])
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
