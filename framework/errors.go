package framework

import (
	errors "github.com/dropoutlabs/cape/partyerrors"
)

var (
	// AuthenticationFailure is caused by authentication failing
	AuthenticationFailure = errors.NewCause(errors.UnauthorizedCategory, "authentication_failure")

	// ErrAuthentication is the error wrapping the AuthenticationFailure cause
	ErrAuthentication = errors.New(AuthenticationFailure, "Failed to authenticate")
)
