// +build integration

package integration

import (
	"context"
	"github.com/capeprivacy/cape/coordinator/harness"
	"github.com/capeprivacy/cape/primitives"
	gm "github.com/onsi/gomega"
	"testing"
)

func TestProjectsCreate(t *testing.T) {
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

	t.Run("Can create a project", func(t *testing.T) {
		p, err := client.CreateProject(ctx, "My Project", nil, "This project does great things")
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(p.Name).To(gm.Equal(primitives.DisplayName("My Project")))
		gm.Expect(p.Label).To(gm.Equal(primitives.Label("my-project")))
		gm.Expect(p.Description).To(gm.Equal(primitives.Description("This project does great things")))
	})

	t.Run("Cannot create a project with the same name", func(t *testing.T) {
		_, err := client.CreateProject(ctx, "Duplicate Me", nil, "This project does great things")
		gm.Expect(err).To(gm.BeNil())

		_, err = client.CreateProject(ctx, "Duplicate Me", nil, "This project does great things")
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(err.Error()).To(gm.Equal("unknown_cause: entity already exists"))
	})
}

func TestProjectsList(t *testing.T) {
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

	seedProjects := []struct {
		Name primitives.DisplayName
	}{
		{Name: "Project One"},
		{Name: "Project Two"},
		{Name: "Project Three"},
		{Name: "Project Four"},
		{Name: "Project Five"},
	}

	for _, p := range seedProjects {
		_, err := client.CreateProject(ctx, p.Name, nil, "This is just a test")
		gm.Expect(err).To(gm.BeNil())
	}

	t.Run("Can list projects", func(t *testing.T) {
		projects, err := client.ListProjects(ctx, nil)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(len(projects)).To(gm.Equal(len(seedProjects)))
	})
}
