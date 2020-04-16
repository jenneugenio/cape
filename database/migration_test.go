// +build integration

package database

import (
	"context"
	"os"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/database/dbtest"

	"github.com/jackc/pgx/v4"
)

func TestMigrateUpAndDown(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()

	testDB, err := dbtest.NewTestPostgres(os.Getenv("CAPE_DB_URL"))
	gm.Expect(err).To(gm.BeNil())

	err = testDB.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer testDB.Teardown(ctx) // nolint: errcheck

	conn, err := pgx.Connect(ctx, testDB.URL().String())
	gm.Expect(err).To(gm.BeNil())

	defer conn.Close(ctx)

	migrator, err := NewMigrator(testDB.URL(), os.Getenv("CAPE_DB_MIGRATIONS"))
	gm.Expect(err).To(gm.BeNil())

	err = migrator.Up(ctx)
	gm.Expect(err).To(gm.BeNil())

	// Checks to see if anything exists in the database, ignores pg_toast
	// tables created by postgres and other system tables
	tableCountQuery := `SELECT COUNT(*) FROM pg_class c
						  JOIN pg_namespace s ON s.oid = c.relnamespace
	 					WHERE s.nspname NOT IN ('pg_catalog', 'information_schema')
						  AND s.nspname NOT LIKE 'pg_temp%' AND c.relname NOT LIKE 'pg_toast%';`

	rows, err := conn.Query(ctx, tableCountQuery)
	gm.Expect(err).To(gm.BeNil())
	defer rows.Close() // nolint: errcheck

	var count int
	rows.Next()
	err = rows.Scan(&count)
	gm.Expect(err).To(gm.BeNil())
	gm.Expect(count).ToNot(gm.Equal(0))

	// need to explicitly Close this here or the connection will remain busy
	rows.Close()

	err = migrator.Down(ctx)
	gm.Expect(err).To(gm.BeNil())

	rows, err = conn.Query(ctx, tableCountQuery)
	gm.Expect(err).To(gm.BeNil())
	defer rows.Close() // nolint: errcheck

	rows.Next()
	err = rows.Scan(&count)
	gm.Expect(err).To(gm.BeNil())

	// One is left here from the migrations tool
	gm.Expect(count).To(gm.Equal(1))
}
