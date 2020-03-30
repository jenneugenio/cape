package connector

import (
	errors "github.com/dropoutlabs/cape/partyerrors"
)

var (
	MissingAuthToken   = errors.NewCause(errors.BadRequestCategory, "missing_auth_token")
	InvalidConfigCause = errors.NewCause(errors.BadRequestCategory, "invalid_config")
)
