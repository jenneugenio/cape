// +build integration

package database

import (
	"context"
	"os"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/coordinator/database/crypto"
	"github.com/capeprivacy/cape/coordinator/database/dbtest"
	"github.com/capeprivacy/cape/coordinator/database/types"
	errors "github.com/capeprivacy/cape/partyerrors"
)

func TestPostgresBackend(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()

	testDB, err := dbtest.New(os.Getenv("CAPE_DB_URL"))
	gm.Expect(err).To(gm.BeNil())

	migrations := []string{
		os.Getenv("CAPE_DB_MIGRATIONS"),
		os.Getenv("CAPE_DB_TEST_MIGRATIONS"),
	}

	migrator, err := NewMigrator(testDB.URL(), migrations...)
	gm.Expect(err).To(gm.BeNil())

	err = testDB.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	err = migrator.Up(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer testDB.Teardown(ctx) // nolint: errcheck

	kURL, err := crypto.NewBase64KeyURL(nil)
	gm.Expect(err).To(gm.BeNil())

	kms, err := crypto.LoadKMS(kURL)
	gm.Expect(err).To(gm.BeNil())

	defer kms.Close()

	codec := crypto.NewSecretBoxCodec(kms)

	t.Run("can create/retrieve an immutable entity", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB, codec)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		e, err := NewTestEntity("hello")
		gm.Expect(err).To(gm.BeNil())

		err = db.Create(ctx, e)
		gm.Expect(err).To(gm.BeNil())

		target := &TestEntity{}
		err = db.Get(ctx, e.GetID(), target)
		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("can create/retrieve an encryptable entity", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB, codec)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		e, err := NewTestEncryptionEntity("boo")
		gm.Expect(err).To(gm.BeNil())

		err = db.Create(ctx, e)
		gm.Expect(err).To(gm.BeNil())

		target := &TestEncryptionEntity{}
		err = db.Get(ctx, e.GetID(), target)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(target.ID).To(gm.Equal(e.ID))
		gm.Expect(target.Data).To(gm.Equal(e.Data))
	})

	t.Run("error if not encrypted", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB, codec)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		e, err := NewTestEncryptionEntity("boo")
		gm.Expect(err).To(gm.BeNil())

		db.(*PostgresBackend).codec = nil

		err = db.Create(ctx, e)
		gm.Expect(err).ToNot(gm.BeNil())
	})

	t.Run("can create multiple of same entity", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB, codec)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		eA, err := NewTestEntity("helloa")
		gm.Expect(err).To(gm.BeNil())

		eB, err := NewTestEntity("yod")
		gm.Expect(err).To(gm.BeNil())

		err = db.Create(ctx, eA, eB)
		gm.Expect(err).To(gm.BeNil())

		entities := []TestEntity{}
		f := Filter{Where: Where{"id": In{eA.ID.String(), eB.ID.String()}}}
		err = db.Query(ctx, &entities, f)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(entities).To(gm.Equal([]TestEntity{*eA, *eB}))
	})

	t.Run("can query encryptable entities", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB, codec)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		eA, err := NewTestEncryptionEntity("doo")
		gm.Expect(err).To(gm.BeNil())

		eB, err := NewTestEncryptionEntity("who")
		gm.Expect(err).To(gm.BeNil())

		err = db.Create(ctx, eA, eB)
		gm.Expect(err).To(gm.BeNil())

		var entities []*TestEncryptionEntity
		f := Filter{Where: Where{"id": In{eA.ID.String(), eB.ID.String()}}}
		err = db.Query(ctx, &entities, f)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(entities).To(gm.Equal([]*TestEncryptionEntity{eA, eB}))
	})

	t.Run("can't insert same entity twice", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB, codec)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		e, err := NewTestEntity("sup")
		gm.Expect(err).To(gm.BeNil())

		err = db.Create(ctx, e)
		gm.Expect(err).To(gm.BeNil())

		err = db.Create(ctx, e)
		gm.Expect(errors.FromCause(err, DuplicateCause)).To(gm.BeTrue())
	})

	t.Run("can't delete without providing an id", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB, codec)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		err = db.Delete(ctx, types.Test)
		gm.Expect(err).ToNot(gm.BeNil())
	})

	t.Run("can delete an entity", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB, codec)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		e, err := NewTestEntity("hi")
		gm.Expect(err).To(gm.BeNil())

		err = db.Create(ctx, e)
		gm.Expect(err).To(gm.BeNil())

		err = db.Delete(ctx, types.Test, e.GetID())
		gm.Expect(err).To(gm.BeNil())

		target := &TestEntity{}
		err = db.Get(ctx, e.GetID(), target)
		gm.Expect(errors.FromCause(err, NotFoundCause)).To(gm.BeTrue())
	})

	t.Run("can delete multiple entities", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB, codec)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		e, err := NewTestEntity("hi")
		gm.Expect(err).To(gm.BeNil())

		e2, err := NewTestEntity("hi2")
		gm.Expect(err).To(gm.BeNil())

		err = db.Create(ctx, e, e2)
		gm.Expect(err).To(gm.BeNil())

		err = db.Delete(ctx, types.Test, e.GetID(), e2.GetID())
		gm.Expect(err).To(gm.BeNil())

		target := []TestEntity{}
		err = db.Query(ctx, &target, NewFilter(Where{
			"id": In{e.GetID(), e2.GetID()},
		}, nil, nil))
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(len(target)).To(gm.Equal(0))
	})

	t.Run("cannot delete an entity from wrong table", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB, codec)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		e, err := NewTestEntity("blah blah")
		gm.Expect(err).To(gm.BeNil())

		err = db.Create(ctx, e)
		gm.Expect(err).To(gm.BeNil())

		// TestEntity has type Test so shouldn't be able to delete a TestEntity when the
		// specified table is TestMutable
		err = db.Delete(ctx, types.TestMutable, e.GetID())
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(errors.CausedBy(err, TypeMismatchCause)).To(gm.BeTrue())
	})

	t.Run("can't retrieve an unknown entity", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB, codec)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		e, err := NewTestEntity("hi")
		gm.Expect(err).To(gm.BeNil())

		target := &TestEntity{}
		err = db.Get(ctx, e.GetID(), target)
		gm.Expect(errors.FromCause(err, NotFoundCause)).To(gm.BeTrue())
	})

	t.Run("can update an entity", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB, codec)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		e, err := NewTestMutableEntity("sup")
		gm.Expect(err).To(gm.BeNil())

		err = db.Create(ctx, e)
		gm.Expect(err).To(gm.BeNil())

		previous := e.GetUpdatedAt()

		e.Data = "hello"
		err = db.Update(ctx, e)
		gm.Expect(err).To(gm.BeNil())

		target := &TestMutableEntity{}
		err = db.Get(ctx, e.GetID(), target)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(target.Data).To(gm.Equal(e.Data))
		gm.Expect(target.GetUpdatedAt().After(previous)).To(gm.BeTrue())
	})

	t.Run("can update an encryptable entity", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB, codec)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		e, err := NewTestEncryptionEntity("hehehe")
		gm.Expect(err).To(gm.BeNil())

		err = db.Create(ctx, e)
		gm.Expect(err).To(gm.BeNil())

		previous := e.GetUpdatedAt()

		e.Data = "hahaha"
		err = db.Update(ctx, e)
		gm.Expect(err).To(gm.BeNil())

		target := &TestEncryptionEntity{}
		err = db.Get(ctx, e.GetID(), target)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(target.Data).To(gm.Equal(e.Data))
		gm.Expect(target.GetUpdatedAt().After(previous)).To(gm.BeTrue())
	})

	t.Run("can't update an immutable entity", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB, codec)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		e, err := NewTestEntity("sup")
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(func() {
			err = db.Update(ctx, e)
		}).To(gm.Panic())
	})

	t.Run("can't update a non-existent entity", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB, codec)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		e, err := NewTestMutableEntity("sup")

		err = db.Update(ctx, e)
		gm.Expect(errors.FromCause(err, NotFoundCause)).To(gm.BeTrue())
	})

	t.Run("can run commands in a transaction", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB, codec)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		tx, err := db.Transaction(ctx)
		gm.Expect(err).To(gm.BeNil())
		defer tx.Rollback(ctx)

		e, err := NewTestMutableEntity("jack")
		gm.Expect(err).To(gm.BeNil())

		err = tx.Create(ctx, e)
		gm.Expect(err).To(gm.BeNil())

		e.Data = "joe"
		err = tx.Update(ctx, e)
		gm.Expect(err).To(gm.BeNil())

		err = tx.Commit(ctx)
		gm.Expect(err).To(gm.BeNil())

		target := &TestMutableEntity{}
		err = db.Get(ctx, e.GetID(), target)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(target.Data).To(gm.Equal("joe"))
	})

	t.Run("can rollback", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB, codec)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		e, err := NewTestMutableEntity("jack")
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

		target := &TestMutableEntity{}
		err = db.Get(ctx, e.GetID(), target)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(target.Data).To(gm.Equal("jack"))
	})

	t.Run("rollback after commit causes an error", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB, codec)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		tx, err := db.Transaction(ctx)
		gm.Expect(err).To(gm.BeNil())
		defer tx.Rollback(ctx)

		e, err := NewTestMutableEntity("jack")
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

		target := &TestMutableEntity{}
		err = db.Get(ctx, e.GetID(), target)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(target.Data).To(gm.Equal("joe"))
	})

	t.Run("can query a single entity", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB, codec)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		eA, err := NewTestEntity("a")
		gm.Expect(err).To(gm.BeNil())

		eB, err := NewTestEntity("b")
		gm.Expect(err).To(gm.BeNil())

		err = db.Create(ctx, eA)
		gm.Expect(err).To(gm.BeNil())

		err = db.Create(ctx, eB)
		gm.Expect(err).To(gm.BeNil())

		target := &TestEntity{}
		err = db.QueryOne(ctx, target, NewFilter(Where{"data": "a"}, nil, nil))
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(eA).To(gm.Equal(target))

		targetTwo := &TestEntity{}
		filter := NewFilter(Where{"id": eB.GetID().String()}, nil, nil)
		err = db.QueryOne(ctx, targetTwo, filter)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(eB).To(gm.Equal(targetTwo))
	})

	t.Run("QueryOne returns not found if in is empty", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB, codec)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		target := &TestEntity{}
		err = db.QueryOne(ctx, target, NewFilter(Where{"id": In{}}, nil, nil))
		gm.Expect(errors.FromCause(err, NotFoundCause)).To(gm.BeTrue())
	})

	t.Run("QueryOne returns a not-found error if the entity can't be found", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB, codec)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		eA, err := NewTestEntity("idontexistok")
		gm.Expect(err).To(gm.BeNil())

		target := &TestEntity{}
		err = db.QueryOne(ctx, target, NewFilter(Where{"id": eA.ID}, nil, nil))
		gm.Expect(errors.FromCause(err, NotFoundCause)).To(gm.BeTrue())
	})

	t.Run("Query does no-op if in is empty", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB, codec)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		target := []*TestEntity{}
		err = db.Query(ctx, &target, NewFilter(Where{"id": In{}}, nil, nil))
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(len(target)).To(gm.Equal(0))
	})

	t.Run("can query using slice of pointers", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB, codec)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		eA, err := NewTestEntity("a1b")
		gm.Expect(err).To(gm.BeNil())

		eB, err := NewTestEntity("b1c")
		gm.Expect(err).To(gm.BeNil())

		eC, err := NewTestEntity("c1d")
		gm.Expect(err).To(gm.BeNil())

		err = db.Create(ctx, eA, eB, eC)
		gm.Expect(err).To(gm.BeNil())

		filter := NewFilter(Where{"id": In{
			eA.GetID().String(),
			eB.GetID().String(),
			eC.GetID().String(),
		}}, nil, nil)
		result := []*TestEntity{}

		err = db.Query(ctx, &result, filter)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(len(result)).To(gm.Equal(3))
		gm.Expect(result[0].ID).To(gm.Equal(eA.ID))
		gm.Expect(result[1].ID).To(gm.Equal(eB.ID))
		gm.Expect(result[2].ID).To(gm.Equal(eC.ID))
	})

	t.Run("can query multiple entities of the same type", func(t *testing.T) {
		db, err := dbConnect(ctx, testDB, codec)
		gm.Expect(err).To(gm.BeNil())
		defer db.Close()

		eA, err := NewTestEntity("a1")
		gm.Expect(err).To(gm.BeNil())

		eB, err := NewTestEntity("b1")
		gm.Expect(err).To(gm.BeNil())

		eC, err := NewTestEntity("c1")
		gm.Expect(err).To(gm.BeNil())

		err = db.Create(ctx, eA, eB, eC)
		gm.Expect(err).To(gm.BeNil())

		tests := map[string]struct {
			f   Filter
			out []TestEntity
		}{
			"can pull back single by id": {
				NewFilter(Where{"id": eA.GetID().String()}, nil, nil),
				[]TestEntity{*eA},
			},
			"can pull back using comparison": {
				NewFilter(Where{"data": "a1"}, nil, nil),
				[]TestEntity{*eA},
			},
			"can pull back using IN operator": {
				NewFilter(Where{"id": In{eA.ID.String(), eB.ID.String()}}, nil, nil),
				[]TestEntity{*eA, *eB},
			},
			"can order via a field": {
				NewFilter(Where{"id": In{eA.ID.String(), eB.ID.String(), eC.ID.String()}},
					&Order{Desc, "data"}, nil),
				[]TestEntity{*eC, *eB, *eA},
			},
			"can order and paginate": {
				NewFilter(Where{"id": In{eA.ID.String(), eB.ID.String(), eC.ID.String()}},
					&Order{Desc, "data"}, &Page{1, 1}),
				[]TestEntity{*eB},
			},
		}

		for d, tc := range tests {
			t.Run(d, func(t *testing.T) {
				results := []TestEntity{}
				err = db.Query(ctx, &results, tc.f)
				gm.Expect(err).To(gm.BeNil())
				gm.Expect(results).To(gm.BeEquivalentTo(tc.out))
			})
		}
	})
}

func dbConnect(ctx context.Context, t dbtest.TestDatabase, codec crypto.EncryptionCodec) (Backend, error) {
	db, err := New(t.URL(), "testing")
	if err != nil {
		return nil, err
	}

	db.SetEncryptionCodec(codec)

	err = db.Open(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
