package framework

import (
	errors "github.com/dropoutlabs/cape/partyerrors"
)

var (
	// AuthenticationFailure is caused by authentication failing
	AuthenticationFailure = errors.NewCause(errors.UnauthorizedCategory, "authentication_failure")

	// AuthorizationFailure is caused by authorization not being given
	AuthorizationFailure = errors.NewCause(errors.UnauthorizedCategory, "authorization_failure")

	// ErrAuthentication is the error wrapping the AuthenticationFailure cause
	ErrAuthentication = errors.New(AuthenticationFailure, "Failed to authenticate")

	// ErrAuthorization is the error wrapping the AuthorizationFailure cause
	ErrAuthorization = errors.New(AuthorizationFailure, "Access denied")
)
