package auth

import (
	"context"
	"fmt"
	"testing"

	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/types"
	"github.com/capeprivacy/cape/models"
	"github.com/capeprivacy/cape/primitives"
)

var ErrGenericError = fmt.Errorf("error")

var (
	NilCanner = TestCanner{nil}
	ErrCanner = TestCanner{ErrGenericError}
)

func ErrNil(e error) bool    { return e == nil }
func ErrNotNil(e error) bool { return e != nil }

type TestCanner struct {
	e error
}

func (tc TestCanner) Can(_ models.Action, _ types.Type) error { return tc.e }

type TestQuerier struct {
	err error
}

func (t TestQuerier) Create(context.Context, ...database.Entity) error         { return t.err }
func (t TestQuerier) Get(context.Context, database.ID, database.Entity) error  { return t.err }
func (t TestQuerier) Delete(context.Context, types.Type, ...database.ID) error { return t.err }
func (t TestQuerier) Upsert(context.Context, database.Entity) error            { return t.err }
func (t TestQuerier) Update(context.Context, database.Entity) error            { return t.err }
func (t TestQuerier) QueryOne(context.Context, database.Entity, database.Filter) error {
	return t.err
}
func (t TestQuerier) Query(context.Context, interface{}, database.Filter) error { return t.err }

var TestEnt = func() database.Entity {
	e, _ := database.NewTestEntity("r2d2")
	return e
}()

var SingleTestEntList = []database.Entity{TestEnt}

