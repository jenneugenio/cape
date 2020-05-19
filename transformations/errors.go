package transformations

import (
	errors "github.com/capeprivacy/cape/partyerrors"
)

var (
	// UnsupportedType happens when a transform encounters a type it does not support
	UnsupportedType = errors.NewCause(errors.BadRequestCategory, "unsupported_type")
)
