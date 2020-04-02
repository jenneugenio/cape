package sources

import (
	"context"

	"google.golang.org/grpc"

	"github.com/dropoutlabs/cape/connector/proto"
	"github.com/dropoutlabs/cape/primitives"
)

// Stream is an interface that represents the response to a query which streams
// records from the data source back to the requester.
//
// This interface exists to make it easy to test the sources package in
// isolation of other packages.
type Stream interface {
	grpc.ServerStream
	Send(*proto.Record) error
}

// Query is an interface that represents a query for data from a data source.
//
// This interface exists to make it easy to test the sources package in
// isolation of other packages.
type Query interface {
	// Source returns the data source label being targeted by this label
	Source() primitives.Label

	// Collection returns the name of the collection (e.g. table or path to
	// table) that this query is attempting to query against.
	Collection() string
}

// Source is an interface that any data source provider must
type Source interface {
	Label() primitives.Label
	Type() primitives.SourceType

	// Schema returns the schema for the collection targeted by the given query
	Schema(context.Context, Query) (*proto.Schema, error)

	// Query begins responding to the given query and sending results to the
	// client.
	Query(context.Context, Query, *proto.Schema, Stream) error

	// Close closes all underlying connections to the database. This will close
	// any on-going requests.
	Close(context.Context) error
}
