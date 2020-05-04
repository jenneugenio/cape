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

		gm.Expect(session.IdentityID).To(gm.Equal(m.Admin.User.ID))
		gm.Expect(session.Type).To(gm.Equal(primitives.Authenticated))
		gm.Expect(session.Token).ToNot(gm.BeNil())
	})

	t.Run("test fake user fails", func(t *testing.T) {
		gm.RegisterTestingT(t)

		client, err := h.Client()
		gm.Expect(err).To(gm.BeNil())

		email, err := primitives.NewEmail("fake@fake.com")
		gm.Expect(err).To(gm.BeNil())

		session, err := client.EmailLogin(ctx, email, []byte("newpasswordwhodis"))
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(session).To(gm.BeNil())

		gm.Expect(err.Error()).To(gm.Equal("unknown_cause: Failed to authenticate"))
	})

	t.Run("test incorrect credentials", func(t *testing.T) {
		gm.RegisterTestingT(t)

		client, err := h.Client()
		gm.Expect(err).To(gm.BeNil())

		// fail because credentials inside login won't be right
		session, err := client.EmailLogin(ctx, m.Admin.User.Email, []byte("idontknowmypassword"))
		gm.Expect(session).To(gm.BeNil())

		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(err.Error()).To(gm.Equal("unknown_cause: Failed to authenticate"))
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

	t.Run("login user can retrieve their identity", func(t *testing.T) {
		gm.RegisterTestingT(t)

		client, err := h.Client()
		gm.Expect(err).To(gm.BeNil())

		session, err := client.EmailLogin(ctx, m.Admin.User.Email, m.Admin.Password)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(session.IdentityID).To(gm.Equal(m.Admin.User.ID))
		gm.Expect(session.Type).To(gm.Equal(primitives.Authenticated))
		gm.Expect(session.Token).ToNot(gm.BeNil())

		identity, err := client.Me(ctx)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(identity.GetID()).To(gm.Equal(session.IdentityID))
		gm.Expect(identity.GetEmail()).To(gm.Equal(m.Admin.User.Email))
	})

	t.Run("non-auth'd user cannot retrieve their identity", func(t *testing.T) {
		gm.RegisterTestingT(t)

		client, err := h.Client()
		gm.Expect(err).To(gm.BeNil())

		_, err = client.Me(ctx)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(err.Error()).To(gm.Equal("unknown_cause: Failed to authenticate"))
	})
}
