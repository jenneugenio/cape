// +build integration

package database

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"testing"

	gm "github.com/onsi/gomega"
)

// Integration test, testing a real query to a real postgres database
func TestPostgresQuery(t *testing.T) {
	gm.RegisterTestingT(t)

	host := os.Getenv("DBHOST")
	addr := fmt.Sprintf("postgres://postgres:dev@%s/postgres?sslmode=disable", host)
	u, err := url.Parse(addr)
	gm.Expect(err).To(gm.BeNil())

	pg := &PostgresBackend{dbURL: u}

	err = pg.Open(context.Background())
	gm.Expect(err).To(gm.BeNil())

	_, err = pg.db.Query(`
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
