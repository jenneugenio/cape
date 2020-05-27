package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"

	"github.com/manifoldco/go-base64"
	"golang.org/x/crypto/argon2"

	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

var (
	DefaultSHA256Producer   = &SHA256Producer{}
	DefaultArgon2IDProducer = &Argon2IDProducer{
		Time:      1,
		Memory:    64 * 1024,
		Threads:   4,
		KeyLength: primitives.SecretLength,
	}
)

// CredentialProducer represents an interface for generating credentials and
// comparing credentials based on a pre-shared key.
type CredentialProducer interface {
	Generate(primitives.Password) (*primitives.Credentials, error)
	Compare(primitives.Password, *primitives.Credentials) error
	Alg() primitives.CredentialsAlgType
}

// Argon2IDProducer implements the CredentialProducer interface.
//
// This producer is designed for production usage as it focuses on a
// memory-hard hashing algorithm which makes it harder to bruteforce in the
// case of a lost database.
type Argon2IDProducer struct {
	Time      uint32
	Memory    uint32
	Threads   uint8
	KeyLength uint32
}

func (a *Argon2IDProducer) Generate(secret primitives.Password) (*primitives.Credentials, error) {
	if err := secret.Validate(); err != nil {
		return nil, err
	}

	salt := make([]byte, primitives.SaltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	value := argon2.IDKey([]byte(secret), salt, a.Time, a.Memory, a.Threads, a.KeyLength)
	creds := &primitives.Credentials{
		Alg:    a.Alg(),
		Secret: base64.New(value),
		Salt:   base64.New(salt),
	}

	return creds, creds.Validate()
}

func (a *Argon2IDProducer) Compare(secret primitives.Password, creds *primitives.Credentials) error {
	if err := secret.Validate(); err != nil {
		return err
	}

	if err := creds.Validate(); err != nil {
		return err
	}

	if creds.Alg != a.Alg() {
		return errors.New(UnsupportedAlgorithm, "Algorithm %s is not supported, requires %s", creds.Alg, a.Alg())
	}

	value := argon2.IDKey([]byte(secret), []byte(*creds.Salt), a.Time, a.Memory, a.Threads, a.KeyLength)
	if subtle.ConstantTimeCompare(value, []byte(*creds.Secret)) == 0 {
		return ErrBadCredentials
	}

	return nil
}

func (a *Argon2IDProducer) Alg() primitives.CredentialsAlgType {
	return primitives.Argon2ID
}

// SHA256Producer implements the CredentialProducer interface. The
// SHA256Producer is designed for _fast_ hashing scenarios and thus should
// only ever be used in development situations
type SHA256Producer struct{}

func (s *SHA256Producer) Generate(secret primitives.Password) (*primitives.Credentials, error) {
	if err := secret.Validate(); err != nil {
		return nil, err
	}

	salt := make([]byte, primitives.SaltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	value := sha256.Sum256(append([]byte(secret), salt...))
	creds := &primitives.Credentials{
		Alg:    s.Alg(),
		Secret: base64.New(value[:]),
		Salt:   base64.New(salt),
	}

	return creds, creds.Validate()
}

func (s *SHA256Producer) Compare(secret primitives.Password, creds *primitives.Credentials) error {
	if err := secret.Validate(); err != nil {
		return err
	}

	if err := creds.Validate(); err != nil {
		return err
	}

	if creds.Alg != s.Alg() {
		return errors.New(UnsupportedAlgorithm, "Algorithm %s is not supported, requires %s", creds.Alg, s.Alg())
	}

	value := sha256.Sum256(append([]byte(secret), []byte(*creds.Salt)...))
	if subtle.ConstantTimeCompare(value[:], []byte(*creds.Secret)) == 0 {
		return ErrBadCredentials
	}

	return nil
}

func (s *SHA256Producer) Alg() primitives.CredentialsAlgType {
	return primitives.SHA256
}
