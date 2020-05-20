// +build integration

package sources

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/connector/proto"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/dbtest"
	"github.com/capeprivacy/cape/framework"
	"github.com/capeprivacy/cape/primitives"
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
		schema, err := source.QuerySchema(ctx, query)
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

	tests := []struct {
		name             string
		limit            int64
		offset           int64
		expectedRowCount int
	}{
		{name: "can stream rows with default limit", limit: 50, expectedRowCount: 9, offset: 0},
		{name: "can stream rows with custom limit", limit: 4, expectedRowCount: 4, offset: 0},
		{name: "can stream rows with custom limit and offset", limit: 3, expectedRowCount: 3, offset: 4},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			source, err := NewPostgresSource(ctx, cfg, src)
			gm.Expect(err).To(gm.BeNil())

			defer source.Close(ctx) // nolint: errcheck

			q := &testQuery{}

			stream := &testStream{}

			err = source.Query(ctx, stream, q, tc.limit, tc.offset)
			gm.Expect(err).To(gm.BeNil())

			gm.Expect(len(stream.Buffer)).To(gm.Equal(tc.expectedRowCount))

			query, params := q.Raw()
			query = fmt.Sprintf("%s LIMIT %d OFFSET %d", query, tc.limit, tc.offset)
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

	t.Run("test record to strings", func(t *testing.T) {
		source, err := NewPostgresSource(ctx, cfg, src)
		gm.Expect(err).To(gm.BeNil())

		defer source.Close(ctx) // nolint: errcheck

		q := &testQuery{}

		stream := &testStream{}

		err = source.Query(ctx, stream, q, 1, 0)
		gm.Expect(err).To(gm.BeNil())
		record, err := NewRecord(stream.Buffer[0].Schema, stream.Buffer[0].Fields)
		gm.Expect(err).To(gm.BeNil())

		strs, err := record.ToStrings()
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(len(strs)).To(gm.Equal(len(stream.Buffer[0].Fields)))

		expectedTime := record.Values()[8].(time.Time).Format(time.RFC3339Nano)
		expectedStrs := []string{
			"2", "4", "8", "8.8", "4.4", "hello", "thisisatest         ", "andthis",
			expectedTime, "false", "deadbeef",
		}

		gm.Expect(strs).To(gm.Equal(expectedStrs))
	})
}
