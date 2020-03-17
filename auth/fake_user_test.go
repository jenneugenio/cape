package auth

import (
	"testing"

	gm "github.com/onsi/gomega"
)

func TestNewFakeUser(t *testing.T) {
	t.Run("test new fake user ", func(t *testing.T) {
		gm.RegisterTestingT(t)

		user, err := NewFakeUser("fake@fake.com")
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(user.Credentials.Salt).NotTo(gm.BeNil())
	})

	t.Run("test fake user returns same data for email ", func(t *testing.T) {
		gm.RegisterTestingT(t)

		user, err := NewFakeUser("fake@fake.com")
		gm.Expect(err).To(gm.BeNil())

		otherUser, err := NewFakeUser("fake@fake.com")
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(user.Credentials.Salt).To(gm.Equal(otherUser.Credentials.Salt))
	})
}
