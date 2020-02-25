// +build integration

package database

import (
	"context"
	"fmt"
	"os"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/dropoutlabs/privacyai/database/dbtest"
	errors "github.com/dropoutlabs/privacyai/partyerrors"
	"github.com/dropoutlabs/privacyai/primitives"
)

func TestPostgresBackend(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()

	testDB, err := dbtest.New(os.Getenv("CAPE_DB_URL"))
	gm.Expect(err).To(gm.BeNil())

	err = testDB.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer testDB.Teardown(ctx) // nolint: errcheck

	err = setupMigrations(ctx, testDB)
	gm.Expect(err).To(gm.BeNil())

	t.Run("can create/retrieve an immutable entity", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		e, err := primitives.NewTestEntity("hello")
		gm.Expect(err).To(gm.BeNil())

		err = db.Create(ctx, e)
		gm.Expect(err).To(gm.BeNil())

		target := &primitives.TestEntity{}
		err = db.Get(ctx, e.GetID(), target)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(target).To(gm.Equal(e))
	})

	t.Run("can't insert same entity twice", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		e, err := primitives.NewTestEntity("sup")
		gm.Expect(err).To(gm.BeNil())

		err = db.Create(ctx, e)
		gm.Expect(err).To(gm.BeNil())

		err = db.Create(ctx, e)
		gm.Expect(errors.FromCause(err, DuplicateCause)).To(gm.BeTrue())
	})

	t.Run("can delete an entity", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		e, err := primitives.NewTestEntity("hi")
		gm.Expect(err).To(gm.BeNil())

		err = db.Create(ctx, e)
		gm.Expect(err).To(gm.BeNil())

		err = db.Delete(ctx, e.GetID())
		gm.Expect(err).To(gm.BeNil())

		target := &primitives.TestEntity{}
		err = db.Get(ctx, e.GetID(), target)
		gm.Expect(errors.FromCause(err, NotFoundCause)).To(gm.BeTrue())
	})

	t.Run("can't retrieve an unknown entity", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		e, err := primitives.NewTestEntity("hi")
		gm.Expect(err).To(gm.BeNil())

		target := &primitives.TestEntity{}
		err = db.Get(ctx, e.GetID(), target)
		gm.Expect(errors.FromCause(err, NotFoundCause)).To(gm.BeTrue())
	})

	t.Run("can update an entity", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		e, err := primitives.NewTestMutableEntity("sup")
		gm.Expect(err).To(gm.BeNil())

		err = db.Create(ctx, e)
		gm.Expect(err).To(gm.BeNil())

		previous := e.GetUpdatedAt()

		e.Data = "hello"
		err = db.Update(ctx, e)
		gm.Expect(err).To(gm.BeNil())

		target := &primitives.TestMutableEntity{}
		err = db.Get(ctx, e.GetID(), target)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(target.Data).To(gm.Equal(e.Data))
		gm.Expect(target.GetUpdatedAt().After(previous)).To(gm.BeTrue())
	})

	t.Run("can't update an immutable entity", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		e, err := primitives.NewTestEntity("sup")
		gm.Expect(err).To(gm.BeNil())

		err = db.Update(ctx, e)
		gm.Expect(errors.FromCause(err, NotMutableCause)).To(gm.BeTrue())
	})

	t.Run("can't update a non-existent entity", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		e, err := primitives.NewTestMutableEntity("sup")

		err = db.Update(ctx, e)
		gm.Expect(errors.FromCause(err, NotFoundCause)).To(gm.BeTrue())
	})
}

func dbConnect(ctx context.Context, t dbtest.TestDatabase) (Backend, error) {
	db, err := New(t.URL(), "testing")
	if err != nil {
		return nil, err
	}

	err = db.Open(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func setupMigrations(ctx context.Context, db dbtest.TestDatabase) error {
	pg, ok := db.(*dbtest.Wrapper).Database().(*dbtest.TestPostgres) // this is all throw away once we have migrations
	if !ok {
		return errors.New(errors.UnsupportedErrorCause, "dbtest must be a TestPostgres")
	}

	_, err := pg.Exec(ctx, baseMigrationSQL())
	if err != nil {
		return err
	}

	_, err = pg.Exec(ctx, deriveMigrationSQL("test"))
	if err != nil {
		return err
	}

	_, err = pg.Exec(ctx, deriveMigrationSQL("test_mutable"))
	return err
}

func baseMigrationSQL() string {
	return `
		CREATE EXTENSION IF NOT EXISTS hstore;

		CREATE FUNCTION hoist_values() RETURNS TRIGGER AS $$
			DECLARE
				value hstore;
				paths text[];
				path text;
				segments text[];
				segment text;
			BEGIN
				value = hstore(NEW);
				paths = TG_ARGV;

				FOREACH path IN ARRAY paths LOOP
					segments = string_to_array(path, '.')::text[];
					segment = segments[array_upper(segments, 1)];

					value := value || hstore(segment, NEW.data::jsonb#>>segments);
					NEW := NEW #= value;
				END LOOP;

				RETURN NEW;
			END;
		$$ LANGUAGE plpgsql;
	`
}

func deriveMigrationSQL(name string) string {
	return fmt.Sprintf(`
		CREATE TABLE %s (
			id char(29) not null primary key,
			data jsonb not null,

			CONSTRAINT %s_id_equals CHECK (data::jsonb#>>'{id}' = id)
		);

		CREATE TRIGGER %s_hoist_tgr
			BEFORE INSERT ON %s
			FOR EACH ROW EXECUTE PROCEDURE hoist_values('id');
	`, name, name, name, name)
}
