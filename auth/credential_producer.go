package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"

	"github.com/manifoldco/go-base64"
	"golang.org/x/crypto/argon2"

	"github.com/capeprivacy/cape/models"
	errors "github.com/capeprivacy/cape/partyerrors"
)

var (
	DefaultSHA256Producer   = &SHA256Producer{}
	DefaultArgon2IDProducer = &Argon2IDProducer{
		Time:      1,
		Memory:    64 * 1024,
		Threads:   4,
		KeyLength: models.SecretLength,
	}
)

// CredentialProducer represents an interface for generating credentials and
// comparing credentials based on a pre-shared key.
type CredentialProducer interface {
	Generate(models.Password) (*models.Credentials, error)
	Compare(models.Password, *models.Credentials) error
	Alg() models.CredentialsAlgType
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

func (a *Argon2IDProducer) Generate(secret models.Password) (*models.Credentials, error) {
	if err := secret.Validate(); err != nil {
		return nil, err
	}

	salt := make([]byte, models.SaltLength)
	_, err := randRead(salt)
	if err != nil {
		return nil, errors.Wrap(models.SystemErrorCause, err)
	}

	value := argon2.IDKey([]byte(secret), salt, a.Time, a.Memory, a.Threads, a.KeyLength)
	creds := &models.Credentials{
		Alg:    a.Alg(),
		Secret: base64.New(value),
		Salt:   base64.New(salt),
	}

	return creds, nil
}

func (a *Argon2IDProducer) Compare(secret models.Password, creds *models.Credentials) error {
	if err := secret.Validate(); err != nil {
		return err
	}

	if creds.Alg != a.Alg() {
		return errors.New(UnsupportedAlgorithm, "Algorithm %s is not supported, requires %s", creds.Alg, a.Alg())
	}

	value := argon2.IDKey([]byte(secret), *creds.Salt, a.Time, a.Memory, a.Threads, a.KeyLength)
	if subtle.ConstantTimeCompare(value, *creds.Secret) == 0 {
		return ErrBadCredentials
	}

	return nil
}

func (a *Argon2IDProducer) Alg() models.CredentialsAlgType {
	return models.Argon2ID
}

// SHA256Producer implements the CredentialProducer interface. The
// SHA256Producer is designed for _fast_ hashing scenarios and thus should
// only ever be used in development situations
type SHA256Producer struct{}

func (s *SHA256Producer) Generate(secret models.Password) (*models.Credentials, error) {
	if err := secret.Validate(); err != nil {
		return nil, err
	}

	salt := make([]byte, models.SaltLength)
	_, err := randRead(salt)
	if err != nil {
		return nil, errors.Wrap(models.SystemErrorCause, err)
	}

	value := sha256.Sum256(append([]byte(secret), salt...))
	creds := &models.Credentials{
		Alg:    s.Alg(),
		Secret: base64.New(value[:]),
		Salt:   base64.New(salt),
	}

	return creds, nil
}

func (s *SHA256Producer) Compare(secret models.Password, creds *models.Credentials) error {
	if err := secret.Validate(); err != nil {
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

func (s *SHA256Producer) Alg() models.CredentialsAlgType {
	return models.SHA256
}

var randRead = rand.Read
