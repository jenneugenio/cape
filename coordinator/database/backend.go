package database

import (
	"context"
	"encoding/json"
	"github.com/capeprivacy/cape/prims"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/url"

	sq "github.com/Masterminds/squirrel"
	"github.com/capeprivacy/cape/coordinator/database/crypto"
	errors "github.com/capeprivacy/cape/partyerrors"
)

type Database interface {
	Connect(context.Context) error

	CreateProject(context.Context, *prims.Project) error
}

func NewPgxDatabase(connStr string) Database {
	return &PgxDatabase{
		connStr: connStr,
	}
}

type PgxDatabase struct {
	connStr string
	pool *pgxpool.Pool
}

func (p *PgxDatabase) Connect(ctx context.Context) error {
	pool, err := pgxpool.Connect(ctx, p.connStr)
	if err != nil {
		return err
	}

	p.pool = pool
	return nil
}

func (p *PgxDatabase) CreateProject(ctx context.Context, project *prims.Project) error {
	data, err := json.Marshal(project)
	if err != nil {
		return err
	}

	sql, args, err := sq.Insert("projects").
		Columns("data").
		Values(data).
		ToSql()

	if err != nil {
		return err
	}

	_, err = p.pool.Exec(ctx, sql, args)
	return err
}

// Backend represents a storage backend (e.g. Postgres, MySQL, etc).
// See a concrete implementation of this interface (e.g. PostgresBackend) for
// more details.
type Backend interface {
	Querier
	Open(context.Context) error
	Close() error
	Transaction(context.Context) (Transaction, error)
	URL() *url.URL
	SetEncryptionCodec(crypto.EncryptionCodec)
}

// NewBackendFunc represents a constructor of a Backend implementation
type NewBackendFunc func(*url.URL, string) (Backend, error)

var validDBs = map[string]NewBackendFunc{
	"postgres": NewPostgresBackend,
}

// New returns a new backend for the given application name
func New(dbURL *url.URL, appName string) (Backend, error) {
	ctor, ok := validDBs[dbURL.Scheme]
	if !ok {
		return nil, errors.New(errors.NotImplementedCause, "database not supported")
	}

	return ctor(dbURL, appName)
}