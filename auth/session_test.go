package auth

import (
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/models"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

func TestCan(t *testing.T) {
	gm.RegisterTestingT(t)

	email := models.Email("jerry@jerry.berry")

	password := primitives.GeneratePassword()

	creds, err := DefaultSHA256Producer.Generate(password)
	gm.Expect(err).To(gm.BeNil())

	t.Run("GetID returns the user id", func(t *testing.T) {
		_, user := models.GenerateUser("hiho", "jerry@berry.jerry")

		session, err := NewSession(&user, &primitives.Session{}, []*models.Policy{}, []*primitives.Role{}, &user)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(session.GetID()).To(gm.Equal(user.ID))
	})

	t.Run("denied no rules", func(t *testing.T) {
		user := models.NewUser("Jerry Berry", email, creds)

		session, err := NewSession(&user, &primitives.Session{}, []*models.Policy{}, []*primitives.Role{}, &user)
		gm.Expect(err).To(gm.BeNil())

		err = session.Can(models.Create, primitives.UserType)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(errors.CausedBy(err, AuthorizationFailure)).To(gm.BeTrue())
	})

	t.Run("denied deny rule exists", func(t *testing.T) {
		user := models.NewUser("Jerry Berry", email, creds)

		spec := &models.PolicySpec{
			Version: 1,
			Label:   "my-policy",
			Rules: []*models.Rule{
				{
					Target: "users:*",
					Action: models.Create,
					Effect: models.Deny,
				},
			},
		}

		p := models.NewPolicy("my-policy", spec)

		session, err := NewSession(&user, &primitives.Session{}, []*models.Policy{&p}, []*primitives.Role{}, &user)
		gm.Expect(err).To(gm.BeNil())

		err = session.Can(models.Create, primitives.UserType)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(errors.CausedBy(err, AuthorizationFailure)).To(gm.BeTrue())
	})

	t.Run("allowed rules", func(t *testing.T) {
		user := models.NewUser("Jerry Berry", email, creds)
		gm.Expect(err).To(gm.BeNil())

		spec := &models.PolicySpec{
			Version: 1,
			Label:   "my-policy",
			Rules: []*models.Rule{
				{
					Target: "users:*",
					Action: models.Create,
					Effect: models.Allow,
				},
			},
		}

		p := models.NewPolicy("my-policy", spec)

		session, err := NewSession(&user, &primitives.Session{}, []*models.Policy{&p}, []*primitives.Role{}, &user)
		gm.Expect(err).To(gm.BeNil())

		err = session.Can(models.Create, primitives.UserType)
		gm.Expect(err).To(gm.BeNil())
	})
}
