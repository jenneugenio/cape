package capepg

import (
	"context"
	"fmt"
	"reflect"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgproto3/v2"
	"github.com/jackc/pgx/v4"
)

var (
	ErrGenericDBError = fmt.Errorf("generic db error")
)

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
	if len(t.obj) == 0 {
		return pgx.ErrNoRows
	}

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
