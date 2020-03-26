package main

import (
	"fmt"
	"github.com/dropoutlabs/cape/controller"
	"github.com/dropoutlabs/cape/primitives"
	"github.com/manifoldco/go-base64"
	"github.com/urfave/cli/v2"
	"net/url"
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
			Action: sourcesAdd,
			Flags:  []cli.Flag{clusterFlag()},
		},
	}

	sourcesListCmd := &Command{
		Usage: "Lists all of your data sources",
		Command: &cli.Command{
			Name:   "list",
			Action: sourcesList,
			Flags:  []cli.Flag{clusterFlag()},
		},
	}

	sourcesCmd := &Command{
		Usage: "Commands for adding, deleting, and modifying data sources",
		Command: &cli.Command{
			Name: "sources",
			Subcommands: []*cli.Command{
				sourcesAddCmd.Package(),
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

	clusterURL, err := url.Parse(cluster.URL)
	if err != nil {
		return err
	}

	token, err := base64.NewFromString(cluster.AuthToken)
	if err != nil {
		return err
	}

	client := controller.NewClient(clusterURL, token)
	source, err := client.AddSource(c.Context, label, credentials)
	if err != nil {
		return err
	}

	fmt.Printf("Added source %s to Cape\n", source.Label)
	return nil
}

func sourcesList(c *cli.Context) error {
	cfgSession := Session(c.Context)
	cluster, err := cfgSession.Cluster()

	if err != nil {
		return err
	}

	clusterURL, err := url.Parse(cluster.URL)
	if err != nil {
		return err
	}

	token, err := base64.NewFromString(cluster.AuthToken)
	if err != nil {
		return err
	}

	client := controller.NewClient(clusterURL, token)
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
