package capepg

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/models"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

func TestUsersCreate(t *testing.T) {
	tests := []struct {
		user    models.User
		wantErr error
		err     error
	}{
		{
			user:    models.User{},
			wantErr: nil,
			err:     nil,
		},
		{
			user:    models.User{},
			wantErr: fmt.Errorf("error creating user: %w", ErrGenericDBError),
			err:     ErrGenericDBError,
		},
	}

	pool := &testPgPool{}
	for i, test := range tests {
		pool.err = test.err

		userDB := pgUser{pool, 0}

		gotErr := userDB.Create(context.TODO(), test.user)
		if (test.wantErr == nil && gotErr != nil) ||
			(test.wantErr != nil && gotErr.Error() != test.wantErr.Error()) {
			t.Errorf("unexpected error on Create() test %d of %d: got %v want %v", i+1, len(tests), gotErr, test.wantErr)
		}
	}
}

func TestUsersUpdate(t *testing.T) {
	tests := []struct {
		user    models.User
		id      string
		wantErr error
		err     error
	}{
		{
			user:    models.User{},
			id:      "idididid",
			wantErr: nil,
			err:     nil,
		},
		{
			user:    models.User{},
			id:      "idididid",
			wantErr: fmt.Errorf("error updating user: %w", ErrGenericDBError),
			err:     ErrGenericDBError,
		},
	}

	pool := &testPgPool{}
	for i, test := range tests {
		pool.err = test.err

		userDB := pgUser{pool, 0}

		gotErr := userDB.Update(context.TODO(), test.id, test.user)
		if (test.wantErr == nil && gotErr != nil) ||
			(test.wantErr != nil && gotErr.Error() != test.wantErr.Error()) {
			t.Errorf("unexpected error on Create() test %d of %d: got %v want %v", i+1, len(tests), gotErr, test.wantErr)
		}
	}
}

func TestUsersDelete(t *testing.T) {
	tagDeleted := pgconn.CommandTag("DELETE 1")
	tests := []struct {
		pool    *testPgPool
		email   models.Email
		wantErr error
	}{
		{
			pool: &testPgPool{
				ct: &tagDeleted,
			},
			email:   models.Email("foo"),
			wantErr: nil,
		},
		{
			pool: &testPgPool{
				ct:  &tagDeleted,
				err: ErrGenericDBError,
			},
			email:   models.Email("foo"),
			wantErr: fmt.Errorf("error deleting user: %w", ErrGenericDBError),
		},
	}

	for i, test := range tests {
		userDB := pgUser{test.pool, 0}

		gotErr := userDB.Delete(context.TODO(), test.email)
		if (test.wantErr == nil && gotErr != nil) ||
			(test.wantErr != nil && gotErr.Error() != test.wantErr.Error()) {
			t.Errorf("unexpected error on Create() test %d of %d: got %v want %v", i+1, len(tests), gotErr, test.wantErr)
		}
	}
}

var EmptyUser = models.User{
	ID:      "foo",
	Version: 1,
	Email:   models.Email("foo"),
}

func TestUserGet(t *testing.T) {
	tests := []struct {
		email    models.Email
		wantUser *models.User
		wantErr  error
		row      pgx.Row
		err      error
	}{
		{
			email:    models.Email("foo"),
			wantUser: &EmptyUser,
			wantErr:  nil,
			row: testRow{
				obj: []interface{}{EmptyUser},
			},
			err: nil,
		},
		{
			email:    models.Email("foo"),
			wantUser: &EmptyUser,
			wantErr:  fmt.Errorf("error retrieving user: %w", ErrGenericDBError),
			row: testRow{
				obj: []interface{}{EmptyUser},
				err: nil,
			},
			err: ErrGenericDBError,
		},
		{
			email:    models.Email("foo"),
			wantUser: &EmptyUser,
			wantErr:  db.ErrNoRows,
			row: testRow{
				obj: []interface{}{EmptyUser},
				err: nil,
			},
			err: pgx.ErrNoRows,
		},
	}

	pool := &testPgPool{}
	for i, test := range tests {
		pool.row = test.row
		pool.err = test.err

		userDB := pgUser{pool, 0}

		gotPol, gotErr := userDB.Get(context.TODO(), test.email)
		if (test.wantErr == nil && gotErr != nil) ||
			(test.wantErr != nil && gotErr != nil && gotErr.Error() != test.wantErr.Error()) {
			t.Errorf("unexpected error on Create() test %d of %d: got %v want %v", i+1, len(tests), gotErr, test.wantErr)
		}
		got, want := *gotPol, *(test.wantUser)
		if got != want {
			t.Errorf("incorrect user returned on Get() test %d of %d: got %v want %v", i+1, len(tests), got, want)
		}
	}
}

