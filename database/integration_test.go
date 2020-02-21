// +build integration

package database

import (
	"context"
	"os"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/dropoutlabs/privacyai/database/dbtest"
)

// Integration test, testing a real query to a real postgres database
func TestPostgresQuery(t *testing.T) {
	gm.RegisterTestingT(t)
	ctx := context.Background()

	testDB, err := dbtest.New(os.Getenv("CAPE_DB_URL"))
	gm.Expect(err).To(gm.BeNil())

	err = testDB.Setup(ctx)
	defer testDB.Teardown(ctx)
	gm.Expect(err).To(gm.BeNil())

	// TODO: Refactor this once we support _any_ database functionality at all
	_, err = testDB.(*dbtest.Wrapper).Database().(*dbtest.TestPostgres).RawQuery(ctx, `
		SELECT
			*
		FROM
			pg_catalog.pg_tables
		WHERE
			schemaname != 'pg_catalog'
		AND schemaname != 'information_schema';
	`)
	gm.Expect(err).To(gm.BeNil())
}
