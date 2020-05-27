package main

import (
	"github.com/urfave/cli/v2"

	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/primitives"
)

func updateCmd(c *cli.Context) error {
	dbURL := EnvVariables(c.Context, capeDBURL).(*primitives.DBURL)
	migrator, err := database.NewMigrator(dbURL.URL, "coordinator/migrations")
	if err != nil {
		return err
	}

	return migrator.Up(c.Context)
}

func init() {
	updateCmd := &Command{
		Usage:     "Update a running Cape coordinator to a new version",
		Variables: []*EnvVar{capeDBURL},
		Command: &cli.Command{
			Name:   "update",
			Action: updateCmd,
		},
	}

	commands = append(commands, updateCmd.Package())
}
