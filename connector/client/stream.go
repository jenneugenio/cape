package client

import (
	"context"
	"io"

	spb "github.com/golang/protobuf/ptypes/struct"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	pb "github.com/capeprivacy/cape/connector/proto"
	"github.com/capeprivacy/cape/connector/sources"
	errors "github.com/capeprivacy/cape/partyerrors"
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

	if s.schema == nil {
		s.schema = record.GetSchema()
	}

	r, err := sources.NewRecord(s.schema, record.Fields)
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
		switch info := details[0].(type) {
		case *spb.Struct:
			by, err := protojson.Marshal(info)
			if err != nil {
				return err
			}

			return errors.ErrorFromBytes(by)
		}
	}

	return st.Err()
}
