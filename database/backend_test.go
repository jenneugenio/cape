package database

import (
	"context"
	"net/url"
	"testing"

	errors "github.com/dropoutlabs/privacyai/partyerrors"
	gm "github.com/onsi/gomega"
)

func TestBackend(t *testing.T) {
	gm.RegisterTestingT(t)
	t.Run("Invalid backend specified", func(t *testing.T) {
		u, err := url.Parse("fakedb://fake.db")
		gm.Expect(err).To(gm.BeNil())

		_, err = New(u)
		gm.Expect(errors.FromCause(err, NotImplementedDBCause)).To(gm.BeTrue())
	})

	t.Run("Valid backend specified", func(t *testing.T) {
		u, err := url.Parse("postgres://fake.db")
		gm.Expect(err).To(gm.BeNil())

		_, err = New(u)
		gm.Expect(err).To(gm.BeNil())
	})
}

// Integration test, testing a real query to a real postgres database
func TestPostgresQuery(t *testing.T) {
	gm.RegisterTestingT(t)

	addr := "postgres://postgres:dev@privacy-db-postgresql.default.svc.cluster.local/postgres?sslmode=disable"
	u, err := url.Parse(addr)
	gm.Expect(err).To(gm.BeNil())

	pg := &PostgresBackend{dbURL: u}

	err = pg.Open(context.Background())
	gm.Expect(err).To(gm.BeNil())

	_, err := pg.db.Query(`
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
