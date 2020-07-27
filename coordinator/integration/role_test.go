// +build integration

package integration

import (
	"context"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/coordinator/harness"
	"github.com/capeprivacy/cape/models"
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

	t.Run("Test getting your own role", func(t *testing.T) {
		r, err := client.MyRole(ctx)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(r.Label).To(gm.Equal(models.AdminRole))
	})

	t.Run("Test getting your own role in a project", func(t *testing.T) {
		_, err := client.CreateProject(ctx, "My Project", nil, "Who cares")
		gm.Expect(err).To(gm.BeNil())

		r, err := client.MyProjectRole(ctx, "my-project")
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(r.Label).To(gm.Equal(models.ProjectOwnerRole))
	})

	t.Run("Can't get a role for a project you don't belong to", func(t *testing.T) {
		_, err := client.MyProjectRole(ctx, "fake-project")
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(err.Error()).To(gm.Equal("unknown_cause: you are not a member of the requested project: fake-project"))
	})
}

//func TestListRoles(t *testing.T) {
//	gm.RegisterTestingT(t)
//
//	ctx := context.Background()
//	cfg, err := harness.NewConfig()
//	gm.Expect(err).To(gm.BeNil())
//
//	h, err := harness.NewHarness(cfg)
//	gm.Expect(err)
//
//	err = h.Setup(ctx)
//	gm.Expect(err).To(gm.BeNil())
//
//	defer h.Teardown(ctx) // nolint: errcheck
//
//	m := h.Manager()
//	client, err := m.Setup(ctx)
//	gm.Expect(err).To(gm.BeNil())
//
//	dsRole, err := client.CreateRole(ctx, "data-scientist", nil)
//	gm.Expect(err).To(gm.BeNil())
//
//	ctoRole, err := client.CreateRole(ctx, "ctoo", nil)
//	gm.Expect(err).To(gm.BeNil())
//
//	roles, err := client.ListRoles(ctx)
//	gm.Expect(err).To(gm.BeNil())
//
//	// create two roles + the system roles
//	gm.Expect(len(roles)).To(gm.Equal(2 + len(primitives.SystemRoles)))
//	gm.Expect(roles).To(gm.ContainElements(dsRole, ctoRole))
//}
//
//func TestAssignments(t *testing.T) {
//	gm.RegisterTestingT(t)
//
//	ctx := context.Background()
//	cfg, err := harness.NewConfig()
//	gm.Expect(err).To(gm.BeNil())
//
//	h, err := harness.NewHarness(cfg)
//	gm.Expect(err)
//
//	err = h.Setup(ctx)
//	gm.Expect(err).To(gm.BeNil())
//
//	defer h.Teardown(ctx) // nolint: errcheck
//
//	m := h.Manager()
//	client, err := m.Setup(ctx)
//	gm.Expect(err).To(gm.BeNil())
//
//	t.Run("assign role", func(t *testing.T) {
//		gm.RegisterTestingT(t)
//
//		label, err := primitives.NewLabel("data-scientist")
//		gm.Expect(err).To(gm.BeNil())
//		role, err := client.CreateRole(ctx, label, nil)
//		gm.Expect(err).To(gm.BeNil())
//
//		assignment, err := client.AssignRole(ctx, m.Admin.User.ID, role.ID)
//		gm.Expect(err).To(gm.BeNil())
//		gm.Expect(assignment).NotTo(gm.BeNil())
//		gm.Expect(assignment.User.ID).To(gm.Equal(m.Admin.User.ID))
//		gm.Expect(assignment.Role.Label).To(gm.Equal(label))
//
//		users, err := client.GetMembersRole(ctx, role.ID)
//		gm.Expect(err).To(gm.BeNil())
//
//		gm.Expect(len(users)).To(gm.Equal(1))
//
//		gm.Expect(users[0].ID).To(gm.Equal(m.Admin.User.ID))
//	})
//
//	t.Run("unassign role", func(t *testing.T) {
//		gm.RegisterTestingT(t)
//
//		label, err := primitives.NewLabel("iamarole")
//		gm.Expect(err).To(gm.BeNil())
//		role, err := client.CreateRole(ctx, label, nil)
//		gm.Expect(err).To(gm.BeNil())
//
//		assignment, err := client.AssignRole(ctx, m.Admin.User.ID, role.ID)
//		gm.Expect(err).To(gm.BeNil())
//		gm.Expect(assignment).NotTo(gm.BeNil())
//
//		err = client.UnassignRole(ctx, m.Admin.User.ID, role.ID)
//		gm.Expect(err).To(gm.BeNil())
//
//		users, err := client.GetMembersRole(ctx, role.ID)
//		gm.Expect(err).To(gm.BeNil())
//
//		gm.Expect(len(users)).To(gm.Equal(0))
//	})
//}
