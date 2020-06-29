package main

import (
	"fmt"

	"github.com/capeprivacy/cape/cmd/ui"
	"github.com/urfave/cli/v2"

	"github.com/capeprivacy/cape/primitives"
)

func init() {
	sourcesAddCmd := &Command{
		Usage:       "Adds a data source to Cape",
		Description: "Adds a data source to Cape",
		Arguments:   []*Argument{SourceLabelArg, SourcesCredentialsArg},
		Examples: []*Example{
			{
				Example: "cape sources add transactions postgres://username:password@location.of.database:5432/mydb",
				Description: "Adds the database `mydb` from the postgres database running at `location.of.database` " +
					"and gives it the label transactions:",
			},
		},
		Command: &cli.Command{
			Name:   "add",
			Action: handleSessionOverrides(sourcesAdd),
			Flags: []cli.Flag{
				clusterFlag(),
				linkFlag(),
			},
		},
	}

	sourcesRemoveCmd := &Command{
		Usage:       "Removes a data source to Cape",
		Description: "Removes a data source to Cape",
		Arguments:   []*Argument{SourceLabelArg},
		Examples: []*Example{
			{
				Example:     "cape sources remove transactions",
				Description: "Removes the database labelled `transactions` from cape",
			},
		},
		Command: &cli.Command{
			Name:   "remove",
			Action: handleSessionOverrides(sourcesRemove),
			Flags:  []cli.Flag{clusterFlag(), yesFlag()},
		},
	}

	sourcesListCmd := &Command{
		Usage: "Lists all of your data sources",
		Command: &cli.Command{
			Name:   "list",
			Action: handleSessionOverrides(sourcesList),
			Flags:  []cli.Flag{clusterFlag()},
		},
	}

	sourcesDescribeCmd := &Command{
		Usage:       "Describes a data source",
		Description: "Provides addition information about a data source, such as its schema",
		Arguments:   []*Argument{SourceLabelArg},
		Examples: []*Example{
			{
				Example:     "cape sources describe transactions",
				Description: "Describes the source labelled `transactios`",
			},
		},
		Command: &cli.Command{
			Name:   "describe",
			Action: handleSessionOverrides(sourcesDescribe),
			Flags:  []cli.Flag{clusterFlag()},
		},
	}

	sourcesCmd := &Command{
		Usage: "Commands for adding, deleting, and modifying data sources",
		Command: &cli.Command{
			Name: "sources",
			Subcommands: []*cli.Command{
				sourcesAddCmd.Package(),
				sourcesRemoveCmd.Package(),
				sourcesListCmd.Package(),
				sourcesDescribeCmd.Package(),
			},
		},
	}

	commands = append(commands, sourcesCmd.Package())
}

func sourcesAdd(c *cli.Context) error {
	provider := GetProvider(c.Context)
	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	label := Arguments(c.Context, SourceLabelArg).(primitives.Label)
	credentials := Arguments(c.Context, SourcesCredentialsArg).(*primitives.DBURL)

	source, err := client.AddSource(c.Context, label, credentials)
	if err != nil {
		return err
	}

	u := provider.UI(c.Context)
	return u.Template("Added source {{ . | bold }} to Cape\n", source.Label.String())
}

func sourcesRemove(c *cli.Context) error {
	skipConfirm := c.Bool("yes")
	provider := GetProvider(c.Context)
	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}
	u := provider.UI(c.Context)

	label := Arguments(c.Context, SourceLabelArg).(primitives.Label)

	if !skipConfirm {
		err := u.Confirm(fmt.Sprintf("Do you really want to delete the source %s", label))
		if err != nil {
			return err
		}
	}

	err = client.RemoveSource(c.Context, label)
	if err != nil {
		return err
	}

	return u.Template("Removed source {{ . | bold }} from Cape\n", label.String())
}

func sourcesList(c *cli.Context) error {
	provider := GetProvider(c.Context)
	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	sources, err := client.ListSources(c.Context)
	if err != nil {
		return err
	}

	u := provider.UI(c.Context)

	if len(sources) > 0 {
		header := []string{"Name", "Type", "Host"}
		body := make([][]string, len(sources))
		for i, s := range sources {
			body[i] = []string{s.Label.String(), s.Endpoint.Scheme, s.Endpoint.String()}
		}

		err = u.Table(header, body)
		if err != nil {
			return err
		}
	}

	return u.Template("\nFound {{ . | toString | faded }} source{{ . | pluralize \"s\"}}\n", len(sources))
}

func sourcesDescribe(c *cli.Context) error {
	provider := GetProvider(c.Context)
	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	label := Arguments(c.Context, SourceLabelArg).(primitives.Label)
	s, err := client.GetSourceByLabel(c.Context, label)
	if err != nil {
		return err
	}

	u := provider.UI(c.Context)
	details := ui.Details{
		"Name": s.Label.String(),
		"Type": s.Endpoint.Scheme,
		"Host": s.Endpoint.String(),
	}

	return u.Details(details)
}
