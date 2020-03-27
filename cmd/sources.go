package main

import (
	"github.com/dropoutlabs/cape/primitives"
	"github.com/urfave/cli/v2"
	"net/url"
)

func init() {
	sourcesListCmd := &Command{
		Usage: "Lists all of your data sources",
		Command: &cli.Command{
			Name:   "list",
			Action: sourcesList,
		},
	}

	sourcesCmd := &Command{
		Usage: "Commands for adding, deleting, and modifying data sources",
		Command: &cli.Command{
			Name: "sources",
			Subcommands: []*cli.Command{
				sourcesListCmd.Package(),
			},
		},
	}

	commands = append(commands, sourcesCmd.Package())
}

func sourcesList(c *cli.Context) error {
	bigdata, err := url.Parse("postgres://bigdata.com/mydb")
	if err != nil {
		return err
	}

	smalldata, err := url.Parse("postgres://smalldata.com/mydb")
	if err != nil {
		return err
	}

	s1, err := primitives.NewSource("BIGDATA", *bigdata)
	if err != nil {
		return err
	}

	s2, err := primitives.NewSource("SMALLDATA", *smalldata)
	if err != nil {
		return err
	}

	ui := UI(c.Context)
	sources := []*primitives.Source{s1, s2}

	header := []string{"Label", "Credentials"}
	body := make([][]string, len(sources))
	for i, s := range sources {
		body[i] = []string{s.Label.String(), s.Credentials.String()}
	}

	return ui.Table(header, body)
}
