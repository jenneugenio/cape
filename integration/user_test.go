// +build integration

package integration

import (
	"context"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/dropoutlabs/cape/auth"
	"github.com/dropoutlabs/cape/controller"
	"github.com/dropoutlabs/cape/primitives"
)

func TestUsers(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()

	tc, err := controller.NewTestController()
	gm.Expect(err).To(gm.BeNil())

	_, err = tc.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer tc.Teardown(ctx) // nolint: errcheck

	client, err := tc.Client()
	gm.Expect(err).To(gm.BeNil())

	t.Run("create user", func(t *testing.T) {
		creds, err := auth.NewCredentials([]byte("jerryberrybuddyboy"), nil)
		gm.Expect(err).To(gm.BeNil())

		user, err := primitives.NewUser("Jerry Berry", "jerry@jerry.berry", creds.Package())
		gm.Expect(err).To(gm.BeNil())

		otherUser, err := client.CreateUser(ctx, user)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(user.Name).To(gm.Equal(otherUser.Name))
		gm.Expect(user.Email).To(gm.Equal(otherUser.Email))
	})

	t.Run("cannot create multiple users with same email", func(t *testing.T) {
		creds, err := auth.NewCredentials([]byte("jerryberrybuddyboy"), nil)
		gm.Expect(err).To(gm.BeNil())

		user, err := primitives.NewUser("Lenny Bonedog", "bones@tails.com", creds.Package())
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(user).NotTo(gm.BeNil())

		otherUser, err := client.CreateUser(ctx, user)
		gm.Expect(err).To(gm.BeNil())

		user, err = primitives.NewUser("Julio Tails", "bones@tails.com", creds.Package())
		gm.Expect(err).To(gm.BeNil())

		otherUser, err = client.CreateUser(ctx, user)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(otherUser).To(gm.BeNil())
	})
}
