package main

import (
	"net/url"

	"github.com/urfave/cli/v2"

	"github.com/dropoutlabs/cape/database"
)

func updateCmd(c *cli.Context) error {
	dbAddr := c.String("db-url")

	dbURL, err := url.Parse(dbAddr)
	if err != nil {
		return err
	}
	migrator, err := database.NewMigrator(dbURL, "migrations")
	if err != nil {
		return err
	}

	return migrator.Up(c.Context)
}

func init() {
	updateCmd := &Command{
		Usage: "Update a running Cape controller to a new version",
		Command: &cli.Command{
			Name:   "update",
			Action: updateCmd,
			Flags:  []cli.Flag{dbURLFlag()},
		},
	}

	commands = append(commands, updateCmd.Package())
}
