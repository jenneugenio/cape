// +build integration

package integration

import (
	"context"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/coordinator/harness"
	"github.com/capeprivacy/cape/primitives"
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
		email, err := primitives.NewEmail("jerry@jerry.berry")
		gm.Expect(err).To(gm.BeNil())

		name, err := primitives.NewName("Jerry Berry")
		gm.Expect(err).To(gm.BeNil())

		result, _, err := client.CreateUser(ctx, name, email)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(name).To(gm.Equal(result.Name))
		gm.Expect(email).To(gm.Equal(result.Email))

		resp, err := client.GetUser(ctx, result.ID)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(len(resp.Roles)).To(gm.Equal(1))
		gm.Expect(resp.Roles[0].Label).To(gm.Equal(primitives.GlobalRole))
	})

	t.Run("cannot create multiple users with same email", func(t *testing.T) {
		n, err := primitives.NewName("Lenny Bonedog")
		gm.Expect(err).To(gm.BeNil())

		e, err := primitives.NewEmail("bones@tails.com")
		gm.Expect(err).To(gm.BeNil())

		user, _, err := client.CreateUser(ctx, n, e)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(user).ToNot(gm.BeNil())

		nTwo, err := primitives.NewName("Julio Tails")
		gm.Expect(err).To(gm.BeNil())

		secondUser, _, err := client.CreateUser(ctx, nTwo, e)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(secondUser).To(gm.BeNil())
	})

	t.Run("Can query ME and get my name", func(t *testing.T) {
		me, err := client.Me(ctx)
		gm.Expect(err).To(gm.BeNil())
		name := me.GetName()
		gm.Expect(name).To(gm.Equal(primitives.Name("admin")))
	})
}

func TestListUsers(t *testing.T) {
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

	n, err := primitives.NewName("Lenny Bonedog")
	gm.Expect(err).To(gm.BeNil())

	e, err := primitives.NewEmail("bones@tails.com")
	gm.Expect(err).To(gm.BeNil())

	user, _, err := client.CreateUser(ctx, n, e)
	gm.Expect(err).To(gm.BeNil())
	gm.Expect(user).ToNot(gm.BeNil())

	nTwo, err := primitives.NewName("Julio Tails")
	gm.Expect(err).To(gm.BeNil())

	e2, err := primitives.NewEmail("bones2@tails.com")
	gm.Expect(err).To(gm.BeNil())

	_, _, err = client.CreateUser(ctx, nTwo, e2)
	gm.Expect(err).To(gm.BeNil())

	users, err := client.ListUsers(ctx)
	gm.Expect(err).To(gm.BeNil())

	// created two here plus admin
	gm.Expect(len(users)).To(gm.Equal(3))
}
