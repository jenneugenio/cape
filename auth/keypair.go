package auth

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"

	"github.com/manifoldco/go-base64"
	"golang.org/x/crypto/scrypt"

	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

const (
	SecretLength              = 8
	GeneratedSecretByteLength = 8

	// N , R, P are used by scrypt, see scrypt docs below
	// https://godoc.org/golang.org/x/crypto/scrypt#Key
	N        = 32768
	R        = 8
	P        = 1
	SeedSize = 32
)

// Keypair represents a ed25519 Private and Public Keypair. This struct is used
// by our TokenAuthority and Credentials struct to derive and package up
// Keypairs.
//
// It is a wrapper around ed25519 and implements functionality for packaging up
// and rederiving keypairs. This is particularly useful for creating and
// recreating Token Authorities.
type Keypair struct {
	secret     []byte
	salt       []byte
	PrivateKey ed25519.PrivateKey `json:"private_key"`

	PublicKey ed25519.PublicKey             `json:"public_key"`
	Alg       primitives.CredentialsAlgType `json:"alg"`
}

// KeypairPackage represents a packaged keypair that can be shared outside of
// the Auth package.
type KeypairPackage struct {
	Secret *base64.Value                 `json:"secret"`
	Salt   *base64.Value                 `json:"salt"`
	Alg    primitives.CredentialsAlgType `json:"alg"`
}

// Validate returns an error if the packaged keypair is invalid
func (kp *KeypairPackage) Validate() error {
	if kp.Salt == nil || kp.Secret == nil {
		return errors.New(BadPackagedKeypair, "Salt or Secret missing")
	}

	salt := []byte(*kp.Salt)
	secret := []byte(*kp.Secret)

	if len(salt) != primitives.SaltLength {
		return errors.New(BadSaltLength, "Salt must be at least %d bytes long, saw %d",
			primitives.SaltLength, len(salt))
	}

	if len(secret) < SecretLength {
		return errors.New(BadSecretLength, "Secret must be at least %d bytes long, saw %d", SecretLength, len(secret))
	}

	if kp.Alg != primitives.EDDSA {
		return errors.New(BadAlgType, "Algorithm %s not recognized", kp.Alg)
	}

	return nil
}

// Unpackage returns a Keypair from the packaging
func (kp *KeypairPackage) Unpackage() (*Keypair, error) {
	if err := kp.Validate(); err != nil {
		return nil, err
	}

	return DeriveKeypair([]byte(*kp.Secret), []byte(*kp.Salt))
}

// Package returns a KeypairPackage which can be serialized to JSON
func (k *Keypair) Package() KeypairPackage {
	return KeypairPackage{
		Secret: base64.New(k.secret),
		Salt:   base64.New(k.salt),
		Alg:    k.Alg,
	}
}

// DeriveKeypair returns a new keypair for
func DeriveKeypair(secret []byte, salt []byte) (*Keypair, error) {
	if len(salt) != primitives.SaltLength {
		return nil, errors.New(BadSaltLength, "Salt must be at least %d bytes long, saw %d",
			primitives.SaltLength, len(salt))
	}

	if len(secret) < SecretLength {
		return nil, errors.New(BadSecretLength, "Secret must be at least %d bytes long, saw %d",
			SecretLength, len(secret))
	}

	// Derive a key by stretching & salting the secret into 32 bytes
	key, err := scrypt.Key(secret, salt, N, R, P, SeedSize)
	if err != nil {
		return nil, err
	}

	// Generate a ed25519 keypair from the derived key
	pk, sk, err := ed25519.GenerateKey(bytes.NewBuffer(key))
	if err != nil {
		return nil, err
	}

	return &Keypair{
		secret:     secret,
		salt:       salt,
		PrivateKey: sk,
		PublicKey:  pk,
		Alg:        primitives.EDDSA,
	}, nil
}

// NewKeypair returns a keypair generated from a newly created secret and salt!
func NewKeypair() (*Keypair, error) {
	secret := make([]byte, SecretLength)
	salt := make([]byte, primitives.SaltLength)

	_, err := rand.Read(secret)
	if err != nil {
		return nil, err
	}

	_, err = rand.Read(salt)
	if err != nil {
		return nil, err
	}

	return DeriveKeypair(secret, salt)
}

// NewKeypairWithSecret returns a keypair generated from the provided secret
// and a generated salt!
func NewKeypairWithSecret(secret []byte) (*Keypair, error) {
	salt := make([]byte, primitives.SaltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	return DeriveKeypair(secret, salt)
}
