package main

import (
	"github.com/capeprivacy/cape/models"
	"testing"

	"github.com/capeprivacy/cape/cmd/ui"
	"github.com/capeprivacy/cape/coordinator"
	gm "github.com/onsi/gomega"
)

func TestProjectsCreate(t *testing.T) {
	gm.RegisterTestingT(t)

	p := models.NewProject("My Project", "my-project", "What is this project even about")

	resp := coordinator.CreateProjectResponse{
		Project: &p,
	}

	t.Run("Can create a project", func(t *testing.T) {
		app, u := NewHarness([]*coordinator.MockResponse{
			{
				Value: resp,
			},
		})
		err := app.Run([]string{"cape", "projects", "create", p.Name.String(), p.Description.String()})
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(len(u.Calls)).To(gm.Equal(2))
		gm.Expect(u.Calls[0].Name).To(gm.Equal("template"))
		gm.Expect(u.Calls[0].Args[1]).To(gm.Equal(p.Name.String()))
	})

	t.Run("Can create a project without a description", func(t *testing.T) {
		app, u := NewHarness([]*coordinator.MockResponse{
			{
				Value: resp,
			},
		})
		err := app.Run([]string{"cape", "projects", "create", p.Name.String()})
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(len(u.Calls)).To(gm.Equal(2))
		gm.Expect(u.Calls[0].Name).To(gm.Equal("template"))
		gm.Expect(u.Calls[0].Args[1]).To(gm.Equal(p.Name.String()))
	})

	t.Run("Must pass at least a name", func(t *testing.T) {
		app, _ := NewHarness([]*coordinator.MockResponse{
			{
				Value: resp,
			},
		})
		err := app.Run([]string{"cape", "projects", "create"})
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(err.Error()).To(gm.Equal("missing_argument: The argument name is required, but was not provided"))
	})

	t.Run("Can update a project", func(t *testing.T) {
		updatedP := models.NewProject("Updated Project", "my-project", "Update ME")
		updateResp := coordinator.UpdateProjectResponse{
			Project: &updatedP,
		}

		app, u := NewHarness([]*coordinator.MockResponse{
			{
				Value: updateResp,
			},
		})
		err := app.Run([]string{
			"cape", "projects", "update",
			"--name", updatedP.Name.String(), "--description", updatedP.Description.String(),
			p.Label.String(),
		})
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(len(u.Calls)).To(gm.Equal(2))
		gm.Expect(u.Calls[0].Name).To(gm.Equal("template"))
		gm.Expect(u.Calls[0].Args[1]).To(gm.Equal(updatedP.Name.String()))
		gm.Expect(u.Calls[1].Args[0].(ui.Details)["Description"]).To(gm.Equal(updatedP.Description.String()))
	})

	t.Run("Can update a project without a description", func(t *testing.T) {
		updatedP := models.NewProject("Updated Project", "my-project", "What is this project even about")
		updateResp := coordinator.UpdateProjectResponse{
			Project: &updatedP,
		}

		app, u := NewHarness([]*coordinator.MockResponse{
			{
				Value: updateResp,
			},
		})
		err := app.Run([]string{
			"cape", "projects", "update",
			"--name", updatedP.Name.String(),
			p.Label.String(),
		})
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(len(u.Calls)).To(gm.Equal(2))
		gm.Expect(u.Calls[0].Name).To(gm.Equal("template"))
		gm.Expect(u.Calls[0].Args[1]).To(gm.Equal(updatedP.Name.String()))
		gm.Expect(u.Calls[1].Args[0].(ui.Details)["Description"]).To(gm.Equal(updatedP.Description.String()))
	})

	t.Run("Can update a project spec", func(t *testing.T) {
		project := models.NewProject("Project", "my-project", "What is this project even about")
		spec := &models.Policy{}
		spec.ID = "my-spec"

		respBody := coordinator.UpdateProjectSpecResponseBody{
			Project:     &project,
			ProjectSpec: spec,
		}

		resp := coordinator.UpdateProjectSpecResponse{UpdateProjectSpecResponseBody: respBody}

		app, u := NewHarness([]*coordinator.MockResponse{
			{
				Value: resp,
			},
		})
		err := app.Run([]string{
			"cape", "projects", "update",
			"--from-spec", "./testdata/project_spec.yaml",
			p.Label.String(),
		})

		gm.Expect(err).To(gm.BeNil())
		gm.Expect(len(u.Calls)).To(gm.Equal(1))
		gm.Expect(u.Calls[0].Name).To(gm.Equal("template"))
	})

	t.Run("Can suggest a policy", func(t *testing.T) {
		resp := coordinator.SuggestPolicyResponse{
			Suggestion: models.Suggestion{
				ID: "123",
			},
		}

		app, _ := NewHarness([]*coordinator.MockResponse{
			{
				Value: resp,
			},
		})
		err := app.Run([]string{
			"cape", "projects", "policy", "create",
			"--from-spec", "./testdata/project_spec.yaml",
			"my-project",
			"\"My Suggestion\"", "\"Rocks\"",
			p.Label.String(),
		})

		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("Can list suggestions", func(t *testing.T) {
		resp := coordinator.GetProjectSuggestionsResponse{
			Suggestions: []models.Suggestion{
				{
					ID: "123",
				},
			},
		}

		app, _ := NewHarness([]*coordinator.MockResponse{
			{
				Value: resp,
			},
		})
		err := app.Run([]string{
			"cape", "projects", "policy", "list-suggestions",
			"my-project",
			p.Label.String(),
		})

		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("Can get a suggestions", func(t *testing.T) {
		app, _ := NewHarness([]*coordinator.MockResponse{
			{
				Value: coordinator.GetProjectResponse{
					GetProject: coordinator.GetProject{
						Project: &models.Project{
							ID: "1234",
						},
					},
				},
			},
			{
				Value: coordinator.GetProjectSuggestionResponse{
					SuggestionResponse: coordinator.ProjectSuggestion{
						Suggestion: &models.Suggestion{
							ID: "123",
						},
					},
				},
			},
		})
		err := app.Run([]string{
			"cape", "projects", "policy", "get", "abc123",
			p.Label.String(),
		})

		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("Can reject suggestions", func(t *testing.T) {
		resp := coordinator.RejectSuggestionResponse{}
		app, _ := NewHarness([]*coordinator.MockResponse{
			{
				Value: resp,
			},
		})
		err := app.Run([]string{
			"cape", "projects", "policy", "reject", "123",
			p.Label.String(),
		})

		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("Can accept suggestions", func(t *testing.T) {
		resp := coordinator.ApproveSuggestionResponse{}
		app, _ := NewHarness([]*coordinator.MockResponse{
			{
				Value: resp,
			},
		})
		err := app.Run([]string{
			"cape", "projects", "policy", "approve", "123",
			p.Label.String(),
		})

		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("Can list contributors", func(t *testing.T) {
		projectResponse := coordinator.GetProjectResponse{
			GetProject: coordinator.GetProject{
				Project: &models.Project{},
			},
		}
		contribResponse := coordinator.ListContributorsResponse{
			Contributors: []coordinator.GQLContributor{},
		}
		app, _ := NewHarness([]*coordinator.MockResponse{
			{
				Value: projectResponse,
			},
			{
				Value: contribResponse,
			},
		})
		err := app.Run([]string{
			"cape", "projects", "contributors", "list", "my-project",
			p.Label.String(),
		})

		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("Can add contributors", func(t *testing.T) {
		projectResponse := coordinator.GetProjectResponse{
			GetProject: coordinator.GetProject{
				Project: &models.Project{},
			},
		}
		contribResponse := coordinator.UpdateContributorResponse{}
		app, _ := NewHarness([]*coordinator.MockResponse{
			{
				Value: projectResponse,
			},
			{
				Value: contribResponse,
			},
		})
		err := app.Run([]string{
			"cape", "projects", "contributors", "add", "friend@cape.com", "my-project", "project-owner",
			p.Label.String(),
		})

		gm.Expect(err).To(gm.BeNil())
	})
}

func TestProjectsList(t *testing.T) {
	gm.RegisterTestingT(t)

	t.Run("Can list some projects", func(t *testing.T) {
		p1 := models.NewProject("My Project One", "my-project-one", "What is this project even about")
		p2 := models.NewProject("My Project Two", "my-project-two", "What is this project even about")

		resp := coordinator.ListProjectsResponse{
			Projects: []*models.Project{&p1, &p2},
		}

		app, u := NewHarness([]*coordinator.MockResponse{
			{
				Value: resp,
			},
		})
		err := app.Run([]string{"cape", "projects", "list"})
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(len(u.Calls)).To(gm.Equal(2))
		gm.Expect(u.Calls[0].Name).To(gm.Equal("table"))
		gm.Expect(u.Calls[1].Name).To(gm.Equal("template"))
	})

	t.Run("Works when there are no projects", func(t *testing.T) {
		resp := coordinator.ListProjectsResponse{
			Projects: []*models.Project{},
		}

		app, u := NewHarness([]*coordinator.MockResponse{
			{
				Value: resp,
			},
		})
		err := app.Run([]string{"cape", "projects", "list"})
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(len(u.Calls)).To(gm.Equal(1))
		gm.Expect(u.Calls[0].Name).To(gm.Equal("template"))
	})
}
