// +build integration

package sources

import (
	"context"
	"net/url"
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	gm "github.com/onsi/gomega"

	"github.com/dropoutlabs/cape/connector/proto"
	"github.com/dropoutlabs/cape/database"
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

	seedMigrations := os.Getenv("CAPE_DB_SEED_MIGRATIONS")

	err = db.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	migrator, err := database.NewMigrator(db.URL(), seedMigrations)
	gm.Expect(err).To(gm.BeNil())

	defer func() {
		migrator.Down(ctx) // nolint: errcheck
		db.Teardown(ctx)
	}()

	err = migrator.Up(ctx)
	gm.Expect(err).To(gm.BeNil())

	cfg := &Config{
		InstanceID: primitives.Label("cape-source-tester"),
		Logger:     framework.TestLogger(),
	}

	dbURL, err := primitives.DBURLFromURL(db.URL())
	gm.Expect(err).To(gm.BeNil())

	src, err := primitives.NewSource(primitives.Label("test"), dbURL, nil)
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

		defer source.Close(ctx) // nolint: errcheck

		query := &testQuery{}
		schema, err := source.Schema(ctx, query)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(schema).ToNot(gm.BeNil())

		gm.Expect(schema.DataSource).To(gm.Equal(src.Label.String()))
		gm.Expect(schema.Target).To(gm.Equal(query.Collection()))
		gm.Expect(schema.Type).To(gm.Equal(proto.RecordType_DOCUMENT))

		gm.Expect(len(schema.Fields)).To(gm.Equal(8))

		expectedFields := []*proto.Field{
			&proto.Field{
				Field: proto.FieldType_INT,
				Name:  "id",
				Size:  4,
			},
			&proto.Field{
				Field: proto.FieldType_TEXT,
				Name:  "processor",
				Size:  VariableSize,
			},
			&proto.Field{
				Field: proto.FieldType_TIMESTAMP,
				Name:  "timestamp",
				Size:  8,
			},
			&proto.Field{
				Field: proto.FieldType_INT,
				Name:  "card_id",
				Size:  4,
			},
			&proto.Field{
				Field: proto.FieldType_BIGINT,
				Name:  "card_number",
				Size:  8,
			},
			&proto.Field{
				Field: proto.FieldType_DOUBLE,
				Name:  "value",
				Size:  8,
			},
			&proto.Field{
				Field: proto.FieldType_INT,
				Name:  "ssn",
				Size:  4,
			},
			&proto.Field{
				Field: proto.FieldType_TEXT,
				Name:  "vendor",
				Size:  VariableSize,
			},
		}

		for i, field := range schema.Fields {
			gm.Expect(field).To(gm.Equal(expectedFields[i]))
		}
	})

	t.Run("can stream rows back for query", func(t *testing.T) {
		source, err := NewPostgresSource(ctx, cfg, src)
		gm.Expect(err).To(gm.BeNil())

		defer source.Close(ctx) // nolint: errcheck

		limit := 5
		q := &testQuery{Limit: limit}

		stream := &testStream{}

		err = source.Query(ctx, q, stream)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(len(stream.Buffer)).To(gm.Equal(limit))

		expectedRows, err := getExpectedRows(ctx, db.URL(), q.Raw())
		gm.Expect(err).To(gm.BeNil())
		for i, row := range stream.Buffer {
			vals, err := Decode(stream.Buffer[0].Schema, row.Value)
			gm.Expect(err).To(gm.BeNil())

			// could check row to row but this is easier to see
			// if there are any errors
			for j, val := range vals {
				gm.Expect(val).To(gm.Equal(expectedRows[i][j]))
			}
		}
	})
}

func getExpectedRows(ctx context.Context, dbURL *url.URL, query string) ([][]interface{}, error) {
	pool, err := pgxpool.Connect(context.Background(), dbURL.String())
	if err != nil {
		return nil, err
	}

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var outVals [][]interface{}
	for rows.Next() {
		vals, err := rows.Values()
		if err != nil {
			return nil, err
		}
		outVals = append(outVals, vals)
	}

	return outVals, nil
}
