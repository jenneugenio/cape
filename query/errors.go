package query

import (
	errors "github.com/dropoutlabs/cape/partyerrors"
)

var (
	InvalidQueryCause = errors.NewCause(errors.BadRequestCategory, "invalid_query")
)
