package main

import (
	"github.com/capeprivacy/cape/coordinator"
	"github.com/capeprivacy/cape/models"
	gm "github.com/onsi/gomega"
	"testing"
)

func TestRolesMe(t *testing.T) {
	gm.RegisterTestingT(t)

	resp := coordinator.MyRoleResponse{
		Role: models.Role{
			ID:    "1234",
			Label: models.AdminRole,
		},
	}

	t.Run("Global Check", func(t *testing.T) {
		app, _ := NewHarness([]*coordinator.MockResponse{
			{
				Value: resp,
			},
		})
		err := app.Run([]string{"cape", "roles", "me"})
		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("Project Check", func(t *testing.T) {
		app, _ := NewHarness([]*coordinator.MockResponse{
			{
				Value: resp,
			},
		})
		err := app.Run([]string{"cape", "roles", "me", "--project", "my-project"})
		gm.Expect(err).To(gm.BeNil())
	})
}

func TestRolesSet(t *testing.T) {
	gm.RegisterTestingT(t)

	t.Run("Global Set", func(t *testing.T) {
		app, _ := NewHarness([]*coordinator.MockResponse{
			{
				Value: nil,
			},
		})
		err := app.Run([]string{"cape", "roles", "set", "person@website.com", "admin"})
		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("Global set with bad args", func(t *testing.T) {
		app, _ := NewHarness([]*coordinator.MockResponse{
			{
				Value: nil,
			},
		})
		err := app.Run([]string{"cape", "roles", "set"})
		gm.Expect(err).ToNot(gm.BeNil())
	})

	t.Run("Global set with fake role", func(t *testing.T) {
		app, _ := NewHarness([]*coordinator.MockResponse{
			{
				Value: nil,
			},
		})
		err := app.Run([]string{"cape", "roles", "set", "whahahdhshdashdsajkdhsa"})
		gm.Expect(err).ToNot(gm.BeNil())
	})

	t.Run("Project set", func(t *testing.T) {
		app, _ := NewHarness([]*coordinator.MockResponse{
			{
				Value: nil,
			},
		})
		err := app.Run([]string{"cape", "roles", "set", "--project", "my-project", "person@website.com", "project-owner"})
		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("Project set with bad args", func(t *testing.T) {
		app, _ := NewHarness([]*coordinator.MockResponse{
			{
				Value: nil,
			},
		})
		err := app.Run([]string{"cape", "roles", "set", "--project", "my-project", "person@website.com"})
		gm.Expect(err).ToNot(gm.BeNil())
	})

	t.Run("Project set with fake role", func(t *testing.T) {
		app, _ := NewHarness([]*coordinator.MockResponse{
			{
				Value: nil,
			},
		})
		err := app.Run([]string{"cape", "roles", "set", "--project", "my-project", "person@website.com", "projectzzzzzz-owner"})
		gm.Expect(err).ToNot(gm.BeNil())
	})
}
