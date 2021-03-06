package jwk

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	"github.com/pkg/errors"
	"github.com/square/go-jose"
)

type RS256Generator struct {
	KeyLength int
}

func (g *RS256Generator) Generate(id string) (*jose.JSONWebKeySet, error) {
	if g.KeyLength < 4096 {
		g.KeyLength = 4096
	}

	key, err := rsa.GenerateKey(rand.Reader, g.KeyLength)
	if err != nil {
		return nil, errors.Errorf("Could not generate key because %s", err)
	} else if err = key.Validate(); err != nil {
		return nil, errors.Errorf("Validation failed because %s", err)
	}

	// jose does not support this...
	key.Precomputed = rsa.PrecomputedValues{}
	return &jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{
				Algorithm:    "RS256",
				Key:          key,
				KeyID:        ider("private", id),
				Certificates: []*x509.Certificate{},
			},
			{
				Algorithm:    "RS256",
				Key:          &key.PublicKey,
				KeyID:        ider("public", id),
				Certificates: []*x509.Certificate{},
			},
		},
	}, nil
}

func ider(typ, id string) string {
	if id != "" {
		return fmt.Sprintf("%s:%s", typ, id)
	}
	return typ
}
