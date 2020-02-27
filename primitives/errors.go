package primitives

import (
	errors "github.com/dropoutlabs/privacyai/partyerrors"
)

var (
	// InvalidTimeCause occurs when a time is provided that outside the range
	// of possible times
	InvalidTimeCause = errors.NewCause(errors.BadRequestCategory, "invalid_time")
)
