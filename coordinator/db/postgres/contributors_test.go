package capepg

import (
	"context"
	"fmt"
	"testing"

	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/models"
	"github.com/jackc/pgx/v4"
	gm "github.com/onsi/gomega"
)

func TestAddContributor(t *testing.T) {
	gm.RegisterTestingT(t)

	tests := []struct {
		name    string
		c       models.Contributor
		wantErr error
		err     error
		row     pgx.Row
	}{
		{
			name:    "Can create without any errors",
			c:       models.Contributor{},
			wantErr: nil,
			err:     nil,
			row: &testRows{
				obj: [][]interface{}{{contributorAddIDs{"123", "123", "123"}}},
				err: nil,
			},
		},
		{
			name:    "Contributor doesn't get created if the DB errors",
			c:       models.Contributor{},
			wantErr: fmt.Errorf("error creating contributor: %w", ErrGenericDBError),
			err:     ErrGenericDBError,
			row: &testRows{
				obj: [][]interface{}{{contributorAddIDs{"123", "123", "123"}}},
				err: nil,
			},
		},
	}

	pool := &testPgPool{}
	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pool.err = test.err
			pool.row = test.row
			contributorDB := pgContributor{pool, 0}

			_, gotErr := contributorDB.Add(context.TODO(), "my-project", "me.com")
			if (test.wantErr == nil && gotErr != nil) ||
				(test.wantErr != nil && gotErr.Error() != test.wantErr.Error()) {
				t.Errorf("unexpected error on Add() test %d of %d: got %v want %v", i+1, len(tests), gotErr, test.wantErr)
			}
		})
	}
}

func TestGetContributor(t *testing.T) {
	gm.RegisterTestingT(t)
	tests := []struct {
		name        string
		c           models.Contributor
		expectedErr error
		row         pgx.Row
	}{
		{
			name:        "querying a contributor that doesn't exist",
			expectedErr: db.ErrCannotFindContributor,
			row: &testRow{
				obj: []interface{}{},
				err: nil,
			},
		},

		{
			name:        "can get a contributor",
			expectedErr: nil,
			row: &testRow{
				obj: []interface{}{models.Contributor{}},
			},
		},
	}

	pool := &testPgPool{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pool.row = test.row
			contributorDB := pgContributor{pool, 0}

			c, err := contributorDB.Get(context.TODO(), "my-project", "me@me.com")
			if test.expectedErr == nil {
				gm.Expect(c).ToNot(gm.BeNil())
				gm.Expect(err).To(gm.BeNil())
			} else {
				gm.Expect(err).To(gm.Equal(test.expectedErr))
				gm.Expect(c).To(gm.BeNil())
			}
		})
	}
}

func TestListContributors(t *testing.T) {
	gm.RegisterTestingT(t)
	tests := []struct {
		name string
		c    models.Contributor
		rows pgx.Rows
		err  error
	}{
		{
			name: "returns db errors",
			err:  fmt.Errorf("Error reading from the db"),
		},

		{
			name: "can list contributors",
			err:  nil,
			rows: &testRows{
				obj: [][]interface{}{{models.Contributor{}}},
			},
		},
	}

	pool := &testPgPool{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pool.rows = test.rows
			pool.err = test.err

			contributorDB := pgContributor{pool, 0}
			contribs, err := contributorDB.List(context.TODO(), "my-project")
			if err != nil {
				gm.Expect(err).ToNot(gm.BeNil())
				gm.Expect(contribs).To(gm.BeNil())
			} else {
				gm.Expect(err).To(gm.BeNil())
				gm.Expect(contribs).ToNot(gm.BeNil())
			}
		})
	}
}

func TestDeleteContributor(t *testing.T) {
	gm.RegisterTestingT(t)
	gm.RegisterTestingT(t)
	tests := []struct {
		name string
		c    models.Contributor
		rows pgx.Rows
		err  error
	}{
		{
			name: "returns db errors",
			err:  fmt.Errorf("Error reading from the db"),
		},

		{
			name: "can list contributors",
			err:  nil,
			rows: &testRows{
				obj: [][]interface{}{{models.Contributor{}}},
			},
		},
	}

	pool := &testPgPool{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pool.rows = test.rows
			pool.err = test.err

			contributorDB := pgContributor{pool, 0}
			contribs, err := contributorDB.Delete(context.TODO(), "my-project", "me@me.com")
			if err != nil {
				gm.Expect(err).ToNot(gm.BeNil())
				gm.Expect(contribs).To(gm.BeNil())
			} else {
				gm.Expect(err).To(gm.BeNil())
				gm.Expect(contribs).ToNot(gm.BeNil())
			}
		})
	}
}
