package connector

import (
	"context"

	"github.com/capeprivacy/cape/connector/proto"
	"google.golang.org/grpc/metadata"
)

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
