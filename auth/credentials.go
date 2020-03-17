package auth

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"

	"github.com/manifoldco/go-base64"
	"golang.org/x/crypto/scrypt"

	errors "github.com/dropoutlabs/cape/partyerrors"
	"github.com/dropoutlabs/cape/primitives"
)

const (
	SaltLength   = 16
	SecretLength = 8

	// N , R, P are used by scrypt, see scrypt docs below
	// https://godoc.org/golang.org/x/crypto/scrypt#Key
	N        = 32768
	R        = 8
	P        = 1
	SeedSize = 32
)

// Credentials holds the public key and salt for a user
type Credentials struct {
	// privateKey not saved in the database or sent anywhere!
	privateKey ed25519.PrivateKey
	PublicKey  *base64.Value                 `json:"public_key"`
	Salt       *base64.Value                 `json:"salt"`
	Alg        primitives.CredentialsAlgType `json:"alg"`
}

//Sign signs a message token
func (c *Credentials) Sign(token *base64.Value) (*base64.Value, error) {
	if token == nil {
		return nil, errors.New(RequiredMsgCause, "Must provide message to sign")
	}

	if c.privateKey == nil {
		return nil, errors.New(RequiredPrivateKeyCause,
			"Cannot sign token, missing private key")
	}

	return base64.New(ed25519.Sign(c.privateKey, []byte(*token))), nil
}

// Verify verifies a given signature based on the token and the public key
func (c *Credentials) Verify(token, sig *base64.Value) error {
	if token == nil {
		return errors.New(RequiredMsgCause, "Must provide message to verify")
	}

	if sig == nil {
		return errors.New(RequiredSigCause, "Must provide signature to verify")
	}

	if c.PublicKey == nil {
		return errors.New(RequiredPublicKeyCause,
			"Credentials must have public key to verify with")
	}

	if !ed25519.Verify([]byte(*c.PublicKey), []byte(*token), []byte(*sig)) {
		return errors.New(SignatureNotValid,
			"Signature does not match the public key")
	}

	return nil
}

// Package packages up a auth.Credentials into a primitives.Credentials
func (c *Credentials) Package() *primitives.Credentials {
	return &primitives.Credentials{
		PublicKey: c.PublicKey,
		Salt:      c.Salt,
		Alg:       c.Alg,
	}
}

// NewCredentials returns a Credential struct for creating credentials & re-deriving credentials.
// To rederive credentials you must provide the same secret & salt with the same Alg parameters as used previously.
func NewCredentials(secret []byte, salt *base64.Value) (*Credentials, error) {
	var saltBytes []byte
	if salt == nil {
		saltBytes = make([]byte, SaltLength)
		_, err := rand.Read(saltBytes)
		if err != nil {
			return nil, err
		}
	} else {
		saltBytes = []byte(*salt)
	}

	pub, priv, err := deriveKey(secret, saltBytes)
	if err != nil {
		return nil, err
	}

	return &Credentials{
		privateKey: priv,
		PublicKey:  base64.New([]byte(pub)),
		Salt:       base64.New(saltBytes),
		Alg:        primitives.EDDSA,
	}, nil
}

// LoadCredentials loads credentials
func LoadCredentials(publicKey, salt *base64.Value) (*Credentials, error) {
	if len(*publicKey) != ed25519.PublicKeySize {
		return nil, errors.New(BadPublicKeyLength, "Public key must be %d bytes long, saw %d",
			ed25519.PublicKeySize, len(*publicKey))
	}

	if len(*salt) != SaltLength {
		return nil, errors.New(BadSaltLength, "Salt must be at least %d bytes long, saw %d",
			SaltLength, len(*salt))
	}

	return &Credentials{
		PublicKey: publicKey,
		Salt:      salt,
		Alg:       primitives.EDDSA,
	}, nil
}

func deriveKey(secret []byte, salt []byte) (ed25519.PublicKey, ed25519.PrivateKey, error) {
	if len(salt) != SaltLength {
		return nil, nil, errors.New(BadSaltLength, "Salt must be at least %d bytes long, saw %d",
			SaltLength, len(salt))
	}

	if len(secret) < SecretLength {
		return nil, nil, errors.New(BadSecretLength, "Secret must be at least %d bytes long, saw %d",
			SecretLength, len(secret))
	}

	// Derive a key by stretching & salting the secret into 32 bytes
	key, err := scrypt.Key(secret, salt, N, R, P, SeedSize)
	if err != nil {
		return nil, nil, err
	}

	// Generate a ed25519 keypair from the derived key
	return ed25519.GenerateKey(bytes.NewBuffer(key))
}

func newKey() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	return ed25519.GenerateKey(nil)
}
