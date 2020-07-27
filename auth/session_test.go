package auth

import (
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/models"
	"github.com/capeprivacy/cape/primitives"
)

func TestCan(t *testing.T) {
	gm.RegisterTestingT(t)

	t.Run("GetID returns the user id", func(t *testing.T) {
		_, user := models.GenerateUser("hiho", "jerry@berry.jerry")

		session, err := NewSession(&user, &primitives.Session{}, models.UserRoles{}, &user)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(session.GetID()).To(gm.Equal(user.ID))
	})
}
