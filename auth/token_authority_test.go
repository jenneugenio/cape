package auth

import (
	"testing"
	"time"

	gm "github.com/onsi/gomega"
	"gopkg.in/square/go-jose.v2"

	"github.com/dropoutlabs/cape/primitives"
)

func TestTokenAuthority(t *testing.T) {
	gm.RegisterTestingT(t)

	keypair, err := NewKeypair()
	gm.Expect(err).To(gm.BeNil())

	t.Run("Generate Token", func(t *testing.T) {
		tokenAuth, err := NewTokenAuthority(keypair, "controller@controller.ai")
		gm.Expect(err).To(gm.BeNil())

		sig, expiresIn, err := tokenAuth.Generate(primitives.Login)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(sig).ToNot(gm.BeNil())
		gm.Expect(expiresIn).To(gm.BeTemporally(">", time.Now().UTC()))
	})

	t.Run("Can Verify Token", func(t *testing.T) {
		tokenAuth, err := NewTokenAuthority(keypair, "controller@controller.ai")
		gm.Expect(err).To(gm.BeNil())

		sig, expiresIn, err := tokenAuth.Generate(primitives.Login)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(sig).ToNot(gm.BeNil())
		gm.Expect(expiresIn).To(gm.BeTemporally(">", time.Now().UTC()))

		err = tokenAuth.Verify(sig)
		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("Can verify from another authority - same keypair", func(t *testing.T) {
		tokenAuth, err := NewTokenAuthority(keypair, "controller@controller.ai")
		gm.Expect(err).To(gm.BeNil())

		sig, expiresIn, err := tokenAuth.Generate(primitives.Login)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(sig).ToNot(gm.BeNil())
		gm.Expect(expiresIn).To(gm.BeTemporally(">", time.Now().UTC()))

		err = tokenAuth.Verify(sig)
		gm.Expect(err).To(gm.BeNil())

		other, err := NewTokenAuthority(keypair, "controller@controller.ai")
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(other.Verify(sig)).To(gm.BeNil())
	})

	t.Run("Can't verify from another authority with different keypair", func(t *testing.T) {
		tokenAuth, err := NewTokenAuthority(keypair, "controller@controller.ai")
		gm.Expect(err).To(gm.BeNil())

		otherKeypair, err := NewKeypair()
		gm.Expect(err).To(gm.BeNil())

		otherTokenAuth, err := NewTokenAuthority(otherKeypair, "controller2@controller.ai")
		gm.Expect(err).To(gm.BeNil())

		sig, expiresIn, err := tokenAuth.Generate(primitives.Login)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(sig).ToNot(gm.BeNil())
		gm.Expect(expiresIn).To(gm.BeTemporally(">", time.Now().UTC()))

		err = otherTokenAuth.Verify(sig)
		gm.Expect(err).To(gm.Equal(jose.ErrCryptoFailure))
	})
}
