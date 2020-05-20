package auth

import (
	"crypto/ed25519"

	"github.com/manifoldco/go-base64"

	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

// Credentials holds the public key and salt for a user
type Credentials struct {
	// privateKey not saved in the database or sent anywhere!
	privateKey ed25519.PrivateKey
	PublicKey  *base64.Value                 `json:"public_key"`
	Salt       *base64.Value                 `json:"salt"`
	Alg        primitives.CredentialsAlgType `json:"alg"`
}

// Sign signs a message token
func (c *Credentials) Sign(token *base64.Value) (*base64.Value, error) {
	if token == nil {
		return nil, errors.New(RequiredTokenCause, "Must provide token to sign")
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
		return errors.New(RequiredTokenCause, "Must provide token to verify")
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
func (c *Credentials) Package() (*primitives.Credentials, error) {
	return primitives.NewCredentials(c.PublicKey, c.Salt)
}

// NewCredentials returns a Credential struct for creating credentials & re-deriving credentials.
// To rederive credentials you must provide the same secret & salt with the same Alg parameters as used previously.
func NewCredentials(secret []byte, salt *base64.Value) (*Credentials, error) {
	var keypair *Keypair
	var err error
	if salt == nil {
		keypair, err = NewKeypairWithSecret(secret)
		if err != nil {
			return nil, err
		}
	} else {
		keypair, err = DeriveKeypair(secret, []byte(*salt))
		if err != nil {
			return nil, err
		}
	}

	return &Credentials{
		privateKey: keypair.PrivateKey,
		PublicKey:  base64.New(keypair.PublicKey),
		Salt:       base64.New(keypair.salt),
		Alg:        keypair.Alg,
	}, nil
}

// LoadCredentials loads credentials
func LoadCredentials(publicKey, salt *base64.Value) (*Credentials, error) {
	if len(*publicKey) != ed25519.PublicKeySize {
		return nil, errors.New(BadPublicKeyLength, "Public key must be %d bytes long, saw %d",
			ed25519.PublicKeySize, len(*publicKey))
	}

	if len(*salt) != primitives.SaltLength {
		return nil, errors.New(BadSaltLength, "Salt must be at least %d bytes long, saw %d",
			primitives.SaltLength, len(*salt))
	}

	return &Credentials{
		PublicKey: publicKey,
		Salt:      salt,
		Alg:       primitives.EDDSA,
	}, nil
}
