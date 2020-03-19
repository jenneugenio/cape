package framework

import (
	errors "github.com/dropoutlabs/cape/partyerrors"
)

var (
	// AuthenticationFailure is caused by authentication failing
	AuthenticationFailure = errors.NewCause(errors.UnauthorizedCategory, "authentication_failure")

	// ErrAuthentication is the error wrapping the AuthenticationFailure cause
	ErrAuthentication = errors.New(AuthenticationFailure, "Failed to authenticate")

	// InvalidAuthHeader occurs when the auth header is in the wrong format
	InvalidAuthHeader = errors.NewCause(errors.BadRequestCategory, "invalid_auth_header")
)
