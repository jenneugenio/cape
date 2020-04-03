package sources

import (
	"context"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	"github.com/dropoutlabs/cape/connector/proto"
	errors "github.com/dropoutlabs/cape/partyerrors"
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
	//
	// TODO: Update this type to reflect the target that's returned for the
	// query
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

// NewSourceFunc represents a function constructor for a Source
type NewSourceFunc func(context.Context, *Config, *primitives.Source) (Source, error)

// Config represents configuration thats common across the Cache and Sources
type Config struct {
	InstanceID primitives.Label
	Logger     *zerolog.Logger
}

// Validate returns an error if the given Config struct is invalid
func (c *Config) Validate() error {
	if err := c.InstanceID.Validate(); err != nil {
		return err
	}

	if c.Logger == nil {
		return errors.New(InvalidConfig, "Missing logger from Config")
	}

	return nil
}
