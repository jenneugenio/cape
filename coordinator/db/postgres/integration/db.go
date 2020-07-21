package integration

import (
	"context"
	"fmt"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/dbtest"
	capepg "github.com/capeprivacy/cape/coordinator/db/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/url"
	"os"
)

type TestDB struct {
	Pool      capepg.Pool
	RootDBURL *url.URL
	TestDBURL *url.URL
	DBName    string
}

func (t *TestDB) Setup(ctx context.Context) error {
	db, err := pgxpool.Connect(ctx, t.RootDBURL.String())
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", t.DBName))
	if err != nil {
		return err
	}

	// okay now we can create our long lasting connection
	db, err = pgxpool.Connect(ctx, t.TestDBURL.String())
	if err != nil {
		return err
	}

	t.Pool = db

	migrations := []string{
		os.Getenv("CAPE_DB_MIGRATIONS"),
		os.Getenv("CAPE_DB_TEST_MIGRATIONS"),
	}

	migrator, err := database.NewMigrator(t.TestDBURL, migrations...)
	if err != nil {
		return err
	}

	return migrator.Up(ctx)
}

func (t *TestDB) Teardown(ctx context.Context) error {
	// TODO -- what happens to t.Pool?
	db, err := pgxpool.Connect(ctx, t.RootDBURL.String())
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec(ctx, fmt.Sprintf("DROP DATABASE %s", t.DBName))
	return err
}

func CreateTestDB() (*TestDB, error) {
	rootURL, err := url.Parse(os.Getenv("CAPE_DB_URL"))
	if err != nil {
		return nil, err
	}

	dbURL, err := url.Parse(os.Getenv("CAPE_DB_URL"))
	if err != nil {
		return nil, err
	}

	dbName, err := dbtest.GenerateName()
	if err != nil {
		return nil, err
	}

	dbURL.Path = dbName

	return &TestDB{
		RootDBURL: rootURL,
		TestDBURL: dbURL,
		DBName:    dbName,
	}, nil
}
