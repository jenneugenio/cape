package database

import "context"

// Transaction is an atomic operation inside a database
type Transaction interface {
	Querier
	Commit(context.Context) error
	Rollback(context.Context) error
}
