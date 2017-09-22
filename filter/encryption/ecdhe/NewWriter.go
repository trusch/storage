package ecdhe

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"io"

	"github.com/trusch/storage/filter/encryption/aes"
)

// NewWriter returns a new ecdhe writer
func NewWriter(base io.Writer, pubkey *ecdsa.PublicKey) (io.WriteCloser, error) {
	// create ephemeral ec key
	ephemeral, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	// compute shared secret by multiplying the public key with ephemeral private key
	pub := pubkey
	x, _ := pub.Curve.ScalarMult(pub.X, pub.Y, ephemeral.D.Bytes())
	if x == nil {
		return nil, errors.New("Failed to generate encryption key")
	}

	// write ephemeral public key
	ephPub := elliptic.Marshal(pub.Curve, ephemeral.PublicKey.X, ephemeral.PublicKey.Y)
	_, err = base.Write([]byte{byte(len(ephPub))})
	if err != nil {
		return nil, err
	}
	_, err = base.Write(ephPub)
	if err != nil {
		return nil, err
	}
	return aes.NewWriter(base, x.String())
}
