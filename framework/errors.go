package framework

import (
	errors "github.com/dropoutlabs/cape/partyerrors"
)

var (
	// InvalidAuthHeader occurs when the auth header is in the wrong format
	InvalidAuthHeader = errors.NewCause(errors.BadRequestCategory, "invalid_auth_header")
)
