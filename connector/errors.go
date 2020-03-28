package connector

import (
	errors "github.com/dropoutlabs/cape/partyerrors"
)

var (
	InvalidConfigCause = errors.NewCause(errors.BadRequestCategory, "invalid_config")
)
