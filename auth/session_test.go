package auth

import (
	"testing"

	gm "github.com/onsi/gomega"

	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

func TestCan(t *testing.T) {
	gm.RegisterTestingT(t)

	email, err := primitives.NewEmail("jerry@jerry.berry")
	gm.Expect(err).To(gm.BeNil())

	password, err := primitives.GeneratePassword()
	gm.Expect(err).To(gm.BeNil())

	creds, err := DefaultSHA256Producer.Generate(password)
	gm.Expect(err).To(gm.BeNil())

	t.Run("GetID returns the identity id", func(t *testing.T) {
		_, user, err := primitives.GenerateUser("hiho", "jerry@berry.jerry")
		gm.Expect(err).To(gm.BeNil())

		session, err := NewSession(user, &primitives.Session{}, []*primitives.Policy{}, []*primitives.Role{}, user)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(session.GetID()).To(gm.Equal(user.ID))
	})

	t.Run("denied no rules", func(t *testing.T) {
		user, err := primitives.NewUser("Jerry Berry", email, creds)
		gm.Expect(err).To(gm.BeNil())

		session, err := NewSession(user, &primitives.Session{}, []*primitives.Policy{}, []*primitives.Role{}, user)
		gm.Expect(err).To(gm.BeNil())

		err = session.Can(primitives.Create, primitives.UserType)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(errors.CausedBy(err, AuthorizationFailure)).To(gm.BeTrue())
	})

	t.Run("denied deny rule exists", func(t *testing.T) {
		user, err := primitives.NewUser("Jerry Berry", email, creds)
		gm.Expect(err).To(gm.BeNil())

		spec := &primitives.PolicySpec{
			Version: 1,
			Label:   "my-policy",
			Rules: []*primitives.Rule{
				{
					Target: "users:*",
					Action: primitives.Create,
					Effect: primitives.Deny,
				},
			},
		}

		p, err := primitives.NewPolicy("my-policy", spec)
		gm.Expect(err).To(gm.BeNil())

		session, err := NewSession(user, &primitives.Session{}, []*primitives.Policy{p}, []*primitives.Role{}, user)
		gm.Expect(err).To(gm.BeNil())

		err = session.Can(primitives.Create, primitives.UserType)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(errors.CausedBy(err, AuthorizationFailure)).To(gm.BeTrue())
	})

	t.Run("allowed rules", func(t *testing.T) {
		user, err := primitives.NewUser("Jerry Berry", email, creds)
		gm.Expect(err).To(gm.BeNil())

		spec := &primitives.PolicySpec{
			Version: 1,
			Label:   "my-policy",
			Rules: []*primitives.Rule{
				{
					Target: "users:*",
					Action: primitives.Create,
					Effect: primitives.Allow,
				},
			},
		}

		p, err := primitives.NewPolicy("my-policy", spec)
		gm.Expect(err).To(gm.BeNil())

		session, err := NewSession(user, &primitives.Session{}, []*primitives.Policy{p}, []*primitives.Role{}, user)
		gm.Expect(err).To(gm.BeNil())

		err = session.Can(primitives.Create, primitives.UserType)
		gm.Expect(err).To(gm.BeNil())
	})
}
