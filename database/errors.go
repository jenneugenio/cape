package database

import errors "github.com/dropoutlabs/privacyai/partyerrors"

var (
	// InvalidDBURLCause happens you pass a poorly formatted request string
	InvalidDBURLCause = errors.NewCause(errors.BadRequestCategory, "invalid_db_url")

	// NotImplementedDBCause happens when you try to connect to a db we do not support
	NotImplementedDBCause = errors.NewCause(errors.NotImplementedCategory, "not_implemented")
)
