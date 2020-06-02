package transformations

import (
	errors "github.com/capeprivacy/cape/partyerrors"
)

var (
	// UnsupportedType happens when a transform encounters a type it does not support
	UnsupportedType = errors.NewCause(errors.BadRequestCategory, "unsupported_type")
	// MissingArgument happens when an argument is missing to apply the transformation
	MissingArgument = errors.NewCause(errors.BadRequestCategory, "missing_arg")
	// WrongArgument happens when a wrong combination of argument is apply to the transformation
	WrongArgument = errors.NewCause(errors.BadRequestCategory, "wrong_arg")
)
