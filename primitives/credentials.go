package primitives

import (
	"crypto/rand"

	"github.com/manifoldco/go-base64"

	errors "github.com/capeprivacy/cape/partyerrors"
)

const (
	SecretLength = 32
	SaltLength   = 16
)

type Credentials struct {
	// Secret is the "scrypt'ed" secret which we store in the database.
	Secret *base64.Value      `json:"secret"`
	Salt   *base64.Value      `json:"salt"`
	Alg    CredentialsAlgType `json:"alg"`
}

func (c *Credentials) Validate() error {
	if c.Secret == nil {
		return errors.New(InvalidCredentialsCause, "Credentials secret must be non-nil")
	}

	if c.Salt == nil {
		return errors.New(InvalidCredentialsCause, "Credentials salt must be non-nil")
	}

	if len(*c.Secret) != SecretLength {
		return errors.New(InvalidCredentialsCause, "Credentials secret length must be %d", SecretLength)
	}

	if len(*c.Salt) != SaltLength {
		return errors.New(InvalidCredentialsCause, "Credentials salt length must be %d", SaltLength)
	}

	if err := c.Alg.Validate(); err != nil {
		return errors.Wrap(InvalidCredentialsCause, err)
	}

	return nil
}

func NewCredentials(secret, salt *base64.Value, alg CredentialsAlgType) (*Credentials, error) {
	c := &Credentials{
		Secret: secret,
		Salt:   salt,
		Alg:    alg,
	}

	return c, c.Validate()
}

// GenerateCredentials returns an instantiated Credentials for use in unit testing.
//
// This function _should only ever_ be used inside of a test.
func GenerateCredentials() (*Credentials, error) {
	secret := make([]byte, SecretLength)
	salt := make([]byte, SaltLength)

	_, err := rand.Read(secret)
	if err != nil {
		return nil, err
	}

	_, err = rand.Read(salt)
	if err != nil {
		return nil, err
	}

	return &Credentials{
		Secret: base64.New(secret),
		Salt:   base64.New(salt),
		Alg:    SHA256,
	}, nil
}
