// +build integration

package dbtest

import (
	"context"
	"os"
	"testing"

	gm "github.com/onsi/gomega"
)

func TestTestPostgres(t *testing.T) {

	ctx := context.Background()
	t.Run("creates a database that can be queried (and destroyed)", func(t *testing.T) {
		gm.RegisterTestingT(t)

		testDB, err := NewTestPostgres(os.Getenv("DB_URL"))
		gm.Expect(err).To(gm.BeNil())

		err = testDB.Setup(ctx)
		gm.Expect(err).To(gm.BeNil())

		_, err = testDB.(*TestPostgres).Exec(ctx, "CREATE TABLE hi (name char)")
		gm.Expect(err).To(gm.BeNil())

		err = testDB.Teardown(ctx)
		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("will truncate tables", func(t *testing.T) {
		gm.RegisterTestingT(t)

		testDB, err := NewTestPostgres(os.Getenv("DB_URL"))
		gm.Expect(err).To(gm.BeNil())

		err = testDB.Setup(ctx)
		gm.Expect(err).To(gm.BeNil())

		defer testDB.Teardown(ctx) //nolint: errcheck

		_, err = testDB.(*TestPostgres).Exec(ctx, "CREATE TABLE hi (name char(20))")
		gm.Expect(err).To(gm.BeNil())

		_, err = testDB.(*TestPostgres).Exec(ctx, "INSERT INTO hi (name) VALUES ($1)", "joe")
		gm.Expect(err).To(gm.BeNil())

		err = testDB.Truncate(ctx)
		gm.Expect(err).To(gm.BeNil())

		rows, err := testDB.(*TestPostgres).Query(ctx, "SELECT Count(1) FROM hi")
		defer rows.Close() //nolint: errcheck
		gm.Expect(err).To(gm.BeNil())

		count := 2
		for rows.Next() {
			err = rows.Scan(&count)
			gm.Expect(err).To(gm.BeNil())
		}

		gm.Expect(count).To(gm.Equal(0))
	})
}