func TestUserGetByID(t *testing.T) {
	tests := []struct {
		id       string
		wantUser *models.User
		wantErr  error
		row      pgx.Row
		err      error
	}{
		{
			id:       "idididid",
			wantUser: &EmptyUser,
			wantErr:  nil,
			row: testRow{
				obj: []interface{}{EmptyUser},
			},
			err: nil,
		},
		{
			id:       "idididid",
			wantUser: &EmptyUser,
			wantErr:  fmt.Errorf("error retrieving user: %w", ErrGenericDBError),
			row: testRow{
				obj: []interface{}{EmptyUser},
				err: nil,
			},
			err: ErrGenericDBError,
		},
		{
			id:       "idididid",
			wantUser: &EmptyUser,
			wantErr:  db.ErrNoRows,
			row: testRow{
				obj: []interface{}{EmptyUser},
				err: nil,
			},
			err: pgx.ErrNoRows,
		},
	}

	pool := &testPgPool{}
	for i, test := range tests {
		pool.row = test.row
		pool.err = test.err

		userDB := pgUser{pool, 0}

		gotUser, gotErr := userDB.GetByID(context.TODO(), test.id)
		if (test.wantErr == nil && gotErr != nil) ||
			(test.wantErr != nil && gotErr != nil && gotErr.Error() != test.wantErr.Error()) {
			t.Errorf("unexpected error on Create() test %d of %d: got %v want %v", i+1, len(tests), gotErr, test.wantErr)
		}
		got, want := *gotUser, *(test.wantUser)
		if got != want {
			t.Errorf("incorrect user returned on Get() test %d of %d: got %v want %v", i+1, len(tests), got, want)
		}
	}
}

func TestUserstList(t *testing.T) {
	tests := []struct {
		opt       *db.ListUserOptions
		wantUsers []models.User
		wantErr   error
		rows      pgx.Rows
		err       error
	}{
		{
			opt:       nil,
			wantUsers: []models.User{{}},
			wantErr:   nil,
			rows: &testRows{
				obj: [][]interface{}{{models.User{}}},
				err: nil,
			},
			err: nil,
		},
		{
			opt: &db.ListUserOptions{
				Options: &struct {
					Offset uint64
					Limit  uint64
				}{
					Offset: 0,
					Limit:  1,
				},
			},
			wantUsers: []models.User{{}},
			wantErr:   nil,
			rows: &testRows{
				obj: [][]interface{}{{models.User{}}},
				err: nil,
			},
			err: nil,
		},
		{
			opt: &db.ListUserOptions{
				FilterIDs: []string{"idididid"},
			},
			wantUsers: []models.User{{
				ID: "idididid",
			}},
			wantErr: nil,
			rows: &testRows{
				obj: [][]interface{}{{models.User{
					ID: "idididid",
				}}},
				err: nil,
			},
			err: nil,
		},
	}

	pool := &testPgPool{}
	for i, test := range tests {
		pool.rows = test.rows
		pool.err = test.err

		userDB := pgUser{pool, 0}

		gotUsers, gotErr := userDB.List(context.TODO(), test.opt)
		if (test.wantErr == nil && gotErr != nil) ||
			(test.wantErr != nil && gotErr != nil && gotErr.Error() != test.wantErr.Error()) {
			t.Errorf("unexpected error on List() test %d of %d: got %v want %v", i+1, len(tests), gotErr, test.wantErr)
		}
		if !reflect.DeepEqual(gotUsers, test.wantUsers) {
			t.Errorf("incorrect user returned on List() test %d of %d: got %v want %v", i+1, len(tests), gotUsers, test.wantUsers)
		}
	}
}