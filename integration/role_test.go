// +build integration

package integration

import (
	"context"
	"testing"

	"github.com/dropoutlabs/cape/controller"
	"github.com/dropoutlabs/cape/database"
	gm "github.com/onsi/gomega"
)

func TestRoles(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()

	tc, err := controller.NewTestController()
	gm.Expect(err).To(gm.BeNil())

	_, err = tc.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer tc.Teardown(ctx) // nolint: errcheck

	client, err := tc.Client()
	gm.Expect(err).To(gm.BeNil())

	_, err = client.Login(ctx, tc.User.Email, tc.UserPassword)
	gm.Expect(err).To(gm.BeNil())

	t.Run("create role", func(t *testing.T) {
		gm.RegisterTestingT(t)

		label := "data-scientist"
		role, err := client.CreateRole(ctx, label, nil)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(role.Label).To(gm.Equal(label))

		// make sure the role exists!!
		role, err = client.GetRole(ctx, role.ID)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(role.Label).To(gm.Equal(label))
	})

	t.Run("delete role", func(t *testing.T) {
		gm.RegisterTestingT(t)

		label := "cio"
		role, err := client.CreateRole(ctx, label, nil)
		gm.Expect(err).To(gm.BeNil())

		err = client.DeleteRole(ctx, role.ID)
		gm.Expect(err).To(gm.BeNil())

		// make sure the role is deleted
		role, err = client.GetRole(ctx, role.ID)
		gm.Expect(err).ToNot(gm.BeNil())
	})

	t.Run("create role with members", func(t *testing.T) {
		gm.RegisterTestingT(t)

		label := "cto"
		role, err := client.CreateRole(ctx, label, []database.ID{tc.User.ID})
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(role).ToNot(gm.BeNil())

		identities, err := client.GetMembersRole(ctx, role.ID)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(len(identities)).To(gm.Equal(1))

		gm.Expect(identities[0].GetID()).To(gm.Equal(tc.User.ID))
	})
}
