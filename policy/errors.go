package policy

import errors "github.com/dropoutlabs/cape/partyerrors"

var (
	AccessDeniedCause = errors.NewCause(errors.BadRequestCategory, "access_denied")
)
