package main

import (
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/prims"
	"testing"

	"github.com/capeprivacy/cape/cmd/ui"
	"github.com/capeprivacy/cape/coordinator"
	"github.com/capeprivacy/cape/primitives"
	gm "github.com/onsi/gomega"
)

func TestProjectsCreate(t *testing.T) {
	gm.RegisterTestingT(t)

	p, err := prims.NewProject("My Project", "my-project", "What is this project even about")
	gm.Expect(err).To(gm.BeNil())

	resp := coordinator.CreateProjectResponse{
		Project: p,
	}

	t.Run("Can create a project", func(t *testing.T) {
		app, u := NewHarness([]*coordinator.MockResponse{
			{
				Value: resp,
			},
		})
		err = app.Run([]string{"cape", "projects", "create", p.Name.String(), p.Description.String()})
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
		err = app.Run([]string{"cape", "projects", "create", p.Name.String()})
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
		err = app.Run([]string{"cape", "projects", "create"})
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(err.Error()).To(gm.Equal("missing_argument: The argument name is required, but was not provided"))
	})

	t.Run("Can update a project", func(t *testing.T) {
		updatedP, err := prims.NewProject("Updated Project", "my-project", "Update ME")
		gm.Expect(err).To(gm.BeNil())

		updateResp := coordinator.UpdateProjectResponse{
			Project: updatedP,
		}

		app, u := NewHarness([]*coordinator.MockResponse{
			{
				Value: updateResp,
			},
		})
		err = app.Run([]string{
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
		updatedP, err := prims.NewProject("Updated Project", "my-project", "What is this project even about")
		gm.Expect(err).To(gm.BeNil())

		updateResp := coordinator.UpdateProjectResponse{
			Project: updatedP,
		}

		app, u := NewHarness([]*coordinator.MockResponse{
			{
				Value: updateResp,
			},
		})
		err = app.Run([]string{
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
		project, err := prims.NewProject("Project", "my-project", "What is this project even about")
		gm.Expect(err).To(gm.BeNil())

		prim, err := database.NewPrimitive(primitives.ProjectSpecType)
		gm.Expect(err).To(gm.BeNil())
		spec := &primitives.ProjectSpec{
			Primitive: prim,
		}
		id, err := database.GenerateID(primitives.ProjectSpecType)
		gm.Expect(err).To(gm.BeNil())
		spec.ID = id

		respBody := coordinator.UpdateProjectSpecResponseBody{
			Project:     project,
			ProjectSpec: spec,
		}

		resp := coordinator.UpdateProjectSpecResponse{UpdateProjectSpecResponseBody: respBody}

		app, u := NewHarness([]*coordinator.MockResponse{
			{
				Value: resp,
			},
		})
		err = app.Run([]string{
			"cape", "projects", "update",
			"--from-spec", "./testdata/project_spec.yaml",
			p.Label.String(),
		})

		gm.Expect(err).To(gm.BeNil())
		gm.Expect(len(u.Calls)).To(gm.Equal(1))
		gm.Expect(u.Calls[0].Name).To(gm.Equal("template"))
	})
}

func TestProjectsList(t *testing.T) {
	gm.RegisterTestingT(t)

	t.Run("Can list some projects", func(t *testing.T) {
		p1, err := prims.NewProject("My Project One", "my-project-one", "What is this project even about")
		gm.Expect(err).To(gm.BeNil())

		p2, err := prims.NewProject("My Project Two", "my-project-two", "What is this project even about")
		gm.Expect(err).To(gm.BeNil())

		resp := coordinator.ListProjectsResponse{
			Projects: []*prims.Project{p1, p2},
		}

		app, u := NewHarness([]*coordinator.MockResponse{
			{
				Value: resp,
			},
		})
		err = app.Run([]string{"cape", "projects", "list"})
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(len(u.Calls)).To(gm.Equal(2))
		gm.Expect(u.Calls[0].Name).To(gm.Equal("table"))
		gm.Expect(u.Calls[1].Name).To(gm.Equal("template"))
	})

	t.Run("Works when there are no projects", func(t *testing.T) {
		resp := coordinator.ListProjectsResponse{
			Projects: []*prims.Project{},
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
