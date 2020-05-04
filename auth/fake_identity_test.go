package auth

import (
	"testing"

	"github.com/capeprivacy/cape/primitives"
	gm "github.com/onsi/gomega"
)

func TestNewFakeIdentity(t *testing.T) {
	t.Run("test new fake identity", func(t *testing.T) {
		gm.RegisterTestingT(t)

		email, err := primitives.NewEmail("fake@fake.com")
		gm.Expect(err).To(gm.BeNil())

		identity, err := NewFakeIdentity(email)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(identity.GetCredentials().Salt).NotTo(gm.BeNil())
	})

	t.Run("test fake identity returns same data for email", func(t *testing.T) {
		gm.RegisterTestingT(t)

		email, err := primitives.NewEmail("fake@fake.com")
		gm.Expect(err).To(gm.BeNil())

		identity, err := NewFakeIdentity(email)
		gm.Expect(err).To(gm.BeNil())

		otherIdentity, err := NewFakeIdentity(email)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(identity.GetCredentials().Salt).To(gm.Equal(otherIdentity.GetCredentials().Salt))
	})
}
