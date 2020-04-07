// +build integration

package integration

import (
	"context"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/dropoutlabs/cape/auth"
	"github.com/dropoutlabs/cape/coordinator/harness"
	"github.com/dropoutlabs/cape/primitives"
)

func TestUsers(t *testing.T) {
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
	client, err := m.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	t.Run("create user", func(t *testing.T) {
		creds, err := auth.NewCredentials([]byte("jerryberrybuddyboy"), nil)
		gm.Expect(err).To(gm.BeNil())

		email, err := primitives.NewEmail("jerry@jerry.berry")
		gm.Expect(err).To(gm.BeNil())

		user, err := primitives.NewUser("Jerry Berry", email, creds.Package())
		gm.Expect(err).To(gm.BeNil())

		result, err := client.CreateUser(ctx, user)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(user.Name).To(gm.Equal(result.Name))
		gm.Expect(user.Email).To(gm.Equal(result.Email))

		resp, err := client.GetUser(ctx, result.ID)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(len(resp.Roles)).To(gm.Equal(1))
		gm.Expect(resp.Roles[0].Label).To(gm.Equal(primitives.GlobalRole))
	})

	t.Run("cannot create multiple users with same email", func(t *testing.T) {
		creds, err := auth.NewCredentials([]byte("jerryberrybuddyboy"), nil)
		gm.Expect(err).To(gm.BeNil())

		n, err := primitives.NewName("Lenny Bonedog")
		gm.Expect(err).To(gm.BeNil())

		e, err := primitives.NewEmail("bones@tails.com")
		gm.Expect(err).To(gm.BeNil())

		user, err := primitives.NewUser(n, e, creds.Package())
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(user).NotTo(gm.BeNil())

		otherUser, err := client.CreateUser(ctx, user)
		gm.Expect(err).To(gm.BeNil())

		n, err = primitives.NewName("Julio Tails")
		gm.Expect(err).To(gm.BeNil())

		e, err = primitives.NewEmail("bones@tails.com")
		gm.Expect(err).To(gm.BeNil())

		user, err = primitives.NewUser(n, e, creds.Package())
		gm.Expect(err).To(gm.BeNil())

		otherUser, err = client.CreateUser(ctx, user)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(otherUser).To(gm.BeNil())
	})
}
