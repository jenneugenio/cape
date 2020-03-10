package cmd

import (
	"context"
	"github.com/dropoutlabs/cape/database"
	"net/url"
	"os"

	"github.com/urfave/cli/v2"
)

func updateCmd(c *cli.Context) error {
	dbAddr := c.String("db-url")

	ctx := context.Background()

	dbURL, err := url.Parse(dbAddr)
	if err != nil {
		return err
	}
	migrator, err := database.NewMigrator(dbURL, os.Getenv("CAPE_DB_MIGRATIONS"))
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
