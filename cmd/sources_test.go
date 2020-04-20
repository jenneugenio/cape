package main

import (
	"github.com/capeprivacy/cape/cmd/ui"
	"github.com/capeprivacy/cape/coordinator"
	"github.com/capeprivacy/cape/primitives"
	gm "github.com/onsi/gomega"
	"testing"
)

func TestListSources(t *testing.T) {
	gm.RegisterTestingT(t)

	t.Run("Can list a single source", func(t *testing.T) {
		gm.RegisterTestingT(t)
		url, err := primitives.NewDBURL("postgres://localhost:5432/mydb")
		gm.Expect(err).To(gm.BeNil())

		resp := coordinator.ListSourcesResponse{
			Sources: []coordinator.SourceResponse{
				{
					Label:    "my-source-1",
					Type:     "postgres",
					Endpoint: url,
				},
			},
		}

		app, u := NewHarness([]interface{}{resp})
		err = app.Run([]string{"cape", "sources", "list"})
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(len(u.Calls)).To(gm.Equal(2))
		gm.Expect(u.Calls[0].Name).To(gm.Equal("table"))
		gm.Expect(u.Calls[0].Args[0]).To(gm.Equal(ui.TableHeader{"Name", "Type", "Host"}))
		gm.Expect(u.Calls[0].Args[1]).To(gm.Equal(ui.TableBody{{"my-source-1", "postgres", url.String()}}))

		gm.Expect(u.Calls[1].Name).To(gm.Equal("template"))
		gm.Expect(u.Calls[1].Args[0]).To(gm.Equal("\nFound {{ . | toString | faded }} source{{ . | pluralize \"s\"}}\n"))
		gm.Expect(u.Calls[1].Args[1]).To(gm.Equal(1))
	})

	t.Run("Can list multiple sources", func(t *testing.T) {
		gm.RegisterTestingT(t)
		url, err := primitives.NewDBURL("postgres://localhost:5432/mydb")
		gm.Expect(err).To(gm.BeNil())

		resp := coordinator.ListSourcesResponse{
			Sources: []coordinator.SourceResponse{
				{
					Label:    "my-source-1",
					Type:     "postgres",
					Endpoint: url,
				},
				{
					Label:    "my-source-2",
					Type:     "postgres",
					Endpoint: url,
				},
				{
					Label:    "my-source-3",
					Type:     "postgres",
					Endpoint: url,
				},
			},
		}

		app, u := NewHarness([]interface{}{resp})
		err = app.Run([]string{"cape", "sources", "list"})
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(len(u.Calls)).To(gm.Equal(2))
		gm.Expect(u.Calls[0].Name).To(gm.Equal("table"))
		gm.Expect(u.Calls[0].Args[0]).To(gm.Equal(ui.TableHeader{"Name", "Type", "Host"}))
		gm.Expect(u.Calls[0].Args[1]).To(gm.Equal(
			ui.TableBody{{"my-source-1", "postgres", url.String()},
				{"my-source-2", "postgres", url.String()},
				{"my-source-3", "postgres", url.String()}},
		))

		gm.Expect(u.Calls[1].Name).To(gm.Equal("template"))
		gm.Expect(u.Calls[1].Args[0]).To(gm.Equal("\nFound {{ . | toString | faded }} source{{ . | pluralize \"s\"}}\n"))
		gm.Expect(u.Calls[1].Args[1]).To(gm.Equal(3))
	})

	t.Run("Doesn't render a table if no sources are returned", func(t *testing.T) {
		gm.RegisterTestingT(t)

		resp := coordinator.ListSourcesResponse{
			Sources: []coordinator.SourceResponse{},
		}

		app, u := NewHarness([]interface{}{resp})
		err := app.Run([]string{"cape", "sources", "list"})
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(len(u.Calls)).To(gm.Equal(1))
		gm.Expect(u.Calls[0].Name).To(gm.Equal("template"))
		gm.Expect(u.Calls[0].Args[0]).To(gm.Equal("\nFound {{ . | toString | faded }} source{{ . | pluralize \"s\"}}\n"))
		gm.Expect(u.Calls[0].Args[1]).To(gm.Equal(0))
	})
}