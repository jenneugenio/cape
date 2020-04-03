package sources

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/dropoutlabs/cape/connector/proto"
	"github.com/dropoutlabs/cape/primitives"
)

// PostgresSource represents a Postgres data source which can be queried
// against by a Cape data connector
//
// PostgresSource completes the Source interface
type PostgresSource struct {
	cfg    *Config
	source *primitives.Source
	pool   *pgxpool.Pool
}

// NewPostgresSource is a constructor for creating a PostgresSource.
//
// Completes the NewSourceFunc interface enabling it to be used by the Registry
func NewPostgresSource(ctx context.Context, cfg *Config, source *primitives.Source) (Source, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	if err := source.Validate(); err != nil {
		return nil, err
	}

	poolCfg, err := pgxpool.ParseConfig(source.Credentials.String())
	if err != nil {
		return nil, err
	}

	poolCfg.ConnConfig.RuntimeParams = map[string]string{
		"application_name": cfg.InstanceID.String(),
	}

	poolCfg.MinConns = 1
	poolCfg.MaxConns = 20
	poolCfg.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.ConnectConfig(ctx, poolCfg)
	if err != nil {
		return nil, err
	}

	return &PostgresSource{
		cfg:    cfg,
		source: source,
		pool:   pool,
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
	if q.Source() != p.source.Label {
		return nil, ErrWrongSource
	}

	// TODO: Do the work with postgres here to extract the schema information
	// we need to pull data out of pg and mutate it into our format.
	return nil, nil
}

// Query issues the given query against the targeted collection from the query.
// The schema should be the schema of the targeted collection.
//
// Results from the query are streamed back to the requester through the
// provided Stream
func (p *PostgresSource) Query(ctx context.Context, q Query, schema *proto.Schema, stream Stream) error {
	// TODO: Using the schema information and the query mutate the data coming
	// out of postgres into our grpc format and then push it over the wire via
	// the Stream.
	source := q.Source().String()
	if source != schema.DataSource {
		return ErrWrongSource
	}

	// XXX: How do you tell a stream that you are done?!
	return stream.Send(&proto.Record{
		Schema: schema,
	})
}

// Close issues a close to cancel all on-going requests and closes any
// connections to the underlying postgres data source
func (p *PostgresSource) Close(ctx context.Context) error {
	if p.pool == nil {
		return nil
	}

	p.pool.Close()
	p.pool = nil

	return nil
}

// init registers the PostgresSource with the global sources registry
func init() {
	err := registry.Register(primitives.PostgresType, NewPostgresSource)
	if err != nil {
		panic("Could not register source: " + err.Error())
	}
}
