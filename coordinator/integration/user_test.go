// +build integration

package integration

import (
	"context"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/coordinator/harness"
	"github.com/capeprivacy/cape/models"
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
		email := models.Email("jerry@jerry.berry")
		name := models.Name("Jerry Berry")

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
		n := models.Name("Lenny Bonedog")
		e := models.Email("lenny@bonedog.com")

		user, _, err := client.CreateUser(ctx, n, e)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(user).ToNot(gm.BeNil())

		nTwo := models.Name("Julio Tails")

		secondUser, _, err := client.CreateUser(ctx, nTwo, e)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(secondUser).To(gm.BeNil())
	})

	t.Run("Can query ME and get my name", func(t *testing.T) {
		me, err := client.Me(ctx)
		gm.Expect(err).To(gm.BeNil())
		name := me.Name
		gm.Expect(name).To(gm.Equal(models.Name("admin")))
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

	n := models.Name("Lenny Bonedog")
	e := models.Email("lenny@bonedog.com")

	user, _, err := client.CreateUser(ctx, n, e)
	gm.Expect(err).To(gm.BeNil())
	gm.Expect(user).ToNot(gm.BeNil())

	nTwo := models.Name("Julio Tails")
	e2 := models.Email("bone2@bonedog.com")

	_, _, err = client.CreateUser(ctx, nTwo, e2)
	gm.Expect(err).To(gm.BeNil())

	users, err := client.ListUsers(ctx)
	gm.Expect(err).To(gm.BeNil())

	// created two here plus admin
	gm.Expect(len(users)).To(gm.Equal(3))
}
