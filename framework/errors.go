package framework

import (
	errors "github.com/capeprivacy/cape/partyerrors"
)

var (
	// InvalidParametersCause happens when you pass invalid input
	InvalidParametersCause = errors.NewCause(errors.BadRequestCategory, "invalid_input_parameters")

	// NotFoundCause occurs when an entity we were looking for could not be found
	NotFoundCause = errors.NewCause(errors.BadRequestCategory, "not_found")

	BadJSONCause = errors.NewCause(errors.BadRequestCategory, "bad_json_cause")
)
