package policy

import errors "github.com/capeprivacy/cape/partyerrors"

var (
	AccessDeniedCause = errors.NewCause(errors.BadRequestCategory, "access_denied")
)
