package ecdhe

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"io"

	"github.com/trusch/storage/filter/encryption/aes"
)

// NewReader returns a new ecdhe reader
func NewReader(base io.Reader, key *ecdsa.PrivateKey) (io.ReadCloser, error) {
	// extract ephemeral public key
	var (
		ephLen [1]byte
		ephPub []byte
	)
	_, err := base.Read(ephLen[:1])
	if err != nil {
		return nil, err
	}
	ephPub = make([]byte, int(ephLen[0]))
	_, err = base.Read(ephPub[:])
	if err != nil {
		return nil, err
	}
	x, y := elliptic.Unmarshal(elliptic.P256(), ephPub)
	ok := elliptic.P256().IsOnCurve(x, y) // Rejects the identity point too.
	if x == nil || !ok {
		return nil, errors.New("Invalid public key")
	}

	// compute shared secret by multiplying ephemeral pubkey with own private key
	x, _ = key.Curve.ScalarMult(x, y, key.D.Bytes())
	if x == nil {
		return nil, errors.New("Failed to generate encryption key")
	}
	return aes.NewReader(base, x.String())
}
