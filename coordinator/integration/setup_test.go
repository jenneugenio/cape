// +build integration

package integration

import (
	"context"
	"fmt"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/coordinator/harness"
	"github.com/capeprivacy/cape/primitives"
)

func TestSetup(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()
	cfg, err := harness.NewConfig()
	gm.Expect(err).To(gm.BeNil())

	h, err := harness.NewHarness(cfg)
	gm.Expect(err)

	err = h.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer h.Teardown(ctx) // nolint: errcheck

	password, err := primitives.NewPassword("jerryberrybuddyboy")
	gm.Expect(err).To(gm.BeNil())

	name, err := primitives.NewName("ben")
	gm.Expect(err).To(gm.BeNil())

	email, err := primitives.NewEmail("ben@capeprivacy.com")
	gm.Expect(err).To(gm.BeNil())

	t.Run("cannot login prior to setup", func(t *testing.T) {
		gm.RegisterTestingT(t)

		client, err := h.Client()
		gm.Expect(err).To(gm.BeNil())

		_, err = client.EmailLogin(ctx, email, password)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(err.Error()).To(gm.Equal("unknown_cause: Failed to authenticate"))
	})

	t.Run("Setup cape", func(t *testing.T) {
		gm.RegisterTestingT(t)

		client, err := h.Client()
		gm.Expect(err).To(gm.BeNil())

		admin, err := client.Setup(ctx, name, email, password)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(admin.Name).To(gm.Equal(name))
		gm.Expect(admin.Email).To(gm.Equal(email))

		_, err = client.EmailLogin(ctx, admin.Email, password)
		gm.Expect(err).To(gm.BeNil())

		roles, err := client.ListRoles(ctx)
		gm.Expect(err).To(gm.BeNil())

		labels := make([]primitives.Label, len(primitives.SystemRoles))
		for i, role := range roles {
			labels[i] = role.Label
			gm.Expect(role.System).To(gm.BeTrue())

			// make sure new user is assigned admin and global roles
			if role.Label == primitives.AdminRole || role.Label == primitives.GlobalRole {
				members, err := client.GetMembersRole(ctx, role.ID)
				gm.Expect(err).To(gm.BeNil())

				gm.Expect(members[0].GetEmail()).To(gm.Equal(admin.Email))
			}
		}

		gm.Expect(labels).To(gm.ContainElements(primitives.SystemRoles))
	})

	t.Run("Setup cannot be called a second time", func(t *testing.T) {
		gm.RegisterTestingT(t)

		client, err := h.Client()
		gm.Expect(err).To(gm.BeNil())

		_, err = client.Setup(ctx, name, email, password)
		gm.Expect(err).ToNot(gm.BeNil())
	})

	t.Run("Setup cannot be called a second time with different information", func(t *testing.T) {
		gm.RegisterTestingT(t)

		client, err := h.Client()
		gm.Expect(err).To(gm.BeNil())

		password, err := primitives.NewPassword("jerryberrybuddythey")
		gm.Expect(err).To(gm.BeNil())

		name, err := primitives.NewName("justin")
		gm.Expect(err).To(gm.BeNil())

		email, err := primitives.NewEmail("justin@capeprivacy.com")
		gm.Expect(err).To(gm.BeNil())

		_, err = client.Setup(ctx, name, email, password)
		gm.Expect(err).ToNot(gm.BeNil())
	})
}

func TestDeleteSystemRoles(t *testing.T) {
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

	systemRoles := []string{"admin", "global"}
	for _, roleLabel := range systemRoles {
		t.Run(fmt.Sprintf("can't delete %s role", roleLabel), func(t *testing.T) {
			admin, err := primitives.NewLabel(roleLabel)
			gm.Expect(err).To(gm.BeNil())

			role, err := client.GetRoleByLabel(ctx, admin)
			gm.Expect(err).To(gm.BeNil())

			err = client.DeleteRole(ctx, role.ID)
			gm.Expect(err).ToNot(gm.BeNil())
		})
	}
}
