package main

import (
	"fmt"
	"github.com/capeprivacy/cape/models"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/capeprivacy/cape/coordinator/db"
)

func updateCmd(c *cli.Context) error {
	dbURL := EnvVariables(c.Context, capeDBURL).(*models.DBURL)

	migrations := c.Args().Slice()
	if len(migrations) == 0 {
		migrations = append(migrations, "coordinator/migrations")
	}

	for _, mPath := range migrations {
		if _, err := os.Stat(mPath); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("path %q not found", mPath)
			}

			return fmt.Errorf("error reading path %q: %s", mPath, err.Error())
		}
	}

	migrator, err := db.NewMigrator(dbURL.URL, migrations...)
	if err != nil {
		return err
	}

	return migrator.Up(c.Context)
}

func init() {
	updateCmd := &Command{
		Usage:     "Update a running Cape coordinator to a new version.",
		Variables: []*EnvVar{capeDBURL},
		Command: &cli.Command{
			Name:      "update",
			Action:    updateCmd,
			ArgsUsage: "[ path/to/migrations [ path/to/more/migrations ... ]]",
		},
	}

	commands = append(commands, updateCmd.Package())
}
