// testing contains functionality that makes it easier to write standalone
// tests for the sources package.
package sources

import (
	"context"
	"net/url"

	"github.com/jackc/pgx/v4"

	"google.golang.org/grpc/metadata"

	"github.com/capeprivacy/cape/connector/proto"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
	"github.com/jackc/pgx/v4/pgxpool"
)

var testSourceType primitives.SourceType = "test"

type testSource struct {
	cfg    *Config
	source *primitives.Source
}

func (t *testSource) Label() primitives.Label {
	return primitives.Label("test")
}
func (t *testSource) Type() primitives.SourceType {
	return testSourceType
}
func (t *testSource) Schema(_ context.Context, _ Query) (*proto.Schema, error) {
	return &proto.Schema{}, nil
}
func (t *testSource) Query(_ context.Context, _ Stream, _ Query, _ int64, _ int64) error {
	return nil
}
func (t *testSource) Close(_ context.Context) error {
	return nil
}

func newTestSource(ctx context.Context, cfg *Config, source *primitives.Source) (Source, error) {
	return &testSource{
		cfg:    cfg,
		source: source,
	}, nil
}

type testClient struct{}

func (t *testClient) GetSourceByLabel(ctx context.Context, source primitives.Label) (*primitives.Source, error) {
	return &primitives.Source{
		Label: source,
		Type:  testSourceType,
	}, nil
}

type errClient struct{}

func (e *errClient) GetSourceByLabel(ctx context.Context, source primitives.Label) (*primitives.Source, error) {
	return nil, errors.New(NotFoundCause, "whoops")
}

// nolint: unused
type testStream struct {
	Buffer []*proto.Record
}

func (t *testStream) Send(r *proto.Record) error {
	t.Buffer = append(t.Buffer, r)
	return nil
}

func (t *testStream) SetHeader(_ metadata.MD) error {
	return nil
}

func (t *testStream) SendHeader(_ metadata.MD) error {
	return nil
}

func (t *testStream) SetTrailer(_ metadata.MD) {}

func (t *testStream) Context() context.Context {
	return context.Background()
}

func (t *testStream) SendMsg(_ interface{}) error {
	return nil
}

func (t *testStream) RecvMsg(_ interface{}) error {
	return nil
}

// nolint: unused
type testQuery struct {
}

func (t *testQuery) Source() primitives.Label {
	return primitives.Label("test")
}

func (t *testQuery) Collection() string {
	return "test"
}

func (t *testQuery) Raw() (string, []interface{}) {
	return "SELECT * FROM test", make([]interface{}, 0)
}

// GetExpectedRows is a testing helper to query the expected rows from the database
func GetExpectedRows(ctx context.Context, dbURL *url.URL, query string, params []interface{}) ([][]interface{}, error) {
	pool, err := pgxpool.Connect(ctx, dbURL.String())
	if err != nil {
		return nil, err
	}

	// If there are no params, pgx will error if you pass anything (even nil!)
	var rows pgx.Rows
	if len(params) == 0 {
		rows, err = pool.Query(ctx, query)
	} else {
		rows, err = pool.Query(ctx, query, params...)
	}

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
