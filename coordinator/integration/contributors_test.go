// +build integration

package integration

import (
	"context"
	"github.com/capeprivacy/cape/coordinator/harness"
	"github.com/capeprivacy/cape/models"
	gm "github.com/onsi/gomega"
	"testing"
)

func TestContributors(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()
	cfg, err := harness.NewConfig()
	gm.Expect(err).To(gm.BeNil())

	h, err := harness.NewHarness(cfg)
	gm.Expect(err).To(gm.BeNil())

	err = h.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer h.Teardown(ctx) // nolint: errcheck

	m := h.Manager()
	client, err := m.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	project, err := client.CreateProject(ctx, "My Project", nil, "This project does great things")
	gm.Expect(err).To(gm.BeNil())

	t.Run("Add a contributor", func(t *testing.T) {
		// Our admin is already the project owner, so we make a new user to test adding contributors
		user, _, err := client.CreateUser(ctx, "Noname Mcgee", "dont@me.com")
		gm.Expect(err).To(gm.BeNil())

		contributor, err := client.AddContributor(ctx, *project, user.Email, models.ProjectContributorRole)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(contributor).ToNot(gm.BeNil())
	})

	t.Run("Can't add a contributor twice", func(t *testing.T) {
		user, _, err := client.CreateUser(ctx, "Double Derry", "dd@cape.com")
		gm.Expect(err).To(gm.BeNil())

		_, err = client.AddContributor(ctx, *project, user.Email, models.ProjectContributorRole)
		gm.Expect(err).To(gm.BeNil())

		_, err = client.AddContributor(ctx, *project, user.Email, models.ProjectContributorRole)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(err.Error()).To(gm.Equal("unknown_cause: duplicate key"))
	})

	t.Run("Can list contributors", func(t *testing.T) {
		contributors, err := client.ListContributors(ctx, *project)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(contributors).ToNot(gm.BeNil())
	})

	t.Run("New projects will have an owner", func(t *testing.T) {
		project, err := client.CreateProject(ctx, "New Project", nil, "This project does good things")
		gm.Expect(err).To(gm.BeNil())

		p, err := client.GetProject(ctx, project.ID, nil)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(p.Contributors).ToNot(gm.BeNil())
		gm.Expect(p.Contributors[0].User.ID).To(gm.Equal(m.Admin.User.ID))
	})

	t.Run("Can remove a contributor", func(t *testing.T) {
		project, err := client.CreateProject(ctx, "Unique Project", nil, "This project does great things")
		gm.Expect(err).To(gm.BeNil())

		user, _, err := client.CreateUser(ctx, "Remove McMee", "rm@cape.com")
		gm.Expect(err).To(gm.BeNil())

		_, err = client.AddContributor(ctx, *project, user.Email, models.ProjectContributorRole)
		gm.Expect(err).To(gm.BeNil())

		contributors, err := client.ListContributors(ctx, *project)
		gm.Expect(err).To(gm.BeNil())

		// Admin and Remove McMee are the contributors
		gm.Expect(len(contributors)).To(gm.Equal(2))

		_, err = client.RemoveContributor(ctx, *user, *project)
		gm.Expect(err).To(gm.BeNil())

		contributors, err = client.ListContributors(ctx, *project)
		gm.Expect(err).To(gm.BeNil())

		// Admin and Remove McMee are the contributors
		gm.Expect(len(contributors)).To(gm.Equal(1))
		gm.Expect(contributors[0].User.ID).To(gm.Equal(m.Admin.User.ID))
	})
}
