# Cape Database

This package encapsulates the Cape database. This is where we store our internal data such as users, tokens, etc.

## Overview

### `Backend`

A database is represented by a `Backend`.  This is probably the most important class and where you should start.
`Backend` is a high level interface that specific databases can implement (e.g. `PostgresBackend`)

### `Entity`

All objects we store in the data layer must satisfy the `Entity` interface. `Entity` represents any primitive data
structure stored inside the Coordinator.

### `Primitive`

Most entities will embed `Primitive`. See `User` for an example. Note that since the database is decoupled from
the app itself, concrete primitives are not declared in this package, but in the `primitives` package (e.g. `User` or `Policy`)

When you are designing in the database you should never leak any other system concepts into this package. Think of this
package that could be pulled out of this repository and shipped on its own in a different system.

### Migrations

#### Postgres

Postgres database migrations are managed through [tern](https://github.com/jackc/tern). You should familiarize yourself with that
package if you are going to migrate the database.  Instructions to write a new migration can be found [here](https://github.com/jackc/tern#migrations).

Postgres is currently the only database we support.

## Testing

You can leverage the `dbtest` package when writing your tests. This gives you access to a "prod looking" database.
Using this will look something like

```go
testDB, err := dbtest.New(os.Getenv("CAPE_DB_URL"))
gm.Expect(err).To(gm.BeNil())

migrations := []string{
    os.Getenv("CAPE_DB_MIGRATIONS"),
    os.Getenv("CAPE_DB_TEST_MIGRATIONS"),
}

migrator, err := NewMigrator(testDB.URL(), migrations...)
gm.Expect(err).To(gm.BeNil())

db, err := dbConnect(ctx, testDB)

// db is a backend!
```

See `postgres_backend_test.go` for more details.
