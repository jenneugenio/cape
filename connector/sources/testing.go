// testing contains functionality that makes it easier to write standalone
// tests for the sources package.
package sources

import (
	"context"

	"google.golang.org/grpc/metadata"

	"github.com/dropoutlabs/cape/connector/proto"
	errors "github.com/dropoutlabs/cape/partyerrors"
	"github.com/dropoutlabs/cape/primitives"
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
func (t *testSource) Query(_ context.Context, _ Query, _ *proto.Schema, _ Stream) error {
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

func (t *testClient) GetSource(ctx context.Context, source primitives.Label) (*primitives.Source, error) {
	return &primitives.Source{
		Label: source,
		Type:  testSourceType,
	}, nil
}

type errClient struct{}

func (e *errClient) GetSource(ctx context.Context, source primitives.Label) (*primitives.Source, error) {
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
type testQuery struct{}

func (t *testQuery) Source() primitives.Label {
	return primitives.Label("test")
}

func (t *testQuery) Collection() string {
	return "transactions"
}
