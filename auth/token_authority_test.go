package auth

import (
	"testing"
	"time"

	gm "github.com/onsi/gomega"
	"gopkg.in/square/go-jose.v2"

	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/primitives"
)

func TestTokenAuthority(t *testing.T) {
	gm.RegisterTestingT(t)

	keypair, err := NewKeypair()
	gm.Expect(err).To(gm.BeNil())

	t.Run("Generate Token", func(t *testing.T) {
		tokenAuth, err := NewTokenAuthority(keypair, "coordinator@coordinator.ai")
		gm.Expect(err).To(gm.BeNil())

		id, err := database.GenerateID(primitives.SessionType)
		gm.Expect(err).To(gm.BeNil())

		sig, expiresIn, err := tokenAuth.Generate(primitives.Login, id)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(sig).ToNot(gm.BeNil())
		gm.Expect(expiresIn).To(gm.BeTemporally(">", time.Now().UTC()))
	})

	t.Run("Can Verify Token", func(t *testing.T) {
		tokenAuth, err := NewTokenAuthority(keypair, "coordinator@coordinator.ai")
		gm.Expect(err).To(gm.BeNil())

		id, err := database.GenerateID(primitives.SessionType)
		gm.Expect(err).To(gm.BeNil())

		sig, expiresIn, err := tokenAuth.Generate(primitives.Login, id)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(sig).ToNot(gm.BeNil())
		gm.Expect(expiresIn).To(gm.BeTemporally(">", time.Now().UTC()))

		otherID, err := tokenAuth.Verify(sig)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(otherID).To(gm.Equal(id))
	})

	t.Run("Can verify from another authority - same keypair", func(t *testing.T) {
		tokenAuth, err := NewTokenAuthority(keypair, "coordinator@coordinator.ai")
		gm.Expect(err).To(gm.BeNil())

		id, err := database.GenerateID(primitives.SessionType)
		gm.Expect(err).To(gm.BeNil())

		sig, expiresIn, err := tokenAuth.Generate(primitives.Login, id)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(sig).ToNot(gm.BeNil())
		gm.Expect(expiresIn).To(gm.BeTemporally(">", time.Now().UTC()))

		otherID, err := tokenAuth.Verify(sig)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(otherID).To(gm.Equal(id))

		other, err := NewTokenAuthority(keypair, "coordinator@coordinator.ai")
		gm.Expect(err).To(gm.BeNil())

		otherOtherID, err := other.Verify(sig)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(otherOtherID).To(gm.Equal(otherID))
	})

	t.Run("Can't verify from another authority with different keypair", func(t *testing.T) {
		tokenAuth, err := NewTokenAuthority(keypair, "coordinator@coordinator.ai")
		gm.Expect(err).To(gm.BeNil())

		otherKeypair, err := NewKeypair()
		gm.Expect(err).To(gm.BeNil())

		otherTokenAuth, err := NewTokenAuthority(otherKeypair, "coordinator2@coordinator.ai")
		gm.Expect(err).To(gm.BeNil())

		id, err := database.GenerateID(primitives.SessionType)
		gm.Expect(err).To(gm.BeNil())

		sig, expiresIn, err := tokenAuth.Generate(primitives.Login, id)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(sig).ToNot(gm.BeNil())
		gm.Expect(expiresIn).To(gm.BeTemporally(">", time.Now().UTC()))

		_, err = otherTokenAuth.Verify(sig)
		gm.Expect(err).To(gm.Equal(jose.ErrCryptoFailure))
	})
}
