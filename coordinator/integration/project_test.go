// +build integration

package integration

import (
	"context"
	"github.com/capeprivacy/cape/models"
	"io/ioutil"
	"testing"
	"time"

	"github.com/capeprivacy/cape/coordinator/harness"
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
		gm.Expect(p.Name).To(gm.Equal(models.ProjectDisplayName("My Project")))
		gm.Expect(p.Label).To(gm.Equal(models.Label("my-project")))
		gm.Expect(p.Description).To(gm.Equal(models.ProjectDescription("This project does great things")))
	})

	t.Run("Can get a project by id", func(t *testing.T) {
		// time out of the database has ms rounded off of it, so testStart can be ahead of the db by a tiny amount of time
		testStart := time.Now().AddDate(0, 0, -1)

		p1, err := client.CreateProject(ctx, "Make Me", nil, "This project does great things")
		gm.Expect(err).To(gm.BeNil())

		p2, err := client.GetProject(ctx, p1.ID, nil)
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

		p2, err := client.GetProject(ctx, "", &p1.Label)
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

		name := models.ProjectDisplayName("New Name")
		desc := models.ProjectDescription("This project is now updated")

		updatedP, err := client.UpdateProject(ctx, "", &p.Label, &name, &desc)
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

		name := models.ProjectDisplayName("Same Project")
		desc := models.ProjectDescription("This project is now updated")

		_, err = client.UpdateProject(ctx, "", &p.Label, &name, &desc)
		gm.Expect(err).NotTo(gm.BeNil())
	})

	t.Run("Can update name and description by id", func(t *testing.T) {
		p, err := client.CreateProject(ctx, "Another Updatable Project", nil, "This project does great things")
		gm.Expect(err).To(gm.BeNil())

		name := models.ProjectDisplayName("New New Name")
		desc := models.ProjectDescription("This project is now updated")

		updatedP, err := client.UpdateProject(ctx, p.ID, nil, &name, &desc)
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
		Name models.ProjectDisplayName
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
		projects, err := client.ListProjects(ctx, models.Any)
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

	f, err := ioutil.ReadFile("./testdata/project_spec.yaml")
	gm.Expect(err).To(gm.BeNil())

	spec, err := models.ParseProjectSpecFile(f)
	gm.Expect(err).To(gm.BeNil())

	t.Run("Can create a spec", func(t *testing.T) {
		gm.Expect(p.Status).To(gm.Equal(models.ProjectPending))
		p, _, err := client.UpdateProjectSpec(ctx, p.Label, spec)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(p.Status).To(gm.Equal(models.ProjectActive))
	})

	t.Run("New specs become active", func(t *testing.T) {
		p, s, err := client.UpdateProjectSpec(ctx, p.Label, spec)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(p.Status).To(gm.Equal(models.ProjectActive))
		gm.Expect(p.CurrentSpecID).To(gm.Equal(s.ID))
	})

	t.Run("Can suggest a spec", func(t *testing.T) {
		p, err := client.CreateProject(ctx, "suggest-me", nil, "This is my project")
		gm.Expect(err).To(gm.BeNil())

		suggestion, err := client.SuggestPolicy(ctx, p.Label, "Make a change", "it's for the best", spec)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(suggestion.State).To(gm.Equal(models.SuggestionPending))
	})

	t.Run("Can list suggestions", func(t *testing.T) {
		p, err := client.CreateProject(ctx, "suggest-me-lots", nil, "This is my project")
		gm.Expect(err).To(gm.BeNil())

		for i := 0; i < 10; i++ {
			_, err := client.SuggestPolicy(ctx, p.Label, "Make a change", "it's for the best", spec)
			gm.Expect(err).To(gm.BeNil())
		}

		suggestions, err := client.GetProjectSuggestions(ctx, p.Label)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(len(suggestions)).To(gm.Equal(10))
	})

	t.Run("Can approve suggestions", func(t *testing.T) {
		p, err := client.CreateProject(ctx, "approve-me", nil, "This is my project")
		gm.Expect(err).To(gm.BeNil())

		s, err := client.SuggestPolicy(ctx, p.Label, "Make a change", "it's for the best", spec)
		gm.Expect(err).To(gm.BeNil())
		err = client.ApproveSuggestion(ctx, *s)
		gm.Expect(err).To(gm.BeNil())

		projectResp, err := client.GetProject(ctx, p.ID, nil)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(projectResp.Policy.ID).To(gm.Equal(s.PolicyID))

		suggs, err := client.GetProjectSuggestions(ctx, "approve-me")
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(suggs[0].State).To(gm.Equal(models.SuggestionApproved))
	})

	t.Run("Can reject suggestions", func(t *testing.T) {
		p, err := client.CreateProject(ctx, "reject-me", nil, "This is my project")
		gm.Expect(err).To(gm.BeNil())

		s, err := client.SuggestPolicy(ctx, p.Label, "Make a change", "it's for the best", spec)
		gm.Expect(err).To(gm.BeNil())
		err = client.RejectSuggestion(ctx, *s)
		gm.Expect(err).To(gm.BeNil())

		suggs, err := client.GetProjectSuggestions(ctx, "reject-me")
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(suggs[0].State).To(gm.Equal(models.SuggestionRejected))
	})

	t.Run("Can get a specific suggestion", func(t *testing.T) {
		p, err := client.CreateProject(ctx, "get-me", nil, "This is my project")
		gm.Expect(err).To(gm.BeNil())

		s, err := client.SuggestPolicy(ctx, p.Label, "Make a change", "it's for the best", spec)
		gm.Expect(err).To(gm.BeNil())

		resp, err := client.GetProjectSuggestion(ctx, s.ID)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(resp.Project.Label).To(gm.Equal(models.Label("get-me")))
		gm.Expect(resp.State).To(gm.Equal(models.SuggestionPending))
		gm.Expect(len(resp.Policy.Rules) > 0).To(gm.BeTrue())
	})
}
