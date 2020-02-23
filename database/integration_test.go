// +build integration

package database

import (
	"context"
	"net/url"
	"os"
	"testing"

	gm "github.com/onsi/gomega"
)

// Integration test, testing a real query to a real postgres database
func TestPostgresQuery(t *testing.T) {
	gm.RegisterTestingT(t)

	dbURL := os.Getenv("CAPE_DB_URL")
	u, err := url.Parse(dbURL)
	gm.Expect(err).To(gm.BeNil())

	db, err := New(u)
	gm.Expect(err).To(gm.BeNil())

	err = db.Open(context.Background())
	gm.Expect(err).To(gm.BeNil())

	// TODO: Refactor this once we support _any_ database functionality at all
	_, err = db.(*PostgresBackend).db.Query(`
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
