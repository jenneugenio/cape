// +build integration

package sources

import (
	"context"
	"os"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/dropoutlabs/cape/connector/proto"
	"github.com/dropoutlabs/cape/database/dbtest"
	"github.com/dropoutlabs/cape/framework"
	"github.com/dropoutlabs/cape/primitives"
)

// TODO; We need to write the "error" flow tests for everything to do with the
// PostgresSource. For example, what happens if our backend returns an error?
func TestPostgresSource(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()

	db, err := dbtest.New(os.Getenv("CAPE_DB_URL"))
	gm.Expect(err).To(gm.BeNil())

	err = db.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer db.Teardown(ctx) // nolint: errcheck

	cfg := &Config{
		InstanceID: primitives.Label("cape-source-tester"),
		Logger:     framework.TestLogger(),
	}

	src, err := primitives.NewSource(primitives.Label("hello"), *db.URL(), nil)
	gm.Expect(err).To(gm.BeNil())

	t.Run("can create and close", func(t *testing.T) {
		source, err := NewPostgresSource(ctx, cfg, src)
		gm.Expect(err).To(gm.BeNil())

		err = source.Close(ctx)
		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("can get schema back for query", func(t *testing.T) {
		source, err := NewPostgresSource(ctx, cfg, src)
		gm.Expect(err).To(gm.BeNil())

		// TODO: Implement tests for checking whether or not schema is being
		// returned (and properly) for a given set of test data.
		//
		// It's imperative that this Schema is set to use the right types and
		// encoding values otherwise our clients will _not_ be able to decode
		// them.
		defer source.Close(ctx) // nolint: errcheck
	})

	t.Run("can stream rows back for query", func(t *testing.T) {
		source, err := NewPostgresSource(ctx, cfg, src)
		gm.Expect(err).To(gm.BeNil())

		defer source.Close(ctx) // nolint: errcheck

		// TODO: Implement tests for checking whether or not records are being
		// streamed back appropriately _and_ if the records contain the right
		// data / slash that we can unmarshal the data back into a format that
		// works for us.
		//
		// Things we need to assert:
		//
		//  - Schema is returned on the first record thats returned
		//  - The right number of records are returned
		//  - The records contain the data as we expected
		//
		// To do this we will need to setup a test database with test data. We
		// may want a static migration that we can rely on for writing these
		// tests that we store in a `testdata` folder.
		//
		// This test will probably need to call .Schema() first.
		stream := &testStream{}
		q := &testQuery{}
		schema := &proto.Schema{
			DataSource: q.Source().String(),
		}
		err = source.Query(ctx, q, schema, stream)
		gm.Expect(err).To(gm.BeNil())
	})
}
