package capepg

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/models"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgproto3/v2"
	"github.com/jackc/pgx/v4"
)

var (
	ErrGenericDBError = fmt.Errorf("generic db error")
)

func TestCreate(t *testing.T) {
	tests := []struct {
		pol     models.Policy
		wantErr error
		err     error
	}{
		{
			pol:     models.Policy{},
			wantErr: nil,
			err:     nil,
		},
		{
			pol:     models.Policy{},
			wantErr: fmt.Errorf("error creating policy: %w", ErrGenericDBError),
			err:     ErrGenericDBError,
		},
	}

	pool := &testPgPool{}
	for i, test := range tests {
		pool.err = test.err

		policyDB := pgPolicy{pool, 0}

		gotErr := policyDB.Create(context.TODO(), test.pol)
		if (test.wantErr == nil && gotErr != nil) ||
			(test.wantErr != nil && gotErr.Error() != test.wantErr.Error()) {
			t.Errorf("unexpected error on Create() test %d of %d: got %v want %v", i+1, len(tests), gotErr, test.wantErr)
		}
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		label   models.Label
		wantErr error
		err     error
	}{
		{
			label:   models.Label("foo"),
			wantErr: nil,
			err:     nil,
		},
		{
			label:   models.Label("foo"),
			wantErr: fmt.Errorf("error deleting policy: %w", ErrGenericDBError),
			err:     ErrGenericDBError,
		},
	}

	pool := &testPgPool{}
	for i, test := range tests {
		pool.err = test.err

		policyDB := pgPolicy{pool, 0}

		gotErr := policyDB.Delete(context.TODO(), test.label)
		if (test.wantErr == nil && gotErr != nil) ||
			(test.wantErr != nil && gotErr.Error() != test.wantErr.Error()) {
			t.Errorf("unexpected error on Create() test %d of %d: got %v want %v", i+1, len(tests), gotErr, test.wantErr)
		}
	}
}

var EmptyPolicy = models.Policy{
	ID:      "foo",
	Version: 1,
	Label:   models.Label("foo"),
}

func TestGet(t *testing.T) {
	tests := []struct {
		label   models.Label
		wantPol *models.Policy
		wantErr error
		row     pgx.Row
		err     error
	}{
		{
			label:   models.Label("foo"),
			wantPol: &EmptyPolicy,
			wantErr: nil,
			row: testRow{
				obj: []interface{}{EmptyPolicy},
			},
			err: nil,
		},
		{
			label:   models.Label("foo"),
			wantPol: &EmptyPolicy,
			wantErr: fmt.Errorf("error retrieving policy: %w", ErrGenericDBError),
			row: testRow{
				obj: []interface{}{EmptyPolicy},
				err: nil,
			},
			err: ErrGenericDBError,
		},
	}

	pool := &testPgPool{}
	for i, test := range tests {
		pool.row = test.row
		pool.err = test.err

		policyDB := pgPolicy{pool, 0}

		gotPol, gotErr := policyDB.Get(context.TODO(), test.label)
		if (test.wantErr == nil && gotErr != nil) ||
			(test.wantErr != nil && gotErr != nil && gotErr.Error() != test.wantErr.Error()) {
			t.Errorf("unexpected error on Create() test %d of %d: got %v want %v", i+1, len(tests), gotErr, test.wantErr)
		}
		got, want := *gotPol, *(test.wantPol)
		if got != want {
			t.Errorf("incorrect policy returned on Get() test %d of %d: got %v want %v", i+1, len(tests), got, want)
		}
	}
}

func TestList(t *testing.T) {
	tests := []struct {
		opt      *db.ListPolicyOptions
		wantPols []models.Policy
		wantErr  error
		rows     pgx.Rows
		err      error
	}{
		{
			opt:      nil,
			wantPols: []models.Policy{{}},
			wantErr:  nil,
			rows: &testRows{
				obj: [][]interface{}{{models.Policy{}}},
				err: nil,
			},
			err: nil,
		},
		{
			opt: &db.ListPolicyOptions{
				Offset: 0,
				Limit:  1,
			},
			wantPols: []models.Policy{{}},
			wantErr:  nil,
			rows: &testRows{
				obj: [][]interface{}{{models.Policy{}}},
				err: nil,
			},
			err: nil,
		},
	}

	pool := &testPgPool{}
	for i, test := range tests {
		pool.rows = test.rows
		pool.err = test.err

		policyDB := pgPolicy{pool, 0}

		gotPols, gotErr := policyDB.List(context.TODO(), test.opt)
		if (test.wantErr == nil && gotErr != nil) ||
			(test.wantErr != nil && gotErr != nil && gotErr.Error() != test.wantErr.Error()) {
			t.Errorf("unexpected error on List() test %d of %d: got %v want %v", i+1, len(tests), gotErr, test.wantErr)
		}
		if !reflect.DeepEqual(gotPols, test.wantPols) {
			t.Errorf("incorrect policy returned on List() test %d of %d: got %v want %v", i+1, len(tests), gotPols, test.wantPols)
		}
	}
}

type testPgPool struct {
	ct   *pgconn.CommandTag
	rows pgx.Rows
	row  pgx.Row
	err  error

	callCount int
	lastSQL   string
	lastArgs  []interface{}
}

func (t *testPgPool) Exec(_ context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	t.callCount++
	t.lastSQL = sql
	t.lastArgs = args
	if t.ct == nil {
		return pgconn.CommandTag(nil), t.err
	}
	return *(t.ct), t.err
}

func (t *testPgPool) Query(_ context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	t.callCount++
	t.lastSQL = sql
	t.lastArgs = args
	return t.rows, t.err
}

func (t *testPgPool) QueryRow(_ context.Context, sql string, args ...interface{}) pgx.Row {
	t.callCount++
	t.lastSQL = sql
	t.lastArgs = args
	return t.row
}

type testRows struct {
	ct  pgconn.CommandTag
	obj [][]interface{}
	err error
}

func (t *testRows) Close()                                         {}
func (t *testRows) Err() error                                     { return t.err }
func (t *testRows) CommandTag() pgconn.CommandTag                  { return t.ct }
func (t *testRows) FieldDescriptions() []pgproto3.FieldDescription { return nil }
func (t *testRows) Next() bool                                     { return len(t.obj) > 0 }
func (t *testRows) Values() ([]interface{}, error)                 { return nil, t.err }
func (t *testRows) RawValues() [][]byte                            { return nil }

func (t *testRows) Scan(dest ...interface{}) error {
	retVal := t.obj[0]
	t.obj = t.obj[1:]

	if len(dest) != len(retVal) {
		panic("unable to scan into dest")
	}

	for i, item := range dest {
		orig := reflect.ValueOf(item)
		replacement := retVal[i]
		reflect.Indirect(orig).Set(reflect.ValueOf(replacement))
	}
	return t.err
}

type testRow struct {
	obj []interface{}
	err error
}

func (t testRow) Scan(dest ...interface{}) error {
	if len(dest) != len(t.obj) {
		panic("unable to scan into dest")
	}

	for i, item := range dest {
		orig := reflect.ValueOf(item)
		replacement := t.obj[i]
		reflect.Indirect(orig).Set(reflect.ValueOf(replacement))
	}
	return t.err
}
