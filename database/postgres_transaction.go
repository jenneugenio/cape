package database

import (
	"context"

	"github.com/jackc/pgx/v4"

	errors "github.com/dropoutlabs/privacyai/partyerrors"
)

// PostgresTransaction implements the Transaction interface
type PostgresTransaction struct {
	*postgresQuerier
	tx pgx.Tx
}

// Commit sends the commit transaction command to postgres
func (t *PostgresTransaction) Commit(ctx context.Context) error {
	// handles returning the connection back to the pool
	return convertPgTxError(t.tx.Commit(ctx))
}

// Rollback sends the rollback transaction command to postgres
func (t *PostgresTransaction) Rollback(ctx context.Context) error {
	// handles returning the connection back to the pool
	return convertPgTxError(t.tx.Rollback(ctx))
}

func convertPgTxError(err error) error {
	switch err {
	case pgx.ErrTxClosed, pgx.ErrTxCommitRollback:
		return errors.New(ClosedCause, "Connection closed")
	default:
		return err
	}
}
