package sources

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/capeprivacy/cape/connector/proto"
	"github.com/capeprivacy/cape/primitives"
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

func (p *PostgresSource) Schema(ctx context.Context, collection primitives.Collection) ([]*proto.Schema, error) {
	sql := "select table_name, column_name, data_type, character_maximum_length from " +
		"information_schema.columns where table_name "
	args := collection.String()

	switch collection {
	case primitives.Collection(primitives.Star):
		sql += "in (select table_name from information_schema.tables where table_schema = $1) order by table_name"
		args = "public"
	default:
		sql += "= $1"
		// args already set
	}

	rows, err := p.pool.Query(ctx, sql, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schemas []*proto.Schema
	var fields []*proto.FieldInfo
	var atTable string
	idx := -1

	for rows.Next() {
		var tableName string
		var columnName string
		var dataType string

		// maxLength is a pointer because it can be null
		var maxLength *int64
		err := rows.Scan(&tableName, &columnName, &dataType, &maxLength)
		if err != nil {
			return nil, err
		}

		if tableName != atTable {
			atTable = tableName
			idx++

			fields = []*proto.FieldInfo{}
			schema := &proto.Schema{
				Type:       proto.RecordType_DOCUMENT,
				DataSource: p.source.Label.String(),
				Target:     tableName,
			}

			schemas = append(schemas, schema)
		}

		f, err := DataTypeToProtoField(p.Type(), dataType)
		if err != nil {
			return nil, err
		}

		size := f.Size
		if maxLength != nil {
			size = *maxLength
		}

		field := &proto.FieldInfo{
			Field: f.Field,
			Size:  size,
			Name:  columnName,
		}

		fields = append(fields, field)
		schemas[idx].Fields = fields
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return schemas, nil
}

// Query issues the given query against the targeted collection from the query.
// The schema should be the schema of the targeted collection.
//
// Results from the query are streamed back to the requester through the
// provided Stream
func (p *PostgresSource) Query(ctx context.Context, stream Stream, q Query, limit int64, offset int64) error {
	qStr, params := q.Raw()

	qStr = fmt.Sprintf("%s LIMIT %d OFFSET %d", qStr, limit, offset)

	// If there are no params, pgx will error if you pass anything (even nil!)
	var rows pgx.Rows
	var err error
	if len(params) == 0 {
		rows, err = p.pool.Query(ctx, qStr)
	} else {
		rows, err = p.pool.Query(ctx, qStr, params...)
	}
	if err != nil {
		return err
	}

	defer rows.Close()

	i := 0
	for rows.Next() {
		record := &proto.Record{}
		if i == 0 {
			fields, err := fieldsFromFieldDescription(rows.FieldDescriptions())
			if err != nil {
				return err
			}

			record.Schema = &proto.Schema{
				DataSource: p.Label().String(),
				Target:     q.Collection(),
				Type:       proto.RecordType_DOCUMENT,
				Fields:     fields,
			}
		}

		pgVals, err := rows.Values()
		if err != nil {
			return err
		}

		fields, err := PostgresEncode(pgVals)
		if err != nil {
			return err
		}

		record.Fields = fields

		// When the grpc connection is closed grpc calls
		// cancel on their context but this Send call also
		// returns an error and the rows get closed properly
		// in the defer above.
		err = stream.Send(record)
		if err != nil {
			return err
		}

		i++
	}

	return rows.Err()
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
