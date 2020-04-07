package client

import (
	"bytes"
	"context"
	"io"

	"github.com/golang/protobuf/jsonpb"
	spb "github.com/golang/protobuf/ptypes/struct"
	"google.golang.org/grpc/status"

	pb "github.com/dropoutlabs/cape/connector/proto"
	"github.com/dropoutlabs/cape/connector/sources"
	errors "github.com/dropoutlabs/cape/partyerrors"
)

// Stream allows data to be streamed easily
// from a grpc stream to the caller
type Stream interface {
	Record() *sources.Record
	NextRecord() bool
	Schema() *pb.Schema
	Error() error

	Context() context.Context
	Close()
}

type stream struct {
	client pb.DataConnector_QueryClient

	ctx    context.Context
	cancel context.CancelFunc

	err error

	schema     *pb.Schema
	nextRecord *sources.Record
}

// NewStream returns a new stream which can pull records off of the given
// grpc client stream
func NewStream(ctx context.Context, client pb.DataConnector_QueryClient) Stream {
	ctx, cancel := context.WithCancel(ctx)
	return &stream{
		schema: &pb.Schema{},

		client: client,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (s *stream) Record() *sources.Record {
	return s.nextRecord
}

func (s *stream) NextRecord() bool {
	record, err := s.client.Recv()
	if err == io.EOF {
		return false
	}
	if err != nil {
		s.err = err
		return false
	}

	if record.Schema != nil {
		*s.schema = *record.Schema
	}

	r, err := sources.NewRecord(s.schema, record.Value)
	if err != nil {
		s.err = err
		return false
	}
	s.nextRecord = r

	return true
}

func (s *stream) Context() context.Context {
	return s.ctx
}

func (s *stream) Close() {
	// closes the grpc stream
	s.cancel()
}

func (s *stream) Schema() *pb.Schema {
	return s.schema
}

func (s *stream) Error() error {
	st := status.Convert(s.err)
	details := st.Details()

	if len(details) > 0 {
		buf := &bytes.Buffer{}

		switch info := details[0].(type) {
		case *spb.Struct:
			marshaler := jsonpb.Marshaler{}
			err := marshaler.Marshal(buf, info)
			if err != nil {
				return err
			}

			return errors.ErrorFromBytes(buf.Bytes())
		}
	}

	return st.Err()
}
