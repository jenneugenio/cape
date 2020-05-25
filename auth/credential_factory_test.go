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

	t.Run("can create credential factory", func(t *testing.T) {
		cf, err := NewCredentialFactory(primitives.SHA256)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(cf.Alg).To(gm.Equal(primitives.SHA256))
	})

	t.Run("returns error on c-tor if unsupported alg", func(t *testing.T) {
		_, err := NewCredentialFactory(primitives.EDDSA)
		gm.Expect(errors.CausedBy(err, ProducerNotFound)).To(gm.BeTrue())
	})

	t.Run("can generate credentials", func(t *testing.T) {
		tests := []struct {
			name   string
			alg    primitives.CredentialsAlgType
			secret primitives.Password
			cause  *errors.Cause
		}{
			{
				name:   "can create with valid password - sha256",
				alg:    primitives.SHA256,
				secret: primitives.Password("helloabcdefgh"),
			},
			{
				name:   "can create with valid password - argon2id",
				alg:    primitives.Argon2ID,
				secret: primitives.Password("helloabcdefgh"),
			},
			{
				name:   "error for invalid password",
				alg:    primitives.SHA256,
				secret: primitives.Password("sdf"),
				cause:  &primitives.InvalidPasswordCause,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				cf, err := NewCredentialFactory(tc.alg)
				gm.Expect(err).To(gm.BeNil())

				creds, err := cf.Generate(tc.secret)
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
		tests := []struct {
			name       string
			alg        primitives.CredentialsAlgType
			initial    primitives.Password
			comparison primitives.Password
			cause      *errors.Cause
		}{
			{
				name:       "matches properly - sha256",
				alg:        primitives.SHA256,
				initial:    primitives.Password("abcdefghijk"),
				comparison: primitives.Password("abcdefghijk"),
			},
			{
				name:       "errors if incorrect - sha256",
				alg:        primitives.SHA256,
				initial:    primitives.Password("sfsdfsfsfsdf"),
				comparison: primitives.Password("sdfsfsfsf231"),
				cause:      &MismatchingCredentials,
			},
			{
				name:       "matches properly - Argon2ID",
				alg:        primitives.Argon2ID,
				initial:    primitives.Password("abcdefghijk"),
				comparison: primitives.Password("abcdefghijk"),
			},
			{
				name:       "errors if incorrect - Argon2ID",
				alg:        primitives.Argon2ID,
				initial:    primitives.Password("sfsdfsfsfsdf"),
				comparison: primitives.Password("sdfsfsfsf231"),
				cause:      &MismatchingCredentials,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				cf, err := NewCredentialFactory(tc.alg)
				gm.Expect(err).To(gm.BeNil())

				creds, err := cf.Generate(tc.initial)
				gm.Expect(err).To(gm.BeNil())

				err = cf.Compare(tc.comparison, creds)
				if tc.cause != nil {
					gm.Expect(errors.FromCause(err, *tc.cause)).To(gm.BeTrue())
					return
				}

				gm.Expect(err).To(gm.BeNil())
			})
		}
	})

	t.Run("can be backwards compatible", func(t *testing.T) {
		sha256, err := NewCredentialFactory(primitives.SHA256)
		gm.Expect(err).To(gm.BeNil())

		argon2id, err := NewCredentialFactory(primitives.Argon2ID)
		gm.Expect(err).To(gm.BeNil())

		pw, err := primitives.GeneratePassword()
		gm.Expect(err).To(gm.BeNil())

		sha256Creds, err := sha256.Generate(pw)
		gm.Expect(err).To(gm.BeNil())

		argon2idCreds, err := argon2id.Generate(pw)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(sha256.Compare(pw, argon2idCreds)).To(gm.BeNil())
		gm.Expect(argon2id.Compare(pw, sha256Creds)).To(gm.BeNil())
	})

	t.Run("cannot compare unsupported credentials", func(t *testing.T) {
		cf, err := NewCredentialFactory(primitives.SHA256)
		gm.Expect(err).To(gm.BeNil())

		creds, err := cf.Generate("sfsfsfasfasffsf")
		gm.Expect(err).To(gm.BeNil())

		creds.Alg = primitives.EDDSA

		err = cf.Compare("sfsfsfasfasffsf", creds)
		gm.Expect(errors.FromCause(err, ProducerNotFound)).To(gm.BeTrue())
	})

	t.Run("cannot compare malignant credentials", func(t *testing.T) {
		cf, err := NewCredentialFactory(primitives.SHA256)
		gm.Expect(err).To(gm.BeNil())

		creds, err := cf.Generate("sfsfsfasfasffsf")
		gm.Expect(err).To(gm.BeNil())

		creds.Secret = base64.New([]byte("sdfasfasfs"))

		err = cf.Compare("sfsfsfasfasffsf", creds)
		gm.Expect(errors.FromCause(err, primitives.InvalidCredentialsCause)).To(gm.BeTrue())
	})
}
