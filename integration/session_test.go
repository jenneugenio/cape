// +build integration

package integration

import (
	"context"
	"testing"

	"github.com/dropoutlabs/cape/controller"
	"github.com/dropoutlabs/cape/primitives"
	gm "github.com/onsi/gomega"
)

func TestSessions(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()

	tc, err := controller.NewTestController()
	gm.Expect(err).To(gm.BeNil())

	_, err = tc.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer tc.Teardown(ctx) // nolint: errcheck

	client, err := tc.Client()
	gm.Expect(err).To(gm.BeNil())

	t.Run("test client login", func(t *testing.T) {
		gm.RegisterTestingT(t)

		session, err := client.Login(ctx, tc.User.Email, tc.UserPassword)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(session.IdentityID).To(gm.Equal(tc.User.ID))
		gm.Expect(session.Type).To(gm.Equal(primitives.Authenticated))
		gm.Expect(session.Token).ToNot(gm.BeNil())
	})

	t.Run("test fake user fails", func(t *testing.T) {
		gm.RegisterTestingT(t)

		client := controller.NewClient(tc.URL(), nil)

		session, err := client.Login(ctx, "fake@fake.com", []byte("newpasswordwhodis"))
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(session).To(gm.BeNil())

		gm.Expect(err.Error()).To(gm.Equal("graphql: authentication_failure: Failed to authenticate"))
	})

	t.Run("test incorrect credentials", func(t *testing.T) {
		gm.RegisterTestingT(t)

		// fail because credentials inside login won't be right
		session, err := client.Login(ctx, tc.User.Email, []byte("idontknowmypassword"))
		gm.Expect(session).To(gm.BeNil())

		gm.Expect(err.Error()).To(gm.Equal("graphql: authentication_failure: Failed to authenticate"))
	})
}
