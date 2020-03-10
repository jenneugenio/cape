package cmd

import (
	"context"
	"github.com/dropoutlabs/cape/database"
	"github.com/urfave/cli/v2"
	"net/url"
)

func updateCmd(c *cli.Context) error {
	dbAddr := c.String("db-url")

	ctx := context.Background()

	dbURL, err := url.Parse(dbAddr)
	if err != nil {
		return err
	}
	migrator, err := database.NewMigrator(dbURL, "migrations")
	if err != nil {
		return err
	}

	return migrator.Up(ctx)
}

func init() {
	updateCmd := &cli.Command{
		Name:        "update",
		Description: "Update Cape Controller database schema version",
		Action:      updateCmd,
		Flags:       []cli.Flag{dbURLFlag()},
	}

	commands = append(commands, updateCmd)
}
