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

// CredentialProducer represents an interface for generating credentials and
// comparing credentials based on a pre-shared key.
type CredentialProducer interface {
	Generate(primitives.Password) (*primitives.Credentials, error)
	Compare(primitives.Password, *primitives.Credentials) error
}

type ProducerRegistry map[primitives.CredentialsAlgType]CredentialProducer

// Get returns the CredentialProducer for the given AlgType. If a producer
// doesn't exist for the algtype an error is returned.
func (pr ProducerRegistry) Get(alg primitives.CredentialsAlgType) (CredentialProducer, error) {
	cp, ok := pr[alg]
	if !ok {
		return nil, errors.New(ProducerNotFound, "Could not find producer: %s", alg.String())
	}

	return cp, nil
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
	salt := make([]byte, primitives.SaltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	value := argon2.IDKey([]byte(secret), salt, a.Time, a.Memory, a.Threads, a.KeyLength)
	creds := &primitives.Credentials{
		Alg:    primitives.Argon2ID,
		Secret: base64.New(value),
		Salt:   base64.New(salt),
	}

	return creds, creds.Validate()
}

func (a *Argon2IDProducer) Compare(secret primitives.Password, creds *primitives.Credentials) error {
	if err := creds.Validate(); err != nil {
		return err
	}

	value := argon2.IDKey([]byte(secret), []byte(*creds.Salt), a.Time, a.Memory, a.Threads, a.KeyLength)
	if subtle.ConstantTimeCompare(value, []byte(*creds.Secret)) == 0 {
		return ErrBadCredentials
	}

	return nil
}

// SHA256Producer implements the CredentialProducer interface. The
// SHA256Producer is designed for _fast_ hashing scenarios and thus should
// only ever be used in development situations
type SHA256Producer struct{}

func (*SHA256Producer) Generate(secret primitives.Password) (*primitives.Credentials, error) {
	salt := make([]byte, primitives.SaltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	value := sha256.Sum256(append([]byte(secret), salt...))
	creds := &primitives.Credentials{
		Alg:    primitives.SHA256,
		Secret: base64.New(value[:]),
		Salt:   base64.New(salt),
	}

	return creds, creds.Validate()
}

func (*SHA256Producer) Compare(secret primitives.Password, creds *primitives.Credentials) error {
	if err := creds.Validate(); err != nil {
		return err
	}

	value := sha256.Sum256(append([]byte(secret), []byte(*creds.Salt)...))
	if subtle.ConstantTimeCompare(value[:], []byte(*creds.Secret)) == 0 {
		return ErrBadCredentials
	}

	return nil
}
