package main

import (
	"fmt"
	"net/url"

	"github.com/urfave/cli/v2"

	"github.com/dropoutlabs/cape/primitives"
)

func init() {
	sourcesAddCmd := &Command{
		Usage:       "Adds a data source to Cape",
		Description: "Adds a data source to Cape",
		Arguments:   []*Argument{LabelArg("source"), SourcesCredentialsArg},
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
			Flags:  []cli.Flag{clusterFlag()},
		},
	}

	sourcesRemoveCmd := &Command{
		Usage:       "Removes a data source to Cape",
		Description: "Removes a data source to Cape",
		Arguments:   []*Argument{LabelArg("source")},
		Examples: []*Example{
			{
				Example:     "cape sources remove transactions",
				Description: "Removes the database labelled `transactions` from cape",
			},
		},
		Command: &cli.Command{
			Name:   "remove",
			Action: sourcesRemove,
			Flags:  []cli.Flag{clusterFlag()},
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

	sourcesCmd := &Command{
		Usage: "Commands for adding, deleting, and modifying data sources",
		Command: &cli.Command{
			Name: "sources",
			Subcommands: []*cli.Command{
				sourcesAddCmd.Package(),
				sourcesRemoveCmd.Package(),
				sourcesListCmd.Package(),
			},
		},
	}

	commands = append(commands, sourcesCmd.Package())
}

func sourcesAdd(c *cli.Context) error {
	cfgSession := Session(c.Context)
	args := Arguments(c.Context)

	cluster, err := cfgSession.Cluster()
	if err != nil {
		return err
	}

	label := args["label"].(primitives.Label)
	credentials := args["connection-string"].(*url.URL)

	client, err := cluster.Client()
	if err != nil {
		return err
	}

	source, err := client.AddSource(c.Context, label, credentials)
	if err != nil {
		return err
	}

	fmt.Printf("Added source %s to Cape\n", source.Label)
	return nil
}

func sourcesRemove(c *cli.Context) error {
	cfgSession := Session(c.Context)
	args := Arguments(c.Context)

	cluster, err := cfgSession.Cluster()
	if err != nil {
		return err
	}

	label := args["label"].(primitives.Label)
	client, err := cluster.Client()
	if err != nil {
		return err
	}

	err = client.RemoveSource(c.Context, label)
	if err != nil {
		return err
	}

	fmt.Printf("Removed source %s from Cape\n", label)
	return nil
}

func sourcesList(c *cli.Context) error {
	cfgSession := Session(c.Context)
	cluster, err := cfgSession.Cluster()
	if err != nil {
		return err
	}

	client, err := cluster.Client()
	if err != nil {
		return err
	}

	sources, err := client.ListSources(c.Context)
	if err != nil {
		return err
	}

	ui := UI(c.Context)

	header := []string{"Name", "Type", "Host"}
	body := make([][]string, len(sources))
	for i, s := range sources {
		body[i] = []string{s.Label.String(), s.Endpoint.Scheme, s.Endpoint.String()}
	}

	return ui.Table(header, body)
}
