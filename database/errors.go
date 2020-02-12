package database

import errors "github.com/dropoutlabs/privacyai/partyerrors"

var (
	// NotImplementedDBCause happens when you try to connect to a db we do not support
	NotImplementedDBCause = errors.NewCause(errors.NotImplementedCategory, "not_implemented")
)
