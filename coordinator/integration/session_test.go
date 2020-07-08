// +build integration

package integration

import (
	"context"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/coordinator/harness"
	"github.com/capeprivacy/cape/primitives"
)

func TestSessions(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()
	cfg, err := harness.NewConfig()
	gm.Expect(err).To(gm.BeNil())

	h, err := harness.NewHarness(cfg)
	gm.Expect(err)

	err = h.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer h.Teardown(ctx) // nolint: errcheck

	m := h.Manager()
	_, err = m.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	t.Run("test client login", func(t *testing.T) {
		gm.RegisterTestingT(t)

		client, err := h.Client()
		gm.Expect(err).To(gm.BeNil())

		session, err := client.EmailLogin(ctx, m.Admin.User.Email, m.Admin.Password)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(session.UserID).To(gm.Equal(m.Admin.User.ID))
		gm.Expect(session.Token).ToNot(gm.BeNil())
	})

	t.Run("test fake user fails", func(t *testing.T) {
		gm.RegisterTestingT(t)

		client, err := h.Client()
		gm.Expect(err).To(gm.BeNil())

		email, err := primitives.NewEmail("fake@fake.com")
		gm.Expect(err).To(gm.BeNil())

		password, err := primitives.NewPassword("newpasswordwhodis")
		gm.Expect(err).To(gm.BeNil())

		session, err := client.EmailLogin(ctx, email, password)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(session).To(gm.BeNil())

		gm.Expect(err.Error()).To(gm.Equal("authentication_failure: Failed to authenticate"))
	})

	t.Run("test incorrect credentials", func(t *testing.T) {
		gm.RegisterTestingT(t)

		client, err := h.Client()
		gm.Expect(err).To(gm.BeNil())

		password, err := primitives.NewPassword("idontknowmypassword")
		gm.Expect(err).To(gm.BeNil())

		// fail because credentials inside login won't be right
		session, err := client.EmailLogin(ctx, m.Admin.User.Email, password)
		gm.Expect(session).To(gm.BeNil())

		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(err.Error()).To(gm.Equal("authentication_failure: Failed to authenticate"))
	})

	t.Run("test delete session", func(t *testing.T) {
		gm.RegisterTestingT(t)

		client, err := h.Client()
		gm.Expect(err).To(gm.BeNil())

		session, err := client.EmailLogin(ctx, m.Admin.User.Email, m.Admin.Password)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(session).ToNot(gm.BeNil())

		err = client.Logout(ctx, nil)
		gm.Expect(err).To(gm.BeNil())

		// Can't do authenticated command after deleting session
		err = client.Logout(ctx, nil)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(err.Error()).To(gm.Equal("unknown_cause: Failed to authenticate"))
	})

	t.Run("login user can retrieve their user", func(t *testing.T) {
		gm.RegisterTestingT(t)

		client, err := h.Client()
		gm.Expect(err).To(gm.BeNil())

		session, err := client.EmailLogin(ctx, m.Admin.User.Email, m.Admin.Password)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(session.UserID).To(gm.Equal(m.Admin.User.ID))
		gm.Expect(session.Token).ToNot(gm.BeNil())

		user, err := client.Me(ctx)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(user.GetID()).To(gm.Equal(session.UserID))
		gm.Expect(user.GetEmail()).To(gm.Equal(m.Admin.User.Email))

		creds, err := user.GetCredentials()
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(creds).To(gm.BeNil())
	})

	t.Run("non-auth'd user cannot retrieve their user", func(t *testing.T) {
		gm.RegisterTestingT(t)

		client, err := h.Client()
		gm.Expect(err).To(gm.BeNil())

		_, err = client.Me(ctx)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(err.Error()).To(gm.Equal("unknown_cause: Failed to authenticate"))
	})
}
