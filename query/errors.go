package query

import (
	errors "github.com/capeprivacy/cape/partyerrors"
)

var (
	InvalidQueryCause     = errors.NewCause(errors.BadRequestCategory, "invalid_query")
	NoPossibleFieldsCause = errors.NewCause(errors.BadRequestCategory, "no_possible_fields")
)
