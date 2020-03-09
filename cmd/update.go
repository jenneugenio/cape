package cmd

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/tern/migrate"
	"github.com/urfave/cli/v2"
)

func updateCmd(c *cli.Context) error {
	dbURL := c.String("db-url")

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		return err
	}

	defer conn.Close(ctx)

	m, err := migrate.NewMigrator(ctx, conn, "migrations")
	if err != nil {
		return err
	}

	err = m.LoadMigrations("migrations")
	if err != nil {
		return err
	}

	return m.Migrate(ctx)
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
