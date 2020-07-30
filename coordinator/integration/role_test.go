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

	t.Run("Getting your own role", func(t *testing.T) {
		r, err := client.MyRole(ctx)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(r.Label).To(gm.Equal(models.AdminRole))
	})

	t.Run("Getting your own role in a project", func(t *testing.T) {
		_, err := client.CreateProject(ctx, "My Project", nil, "Who cares")
		gm.Expect(err).To(gm.BeNil())

		r, err := client.MyProjectRole(ctx, "my-project")
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(r.Label).To(gm.Equal(models.ProjectOwnerRole))
	})

	t.Run("Can't change your own role as a member", func(t *testing.T) {
		_, pw, err := client.CreateUser(ctx, "Cool Guy", "cool@person.com")
		gm.Expect(err).To(gm.BeNil())

		_, err = client.EmailLogin(ctx, "cool@person.com", pw)
		gm.Expect(err).To(gm.BeNil())

		err = client.SetOrgRole(ctx, "cool@person.com", models.AdminRole)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(err.Error()).To(gm.Equal("unknown_cause: invalid permissions to change user role"))

		// Return to the admin user
		_, err = client.EmailLogin(ctx, h.Manager().Admin.User.Email, h.Manager().Admin.Password)
		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("Admin can change member roles", func(t *testing.T) {
		u, _, err := client.CreateUser(ctx, "Best Friend", "bestfriend@person.com")
		gm.Expect(err).To(gm.BeNil())

		err = client.SetOrgRole(ctx, "bestfriend@person.com", models.AdminRole)
		gm.Expect(err).To(gm.BeNil())

		user, err := client.GetUser(ctx, u.ID)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(user.Role.Label, models.AdminRole)
	})

	t.Run("Project owner can change member roles", func(t *testing.T) {
		label := models.Label("epic-project")
		p, err := client.CreateProject(ctx, "Epic Project With My Friends", &label, "Who cares")
		gm.Expect(err).To(gm.BeNil())

		u, _, err := client.CreateUser(ctx, "Abc Def", "alphabet@person.com")
		gm.Expect(err).To(gm.BeNil())

		_, err = client.AddContributor(ctx, *p, *u, models.ProjectContributorRole)
		gm.Expect(err).To(gm.BeNil())

		err = client.SetProjectRole(ctx, "alphabet@person.com", label, models.ProjectReaderRole)
		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("Can't change project roles for non-contributors", func(t *testing.T) {
		label := models.Label("epic-project-two")
		_, err := client.CreateProject(ctx, "Hmmmmmmm", &label, "Who cares")
		gm.Expect(err).To(gm.BeNil())

		_, _, err = client.CreateUser(ctx, "Person Person", "iexist@realperson.com")
		gm.Expect(err).To(gm.BeNil())

		err = client.SetProjectRole(ctx, "iexist@realperson.com", label, models.ProjectReaderRole)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(err.Error()).To(gm.Equal("unknown_cause: provided user iexist@realperson.com not found in project epic-project-two"))
	})
}
