// +build integration

package integration

import (
	"context"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/dropoutlabs/cape/auth"
	"github.com/dropoutlabs/cape/coordinator/harness"
	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/primitives"
)

func TestRoles(t *testing.T) {
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
		role, err := client.CreateRole(ctx, label, []database.ID{m.Admin.User.ID})
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(role).ToNot(gm.BeNil())

		identities, err := client.GetMembersRole(ctx, role.ID)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(len(identities)).To(gm.Equal(1))

		gm.Expect(identities[0].GetID()).To(gm.Equal(m.Admin.User.ID))
	})

	t.Run("create role with multiple members", func(t *testing.T) {
		gm.RegisterTestingT(t)

		label, err := primitives.NewLabel("ceoooo")
		gm.Expect(err).To(gm.BeNil())

		email, err := primitives.NewEmail("jye@jdfjkf.com")
		gm.Expect(err).To(gm.BeNil())

		name, err := primitives.NewName("Jye Jdfjkf")
		gm.Expect(err).To(gm.BeNil())

		creds, err := auth.NewCredentials([]byte("randompassword"), nil)
		gm.Expect(err).To(gm.BeNil())

		user, err := primitives.NewUser(name, email, creds.Package())
		gm.Expect(err).To(gm.BeNil())

		user, err = client.CreateUser(ctx, user)
		gm.Expect(err).To(gm.BeNil())

		role, err := client.CreateRole(ctx, label, []database.ID{m.Admin.User.ID, user.ID})
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(role).ToNot(gm.BeNil())

		identities, err := client.GetMembersRole(ctx, role.ID)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(len(identities)).To(gm.Equal(2))

		gm.Expect(identities[0].GetID()).To(gm.Equal(m.Admin.User.ID))
		gm.Expect(identities[1].GetID()).To(gm.Equal(user.ID))
	})

	t.Run("Roles will not default to system roles", func(t *testing.T) {
		l, err := primitives.NewLabel("coolguy")
		gm.Expect(err).To(gm.BeNil())

		role, err := client.CreateRole(ctx, l, []database.ID{m.Admin.User.ID})

		gm.Expect(err).To(gm.BeNil())
		gm.Expect(role.System).To(gm.BeFalse())
	})

	t.Run("role by label", func(t *testing.T) {
		l, err := primitives.NewLabel("coolguy-five")
		gm.Expect(err).To(gm.BeNil())

		role, err := client.CreateRole(ctx, l, []database.ID{m.Admin.User.ID})
		gm.Expect(err).To(gm.BeNil())

		otherRole, err := client.GetRoleByLabel(ctx, role.Label)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(role).To(gm.Equal(otherRole))
	})

	t.Run("cannot delete system role", func(t *testing.T) {
		admin, err := primitives.NewLabel("admin")
		gm.Expect(err).To(gm.BeNil())

		role, err := client.GetRoleByLabel(ctx, admin)
		gm.Expect(err).To(gm.BeNil())

		err = client.DeleteRole(ctx, role.ID)
		gm.Expect(err).ToNot(gm.BeNil())
	})
}

func TestListRoles(t *testing.T) {
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

	dsRole, err := client.CreateRole(ctx, "data-scientist", nil)
	gm.Expect(err).To(gm.BeNil())

	ctoRole, err := client.CreateRole(ctx, "cto", nil)
	gm.Expect(err).To(gm.BeNil())

	roles, err := client.ListRoles(ctx)
	gm.Expect(err).To(gm.BeNil())

	// create two roles + the system roles
	gm.Expect(len(roles)).To(gm.Equal(5))
	gm.Expect(roles).To(gm.ContainElements(dsRole, ctoRole))
}

func TestAssignments(t *testing.T) {
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

	t.Run("assign role", func(t *testing.T) {
		gm.RegisterTestingT(t)

		label, err := primitives.NewLabel("data-scientist")
		gm.Expect(err).To(gm.BeNil())
		role, err := client.CreateRole(ctx, label, nil)
		gm.Expect(err).To(gm.BeNil())

		assignment, err := client.AssignRole(ctx, m.Admin.User.ID, role.ID)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(assignment).NotTo(gm.BeNil())
		gm.Expect(assignment.Identity.GetID()).To(gm.Equal(m.Admin.User.ID))
		gm.Expect(assignment.Role.Label).To(gm.Equal(label))

		identities, err := client.GetMembersRole(ctx, role.ID)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(len(identities)).To(gm.Equal(1))

		gm.Expect(identities[0].GetID()).To(gm.Equal(m.Admin.User.ID))
	})

	t.Run("unassign role", func(t *testing.T) {
		gm.RegisterTestingT(t)

		label, err := primitives.NewLabel("iamarole")
		gm.Expect(err).To(gm.BeNil())
		role, err := client.CreateRole(ctx, label, nil)
		gm.Expect(err).To(gm.BeNil())

		assignment, err := client.AssignRole(ctx, m.Admin.User.ID, role.ID)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(assignment).NotTo(gm.BeNil())

		err = client.UnassignRole(ctx, m.Admin.User.ID, role.ID)
		gm.Expect(err).To(gm.BeNil())

		identities, err := client.GetMembersRole(ctx, role.ID)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(len(identities)).To(gm.Equal(0))
	})
}
