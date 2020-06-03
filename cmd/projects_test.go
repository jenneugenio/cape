package main

import (
	"github.com/capeprivacy/cape/coordinator"
	"github.com/capeprivacy/cape/primitives"
	gm "github.com/onsi/gomega"
	"testing"
)

func TestProjectsCreate(t *testing.T) {
	gm.RegisterTestingT(t)

	p, err := primitives.NewProject("My Project", "my-project", "What is this project even about")
	gm.Expect(err).To(gm.BeNil())

	resp := coordinator.CreateProjectResponse{
		Project: p,
	}

	t.Run("Can create a project", func(t *testing.T) {
		app, u := NewHarness([]interface{}{resp})
		err = app.Run([]string{"cape", "projects", "create", p.Name.String(), p.Description.String()})
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(len(u.Calls)).To(gm.Equal(2))
		gm.Expect(u.Calls[0].Name).To(gm.Equal("template"))
		gm.Expect(u.Calls[0].Args[1]).To(gm.Equal(p.Name.String()))
	})

	t.Run("Can create a project without a description", func(t *testing.T) {
		app, u := NewHarness([]interface{}{resp})
		err = app.Run([]string{"cape", "projects", "create", p.Name.String()})
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(len(u.Calls)).To(gm.Equal(2))
		gm.Expect(u.Calls[0].Name).To(gm.Equal("template"))
		gm.Expect(u.Calls[0].Args[1]).To(gm.Equal(p.Name.String()))
	})

	t.Run("Must pass at least a name", func(t *testing.T) {
		app, _ := NewHarness([]interface{}{resp})
		err = app.Run([]string{"cape", "projects", "create"})
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(err.Error()).To(gm.Equal("missing_argument: The argument name is required, but was not provided"))
	})
}
