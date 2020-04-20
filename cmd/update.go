package main

import (
	"github.com/urfave/cli/v2"

	"github.com/capeprivacy/cape/database"
)

func updateCmd(c *cli.Context) error {
	dbURL, err := getDBURL(c)
	if err != nil {
		return err
	}

	migrator, err := database.NewMigrator(dbURL.URL, "migrations")
	if err != nil {
		return err
	}

	return migrator.Up(c.Context)
}

func init() {
	updateCmd := &Command{
		Usage:     "Update a running Cape coordinator to a new version",
		Variables: []*EnvVar{capeDBPassword, capeDBURL},
		Command: &cli.Command{
			Name:   "update",
			Action: updateCmd,
		},
	}

	commands = append(commands, updateCmd.Package())
}
