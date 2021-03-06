package auth

import (
	"fmt"
	"testing"

	"github.com/manifoldco/go-base64"
	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/models"
	errors "github.com/capeprivacy/cape/partyerrors"
)

func TestCredentialProducer(t *testing.T) {
	gm.RegisterTestingT(t)

	t.Run("can generate credentials", func(t *testing.T) {
		tests := []struct {
			name     string
			alg      models.CredentialsAlgType
			producer CredentialProducer
			secret   models.Password
			randRead func([]byte) (int, error)
			cause    *errors.Cause
		}{
			{
				name:     "can create with valid password - sha256",
				alg:      models.SHA256,
				producer: DefaultSHA256Producer,
				secret:   models.Password("helloabcdefgh"),
			},
			{
				name:     "can create with valid password - argon2id",
				producer: DefaultArgon2IDProducer,
				alg:      models.Argon2ID,
				secret:   models.Password("helloabcdefgh"),
			},
			{
				name:     "error for invalid password - sha256",
				alg:      models.SHA256,
				producer: DefaultSHA256Producer,
				secret:   models.Password("sdf"),
				cause:    &models.InvalidPasswordCause,
			},
			{
				name:     "error for invalid password - argon2id",
				alg:      models.Argon2ID,
				producer: DefaultArgon2IDProducer,
				secret:   models.Password("sdf"),
				cause:    &models.InvalidPasswordCause,
			},
			{
				name:     "errors on rand.Read error - sha256",
				alg:      models.SHA256,
				producer: DefaultSHA256Producer,
				secret:   models.Password("helloabcdefgh"),
				randRead: errRand,
				cause:    &models.SystemErrorCause,
			},
		}

		origRandRead := randRead
		defer func() { randRead = origRandRead }()

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				if tc.randRead == nil {
					randRead = origRandRead
				} else {
					randRead = tc.randRead
				}
				creds, err := tc.producer.Generate(tc.secret)
				if tc.cause != nil {
					gm.Expect(errors.FromCause(err, *tc.cause)).To(gm.BeTrue())
					return
				}

				gm.Expect(err).To(gm.BeNil())
				gm.Expect(creds.Alg).To(gm.Equal(tc.alg))
			})
		}
	})

	t.Run("can compare credentials", func(t *testing.T) {
		matching := models.GeneratePassword()

		tests := []struct {
			name       string
			producer   CredentialProducer
			alg        models.CredentialsAlgType
			initial    models.Password
			comparison models.Password
			cause      *errors.Cause
			genCause   *errors.Cause
			secret     *base64.Value
			randRead   func([]byte) (int, error)
		}{
			{
				name:       "matches properly - sha256",
				producer:   DefaultSHA256Producer,
				alg:        models.SHA256,
				initial:    models.Password("abcdefghijk"),
				comparison: models.Password("abcdefghijk"),
			},
			{
				name:       "errors if incorrect - sha256",
				producer:   DefaultSHA256Producer,
				alg:        models.SHA256,
				initial:    models.Password("sfsdfsfsfsdf"),
				comparison: models.Password("sdfsfsfsf231"),
				cause:      &MismatchingCredentials,
			},
			{
				name:       "errors if wrong alg - sha256",
				producer:   DefaultSHA256Producer,
				alg:        models.Argon2ID,
				initial:    matching,
				comparison: matching,
				cause:      &UnsupportedAlgorithm,
			},
			{
				name:       "errors if wrong alg - argon2id",
				producer:   DefaultArgon2IDProducer,
				alg:        models.SHA256,
				initial:    matching,
				comparison: matching,
				cause:      &UnsupportedAlgorithm,
			},
			{
				name:       "matches properly - Argon2ID",
				producer:   DefaultArgon2IDProducer,
				alg:        models.Argon2ID,
				initial:    models.Password("abcdefghijk"),
				comparison: models.Password("abcdefghijk"),
			},
			{
				name:       "errors if incorrect - Argon2ID",
				producer:   DefaultArgon2IDProducer,
				alg:        models.Argon2ID,
				initial:    models.Password("sfsdfsfsfsdf"),
				comparison: models.Password("sdfsfsfsf231"),
				cause:      &MismatchingCredentials,
			},
			{
				name:       "error for invalid password - sha256",
				alg:        models.SHA256,
				producer:   DefaultSHA256Producer,
				initial:    models.Password("sdf"),
				comparison: matching,
				genCause:   &models.InvalidPasswordCause,
			},
			{
				name:       "error for invalid password - argon2id",
				alg:        models.Argon2ID,
				producer:   DefaultArgon2IDProducer,
				initial:    models.Password("sdf"),
				comparison: matching,
				genCause:   &models.InvalidPasswordCause,
			},
			{
				name:       "matches properly - sha256",
				producer:   DefaultSHA256Producer,
				alg:        models.SHA256,
				initial:    models.Password("abcdefghijk"),
				comparison: models.Password("abcdefghijk"),
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				creds, err := tc.producer.Generate(tc.initial)
				if tc.genCause != nil {
					gm.Expect(errors.FromCause(err, *tc.genCause)).To(gm.BeTrue())
					return
				}
				gm.Expect(err).To(gm.BeNil())

				creds.Alg = tc.alg
				if tc.secret != nil {
					creds.Secret = tc.secret
				}

				err = tc.producer.Compare(tc.comparison, creds)
				if tc.cause != nil {
					gm.Expect(errors.FromCause(err, *tc.cause)).To(gm.BeTrue())
					return
				}

				gm.Expect(err).To(gm.BeNil())
			})
		}
	})
}

type ErrRandReader struct {
	i   int
	err error
}

func (e ErrRandReader) Read(_ []byte) (int, error) { return e.i, e.err }

var errRand = ErrRandReader{0, fmt.Errorf("bad rand.Read result")}.Read
