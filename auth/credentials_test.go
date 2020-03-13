package auth

import (
	"crypto/ed25519"
	"testing"

	"github.com/manifoldco/go-base64"
	gm "github.com/onsi/gomega"

	errors "github.com/dropoutlabs/cape/partyerrors"
)

func TestNewCredential(t *testing.T) {
	gm.RegisterTestingT(t)

	creds, err := NewCredentials([]byte("my-cool-secret"), nil)
	gm.Expect(err).To(gm.BeNil())
	gm.Expect(len(*creds.Salt)).To(gm.Equal(SaltLength))
	gm.Expect(len(*creds.PublicKey)).To(gm.Equal(ed25519.PublicKeySize))
	gm.Expect(len(creds.privateKey)).To(gm.Equal(ed25519.PrivateKeySize))
	gm.Expect(creds.Alg).To(gm.Equal(EDDSA))
}

func TestSignVerifyChallenge(t *testing.T) {
	gm.RegisterTestingT(t)

	creds, err := NewCredentials([]byte("my-cool-secret"), nil)
	gm.Expect(err).To(gm.BeNil())

	msg := base64.New([]byte("my-awesome-msg"))
	sig, err := creds.Sign(msg)
	gm.Expect(err).To(gm.BeNil())
	gm.Expect(len(*sig)).To(gm.Equal(64))

	err = creds.Verify(msg, sig)
	gm.Expect(err).To(gm.BeNil())
}

func TestRederiveCredentials(t *testing.T) {
	gm.RegisterTestingT(t)

	creds, err := NewCredentials([]byte("my-cool-secret"), nil)
	gm.Expect(err).To(gm.BeNil())

	sameCreds, err := NewCredentials([]byte("my-cool-secret"), creds.Salt)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(sameCreds).To(gm.Equal(creds))
}

func TestLoadCredentials(t *testing.T) {
	gm.RegisterTestingT(t)

	creds, err := NewCredentials([]byte("my-cool-secret"), nil)
	gm.Expect(err).To(gm.BeNil())

	msg := base64.New([]byte("my-cool-msg"))
	sig, err := creds.Sign(msg)
	gm.Expect(err).To(gm.BeNil())

	loadedCreds, err := LoadCredentials(creds.PublicKey, creds.Salt)
	gm.Expect(err).To(gm.BeNil())

	err = loadedCreds.Verify(msg, sig)
	gm.Expect(err).To(gm.BeNil())
}

func TestCredentialErrors(t *testing.T) {
	t.Run("verified by wrong public key", func(t *testing.T) {
		gm.RegisterTestingT(t)

		creds, err := NewCredentials([]byte("my-cool-secret"), nil)
		gm.Expect(err).To(gm.BeNil())

		otherCreds, err := NewCredentials([]byte("my-cool-secret2"), nil)
		gm.Expect(err).To(gm.BeNil())

		msg := base64.New([]byte("my-cool-msg"))
		sig, err := creds.Sign(msg)
		gm.Expect(err).To(gm.BeNil())

		err = otherCreds.Verify(msg, sig)
		gm.Expect(err).ToNot(gm.BeNil())

		gm.Expect(errors.FromCause(err, SignatureNotValid)).To(gm.BeTrue())
	})

	t.Run("required private key sign", func(t *testing.T) {
		gm.RegisterTestingT(t)

		creds, err := NewCredentials([]byte("my-cool-secret"), nil)
		gm.Expect(err).To(gm.BeNil())

		loadedCreds, err := LoadCredentials(creds.PublicKey, creds.Salt)
		gm.Expect(err).To(gm.BeNil())

		msg := base64.New([]byte("my-cool-msg"))
		sig, err := loadedCreds.Sign(msg)
		gm.Expect(sig).To(gm.BeNil())
		gm.Expect(err).ToNot(gm.BeNil())

		gm.Expect(errors.FromCause(err, RequiredPrivateKeyCause)).To(gm.BeTrue())
	})

	t.Run("bad secret length", func(t *testing.T) {
		gm.RegisterTestingT(t)

		creds, err := NewCredentials([]byte("124567"), nil)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(creds).To(gm.BeNil())

		gm.Expect(errors.FromCause(err, BadSecretLength)).To(gm.BeTrue())
	})

	t.Run("bad salt length", func(t *testing.T) {
		gm.RegisterTestingT(t)

		creds, err := NewCredentials([]byte("my-cool-password"), base64.New([]byte("2")))
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(creds).To(gm.BeNil())

		gm.Expect(errors.FromCause(err, BadSaltLength)).To(gm.BeTrue())
	})
}
