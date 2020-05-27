package auth

import (
	"testing"

	"github.com/manifoldco/go-base64"
	gm "github.com/onsi/gomega"

	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

func TestCredentialFactory(t *testing.T) {
	gm.RegisterTestingT(t)

	t.Run("can generate credentials", func(t *testing.T) {
		tests := []struct {
			name     string
			alg      primitives.CredentialsAlgType
			producer CredentialProducer
			secret   primitives.Password
			cause    *errors.Cause
		}{
			{
				name:     "can create with valid password - sha256",
				alg:      primitives.SHA256,
				producer: DefaultSHA256Producer,
				secret:   primitives.Password("helloabcdefgh"),
			},
			{
				name:     "can create with valid password - argon2id",
				producer: DefaultArgon2IDProducer,
				alg:      primitives.Argon2ID,
				secret:   primitives.Password("helloabcdefgh"),
			},
			{
				name:     "error for invalid password - sha256",
				alg:      primitives.SHA256,
				producer: DefaultSHA256Producer,
				secret:   primitives.Password("sdf"),
				cause:    &primitives.InvalidPasswordCause,
			},
			{
				name:     "error for invalid password - argon2id",
				alg:      primitives.Argon2ID,
				producer: DefaultArgon2IDProducer,
				secret:   primitives.Password("sdf"),
				cause:    &primitives.InvalidPasswordCause,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				creds, err := tc.producer.Generate(tc.secret)
				if tc.cause != nil {
					gm.Expect(errors.FromCause(err, *tc.cause)).To(gm.BeTrue())
					return
				}

				gm.Expect(creds.Alg).To(gm.Equal(tc.alg))
				gm.Expect(creds.Validate()).To(gm.BeNil())
			})
		}
	})

	t.Run("can compare credentials", func(t *testing.T) {
		matching, err := primitives.GeneratePassword()
		gm.Expect(err).To(gm.BeNil())

		tests := []struct {
			name       string
			producer   CredentialProducer
			alg        primitives.CredentialsAlgType
			initial    primitives.Password
			comparison primitives.Password
			cause      *errors.Cause
			genCause   *errors.Cause
			secret     *base64.Value
		}{
			{
				name:       "matches properly - sha256",
				producer:   DefaultSHA256Producer,
				alg:        primitives.SHA256,
				initial:    primitives.Password("abcdefghijk"),
				comparison: primitives.Password("abcdefghijk"),
			},
			{
				name:       "errors if incorrect - sha256",
				producer:   DefaultSHA256Producer,
				alg:        primitives.SHA256,
				initial:    primitives.Password("sfsdfsfsfsdf"),
				comparison: primitives.Password("sdfsfsfsf231"),
				cause:      &MismatchingCredentials,
			},
			{
				name:       "errors if wrong alg - sha256",
				producer:   DefaultSHA256Producer,
				alg:        primitives.Argon2ID,
				initial:    matching,
				comparison: matching,
				cause:      &UnsupportedAlgorithm,
			},
			{
				name:       "errors if wrong alg - argon2id",
				producer:   DefaultArgon2IDProducer,
				alg:        primitives.SHA256,
				initial:    matching,
				comparison: matching,
				cause:      &UnsupportedAlgorithm,
			},
			{
				name:       "matches properly - Argon2ID",
				producer:   DefaultArgon2IDProducer,
				alg:        primitives.Argon2ID,
				initial:    primitives.Password("abcdefghijk"),
				comparison: primitives.Password("abcdefghijk"),
			},
			{
				name:       "errors if incorrect - Argon2ID",
				producer:   DefaultArgon2IDProducer,
				alg:        primitives.Argon2ID,
				initial:    primitives.Password("sfsdfsfsfsdf"),
				comparison: primitives.Password("sdfsfsfsf231"),
				cause:      &MismatchingCredentials,
			},
			{
				name:       "cannot compare malignant creds - sha256",
				producer:   DefaultSHA256Producer,
				alg:        primitives.SHA256,
				secret:     base64.New([]byte("hi")),
				initial:    matching,
				comparison: matching,
				cause:      &primitives.InvalidCredentialsCause,
			},
			{
				name:       "cannot compare malignant creds - argon2id",
				producer:   DefaultArgon2IDProducer,
				alg:        primitives.Argon2ID,
				secret:     base64.New([]byte("hi")),
				initial:    matching,
				comparison: matching,
				cause:      &primitives.InvalidCredentialsCause,
			},
			{
				name:       "error for invalid password - sha256",
				alg:        primitives.SHA256,
				producer:   DefaultSHA256Producer,
				initial:    primitives.Password("sdf"),
				comparison: matching,
				genCause:   &primitives.InvalidPasswordCause,
			},
			{
				name:       "error for invalid password - argon2id",
				alg:        primitives.Argon2ID,
				producer:   DefaultArgon2IDProducer,
				initial:    primitives.Password("sdf"),
				comparison: matching,
				genCause:   &primitives.InvalidPasswordCause,
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