func TestCreate(t *testing.T) {
	tests := []struct {
		name     string
		canErr   error
		q        database.Querier
		ctx      context.Context
		ent      []database.Entity
		wantRet  error
		wantFunc func(error) bool
	}{
		{
			name:     "fails if session disallows",
			canErr:   ErrGenericError,
			q:        TestQuerier{nil},
			ctx:      context.TODO(),
			ent:      SingleTestEntList,
			wantFunc: ErrNotNil,
		},
		{
			name:     "fails if entity list is empty",
			canErr:   nil,
			q:        TestQuerier{nil},
			ctx:      context.TODO(),
			ent:      []database.Entity{},
			wantFunc: ErrNotNil,
		},
		{
			name:     "fails if db query fails",
			canErr:   nil,
			q:        TestQuerier{ErrGenericError},
			ctx:      context.TODO(),
			ent:      SingleTestEntList,
			wantFunc: ErrNotNil,
		},
		{
			name:     "succeeds if session can and db query succeeds",
			canErr:   nil,
			q:        TestQuerier{nil},
			ctx:      context.TODO(),
			ent:      SingleTestEntList,
			wantFunc: ErrNil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := NewEnforcer(TestCanner{test.canErr}, test.q)
			gotRet := e.Create(test.ctx, test.ent...)
			wantFunc := test.wantFunc
			if wantFunc == nil {
				wantFunc = ErrNil
			}
			if !wantFunc(gotRet) {
				t.Errorf("failed test %q", test.name)
			}
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name     string
		canErr   error
		q        database.Querier
		ctx      context.Context
		id       database.ID
		ent      database.Entity
		wantFunc func(error) bool
	}{
		{
			name:     "handles nil entity",
			canErr:   nil,
			q:        TestQuerier{nil},
			ctx:      context.TODO(),
			id:       database.EmptyID,
			ent:      database.Entity(nil),
			wantFunc: ErrNotNil,
		},
		{
			name:     "fails if db query fails",
			canErr:   nil,
			q:        TestQuerier{ErrGenericError},
			ctx:      context.TODO(),
			id:       database.EmptyID,
			ent:      TestEnt,
			wantFunc: ErrNotNil,
		},
		{
			name:     "fails if session disallows",
			canErr:   ErrGenericError,
			q:        TestQuerier{nil},
			ctx:      context.TODO(),
			id:       database.EmptyID,
			ent:      TestEnt,
			wantFunc: ErrNotNil,
		},
		{
			name:     "succeeds if allowed and query ok",
			canErr:   nil,
			q:        TestQuerier{nil},
			ctx:      context.TODO(),
			id:       database.EmptyID,
			ent:      TestEnt,
			wantFunc: ErrNil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := NewEnforcer(TestCanner{test.canErr}, test.q)
			gotRet := e.Get(test.ctx, test.id, test.ent)
			wantFunc := test.wantFunc
			if wantFunc == nil {
				wantFunc = ErrNil
			}
			if !wantFunc(gotRet) {
				t.Errorf("failed test %q", test.name)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name     string
		canErr   error
		q        database.Querier
		ctx      context.Context
		id       database.ID
		wantFunc func(error) bool
	}{
		{
			name:     "fails if session disallows",
			canErr:   ErrGenericError,
			q:        TestQuerier{nil},
			ctx:      context.TODO(),
			id:       database.EmptyID,
			wantFunc: ErrNotNil,
		},
		{
			name:     "fails if db query fails",
			canErr:   nil,
			q:        TestQuerier{ErrGenericError},
			ctx:      context.TODO(),
			id:       database.EmptyID,
			wantFunc: ErrNotNil,
		},
		{
			name:     "succeeds if session allows and query succeeds",
			canErr:   nil,
			q:        TestQuerier{nil},
			ctx:      context.TODO(),
			id:       database.EmptyID,
			wantFunc: ErrNil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := NewEnforcer(TestCanner{test.canErr}, test.q)
			gotRet := e.Delete(test.ctx, types.Test, test.id)
			wantFunc := test.wantFunc
			if wantFunc == nil {
				wantFunc = ErrNil
			}
			if !wantFunc(gotRet) {
				t.Errorf("failed test %q", test.name)
			}
		})
	}
}

func TestUpsert(t *testing.T) {
	tests := []struct {
		name     string
		canErr   error
		q        database.Querier
		ctx      context.Context
		ent      database.Entity
		wantFunc func(error) bool
	}{
		{
			name:     "handles nil entity",
			canErr:   nil,
			q:        TestQuerier{nil},
			ctx:      context.TODO(),
			ent:      database.Entity(nil),
			wantFunc: ErrNotNil,
		},
		{
			name:     "fails if db query fails",
			canErr:   nil,
			q:        TestQuerier{ErrGenericError},
			ctx:      context.TODO(),
			ent:      TestEnt,
			wantFunc: ErrNotNil,
		},
		{
			name:     "fails if session disallows",
			canErr:   ErrGenericError,
			q:        TestQuerier{nil},
			ctx:      context.TODO(),
			ent:      TestEnt,
			wantFunc: ErrNotNil,
		},
		{
			name:     "succeeds if allowed and query ok",
			canErr:   nil,
			q:        TestQuerier{nil},
			ctx:      context.TODO(),
			ent:      TestEnt,
			wantFunc: ErrNil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := NewEnforcer(TestCanner{test.canErr}, test.q)
			gotRet := e.Upsert(test.ctx, test.ent)
			wantFunc := test.wantFunc
			if wantFunc == nil {
				wantFunc = ErrNil
			}
			if !wantFunc(gotRet) {
				t.Errorf("failed test %q", test.name)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name     string
		canErr   error
		q        database.Querier
		ctx      context.Context
		ent      database.Entity
		wantFunc func(error) bool
	}{
		{
			name:     "handles nil entity",
			canErr:   nil,
			q:        TestQuerier{nil},
			ctx:      context.TODO(),
			ent:      database.Entity(nil),
			wantFunc: ErrNotNil,
		},
		{
			name:     "fails if db query fails",
			canErr:   nil,
			q:        TestQuerier{ErrGenericError},
			ctx:      context.TODO(),
			ent:      TestEnt,
			wantFunc: ErrNotNil,
		},
		{
			name:     "fails if session disallows",
			canErr:   ErrGenericError,
			q:        TestQuerier{nil},
			ctx:      context.TODO(),
			ent:      TestEnt,
			wantFunc: ErrNotNil,
		},
		{
			name:     "succeeds if allowed and query ok",
			canErr:   nil,
			q:        TestQuerier{nil},
			ctx:      context.TODO(),
			ent:      TestEnt,
			wantFunc: ErrNil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := NewEnforcer(TestCanner{test.canErr}, test.q)
			gotRet := e.Update(test.ctx, test.ent)
			wantFunc := test.wantFunc
			if wantFunc == nil {
				wantFunc = ErrNil
			}
			if !wantFunc(gotRet) {
				t.Errorf("failed test %q", test.name)
			}
		})
	}
}

func TestQueryOne(t *testing.T) {
	tests := []struct {
		name     string
		canErr   error
		q        database.Querier
		ctx      context.Context
		ent      database.Entity
		wantFunc func(error) bool
	}{
		{
			name:     "handles nil entity",
			canErr:   nil,
			q:        TestQuerier{nil},
			ctx:      context.TODO(),
			ent:      database.Entity(nil),
			wantFunc: ErrNotNil,
		},
		{
			name:     "fails if db query fails",
			canErr:   nil,
			q:        TestQuerier{ErrGenericError},
			ctx:      context.TODO(),
			ent:      TestEnt,
			wantFunc: ErrNotNil,
		},
		{
			name:     "fails if session disallows",
			canErr:   ErrGenericError,
			q:        TestQuerier{nil},
			ctx:      context.TODO(),
			ent:      TestEnt,
			wantFunc: ErrNotNil,
		},
		{
			name:     "succeeds if allowed and query ok",
			canErr:   nil,
			q:        TestQuerier{nil},
			ctx:      context.TODO(),
			ent:      TestEnt,
			wantFunc: ErrNil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := NewEnforcer(TestCanner{test.canErr}, test.q)
			gotRet := e.QueryOne(test.ctx, test.ent, database.Filter{})
			wantFunc := test.wantFunc
			if wantFunc == nil {
				wantFunc = ErrNil
			}
			if !wantFunc(gotRet) {
				t.Errorf("failed test %q", test.name)
			}
		})
	}
}

func TestQuery(t *testing.T) {
	tests := []struct {
		name     string
		canErr   error
		q        database.Querier
		ctx      context.Context
		i        interface{}
		wantFunc func(error) bool
	}{
		{
			name:     "handles empty array of valid type",
			canErr:   nil,
			q:        TestQuerier{nil},
			ctx:      context.TODO(),
			i:        &([]*primitives.Role{}),
			wantFunc: ErrNil,
		},
		{
			name:     "fails if db query fails",
			canErr:   nil,
			q:        TestQuerier{ErrGenericError},
			ctx:      context.TODO(),
			i:        &([]*primitives.Role{{}}),
			wantFunc: ErrNotNil,
		},
		{
			name:     "fails if session disallows",
			canErr:   ErrGenericError,
			q:        TestQuerier{nil},
			ctx:      context.TODO(),
			i:        &([]*primitives.Role{{}}),
			wantFunc: ErrNotNil,
		},
		{
			name:     "succeeds if allowed and query ok",
			canErr:   nil,
			q:        TestQuerier{nil},
			ctx:      context.TODO(),
			i:        &([]*primitives.Role{{}}),
			wantFunc: ErrNil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Logf("starting test %q\n", test.name)
			e := NewEnforcer(TestCanner{test.canErr}, test.q)
			gotRet := e.Query(test.ctx, test.i, database.Filter{})
			wantFunc := test.wantFunc
			if wantFunc == nil {
				wantFunc = ErrNil
			}
			if !wantFunc(gotRet) {
				t.Errorf("failed test %q", test.name)
			}
		})
	}
}
