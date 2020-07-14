package capepg

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/models"
	"github.com/jackc/pgx/v4"
)

func TestCreateRBAC(t *testing.T) {
	tests := []struct {
		r       models.RBAC
		wantErr error
		err     error
	}{
		{
			r:       models.RBAC{},
			wantErr: nil,
			err:     nil,
		},
		{
			r:       models.RBAC{},
			wantErr: fmt.Errorf("error creating rbac: %w", ErrGenericDBError),
			err:     ErrGenericDBError,
		},
	}

	pool := &testPgPool{}
	for i, test := range tests {
		pool.err = test.err

		rbacDB := pgRBAC{pool, 0}

		gotErr := rbacDB.Create(context.TODO(), test.r)
		if (test.wantErr == nil && gotErr != nil) ||
			(test.wantErr != nil && gotErr.Error() != test.wantErr.Error()) {
			t.Errorf("unexpected error on Create() test %d of %d: got %v want %v", i+1, len(tests), gotErr, test.wantErr)
		}
	}
}

var EmptyRBAC = models.RBAC{
	ID:      "foo",
	Version: 1,
	Label:   models.Label("foo"),
}

func TestListRBAC(t *testing.T) {
	by, _ := json.Marshal(models.RBAC{ID: "idididid"})

	tests := []struct {
		opt     *db.ListRBACOptions
		wantRs  []models.RBAC
		wantErr error
		rows    pgx.Rows
		err     error
	}{
		{
			opt:     nil,
			wantRs:  []models.RBAC{{ID: "idididid"}},
			wantErr: nil,
			rows: &testRows{
				obj: [][]interface{}{{by}},
				err: nil,
			},
			err: nil,
		},
		{
			opt: &db.ListRBACOptions{
				FilterIDs: []string{"idididid"},
			},
			wantRs:  []models.RBAC{{ID: "idididid"}},
			wantErr: nil,
			rows: &testRows{
				obj: [][]interface{}{{by}},
				err: nil,
			},
			err: nil,
		},
	}

	pool := &testPgPool{}
	for i, test := range tests {
		pool.rows = test.rows
		pool.err = test.err

		rbacDB := pgRBAC{pool, 0}

		gotRs, gotErr := rbacDB.List(context.TODO(), test.opt)
		if (test.wantErr == nil && gotErr != nil) ||
			(test.wantErr != nil && gotErr != nil && gotErr.Error() != test.wantErr.Error()) {
			t.Errorf("unexpected error on List() test %d of %d: got %v want %v", i+1, len(tests), gotErr, test.wantErr)
		}
		if !reflect.DeepEqual(gotRs, test.wantRs) {
			t.Errorf("incorrect rbac returned on List() test %d of %d: got %v want %v", i+1, len(tests), gotRs, test.wantRs)
		}
	}
}