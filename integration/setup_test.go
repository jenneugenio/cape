// +build integration

package integration

import (
	"context"
	"fmt"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/dropoutlabs/cape/auth"
	"github.com/dropoutlabs/cape/controller/harness"
	"github.com/dropoutlabs/cape/primitives"
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

	t.Run("Setup cape", func(t *testing.T) {
		gm.RegisterTestingT(t)

		client, err := h.Client()
		gm.Expect(err).To(gm.BeNil())

		password := []byte("jerryberrybuddyboy")
		creds, err := auth.NewCredentials(password, nil)
		gm.Expect(err).To(gm.BeNil())

		email, err := primitives.NewEmail("ben@capeprivacy.com")
		gm.Expect(err).To(gm.BeNil())

		user, err := primitives.NewUser("ben", email, creds.Package())
		gm.Expect(err).To(gm.BeNil())

		admin, err := client.Setup(ctx, user)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(admin.Name).To(gm.Equal(user.Name))
		gm.Expect(admin.Email).To(gm.Equal(user.Email))

		_, err = client.Login(ctx, admin.Email, password)
		gm.Expect(err).To(gm.BeNil())

		roles, err := client.ListRoles(ctx)
		gm.Expect(err).To(gm.BeNil())

		roleLabels := []primitives.Label{
			primitives.AdminRole,
			primitives.DataConnectorRole,
			primitives.GlobalRole,
		}

		labels := make([]primitives.Label, 3)
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

		gm.Expect(labels).To(gm.ContainElements(roleLabels))
	})

	t.Run("Setup cannot be called a second time", func(t *testing.T) {
		gm.RegisterTestingT(t)

		client, err := h.Client()
		gm.Expect(err).To(gm.BeNil())

		creds, err := auth.NewCredentials([]byte("jerryberrybuddyboy"), nil)
		gm.Expect(err).To(gm.BeNil())

		email, err := primitives.NewEmail("ben@capeprivacy.com")
		gm.Expect(err).To(gm.BeNil())

		user, err := primitives.NewUser("ben", email, creds.Package())
		gm.Expect(err).To(gm.BeNil())

		_, err = client.Setup(ctx, user)
		gm.Expect(err).ToNot(gm.BeNil())
	})

	t.Run("Setup cannot be called a second time with different information", func(t *testing.T) {
		gm.RegisterTestingT(t)

		client, err := h.Client()
		gm.Expect(err).To(gm.BeNil())

		creds, err := auth.NewCredentials([]byte("berryjerrybuddyboy"), nil)
		gm.Expect(err).To(gm.BeNil())

		email, err := primitives.NewEmail("justin@capeprivacy.com")
		gm.Expect(err).To(gm.BeNil())

		user, err := primitives.NewUser("justin", email, creds.Package())
		gm.Expect(err).To(gm.BeNil())

		_, err = client.Setup(ctx, user)
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

	systemRoles := []string{"admin", "global", "data-connector"}
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
