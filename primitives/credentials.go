package primitives

import (
	"crypto/ed25519"

	"github.com/manifoldco/go-base64"

	errors "github.com/capeprivacy/cape/partyerrors"
)

const (
	SaltLength = 16
)

type Credentials struct {
	PublicKey *base64.Value      `json:"public_key"`
	Salt      *base64.Value      `json:"salt"`
	Alg       CredentialsAlgType `json:"alg"`
}

func (c *Credentials) Validate() error {
	if c.PublicKey == nil {
		return errors.New(InvalidCredentialsCause, "Credentials public key must not be nil")
	}

	if len(*c.PublicKey) != ed25519.PublicKeySize {
		return errors.New(InvalidCredentialsCause, "Credentials public key length must be %d", ed25519.PublicKeySize)
	}

	if c.Salt == nil {
		return errors.New(InvalidCredentialsCause, "Credentials salt must not be nil")
	}

	if len(*c.Salt) != SaltLength {
		return errors.New(InvalidCredentialsCause, "Credentials salt must be %d", SaltLength)
	}

	if err := c.Alg.Validate(); err != nil {
		return errors.Wrap(InvalidCredentialsCause, err)
	}

	return nil
}

func NewCredentials(publicKey *base64.Value, salt *base64.Value) (*Credentials, error) {
	creds := &Credentials{
		PublicKey: publicKey,
		Salt:      salt,
		Alg:       EDDSA,
	}

	return creds, creds.Validate()
}
