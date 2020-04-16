package types

import (
	errors "github.com/capeprivacy/cape/partyerrors"
)

var (
	// UnknownTypeCause occurs when a provided type is not registered
	UnknownTypeCause = errors.NewCause(errors.BadRequestCategory, "unknown_type")
)
