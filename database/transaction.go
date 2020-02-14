package database

import "context"

// Transaction is an atomic operation inside a database
type Transaction interface {
	Begin(context.Context) error
	Commit(context.Context) error
	Rollback(context.Context) error
}