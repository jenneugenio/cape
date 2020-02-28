package cmd

import (
	"context"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/tern/migrate"
	"github.com/spf13/cobra"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update PrivacyAI",
	RunE: func(cmd *cobra.Command, args []string) error {
		dbAddr := os.Getenv("CAPE_DB_URL")

		ctx := context.Background()
		conn, err := pgx.Connect(ctx, dbAddr)
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
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
