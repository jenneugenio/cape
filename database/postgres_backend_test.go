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
	})

	t.Run("can create multiple of same entity", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		eA, err := primitives.NewTestEntity("helloa")
		gm.Expect(err).To(gm.BeNil())

		eB, err := primitives.NewTestEntity("yod")
		gm.Expect(err).To(gm.BeNil())

		err = db.Create(ctx, eA, eB)
		gm.Expect(err).To(gm.BeNil())

		entities := []primitives.TestEntity{}
		f := Filter{Where: Where{"id": In{eA.ID.String(), eB.ID.String()}}}
		err = db.Query(ctx, &entities, f)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(entities).To(gm.Equal([]primitives.TestEntity{*eA, *eB}))
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

		gm.Expect(func() {
			err = db.Update(ctx, e)
		}).To(gm.Panic())
	})

	t.Run("can't update a non-existent entity", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		e, err := primitives.NewTestMutableEntity("sup")

		err = db.Update(ctx, e)
		gm.Expect(errors.FromCause(err, NotFoundCause)).To(gm.BeTrue())
	})

	t.Run("can run commands in a transaction", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		tx, err := db.Transaction(ctx)
		gm.Expect(err).To(gm.BeNil())
		defer tx.Rollback(ctx)

		e, err := primitives.NewTestMutableEntity("jack")
		gm.Expect(err).To(gm.BeNil())

		err = tx.Create(ctx, e)
		gm.Expect(err).To(gm.BeNil())

		e.Data = "joe"
		err = tx.Update(ctx, e)
		gm.Expect(err).To(gm.BeNil())

		err = tx.Commit(ctx)
		gm.Expect(err).To(gm.BeNil())

		target := &primitives.TestMutableEntity{}
		err = db.Get(ctx, e.GetID(), target)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(target.Data).To(gm.Equal("joe"))
	})

	t.Run("can rollback", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		e, err := primitives.NewTestMutableEntity("jack")
		gm.Expect(err).To(gm.BeNil())

		err = db.Create(ctx, e)
		gm.Expect(err).To(gm.BeNil())

		tx, err := db.Transaction(ctx)
		gm.Expect(err).To(gm.BeNil())
		defer tx.Rollback(ctx)

		e.Data = "joe"
		err = tx.Update(ctx, e)
		gm.Expect(err).To(gm.BeNil())

		err = tx.Rollback(ctx)
		gm.Expect(err).To(gm.BeNil())

		target := &primitives.TestMutableEntity{}
		err = db.Get(ctx, e.GetID(), target)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(target.Data).To(gm.Equal("jack"))
	})

	t.Run("rollback after commit causes an error", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		tx, err := db.Transaction(ctx)
		gm.Expect(err).To(gm.BeNil())
		defer tx.Rollback(ctx)

		e, err := primitives.NewTestMutableEntity("jack")
		gm.Expect(err).To(gm.BeNil())

		err = tx.Create(ctx, e)
		gm.Expect(err).To(gm.BeNil())

		e.Data = "joe"
		err = tx.Update(ctx, e)
		gm.Expect(err).To(gm.BeNil())

		err = tx.Commit(ctx)
		gm.Expect(err).To(gm.BeNil())

		err = tx.Rollback(ctx)
		gm.Expect(errors.FromCause(err, ClosedCause)).To(gm.BeTrue())

		target := &primitives.TestMutableEntity{}
		err = db.Get(ctx, e.GetID(), target)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(target.Data).To(gm.Equal("joe"))
	})

	t.Run("can query a single entity", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		eA, err := primitives.NewTestEntity("a")
		gm.Expect(err).To(gm.BeNil())

		eB, err := primitives.NewTestEntity("b")
		gm.Expect(err).To(gm.BeNil())

		err = db.Create(ctx, eA)
		gm.Expect(err).To(gm.BeNil())

		err = db.Create(ctx, eB)
		gm.Expect(err).To(gm.BeNil())

		target := &primitives.TestEntity{}
		err = db.QueryOne(ctx, target, NewFilter(Where{"data": "a"}, nil, nil))
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(eA).To(gm.Equal(target))

		targetTwo := &primitives.TestEntity{}
		filter := NewFilter(Where{"id": eB.GetID().String()}, nil, nil)
		err = db.QueryOne(ctx, targetTwo, filter)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(eB).To(gm.Equal(targetTwo))
	})

	t.Run("can query multiple entities of the same type", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		eA, err := primitives.NewTestEntity("a1")
		gm.Expect(err).To(gm.BeNil())

		eB, err := primitives.NewTestEntity("b1")
		gm.Expect(err).To(gm.BeNil())

		eC, err := primitives.NewTestEntity("c1")
		gm.Expect(err).To(gm.BeNil())

		err = db.Create(ctx, eA, eB, eC)
		gm.Expect(err).To(gm.BeNil())

		tests := map[string]struct {
			f   Filter
			out []primitives.TestEntity
		}{
			"can pull back single by id": {
				NewFilter(Where{"id": eA.GetID().String()}, nil, nil),
				[]primitives.TestEntity{*eA},
			},
			"can pull back using comparison": {
				NewFilter(Where{"data": "a1"}, nil, nil),
				[]primitives.TestEntity{*eA},
			},
			"can pull back using IN operator": {
				NewFilter(Where{"id": In{eA.ID.String(), eB.ID.String()}}, nil, nil),
				[]primitives.TestEntity{*eA, *eB},
			},
			"can order via a field": {
				NewFilter(Where{"id": In{eA.ID.String(), eB.ID.String(), eC.ID.String()}},
					&Order{Desc, "data"}, nil),
				[]primitives.TestEntity{*eC, *eB, *eA},
			},
			"can order and paginate": {
				NewFilter(Where{"id": In{eA.ID.String(), eB.ID.String(), eC.ID.String()}},
					&Order{Desc, "data"}, &Page{1, 1}),
				[]primitives.TestEntity{*eB},
			},
		}

		for d, tc := range tests {
			t.Run(d, func(t *testing.T) {
				results := []primitives.TestEntity{}
				err = db.Query(ctx, &results, tc.f)
				gm.Expect(err).To(gm.BeNil())
				gm.Expect(results).To(gm.BeEquivalentTo(tc.out))
			})
		}
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

	_, err := pg.Exec(ctx, deriveMigrationSQL("test"))
	if err != nil {
		return err
	}

	_, err = pg.Exec(ctx, deriveMigrationSQL("test_mutable"))
	return err
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
