// +build integration

package sources

import (
	"context"
	"os"
	"testing"

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

	err = db.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	migrator, err := database.NewMigrator(db.URL(), "testdata")
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

		gm.Expect(len(schema.Fields)).To(gm.Equal(11))

		expectedFields := []*proto.FieldInfo{
			{
				Field: proto.FieldType_SMALLINT,
				Name:  "int2",
				Size:  2,
			},
			{
				Field: proto.FieldType_INT,
				Name:  "int4",
				Size:  4,
			},
			{
				Field: proto.FieldType_BIGINT,
				Name:  "int8",
				Size:  8,
			},
			{
				Field: proto.FieldType_DOUBLE,
				Name:  "float8",
				Size:  8,
			},
			{
				Field: proto.FieldType_REAL,
				Name:  "float4",
				Size:  4,
			},
			{
				Field: proto.FieldType_VARCHAR,
				Name:  "vchar",
				Size:  20,
			},
			{
				Field: proto.FieldType_CHAR,
				Name:  "ch",
				Size:  20,
			},
			{
				Field: proto.FieldType_TEXT,
				Name:  "txt",
				Size:  VariableSize,
			},
			{
				Field: proto.FieldType_TIMESTAMP,
				Name:  "ts",
				Size:  8,
			},
			{
				Field: proto.FieldType_BOOL,
				Name:  "bool",
				Size:  1,
			},
			{
				Field: proto.FieldType_BYTEA,
				Name:  "bytes",
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

		q := &testQuery{}

		stream := &testStream{}

		err = source.Query(ctx, q, stream)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(len(stream.Buffer)).To(gm.Equal(9))

		query, params := q.Raw()
		expectedRows, err := GetExpectedRows(ctx, db.URL(), query, params)
		gm.Expect(err).To(gm.BeNil())
		for i, row := range stream.Buffer {
			vals, err := Decode(stream.Buffer[0].Schema, row.Fields)
			gm.Expect(err).To(gm.BeNil())

			// could check row to row but this is easier to see
			// if there are any errors
			for j, val := range vals {
				gm.Expect(val).To(gm.Equal(expectedRows[i][j]))
			}
		}
	})
}
