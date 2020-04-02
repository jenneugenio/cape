package sources

import (
	"context"

	"github.com/dropoutlabs/cape/connector/proto"
	"github.com/dropoutlabs/cape/primitives"
)

// PostgresSource represents a Postgres data source which can be queried
// against by a Cape data connector
//
// PostgresSource completes the Source interface
type PostgresSource struct {
	source *primitives.Source
}

// NewPostgresSource is a constructor for creating a PostgresSource.
//
// Completes the NewSourceFunc interface enabling it to be used by the Registry
func NewPostgresSource(source *primitives.Source) (Source, error) {
	if err := source.Validate(); err != nil {
		return nil, err
	}

	return &PostgresSource{
		source: source,
	}, nil
}

// Label returns the label for the data source represented by this struct
func (p *PostgresSource) Label() primitives.Label {
	return p.source.Label
}

// Type returns the type of data source supported by this struct
func (p *PostgresSource) Type() primitives.SourceType {
	return primitives.PostgresType
}

// Schema returns the schema for the given Query. This schema can be used to
// issue a query or used to rewrite the query.
func (p *PostgresSource) Schema(ctx context.Context, q Query) (*proto.Schema, error) {
	return nil, nil
}

// Query issues the given query against the targeted collection from the query.
// The schema should be the schema of the targeted collection.
//
// Results from the query are streamed back to the requester through the
// provided Stream
func (p *PostgresSource) Query(ctx context.Context, q Query, schema *proto.Schema, stream Stream) error {
	return nil
}

// Close issues a close to cancel all on-going requests and closes any
// connections to the underlying postgres data source
func (p *PostgresSource) Close(ctx context.Context) error {
	return nil
}

// init registers the PostgresSource with the global sources registry
func init() {
	err := registry.Register(primitives.PostgresType, NewPostgresSource)
	if err != nil {
		panic("Could not register source: " + err.Error())
	}
}
