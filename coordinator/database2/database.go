package database2

import (
	"context"
	"io"
	"math/rand"
	"net/url"
	"sync"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/oklog/ulid"
)

type Database struct {
	url          *url.URL
	cfg          *pgxpool.Config
	Pool         *pgxpool.Pool
	entropy      io.Reader
	entropyMutex sync.Mutex
}

func NewDatabase(dbURL *url.URL, name string) (*Database, error) {
	cfg, err := pgxpool.ParseConfig(dbURL.String())
	if err != nil {
		return nil, err
	}

	// Set the application name which can be used for identifying which service
	// is connecting to postgres
	cfg.ConnConfig.RuntimeParams = map[string]string{
		"application_name": name,
	}

	t := time.Now().UnixNano()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t)), 0)

	return &Database{
		url:     dbURL,
		cfg:     cfg,
		entropy: entropy,
	}, nil
}

func (d *Database) Open(ctx context.Context) error {
	c, err := pgxpool.ConnectConfig(ctx, d.cfg)
	if err != nil {
		return err
	}

	d.Pool = c
	return nil
}

func (d *Database) Close() error {
	if d.Pool == nil {
		return nil
	}

	d.Pool.Close()
	d.Pool = nil

	return nil
}

func (d *Database) GetID() ulid.ULID {
	t := time.Now()

	d.entropyMutex.Lock()
	id := ulid.MustNew(ulid.Timestamp(t), d.entropy)
	d.entropyMutex.Unlock()

	return id
}
