// +build integration

package integration

import (
	"context"
	"github.com/dropoutlabs/cape/primitives"
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

		label, err := primitives.NewLabel("data-scientist")
		gm.Expect(err).To(gm.BeNil())
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

		label, err := primitives.NewLabel("cio-person")
		gm.Expect(err).To(gm.BeNil())
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

		label, err := primitives.NewLabel("cto-person")
		gm.Expect(err).To(gm.BeNil())
		role, err := client.CreateRole(ctx, label, []database.ID{tc.User.ID})
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(role).ToNot(gm.BeNil())

		identities, err := client.GetMembersRole(ctx, role.ID)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(len(identities)).To(gm.Equal(1))

		gm.Expect(identities[0].GetID()).To(gm.Equal(tc.User.ID))
	})

	t.Run("list roles", func(t *testing.T) {
		gm.RegisterTestingT(t)

		roles, err := client.ListRoles(ctx)
		gm.Expect(err).To(gm.BeNil())

		// create three early in the tests and then delete one
		gm.Expect(len(roles)).To(gm.Equal(2))
		gm.Expect(roles[0].Label.String()).To(gm.Equal("data-scientist"))
		gm.Expect(roles[1].Label.String()).To(gm.Equal("cto-person"))
	})

	t.Run("Roles will not default to system roles", func(t *testing.T) {
		l, err := primitives.NewLabel("coolguy")
		role, err := client.CreateRole(ctx, l, []database.ID{tc.User.ID})

		gm.Expect(err).To(gm.BeNil())
		gm.Expect(role.System).To(gm.BeFalse())
	})
}

func TestListRoles(t *testing.T) {
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

	dsRole, err := client.CreateRole(ctx, "data-scientist", nil)
	gm.Expect(err).To(gm.BeNil())

	ctoRole, err := client.CreateRole(ctx, "cto", nil)
	gm.Expect(err).To(gm.BeNil())

	roles, err := client.ListRoles(ctx)
	gm.Expect(err).To(gm.BeNil())

	// create three early in the tests and then delete one
	gm.Expect(len(roles)).To(gm.Equal(2))
	gm.Expect(roles).To(gm.ContainElements(dsRole, ctoRole))
}

func TestAssignments(t *testing.T) {
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

	t.Run("assign role", func(t *testing.T) {
		gm.RegisterTestingT(t)

		label, err := primitives.NewLabel("data-scientist")
		gm.Expect(err).To(gm.BeNil())
		role, err := client.CreateRole(ctx, label, nil)
		gm.Expect(err).To(gm.BeNil())

		assignment, err := client.AssignRole(ctx, tc.User.ID, role.ID)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(assignment).NotTo(gm.BeNil())
		gm.Expect(assignment.Identity.GetID()).To(gm.Equal(tc.User.ID))
		gm.Expect(assignment.Role.Label).To(gm.Equal(label))

		identities, err := client.GetMembersRole(ctx, role.ID)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(len(identities)).To(gm.Equal(1))

		gm.Expect(identities[0].GetID()).To(gm.Equal(tc.User.ID))
	})

	t.Run("unassign role", func(t *testing.T) {
		gm.RegisterTestingT(t)

		label, err := primitives.NewLabel("iamarole")
		gm.Expect(err).To(gm.BeNil())
		role, err := client.CreateRole(ctx, label, nil)
		gm.Expect(err).To(gm.BeNil())

		assignment, err := client.AssignRole(ctx, tc.User.ID, role.ID)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(assignment).NotTo(gm.BeNil())

		err = client.UnassignRole(ctx, tc.User.ID, role.ID)
		gm.Expect(err).To(gm.BeNil())

		identities, err := client.GetMembersRole(ctx, role.ID)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(len(identities)).To(gm.Equal(0))
	})
}
