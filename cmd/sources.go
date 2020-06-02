package main

import (
	"fmt"
	"github.com/capeprivacy/cape/cmd/ui"
	"github.com/capeprivacy/cape/coordinator/client"
	"github.com/urfave/cli/v2"

	"github.com/capeprivacy/cape/coordinator/database"
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

	sourcesUpdateCmd := &Command{
		Usage:       "Modifies the configuration of an existing data source",
		Description: "Modifies the configuration of an existing data source",
		Arguments:   []*Argument{SourceLabelArg},
		Examples: []*Example{
			{
				Example: "cape sources update transactions --set-data-connector service:dc@my-cape.org",
				Description: "Modifies the configuration of the data source labelled `transactions` to link to " +
					"the data connector service identified by `service:dc@my-cape.org`",
			},
		},
		Command: &cli.Command{
			Name:   "update",
			Action: handleSessionOverrides(sourcesUpdate),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "set-data-connector",
					Usage:   "Link the source to the data connector `CONNECTOR`",
					EnvVars: []string{"CAPE_DATA_CONNECTOR"},
				},
				yesFlag(),
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
		Arguments:   []*Argument{SourceLabelArg, CollectionLabelArg},
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
				sourcesUpdateCmd.Package(),
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
	linkEmail := c.String("link")

	var serviceID *database.ID
	if linkEmail != "" {
		email, err := primitives.NewEmail(linkEmail)
		if err != nil {
			return err
		}

		service, err := client.GetServiceByEmail(c.Context, email)
		if err != nil {
			return err
		}

		serviceID = &service.ID
	}

	source, err := client.AddSource(c.Context, label, credentials, serviceID)
	if err != nil {
		return err
	}

	u := provider.UI(c.Context)
	return u.Template("Added source {{ . | bold }} to Cape\n", source.Label.String())
}

func sourcesUpdate(c *cli.Context) error {
	skipConfirm := c.Bool("yes")
	provider := GetProvider(c.Context)
	u := provider.UI(c.Context)
	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	label := Arguments(c.Context, SourceLabelArg).(primitives.Label)

	serviceEmail := c.String("set-data-connector")

	var serviceID *database.ID
	if serviceEmail != "" {
		email, err := primitives.NewEmail(serviceEmail)
		if err != nil {
			return err
		}

		service, err := client.GetServiceByEmail(c.Context, email)
		if err != nil {
			return err
		}

		serviceID = &service.ID
	} else {
		if !skipConfirm {
			err := u.Confirm(fmt.Sprintf("Do you really want to unlink the service for this source %s", label))
			if err != nil {
				return err
			}
		}
	}

	source, err := client.UpdateSource(c.Context, label, serviceID)
	if err != nil {
		return err
	}

	return u.Template("Updated source {{ . | bold }} with data connector {{ . | bold }}\n", source.Label.String())
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
		header := []string{"Name", "Type", "Host", "Data Connector"}
		body := make([][]string, len(sources))
		for i, s := range sources {
			body[i] = []string{s.Label.String(), s.Endpoint.Scheme, s.Endpoint.String(), ""}

			if s.Service != nil {
				body[i][3] = s.Service.Email.String()
			}
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
	connClient, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	label := Arguments(c.Context, SourceLabelArg).(primitives.Label)
	opts := &client.SourceOptions{WithSchema: true}

	collection, ok := Arguments(c.Context, CollectionLabelArg).(primitives.Label)
	if ok {
		opts.SchemaOptions = &client.SchemaOptions{Definition: collection.String()}
	}

	s, err := connClient.GetSourceByLabel(c.Context, label, opts)
	if err != nil {
		return err
	}

	u := provider.UI(c.Context)
	details := ui.Details{
		"Name": s.Label.String(),
		"Type": s.Endpoint.Scheme,
		"Host": s.Endpoint.String(),
	}

	if s.Service != nil {
		details["Data Connector"] = s.Service.Email.String()
	}

	err = u.Details(details)
	if err != nil {
		return err
	}

	template := "\n{{ \"Schema\" | bold }}\n"
	for tableName, table := range s.Schema.Definition {
		template += fmt.Sprintf("%s\n", tableName)
		for columnName, fieldType := range table {
			template += fmt.Sprintf("\t%s:\t%s\n", columnName, fieldType)
		}
		template += "\n"
	}

	return u.Template(template, nil)
}
