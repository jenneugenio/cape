package cmd

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update PrivacyAI",
	RunE: func(cmd *cobra.Command, args []string) error {
		dbAddr := os.Getenv("CAPE_DB_URL")

		db, err := sql.Open("postgres", dbAddr)
		if err != nil {
			return err
		}

		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			return err
		}

		m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
		if err != nil {
			return err
		}

		err = m.Up()

		if err == migrate.ErrNoChange {
			fmt.Println("No change, ignoring")
			return nil
		}

		return err
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
