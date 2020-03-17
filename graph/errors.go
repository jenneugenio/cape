package graph

import (
	errors "github.com/dropoutlabs/cape/partyerrors"
)

var (
	AuthenticationFailure = errors.NewCause(errors.UnauthorizedCategory, "authentication_failure")
	AuthenticationError   = errors.New(AuthenticationFailure, "Failed to authenticate")

	RouteNotImplemented = errors.NewCause(errors.NotImplementedCategory, "route_not_implemented")
)
