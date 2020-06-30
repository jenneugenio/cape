// +build integration

package integration

import (
	"context"
	"io/ioutil"
	"testing"
	"time"

	"github.com/capeprivacy/cape/coordinator/harness"
	"github.com/capeprivacy/cape/primitives"
	gm "github.com/onsi/gomega"
)

func TestProjects(t *testing.T) {
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

	t.Run("Can get a project by id", func(t *testing.T) {
		// time out of the database has ms rounded off of it, so testStart can be ahead of the db by a tiny amount of time
		testStart := time.Now().AddDate(0, 0, -1)

		p1, err := client.CreateProject(ctx, "Make Me", nil, "This project does great things")
		gm.Expect(err).To(gm.BeNil())

		p2, err := client.GetProject(ctx, &p1.ID, nil)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(p2.Name).To(gm.Equal(p1.Name))
		gm.Expect(p2.Description).To(gm.Equal(p1.Description))
		gm.Expect(p2.CreatedAt.After(testStart)).To(gm.BeTrue())
		gm.Expect(p2.UpdatedAt.After(testStart)).To(gm.BeTrue())
	})

	t.Run("Can get a project by label", func(t *testing.T) {
		// time out of the database has ms rounded off of it, so testStart can be ahead of the db by a tiny amount of time
		testStart := time.Now().AddDate(0, 0, -1)

		p1, err := client.CreateProject(ctx, "Make Me Please", nil, "This project does great things")
		gm.Expect(err).To(gm.BeNil())

		p2, err := client.GetProject(ctx, nil, &p1.Label)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(p2.Name).To(gm.Equal(p1.Name))
		gm.Expect(p2.Description).To(gm.Equal(p1.Description))
		gm.Expect(p2.CreatedAt.After(testStart)).To(gm.BeTrue())
		gm.Expect(p2.UpdatedAt.After(testStart)).To(gm.BeTrue())
	})

	t.Run("Cannot create a project with the same name", func(t *testing.T) {
		_, err := client.CreateProject(ctx, "Duplicate Me", nil, "This project does great things")
		gm.Expect(err).To(gm.BeNil())

		_, err = client.CreateProject(ctx, "Duplicate Me", nil, "This project does great things")
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(err.Error()).To(gm.Equal("unknown_cause: entity already exists"))
	})

	t.Run("Can update name and description by label", func(t *testing.T) {
		p, err := client.CreateProject(ctx, "Updatable Project", nil, "This project does great things")
		gm.Expect(err).To(gm.BeNil())

		name := primitives.DisplayName("New Name")
		desc := primitives.Description("This project is now updated")

		updatedP, err := client.UpdateProject(ctx, nil, &p.Label, &name, &desc)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(updatedP.Name).To(gm.Equal(name))
		gm.Expect(updatedP.Label).To(gm.Equal(p.Label))
		gm.Expect(updatedP.Description).To(gm.Equal(desc))
	})

	t.Run("Cannot update a project with the same name", func(t *testing.T) {
		_, err := client.CreateProject(ctx, "Same Project", nil, "This project does great things")
		gm.Expect(err).To(gm.BeNil())

		p, err := client.CreateProject(ctx, "Not Same Project", nil, "This project does great things")
		gm.Expect(err).To(gm.BeNil())

		name := primitives.DisplayName("Same Project")
		desc := primitives.Description("This project is now updated")

		_, err = client.UpdateProject(ctx, nil, &p.Label, &name, &desc)
		gm.Expect(err).NotTo(gm.BeNil())
	})

	t.Run("Can update name and description by id", func(t *testing.T) {
		p, err := client.CreateProject(ctx, "Another Updatable Project", nil, "This project does great things")
		gm.Expect(err).To(gm.BeNil())

		name := primitives.DisplayName("New New Name")
		desc := primitives.Description("This project is now updated")

		updatedP, err := client.UpdateProject(ctx, &p.ID, nil, &name, &desc)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(updatedP.Name).To(gm.Equal(name))
		gm.Expect(updatedP.Label).To(gm.Equal(p.Label))
		gm.Expect(updatedP.Description).To(gm.Equal(desc))
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

func TestProjectSpecCreate(t *testing.T) {
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

	p, err := client.CreateProject(ctx, "My Project", nil, "This is my project")
	gm.Expect(err).To(gm.BeNil())

	// seed a source
	l, err := primitives.NewLabel("transactions")
	gm.Expect(err).To(gm.BeNil())
	dbURL, err := primitives.NewDBURL("postgres://postgres:dev@my.cool.website:5432/mydb")
	gm.Expect(err).To(gm.BeNil())
	_, err = client.AddSource(ctx, l, dbURL)
	gm.Expect(err).To(gm.BeNil())

	f, err := ioutil.ReadFile("./testdata/project_spec.yaml")
	gm.Expect(err).To(gm.BeNil())

	spec, err := primitives.ParseProjectSpecFile(f)
	gm.Expect(err).To(gm.BeNil())

	t.Run("Can create a spec", func(t *testing.T) {
		gm.Expect(p.Status).To(gm.Equal(primitives.ProjectPending))
		p, _, err := client.UpdateProjectSpec(ctx, p.Label, spec)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(p.Status).To(gm.Equal(primitives.ProjectActive))
	})

	t.Run("New specs become active", func(t *testing.T) {
		p, s, err := client.UpdateProjectSpec(ctx, p.Label, spec)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(p.Status).To(gm.Equal(primitives.ProjectActive))
		gm.Expect(p.CurrentSpecID).To(gm.Equal(&s.ID))
	})

	t.Run("Cannot create a spec for a source that doesn't exist", func(t *testing.T) {
		spec.Sources = []primitives.Label{"iamprettysureiwontbeinthedatabase"}
		p, _, err := client.UpdateProjectSpec(ctx, p.Label, spec)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(err.Error()).To(gm.Equal("unknown_cause: One or more sources declared in the project spec do not exist"))
		gm.Expect(p).To(gm.BeNil())
	})

	t.Run("Cannot update a spec with a source with at least one non existant source", func(t *testing.T) {
		spec, err = primitives.ParseProjectSpecFile(f)
		gm.Expect(err).To(gm.BeNil())
		spec.Sources = append(spec.Sources, "imnotreal")
		_, _, err = client.UpdateProjectSpec(ctx, p.Label, spec)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(err.Error()).To(gm.Equal("unknown_cause: One or more sources declared in the project spec do not exist"))
	})
}
